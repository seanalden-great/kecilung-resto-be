package domain

type Menu struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  uint     `json:"category_id"` // Menyimpan ID kategori
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"` // Relasi objek
}

type MenuRepository interface {
	FetchAll() ([]Menu, error)
	FindByID(id uint) (Menu, error)
	Store(m *Menu) error
	Update(m *Menu) error
	Delete(id uint) error
}

type MenuUsecase interface {
	GetAll() ([]Menu, error)
	GetByID(id uint) (Menu, error)
	Create(m *Menu) error
	Update(id uint, m *Menu) error
	Delete(id uint) error
}