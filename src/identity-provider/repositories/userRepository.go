package repositories

import (
	"identity-provider/errors"
	"identity-provider/objects"

	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	Create(loyalty *objects.User) (*objects.User, error)
	Find(login string) (*objects.User, error)
	CheckCredentials(login, password string) (bool, error)
	Update(*objects.User) error
	Delete(login string) error
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repository *PostgresUserRepository) Create(user *objects.User) (*objects.User, error) {
	return user, repository.db.Create(user).Error
}

func (repository *PostgresUserRepository) Find(login string) (*objects.User, error) {
	temp := new(objects.User)
	err := repository.db.
		First(temp, "login = ?", login).
		Error

	switch err {
	case nil:
		return temp, err
	case gorm.ErrRecordNotFound:
		return nil, errors.RecordNotFound
	default:
		return nil, errors.UnknownError
	}
}

func (repository *PostgresUserRepository) CheckCredentials(login, password string) (bool, error) {
	temp := new(objects.User)
	err := repository.db.
		First(temp, "login = ?", login).
		First(temp, "password = ?", password).
		Error

	switch err {
	case nil:
		return true, err
	case gorm.ErrRecordNotFound:
		return false, errors.RecordNotFound
	default:
		return false, errors.UnknownError
	}
}

func (repository *PostgresUserRepository) Update(loyalty *objects.User) error {
	return repository.db.
		Save(loyalty).
		Error
}

func (repository *PostgresUserRepository) Delete(login string) error {
	record, err := repository.Find(login)
	if err != nil {
		return err
	}

	err = repository.db.
		Where(objects.User{Login: login}).
		Delete(record).
		Error

	switch err {
	case nil:
		return err
	case gorm.ErrRecordNotFound:
		return errors.RecordNotFound
	default:
		return errors.UnknownError
	}
}
