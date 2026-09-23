package domain

import "time"

// ContactUs struct untuk tabel database 'contact_us'
type ContactUs struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Description string    `json:"description" gorm:"type:text"`
	Reply       string    `json:"reply" gorm:"type:text"`      // BARU: Isi balasan admin
	IsReplied   bool      `json:"is_replied" gorm:"default:false"` // BARU: Status balasan
	CreatedAt   time.Time `json:"created_at"`
}

// ContactRepository interface
type ContactRepository interface {
	Store(c *ContactUs) error
	FetchAll() ([]ContactUs, error) // Opsional: Untuk panel admin nanti
	// BARU: Fungsi untuk update balasan
	UpdateReply(id uint, replyText string) error
}

// ContactUsecase interface
type ContactUsecase interface {
	Create(c *ContactUs) error
	GetAll() ([]ContactUs, error) // Opsional: Untuk panel admin nanti
	// BARU: Fungsi untuk update balasan
	ReplyMessage(id uint, replyText string) error
}