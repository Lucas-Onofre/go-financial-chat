package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Lucas-Onofre/financial-chat/chat-service/user/dao"
)

type DB interface {
	Create(entity any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(_ context.Context, user dao.User) error {
	tx := r.db.Create(&user)
	return tx.Error
}

func (r *Repository) FindByUsername(_ context.Context, username string) (*dao.User, error) {
	var user dao.User

	tx := r.db.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}
