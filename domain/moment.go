package domain

import "time"

type Moment struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Images      []MomentImage `json:"images" gorm:"foreignKey:MomentID;constraint:OnDelete:CASCADE;"`
}

type MomentImage struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	MomentID uint   `json:"moment_id"`
	ImageURL string `json:"image_url"`
}

type MomentBooking struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MomentID     uint      `json:"moment_id"`
	Moment       Moment    `json:"moment" gorm:"foreignKey:MomentID"`
	CustomerName string    `json:"customer_name"`
	Phone        string    `json:"phone"`
	BookingDate  time.Time `json:"booking_date"`
	Status       string    `json:"status" gorm:"default:'PENDING'"`
}

type MomentRepository interface {
	FetchAllMoments() ([]Moment, error)
	FindMomentByID(id uint) (Moment, error)
	CreateMoment(m *Moment) error
	UpdateMoment(m *Moment) error
	DeleteMoment(id uint) error
	
	FetchBookings() ([]MomentBooking, error)
	CreateBooking(b *MomentBooking) error
	UpdateBookingStatus(id uint, status string) error
	CheckApprovedBookingExists(date time.Time) (bool, error)
}

type MomentUsecase interface {
	GetAllMoments() ([]Moment, error)
	GetMomentByID(id uint) (Moment, error)
	CreateMoment(m *Moment) error
	UpdateMoment(id uint, m *Moment) error
	DeleteMoment(id uint) error
	
	GetAllBookings() ([]MomentBooking, error)
	CreateBooking(b *MomentBooking) error
	ApproveBooking(id uint) error
	RejectBooking(id uint) error
}