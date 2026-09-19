package usecase

import (
	"errors"
	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type momentUsecase struct {
	repo domain.MomentRepository
}

func NewMomentUsecase(r domain.MomentRepository) domain.MomentUsecase {
	return &momentUsecase{repo: r}
}

func (u *momentUsecase) GetAllMoments() ([]domain.Moment, error) { return u.repo.FetchAllMoments() }
func (u *momentUsecase) GetMomentByID(id uint) (domain.Moment, error) { return u.repo.FindMomentByID(id) }
func (u *momentUsecase) CreateMoment(m *domain.Moment) error { return u.repo.CreateMoment(m) }

func (u *momentUsecase) UpdateMoment(id uint, m *domain.Moment) error {
	m.ID = id
	return u.repo.UpdateMoment(m)
}

func (u *momentUsecase) DeleteMoment(id uint) error { return u.repo.DeleteMoment(id) }

func (u *momentUsecase) GetAllBookings() ([]domain.MomentBooking, error) { return u.repo.FetchBookings() }

func (u *momentUsecase) CreateBooking(b *domain.MomentBooking) error {
	exists, err := u.repo.CheckApprovedBookingExists(b.BookingDate)
	if err != nil { return err }
	if exists {
		return errors.New("Mohon maaf, jadwal pada tanggal tersebut sudah penuh dipesan.")
	}
	b.Status = "PENDING"
	return u.repo.CreateBooking(b)
}

func (u *momentUsecase) ApproveBooking(id uint) error { return u.repo.UpdateBookingStatus(id, "APPROVED") }
func (u *momentUsecase) RejectBooking(id uint) error { return u.repo.UpdateBookingStatus(id, "REJECTED") }