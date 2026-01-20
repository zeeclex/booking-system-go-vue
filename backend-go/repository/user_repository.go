package repository

import (
	"backend/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, role, email, password FROM users WHERE email = ?"
	err := r.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, role, email, password FROM users WHERE id = ?"
	err := r.DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users (excluding passwords for security)
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	query := "SELECT id, name, role, email FROM users ORDER BY name ASC"
	err := r.DB.Select(&users, query)
	return users, err
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	query := "INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)"
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Role)
	return err
}

// Update modifies an existing user's data
func (r *UserRepository) Update(user *models.User) error {
	query := "UPDATE users SET name = ?, email = ?, password = ?, role = ? WHERE id = ?"
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Role, user.ID)
	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	return err
}