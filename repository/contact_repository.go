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

// Tambahkan fungsi ini di bagian bawah repository/contact_repository.go
func (r *contactRepo) UpdateReply(id uint, replyText string) error {
	return r.db.Model(&domain.ContactUs{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reply":      replyText,
		"is_replied": true,
	}).Error
}