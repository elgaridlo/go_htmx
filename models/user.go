package models

import (
	"database/sql"
	"fmt"
	"project/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	Name         string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

type ResultUser struct {
	Data       []*User
	Pagination *utils.PaginationResult
}

func (us UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
  SELECT id, password_hash
  FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) CreateUser(user User) (*User, error) {
	email := strings.ToLower(user.Email)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)
	postUser := User{
		Email:        email,
		Name:         user.Name,
		PasswordHash: string(passwordHash),
	}
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3) RETURNING id`, postUser.Email, postUser.PasswordHash, postUser.Name)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) UserMany(name string, pagination utils.PaginationParams) (*ResultUser, error) {
	var indexArg int = 0
	var conditions string
	var valArgs []interface{}
	if name != "" {
		indexArg++
		conditions = fmt.Sprintf(`AND u.name ilike $%d `, indexArg)
		valArgs = append(valArgs, "%"+name+"%")
	}

	tableClause := `users u`

	table := utils.TableService(*us)

	resultPagination, _ := table.GenericPagination(tableClause, conditions, valArgs, &pagination)
	valArgs = append(valArgs, pagination.PageSize, (pagination.Page-1)*pagination.PageSize)

	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.name,
			u.email
		FROM %s
		WHERE true %s
		limit $%d
		offset $%d`, tableClause, conditions, indexArg+1, indexArg+2)

	// Execute the query
	rows, err := us.DB.Query(query, valArgs...)
	if err != nil {
		return nil, fmt.Errorf("User Many: %w", err)
	}
	defer rows.Close()

	var userList []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("User Many: %w", err)
		}
		userList = append(userList, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("User Many: %w", err)
	}

	result := ResultUser{
		Data:       userList,
		Pagination: resultPagination,
	}

	return &result, nil
}

func (us *UserService) UserDetail(id string) (*User, error) {
	row, err := us.DB.Query(`
		SELECT id, name, email
		FROM users
		where id = $1 limit 1`, id)
	if err != nil {
		return nil, fmt.Errorf("UserDetail: %w", err)
	}

	defer row.Close()

	if row.Next() {
		var user User
		errScan := row.Scan(&user.ID, &user.Name, &user.Email)
		if errScan != nil {
			return nil, fmt.Errorf("UserDetail: %w", errScan)
		}

		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (us *UserService) UserUpdate(user User) (*User, error) {
	exist := us.DB.QueryRow(`
		SELECT id from users where id = $1
	`, user.ID)
	var existUser User

	err := exist.Scan(&existUser.ID)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	email := strings.ToLower(user.Email)
	resUser := User{
		Name:  user.Name,
		Email: email,
	}
	row := us.DB.QueryRow(`
		UPDATE users
		SET name=$1, email=$2, updated_at=$3
		WHERE id=$4  RETURNING id`,
		user.Name,
		email,
		time.Now(),
		user.ID,
	)

	err = row.Scan(&resUser.ID)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &resUser, nil
}

func (us *UserService) UserDelete(id int) error {
	exist := us.DB.QueryRow(`
		SELECT id from users where id = $1
	`, id)
	var existUser User

	err := exist.Scan(&existUser.ID)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	resUser := User{
		ID: id,
	}
	row := us.DB.QueryRow(`
		DELETE FROM users
		WHERE id=$1  RETURNING id`,
		id)

	err = row.Scan(&resUser.ID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
