// package domain

// type Menu struct {
// 	ID          uint     `json:"id" gorm:"primaryKey"`
// 	Name        string   `json:"name"`
// 	Description string   `json:"description"`
// 	Price       float64  `json:"price"`
// 	CategoryID  uint     `json:"category_id"` // Menyimpan ID kategori
// 	Category    Category `json:"category" gorm:"foreignKey:CategoryID"` // Relasi objek
// }

// type MenuRepository interface {
// 	FetchAll() ([]Menu, error)
// 	FindByID(id uint) (Menu, error)
// 	Store(m *Menu) error
// 	Update(m *Menu) error
// 	Delete(id uint) error
// }

// type MenuUsecase interface {
// 	GetAll() ([]Menu, error)
// 	GetByID(id uint) (Menu, error)
// 	Create(m *Menu) error
// 	Update(id uint, m *Menu) error
// 	Delete(id uint) error
// }

package domain

type Menu struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"image_url"` // Kolom baru untuk menyimpan link S3
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type MenuRepository interface {
	FetchAll() ([]Menu, error)
	FindByID(id uint) (Menu, error)
	FetchByCategoryID(categoryID uint) ([]Menu, error) // Tambahan baru
	Store(m *Menu) error
	Update(m *Menu) error
	Delete(id uint) error
}

type MenuUsecase interface {
	GetAll() ([]Menu, error)
	GetByID(id uint) (Menu, error)
	GetByCategoryID(categoryID uint) ([]Menu, error) // Tambahan baru
	Create(m *Menu) error
	Update(id uint, m *Menu) error
	Delete(id uint) error
}