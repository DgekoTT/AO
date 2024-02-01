package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"yourTeamAuth/internal/config"
	"yourTeamAuth/internal/domain/models"
	"yourTeamAuth/storage"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(cfg *config.DatabaseConfig) (*Postgres, error) {
	const op = "storage.postgres.New"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %v\n", err)
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Stop() error {
	return p.db.Close()
}

func (p *Postgres) SaveUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	const op = "auth.postgres.saveUser"

	tx, err := p.db.BeginTxx(ctx, nil)

	if err != nil {

		return "0", fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback()
			if errRb != nil {
				err = fmt.Errorf("error during rollback: %s", errRb)
				return
			}
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(ctx, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return "0", fmt.Errorf("set transaction isolation level: %s", err)
	}

	return p.SaveUserTX(ctx, tx, email, hashedPassword)

}

func (p *Postgres) SaveUserTX(ctx context.Context, tx *sqlx.Tx, email string, hashedPassword []byte) (string, error) {
	const op = "auth.postgres.SaveUserTX"

	userExist, err := p.CheckEmail(ctx, tx, email)
	if err != nil {
		return "0", fmt.Errorf("%s: %w", op, err)
	}

	if userExist {
		return "0", storage.ErrUserExists
	}

	userID, err := p.CreateUser(ctx, tx, email, "email", hashedPassword)
	if err != nil {
		return "0", err
	}

	// доделать роли, верификацию и стороннюю регистрацию

	return userID, nil

}

func (p *Postgres) CheckEmail(ctx context.Context, tx *sqlx.Tx, email string) (bool, error) {
	var userExist bool

	err := tx.QueryRowxContext(ctx, `SELECT EXISTS(SELECT ID FROM users WHERE email = $1) AS user_exist`, email).Scan(&userExist)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %s", err)
	}

	return userExist, nil
}

func (p *Postgres) CreateUser(ctx context.Context, tx *sqlx.Tx, email, provider string, hashedPassword []byte) (string, error) {

	var userID string

	err := tx.QueryRowContext(ctx, `INSERT INTO users (email, registrationProvider, hashedPassword) VALUES ($1, $2, $3) RETURNING ID`, email, provider, hashedPassword).Scan(&userID)
	if err != nil {
		return "0", fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil

}

func (p *Postgres) User(ctx context.Context, email string) (models.Users, error) {
	var user models.Users

	err := p.db.GetContext(ctx, &user,
		`SELECT u.*, r.Role, v.verificationToken, v.isVerified, o.accountID 
              FROM users u 
              LEFT JOIN user_roles ur ON u.ID = ur.UserID 
              LEFT JOIN roles r ON ur.RoleID = r.RoleID 
              LEFT JOIN verify_infos v ON u.ID = v.UserID 
              LEFT JOIN oauth_providers o ON u.ID = o.UserID 
              WHERE u.Email = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Users{}, fmt.Errorf("no user found with email: %s", email)
		}
		return models.Users{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func (p *Postgres) Logout(ctx context.Context, accessToken string, refreshToken string) (bool, string, error) {
	return true, accessToken, nil
}
