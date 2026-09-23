package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"golang.org/x/crypto/bcrypt"
)

// GANTI INI DENGAN SECRET KEY YANG LEBIH AMAN (Bisa ditaruh di .env)
var jwtSecretKey = []byte("KecilungRestoSecretKey2026")

type authUsecase struct {
	repo domain.AuthRepository
}

func NewAuthUsecase(r domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{repo: r}
}

func (u *authUsecase) Login(username, password string) (string, domain.Admin, error) {
	admin, err := u.repo.GetAdminByUsername(username)
	if err != nil {
		return "", admin, errors.New("username tidak ditemukan")
	}

	// Bandingkan password input dengan password hash di DB
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", admin, errors.New("password salah")
	}

	// Generate JWT Token (Berlaku 24 Jam)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   admin.ID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	
	// Kosongkan password agar tidak terkirim kembali di JSON response
	admin.Password = "" 
	
	return tokenString, admin, err
}

func (u *authUsecase) GetProfile(id uint) (domain.Admin, error) {
	admin, err := u.repo.GetAdminByID(id)
	admin.Password = "" // Keamanan
	return admin, err
}

func (u *authUsecase) UpdateProfile(id uint, data *domain.Admin) error {
	// Jika ada password baru, Hash terlebih dahulu
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hashedPassword)
	}
	return u.repo.UpdateAdmin(id, data)
}