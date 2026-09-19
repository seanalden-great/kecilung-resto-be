package repository

import (
	"time"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type cateringRepo struct {
	db *gorm.DB
}

func NewCateringRepository(db *gorm.DB) domain.CateringRepository {
	return &cateringRepo{db}
}

// -- Greeting --
func (r *cateringRepo) GetGreeting() (domain.GreetingMessage, error) {
	var msg domain.GreetingMessage
	err := r.db.FirstOrCreate(&msg, domain.GreetingMessage{ID: 1, Message: "Selamat Datang di Layanan Katering Kami"}).Error
	return msg, err
}

func (r *cateringRepo) UpdateGreeting(msg *domain.GreetingMessage) error {
	msg.ID = 1
	return r.db.Save(msg).Error
}

// -- Catering --
func (r *cateringRepo) FetchAllCaterings() ([]domain.Catering, error) {
	var caterings []domain.Catering
	err := r.db.Preload("Images").Find(&caterings).Error
	return caterings, err
}

func (r *cateringRepo) FindCateringByID(id uint) (domain.Catering, error) {
	var catering domain.Catering
	err := r.db.Preload("Images").First(&catering, id).Error
	return catering, err
}

func (r *cateringRepo) CreateCatering(c *domain.Catering) error {
	return r.db.Create(c).Error
}

// === TAMBAHAN BARU: UPDATE CATERING ===
func (r *cateringRepo) UpdateCatering(c *domain.Catering) error {
	// Gunakan transaksi agar jika gagal di tengah jalan, database di-rollback
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update data teks utama
		if err := tx.Model(&domain.Catering{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
			"name":        c.Name,
			"description": c.Description,
		}).Error; err != nil {
			return err
		}

		// 2. Jika ada gambar baru yang diunggah, hapus gambar lama dan masukkan yang baru
		if len(c.Images) > 0 {
			if err := tx.Where("catering_id = ?", c.ID).Delete(&domain.CateringImage{}).Error; err != nil {
				return err
			}
			for i := range c.Images {
				c.Images[i].CateringID = c.ID
			}
			if err := tx.Create(&c.Images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// === TAMBAHAN BARU: DELETE CATERING ===
func (r *cateringRepo) DeleteCatering(id uint) error {
	// Berkat constraint:OnDelete:CASCADE di struct, menghapus Catering otomatis menghapus semua CateringImage miliknya
	return r.db.Delete(&domain.Catering{}, id).Error
}

// -- Booking --
func (r *cateringRepo) FetchBookings() ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Preload("Catering").Order("booking_date desc").Find(&bookings).Error
	return bookings, err
}

func (r *cateringRepo) CreateBooking(b *domain.Booking) error {
	return r.db.Create(b).Error
}

func (r *cateringRepo) UpdateBookingStatus(id uint, status string) error {
	return r.db.Model(&domain.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *cateringRepo) CheckApprovedBookingExists(date time.Time) (bool, error) {
	var count int64
	// Cek apakah ada booking berstatus APPROVED pada hari (tanggal) yang sama
	err := r.db.Model(&domain.Booking{}).
		Where("DATE(booking_date) = DATE(?) AND status = ?", date, "APPROVED").
		Count(&count).Error
	return count > 0, err
}