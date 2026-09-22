package usecase

import (
	"errors"
	"strings"

	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type contactUsecase struct {
	repo domain.ContactRepository
}

func NewContactUsecase(r domain.ContactRepository) domain.ContactUsecase {
	return &contactUsecase{repo: r}
}

func (u *contactUsecase) Create(c *domain.ContactUs) error {
	// Validasi Sederhana
	if strings.TrimSpace(c.FullName) == "" {
		return errors.New("nama lengkap tidak boleh kosong")
	}
	if strings.TrimSpace(c.Description) == "" {
		return errors.New("deskripsi pesan tidak boleh kosong")
	}
	return u.repo.Store(c)
}

func (u *contactUsecase) GetAll() ([]domain.ContactUs, error) {
	return u.repo.FetchAll()
}