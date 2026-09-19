package usecase

import (
	"errors"
	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type cateringUsecase struct {
	repo domain.CateringRepository
}

func NewCateringUsecase(r domain.CateringRepository) domain.CateringUsecase {
	return &cateringUsecase{repo: r}
}

func (u *cateringUsecase) GetGreeting() (domain.GreetingMessage, error) { return u.repo.GetGreeting() }
func (u *cateringUsecase) UpdateGreeting(msg *domain.GreetingMessage) error { return u.repo.UpdateGreeting(msg) }
func (u *cateringUsecase) GetAllCaterings() ([]domain.Catering, error) { return u.repo.FetchAllCaterings() }
func (u *cateringUsecase) GetCateringByID(id uint) (domain.Catering, error) { return u.repo.FindCateringByID(id) }
func (u *cateringUsecase) CreateCatering(c *domain.Catering) error { return u.repo.CreateCatering(c) }

// === TAMBAHAN BARU ===
func (u *cateringUsecase) UpdateCatering(id uint, c *domain.Catering) error {
	c.ID = id
	return u.repo.UpdateCatering(c)
}

func (u *cateringUsecase) DeleteCatering(id uint) error {
	return u.repo.DeleteCatering(id)
}
// =====================

func (u *cateringUsecase) GetAllBookings() ([]domain.Booking, error) { return u.repo.FetchBookings() }

func (u *cateringUsecase) CreateBooking(b *domain.Booking) error {
	// CEK BENTROK: Apakah tanggal ini sudah ada yang APPROVED?
	exists, err := u.repo.CheckApprovedBookingExists(b.BookingDate)
	if err != nil { return err }
	if exists {
		return errors.New("Mohon maaf, jadwal pada tanggal tersebut sudah penuh dipesan.")
	}
	b.Status = "PENDING"
	return u.repo.CreateBooking(b)
}

func (u *cateringUsecase) ApproveBooking(id uint) error {
	// Jika admin meng-approve, bisa saja ditambahkan validasi ulang di sini
	return u.repo.UpdateBookingStatus(id, "APPROVED")
}

func (u *cateringUsecase) RejectBooking(id uint) error {
	return u.repo.UpdateBookingStatus(id, "REJECTED")
}