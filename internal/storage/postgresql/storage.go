package postgresql

import (
	"context"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/model"
)

// NewDB creates new connect to db
func NewDB(ctx context.Context, dsn string) (*Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return newDatabase(pool), nil
}

// GetUser – ...
func (d *Database) GetUser(ctx context.Context, email string) (user model.DBUser, err error) {
	q := `
	SELECT id, email, pass_hash
	FROM users
	WHERE email = $1;`

	err = d.Get(ctx, &user, q, email)
	if err != nil {
		return user, errorpkg.WrapErr(err, "can't get user")
	}

	return user, nil
}

// RegisterUser – ...
func (d *Database) RegisterUser(ctx context.Context, email string, passwordHash []byte) (userID int64, err error) {
	q := `
	INSERT INTO users(email,pass_hash)
	VALUES($1, $2);`

	_, err = d.Exec(ctx, q, email, passwordHash)
	if err != nil {
		return 0, errorpkg.WrapErr(err, "can't register user")
	}

	return userID, nil
}
