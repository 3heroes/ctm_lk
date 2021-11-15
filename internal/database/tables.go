package database

import (
	"context"

	"ctm_lk/pkg/logger"
)

func (s *serverDB) createTables(ctx context.Context) {
	tx, err := s.DB.Begin()
	if err != nil {
		logger.Panic("Ошибка создания таблиц", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (
								    id SERIAL PRIMARY KEY,
									user_name VARCHAR(50) UNIQUE,
									user_password VARCHAR(36),
									user_key VARCHAR(36),
									user_token VARCHAR(36),
									date_add TIMESTAMPTZ(0) default (NOW() at time zone 'UTC+3'))
	`)
	if err != nil {
		logger.Panic("Ошибка создания таблиц", err)
	}

	tx.Commit()
}
