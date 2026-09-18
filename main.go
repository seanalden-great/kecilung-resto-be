package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seanalden-great/kecilung-resto-be/config"
	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/repository"
	"github.com/seanalden-great/kecilung-resto-be/usecase"
)

func main() {
	_ = godotenv.Load()

	db := config.ConnectDatabase()
	
	// MIGRATION: Category harus di atas Menu karena Menu bergantung pada ID Category
	err := db.AutoMigrate(&domain.Category{}, &domain.Menu{})
	if err != nil {
		log.Fatal("Gagal migrasi tabel:", err)
	}

	// WIRING ARCHITECTURE
	catRepo := repository.NewCategoryRepository(db)
	catUseCase := usecase.NewCategoryUsecase(catRepo)

	menuRepo := repository.NewMenuRepository(db)
	menuUseCase := usecase.NewMenuUsecase(menuRepo)

	// INIT ROUTER
	r := gin.Default()

	// MIDDLEWARE CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// REGISTER ENDPOINTS
	httpDelivery.RegisterHandlers(r, catUseCase, menuUseCase)

	r.Run(":8080")
}