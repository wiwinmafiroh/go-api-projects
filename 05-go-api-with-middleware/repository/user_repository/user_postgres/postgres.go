package user_postgres

import (
	"05-go-api-with-middleware/entity"
	"05-go-api-with-middleware/pkg/errs"
	"05-go-api-with-middleware/repository/user_repository"
	"database/sql"

	"github.com/lib/pq"
)

const (
	createUserQuery = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3);
	`

	getUserByEmailQuery = `
		SELECT id, name, email, password, role FROM users
		WHERE email = $1;
	`
)

type userPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) user_repository.UserRepository {
	return &userPostgres{
		db: db,
	}
}

func (u *userPostgres) CreateUser(userEntity entity.User) errs.ErrorMessage {
	_, err := u.db.Exec(createUserQuery, userEntity.Name, userEntity.Email, userEntity.Password)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return errs.NewBadRequestError("Email already exists")
		}

		return errs.NewInternalServerError("Failed to create new user")
	}

	return nil
}

func (u *userPostgres) GetUserByEmail(userEmail string) (*entity.User, errs.ErrorMessage) {
	row := u.db.QueryRow(getUserByEmailQuery, userEmail)

	var retrievedUser entity.User

	err := row.Scan(&retrievedUser.Id, &retrievedUser.Name, &retrievedUser.Email, &retrievedUser.Password, &retrievedUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("User not found")
		}

		return nil, errs.NewInternalServerError("Failed to retrieve user data")
	}

	return &retrievedUser, nil
}
