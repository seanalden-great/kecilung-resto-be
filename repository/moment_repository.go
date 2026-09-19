package repository

import (
	"time"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type momentRepo struct {
	db *gorm.DB
}

func NewMomentRepository(db *gorm.DB) domain.MomentRepository {
	return &momentRepo{db}
}

func (r *momentRepo) FetchAllMoments() ([]domain.Moment, error) {
	var moments []domain.Moment
	err := r.db.Preload("Images").Find(&moments).Error
	return moments, err
}

func (r *momentRepo) FindMomentByID(id uint) (domain.Moment, error) {
	var moment domain.Moment
	err := r.db.Preload("Images").First(&moment, id).Error
	return moment, err
}

func (r *momentRepo) CreateMoment(m *domain.Moment) error {
	return r.db.Create(m).Error
}

func (r *momentRepo) UpdateMoment(m *domain.Moment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Moment{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
			"name":        m.Name,
			"description": m.Description,
		}).Error; err != nil {
			return err
		}
		if len(m.Images) > 0 {
			if err := tx.Where("moment_id = ?", m.ID).Delete(&domain.MomentImage{}).Error; err != nil {
				return err
			}
			for i := range m.Images {
				m.Images[i].MomentID = m.ID
			}
			if err := tx.Create(&m.Images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *momentRepo) DeleteMoment(id uint) error {
	return r.db.Delete(&domain.Moment{}, id).Error
}

func (r *momentRepo) FetchBookings() ([]domain.MomentBooking, error) {
	var bookings []domain.MomentBooking
	err := r.db.Preload("Moment").Order("booking_date desc").Find(&bookings).Error
	return bookings, err
}

func (r *momentRepo) CreateBooking(b *domain.MomentBooking) error {
	return r.db.Create(b).Error
}

func (r *momentRepo) UpdateBookingStatus(id uint, status string) error {
	return r.db.Model(&domain.MomentBooking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *momentRepo) CheckApprovedBookingExists(date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&domain.MomentBooking{}).
		Where("DATE(booking_date) = DATE(?) AND status = ?", date, "APPROVED").
		Count(&count).Error
	return count > 0, err
}