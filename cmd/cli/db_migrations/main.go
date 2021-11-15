package main

import (
	"database/sql"

	"ctm_lk/internal/config"
	"ctm_lk/pkg/logger"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
)

func main() {
	logger.NewLogs()
	p := "Миграции базы данных:"
	logger.Info(p, "Старт")
	config.NewConfig()
	logger.Info(p, "Подключение к БД")
	db, err := sql.Open("postgres", config.Cfg.DBConnString())
	if err != nil {
		logger.Error(p, err)
	}
	//of
	defer db.Close()
	// setup database
	logger.Info(p, "Применение миграций")
	if err := goose.Up(db, "../../../db/migrations"); err != nil {
		logger.Error(p, err)
	}
	logger.Info(p, "Завершение") // run app
}
