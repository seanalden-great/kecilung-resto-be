package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type contactRepo struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) domain.ContactRepository {
	return &contactRepo{db}
}

func (r *contactRepo) Store(c *domain.ContactUs) error {
	return r.db.Create(c).Error
}

func (r *contactRepo) FetchAll() ([]domain.ContactUs, error) {
	var contacts []domain.ContactUs
	err := r.db.Order("created_at desc").Find(&contacts).Error
	return contacts, err
}