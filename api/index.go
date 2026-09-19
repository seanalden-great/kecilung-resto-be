// package api

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/seanalden-great/kecilung-resto-be/config"
// 	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
// 	"github.com/seanalden-great/kecilung-resto-be/domain"
// 	"github.com/seanalden-great/kecilung-resto-be/repository"
// 	"github.com/seanalden-great/kecilung-resto-be/usecase"
// )

// var app *gin.Engine

// // Fungsi init() otomatis dieksekusi pertama kali oleh Vercel sebelum menerima request
// func init() {
// 	// Koneksi Database (menggunakan variabel environment Vercel)
// 	db := config.ConnectDatabase()

// 	// Migrasi Tabel
// 	err := db.AutoMigrate(&domain.Category{}, &domain.Menu{})
// 	if err != nil {
// 		log.Println("Gagal migrasi:", err)
// 	}

// 	// Inisialisasi Arsitektur
// 	catRepo := repository.NewCategoryRepository(db)
// 	catUseCase := usecase.NewCategoryUsecase(catRepo)

// 	menuRepo := repository.NewMenuRepository(db)
// 	menuUseCase := usecase.NewMenuUsecase(menuRepo)

// 	// Set Gin ke mode Release agar log lebih bersih di Vercel
// 	gin.SetMode(gin.ReleaseMode)
// 	app = gin.Default()

// 	// Middleware CORS
// 	app.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	// Daftarkan Routes
// 	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase)
// }

// // Handler adalah fungsi standar yang dicari oleh Vercel
// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// Serahkan request HTTP dari Vercel ke dalam router Gin
// 	app.ServeHTTP(w, r)
// }

package api

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/config"
	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/repository"
	"github.com/seanalden-great/kecilung-resto-be/usecase"
)

// var app *gin.Engine

var (
	app  *gin.Engine
	once sync.Once // 2. Daftarkan variabel penjaga gerbang
)

// Kita buat fungsi inisialisasi terpisah (BUKAN init() bawaan Go)
func initApp() {
	// Koneksi Database
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

	// Endpoint Debug (Untuk memastikan Vercel benar-benar membaca ENV)
	app.GET("/api/debug-env", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"host": os.Getenv("DB_HOST"),
			"port": os.Getenv("DB_PORT"),
			"user": os.Getenv("DB_USER"),
		})
	})

	// Daftarkan Routes
	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase)
}

// // Handler ini adalah pintu masuk utama Vercel
// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// Jika aplikasi belum diinisialisasi (request pertama), maka inisialisasi sekarang
// 	if app == nil {
// 		initApp()
// 	}
// 	// Teruskan request ke Gin Router
// 	app.ServeHTTP(w, r)
// }

// 3. Modifikasi fungsi Handler utama Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// once.Do memastikan fungsi initApp hanya dieksekusi 1x seumur hidup instance
	// Meskipun ratusan request datang menabrak bersamaan, yang lain akan menunggu hingga eksekusi pertama selesai
	once.Do(func() {
		initApp()
	})

	// Teruskan request ke Gin Router
	app.ServeHTTP(w, r)
}