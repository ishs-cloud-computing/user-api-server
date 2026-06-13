package repository

import (
	"context"

	"github.com/ishs-cloud-computing/user-api/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*model.User, int, error) {
	query := `SELECT * FROM user_api.users`
	countQuery := `SELECT COUNT(*) FROM user_api.users`

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	var users []*model.User
	if err := r.db.SelectContext(ctx, &total, query); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint32) (*model.User, error) {
	query := `SELECT id, name, email, age FROM user_api.users WHERE id = ?`

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO user_api.users(name, email, age) VALUES(?, ?, ?)`,
		req.Name, req.Email, req.Age)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, uint32(id))
}

// TODO: Update

// TODO: Delete
