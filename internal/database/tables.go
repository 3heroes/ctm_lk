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

	// Создание таблицы на SQL Server
	var sSqlText = `IF OBJECT_ID(N'dbo.users', N'U') IS NULL
                 BEGIN
                  CREATE TABLE [dbo].[users](
	                           [id] [int] IDENTITY(1,1) NOT NULL,
	                           [user_name] [nchar](50) NOT NULL,
	                           [user_password] [nchar](36) NULL,
	                           [user_key] [nchar](36) NULL,
	                           [user_token] [nchar](36) NULL,
	                           [date_add] [datetime] NULL,
                   CONSTRAINT [PK_users] PRIMARY KEY CLUSTERED (	[id] ASC),
                   CONSTRAINT [AK_TransactionID] UNIQUE NONCLUSTERED ([user_name] ASC)
                   );
                   ALTER TABLE [dbo].[users] ADD  CONSTRAINT [DF_users]  DEFAULT (getdate()) FOR [date_add];
                  END`
	_, err = tx.ExecContext(ctx, sSqlText)

	//Создание таблицы на Postgres
	/*	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (
									    id SERIAL PRIMARY KEY,
										user_name VARCHAR(50) UNIQUE,
										user_password VARCHAR(36),
										user_key VARCHAR(36),
										user_token VARCHAR(36),
										date_add TIMESTAMPTZ(0) default (NOW() at time zone 'UTC+3'))
		`)
	*/
	if err != nil {
		logger.Panic("Ошибка создания таблиц", err)
	}

	tx.Commit()
}
