package user

import (
	"assignment2/entity"
	"context"
	"database/sql"
	"log"
)

type authRepositoryImpl struct {
	DB *sql.DB
}

type AuthRepository interface {
	StoreUser(ctx context.Context, data entity.User) (userID int64, err error)
	GetUserByEmail(ctx context.Context, username string) (data *entity.User, err error)
}

func ProvideAuthRepository(DB *sql.DB) *authRepositoryImpl {
	return &authRepositoryImpl{DB: DB}
}

var (
	INSERT_USER          = "INSERT INTO `user` (username, email, password) VALUES (?, ?, ?)"
	GET_USER_BY_USERNAME = "SELECT email, username, password FROM user WHERE username = ?"
)

func (a authRepositoryImpl) StoreUser(ctx context.Context, data entity.User) (userID int64, err error) {
	stmt, err := a.DB.PrepareContext(ctx, INSERT_USER)
	if err != nil {
		log.Printf("[StoreUser] failed to prepare the statement, err => %v", err)
		return
	}
	res, err := stmt.ExecContext(ctx, data.Username, data.Email, data.Password)
	if err != nil {
		log.Printf("[StoreUser] failed to store user, err => %v", err)
		return
	}
	userID, _ = res.LastInsertId()
	return
}

func (a authRepositoryImpl) GetUserByEmail(ctx context.Context, username string) (data *entity.User, err error) {
	stmt, err := a.DB.PrepareContext(ctx, GET_USER_BY_USERNAME)
	if err != nil {
		log.Printf("[GetUserByEmail] failed to prepare the statement, err => %v", err)
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		log.Printf("[GetUserByEmail] failed to query to the database, err => %v", err)
		return nil, err
	}

	for rows.Next() {
		us := entity.User{}

		err = rows.Scan(
			&us.Email,
			&us.Username,
			&us.Password,
		)
		if err != nil {
			log.Printf("[GetUserByEmail] failed to scan the data, err => %v", err)
			return nil, err
		}
		data = &us
	}
	return data, nil
}
