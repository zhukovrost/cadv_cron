package main

import (
	"github.com/zhukovrost/cadv_cron/internal/app"
	"github.com/zhukovrost/cadv_cron/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	err = app.Run(cfg)
	if err != nil {
		panic(err)
	}
}
