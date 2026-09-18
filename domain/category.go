package domain

type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"size:50;uniqueIndex"` // Tambahkan size:50 di sini
	Name        string `json:"name" gorm:"size:255"`            // Menjadi varchar(255)
	Description string `json:"description"`                     // Tetap longtext
}

type CategoryRepository interface {
	FetchAll() ([]Category, error)
	FindByID(id uint) (Category, error)
	Store(c *Category) error
	Update(c *Category) error
	Delete(id uint) error
}

type CategoryUsecase interface {
	GetAll() ([]Category, error)
	GetByID(id uint) (Category, error)
	Create(c *Category) error
	Update(id uint, c *Category) error
	Delete(id uint) error
}