package usecase

import "github.com/seanalden-great/kecilung-resto-be/domain"

type articleUsecase struct {
	repo domain.ArticleRepository
}

func NewArticleUsecase(repo domain.ArticleRepository) domain.ArticleUsecase {
	return &articleUsecase{repo}
}

func (u *articleUsecase) GetAll() ([]domain.Article, error) {
	return u.repo.GetAll()
}

func (u *articleUsecase) GetByID(id uint) (domain.Article, error) {
	return u.repo.GetByID(id)
}

func (u *articleUsecase) Create(article *domain.Article) error {
	return u.repo.Create(article)
}

func (u *articleUsecase) Update(id uint, article *domain.Article) error {
	return u.repo.Update(id, article)
}

func (u *articleUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}