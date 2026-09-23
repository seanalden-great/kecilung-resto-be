package domain

import "time"

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"` // Akan menyimpan Hashed Password
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthRepository interface {
	GetAdminByUsername(username string) (Admin, error)
	GetAdminByID(id uint) (Admin, error)
	UpdateAdmin(id uint, data *Admin) error
	CreateDefaultAdmin(admin *Admin) error // Opsional untuk Inisialisasi awal
}

type AuthUsecase interface {
	Login(username, password string) (string, Admin, error) // Returns Token JWT
	GetProfile(id uint) (Admin, error)
	UpdateProfile(id uint, data *Admin) error
}