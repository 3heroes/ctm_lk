package database

import (
	"context"
	"database/sql"
	"time"

	"ctm_lk/internal/models"
	"ctm_lk/pkg/encription"
	"ctm_lk/pkg/logger"
)

type DBUserRepo struct {
	serverDB
}

func (db DBUserRepo) Get(ctx context.Context, u *models.User) (bool, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()
	q := `SELECT COALESCE(id, 0), user_name, user_token, user_password FROM users WHERE user_token=@Token OR user_name=@Name`
	row := db.QueryRowContext(ctx, q, sql.Named("Token", u.Token), sql.Named("Name", u.Login))
	if err := row.Scan(&u.ID, &u.Login, &u.Token, &u.Password); err != nil && err != sql.ErrNoRows {
		logger.Info(q, err)
		return false, err
	}
	if u.ID == 0 {
		return false, nil
	}
	if len(u.Token) == 0 {
		u.Token = encription.EncriptStr(u.Login)
		db.update(ctx, u)
	}
	return true, nil
}

func (db DBUserRepo) Add(ctx context.Context, u *models.User) error {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	u.Token = encription.EncriptStr(u.Login)

	q := `INSERT INTO users (user_name, user_password, user_token) OUTPUT INSERTED.ID VALUES (@Name,@Password,@Token)`
	row := db.QueryRowContext(ctx, q, sql.Named("Name", u.Login), sql.Named("Password", u.Password), sql.Named("Token", u.Token))

	if err := row.Scan(&u.ID); err != nil {
		logger.Info(q, err)
		return err
	}

	return nil
}

func (db DBUserRepo) update(ctx context.Context, u *models.User) bool {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	u.Token = encription.EncriptStr(u.Login)

	q := `UPDATE users SET user_name=$2, user_password=$3, user_token=$4 WHERE ID=$4`
	_, err := db.ExecContext(ctx, q, u.ID, u.Login, u.Password, u.Token)

	return err == nil
}

func (db DBUserRepo) Del(ctx context.Context, u *models.User) error {
	return nil
}

func (s serverDB) NewDBUserRepo() models.UsersRepo {
	ur := new(DBUserRepo)
	ur.serverDB = s
	return ur
}
