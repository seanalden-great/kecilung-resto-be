package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/config"
	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/repository"
	"github.com/seanalden-great/kecilung-resto-be/usecase"
)

var app *gin.Engine

// Fungsi init() otomatis dieksekusi pertama kali oleh Vercel sebelum menerima request
func init() {
	// Koneksi Database (menggunakan variabel environment Vercel)
	db := config.ConnectDatabase()
	
	// Migrasi Tabel
	err := db.AutoMigrate(&domain.Category{}, &domain.Menu{})
	if err != nil {
		log.Println("Gagal migrasi:", err)
	}

	// Inisialisasi Arsitektur
	catRepo := repository.NewCategoryRepository(db)
	catUseCase := usecase.NewCategoryUsecase(catRepo)

	menuRepo := repository.NewMenuRepository(db)
	menuUseCase := usecase.NewMenuUsecase(menuRepo)

	// Set Gin ke mode Release agar log lebih bersih di Vercel
	gin.SetMode(gin.ReleaseMode)
	app = gin.Default()

	// Middleware CORS
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Daftarkan Routes
	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase)
}

// Handler adalah fungsi standar yang dicari oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Serahkan request HTTP dari Vercel ke dalam router Gin
	app.ServeHTTP(w, r)
}