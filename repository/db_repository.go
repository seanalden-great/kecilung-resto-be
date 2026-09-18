package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type dbRepo struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *dbRepo {
	return &dbRepo{db}
}

// --- CRUD CATEGORY ---
func (r *dbRepo) FindCategories() ([]domain.Category, error) {
	var cats []domain.Category
	err := r.db.Find(&cats).Error
	return cats, err
}

func (r *dbRepo) FindCategoryByID(id string) (domain.Category, error) {
	var cat domain.Category
	// Preload("Menus") akan otomatis menarik semua menu yang terkait dengan kategori ini
	err := r.db.Preload("Menus").First(&cat, id).Error
	return cat, err
}

func (r *dbRepo) CreateCategory(cat *domain.Category) error {
	return r.db.Create(cat).Error
}

func (r *dbRepo) UpdateCategory(cat *domain.Category) error {
	return r.db.Save(cat).Error
}

func (r *dbRepo) DeleteCategory(id string) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

// --- CRUD MENU ---
func (r *dbRepo) FindMenus() ([]domain.Menu, error) {
	var menus []domain.Menu
	// Preload("Category") otomatis menarik data kategori dari foreign key CategoryID
	err := r.db.Preload("Category").Find(&menus).Error
	return menus, err
}

func (r *dbRepo) FindMenuByID(id string) (domain.Menu, error) {
	var menu domain.Menu
	err := r.db.Preload("Category").First(&menu, id).Error
	return menu, err
}

func (r *dbRepo) CreateMenu(menu *domain.Menu) error {
	return r.db.Create(menu).Error
}

func (r *dbRepo) UpdateMenu(menu *domain.Menu) error {
	return r.db.Save(menu).Error
}

func (r *dbRepo) DeleteMenu(id string) error {
	return r.db.Delete(&domain.Menu{}, id).Error
}