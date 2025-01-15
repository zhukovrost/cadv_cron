package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zhukovrost/cadv_cron/internal/config"
	"github.com/zhukovrost/cadv_cron/internal/service"
	"github.com/zhukovrost/cadv_cron/pkg/mydb"
	logger "github.com/zhukovrost/cadv_logger"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) error {
	log := logger.New("", true)
	log.Info("Starting app...")

	log.Info("Connecting to the database...")
	db, err := mydb.Init(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := service.New(db, log)

	location, err := time.LoadLocation("Europe/Moscow") // Замените "Europe/Moscow" на нужный вам часовой пояс
	if err != nil {
		return fmt.Errorf("Ошибка загрузки часового пояса: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	log.Info("Starting cron...")
	c := cron.New(cron.WithLocation(location))

	id, err := c.AddFunc("0 0 * * *", func() {
		rows, err := srv.ClearTokens(ctx)
		if err != nil {
			log.Error("ClearTokens error", zap.Error(err))
			cancel()
			return
		}
		log.Debug("Database has been cleared", zap.Int64("rows_affected", rows))
	})

	if err != nil {
		return fmt.Errorf("Error occured while adding a new task to cron: %w", err)
	}
	log.Info("Clear tokens function added", zap.String("id", fmt.Sprintf("%v", id)))

	c.Start()
	defer c.Stop()

	select {
	case <-quit:
		log.Info("Stopping app...")
	case <-ctx.Done():
		log.Info("Stopping app because of the error")
	}

	return nil
}
