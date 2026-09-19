package usecase

import "github.com/seanalden-great/kecilung-resto-be/domain"

type menuUsecase struct {
	repo domain.MenuRepository
}

func NewMenuUsecase(r domain.MenuRepository) domain.MenuUsecase {
	return &menuUsecase{repo: r}
}

func (u *menuUsecase) GetAll() ([]domain.Menu, error) {
	return u.repo.FetchAll()
}

func (u *menuUsecase) GetByCategoryID(categoryID uint) ([]domain.Menu, error) {
	return u.repo.FetchByCategoryID(categoryID)
}

func (u *menuUsecase) GetByID(id uint) (domain.Menu, error) {
	return u.repo.FindByID(id)
}

func (u *menuUsecase) Create(m *domain.Menu) error {
	return u.repo.Store(m)
}

func (u *menuUsecase) Update(id uint, m *domain.Menu) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	m.ID = existing.ID
	return u.repo.Update(m)
}

func (u *menuUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}