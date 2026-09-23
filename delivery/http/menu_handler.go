package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

// Routes Menu
func RegisterMenuHandlers(rg *gin.RouterGroup, u domain.MenuUsecase) {
	menu := rg.Group("/menus")
	{
		menu.GET("", fetchMenus(u))
		menu.GET("/category/:id", getMenusByCategory(u))
		menu.POST("", createMenu(u))
		menu.PUT("/:id", updateMenu(u))
		menu.DELETE("/:id", deleteMenu(u))
	}
}

// === HANDLERS MENU ===
func fetchMenus(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func getMenusByCategory(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID Kategori tidak valid"})
			return
		}

		menus, err := u.GetByCategoryID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Pengaman: Jika tidak ada menu di kategori ini, kembalikan array kosong
		if menus == nil {
			menus = []domain.Menu{}
		}

		c.JSON(http.StatusOK, gin.H{"data": menus})
	}
}

func createMenu(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil input teks dari form-data
		name := c.PostForm("name")
		description := c.PostForm("description")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
		categoryID, _ := strconv.Atoi(c.PostForm("category_id"))

		var imageURL string

		// 2. Cek apakah ada file gambar yang dilampirkan
		file, err := c.FormFile("image")
		if err == nil {
			// Jika ada file, unggah ke Clever Cloud S3
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah gambar ke server S3"})
				// return
				// UBAH BARIS INI:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		// 3. Susun data struct
		menu := domain.Menu{
			Name:        name,
			Description: description,
			Price:       price,
			CategoryID:  uint(categoryID),
			ImageURL:    imageURL, // Simpan URL-nya di sini
		}

		// 4. Simpan ke database via Usecase -> Repository
		if err := u.Create(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Menu berhasil dibuat", "data": menu})
	}
}

func updateMenu(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		// 1. Cek apakah menu ada, sekaligus mengambil URL gambar lama
		existingMenu, err := u.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
			return
		}

		// 2. Ambil data teks dari form-data
		name := c.PostForm("name")
		description := c.PostForm("description")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
		categoryID, _ := strconv.Atoi(c.PostForm("category_id"))

		// Set default gambar ke gambar lama
		imageURL := existingMenu.ImageURL

		// 3. Cek apakah ada gambar BARU yang diunggah
		file, err := c.FormFile("image")
		if err == nil {
			// Jika ada, unggah dan timpa URL gambar lama dengan yang baru
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah gambar baru ke S3"})
				// return
				// UBAH BARIS INI:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		// 4. Susun data struct
		menu := domain.Menu{
			Name:        name,
			Description: description,
			Price:       price,
			CategoryID:  uint(categoryID),
			ImageURL:    imageURL,
		}

		// 5. Eksekusi update
		if err := u.Update(uint(id), &menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil diupdate"})
	}
}

func deleteMenu(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus"})
	}
}