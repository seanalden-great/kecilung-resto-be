package domain

import "time"

// Article struct untuk tabel articles
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      []ArticleImage `json:"images" gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ArticleImage struct untuk multi-gambar
type ArticleImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleRepository interface
type ArticleRepository interface {
	GetAll() ([]Article, error)
	GetByID(id uint) (Article, error)
	Create(article *Article) error
	Update(id uint, article *Article) error
	Delete(id uint) error
}

// ArticleUsecase interface
type ArticleUsecase interface {
	GetAll() ([]Article, error)
	GetByID(id uint) (Article, error)
	Create(article *Article) error
	Update(id uint, article *Article) error
	Delete(id uint) error
}