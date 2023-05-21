package repository

import (
	"database/sql"
	"errors"
	"lifeAchieve/entity"
	"lifeAchieve/utils"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func FindUserByEmail(email string) (*entity.User, error) {
	stmt, err := DB.Prepare("select bin_to_uuid(id), email, first_name, last_name, password from user where email = ? order by email asc limit 1;")
	if err != nil {
		return nil, err
	}

	var user entity.User

	err = stmt.QueryRow(email).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserById(id string) (*entity.User, error) {
	stmt, err := DB.Prepare("SELECT user.first_name, user.last_name, user.email FROM user WHERE user.id = uuid_to_bin(?)")
	if err != nil {
		return nil, err
	}

	var user entity.User

	err = stmt.QueryRow(id).Scan(&user.FirstName, &user.LastName, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserNamesById(id string) (*entity.User, error) {
	stmt, err := DB.Prepare("SELECT first_name, last_name from user where id = uuid_to_bin(?)")
	if err != nil {
		return nil, err
	}

	var user entity.User

	err = stmt.QueryRow(id).Scan(&user.FirstName, &user.LastName)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserTypeById(id string) (string, error) {
	stmt, err := DB.Prepare("SELECT type FROM user WHERE id = uuid_to_bin(?)")
	if err != nil {
		return "", err
	}

	var userType string

	err = stmt.QueryRow(id).Scan(&userType)

	if err != nil {
		return "", err
	}

	return userType, nil
}

func InsertUser(user entity.PostUser) error {
	stmt, err := DB.Prepare("INSERT INTO user (id, first_name, last_name, email, password) values (uuid_to_bin(?), ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	id := uuid.New()

	_, err = stmt.Exec(id, user.FirstName, user.LastName, user.Email, user.Password)

	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(id string, user entity.User) error {

	update := "UPDATE user SET"

	userObject := reflect.ValueOf(user)
	userKeys := userObject.Type()

	var key string
	var val interface{}

	var valSlice []interface{}

	for i := 0; i < userObject.NumField(); i++ {
		key = userKeys.Field(i).Name
		key = utils.ToSnakeCase(key)
		val = userObject.Field(i).Interface()
		if !reflect.ValueOf(val).IsZero() && key != "id" {
			update += " " + key + " = ?,"
			valSlice = append(valSlice, val)
		}
	}

	update = strings.TrimSuffix(update, ",")

	update += " WHERE id = uuid_to_bin(?)"
	valSlice = append(valSlice, id)

	stmt, err := DB.Prepare(update)
	if err != nil {
		return err
	}
	var result sql.Result
	result, err = stmt.Exec(valSlice...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if the entity exists
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM user WHERE id = uuid_to_bin(?)", id).Scan(&count)
		if err != nil {
			return err
		}
		if count != 0 {
			return errors.New("no changes")
		} else {
			return errors.New("not found")
		}
	}
	return nil
}
