package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type menuRepo struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domain.MenuRepository {
	return &menuRepo{db}
}

func (r *menuRepo) FetchAll() ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.db.Preload("Category").Find(&menus).Error
	return menus, err
}

func (r *menuRepo) FindByID(id uint) (domain.Menu, error) {
	var menu domain.Menu
	err := r.db.Preload("Category").First(&menu, id).Error
	return menu, err
}

func (r *menuRepo) Store(m *domain.Menu) error {
	return r.db.Create(m).Error
}

func (r *menuRepo) Update(m *domain.Menu) error {
	return r.db.Save(m).Error
}

func (r *menuRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Menu{}, id).Error
}