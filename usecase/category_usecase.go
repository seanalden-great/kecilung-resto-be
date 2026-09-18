package usecase

import "github.com/seanalden-great/kecilung-resto-be/domain"

type categoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(r domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{repo: r}
}

func (u *categoryUsecase) GetAll() ([]domain.Category, error) {
	return u.repo.FetchAll()
}

func (u *categoryUsecase) GetByID(id uint) (domain.Category, error) {
	return u.repo.FindByID(id)
}

func (u *categoryUsecase) Create(c *domain.Category) error {
	return u.repo.Store(c)
}

func (u *categoryUsecase) Update(id uint, c *domain.Category) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	return u.repo.Update(c)
}

func (u *categoryUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}