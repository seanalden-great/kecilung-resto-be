package domain

import "time"

// ContactUs struct untuk tabel database 'contact_us'
type ContactUs struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

// ContactRepository interface
type ContactRepository interface {
	Store(c *ContactUs) error
	FetchAll() ([]ContactUs, error) // Opsional: Untuk panel admin nanti
}

// ContactUsecase interface
type ContactUsecase interface {
	Create(c *ContactUs) error
	GetAll() ([]ContactUs, error) // Opsional: Untuk panel admin nanti
}