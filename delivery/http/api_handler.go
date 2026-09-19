package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

// func RegisterHandlers(r *gin.Engine, cu domain.CategoryUsecase, mu domain.MenuUsecase) {
// 	api := r.Group("/api")

// 	// Routes Kategori
// 	cat := api.Group("/categories")
// 	{
// 		cat.GET("", fetchCategories(cu))
// 		cat.POST("", createCategory(cu))
// 		cat.PUT("/:id", updateCategory(cu))
// 		cat.DELETE("/:id", deleteCategory(cu))
// 	}

// 	// Routes Menu
// 	menu := api.Group("/menus")
// 	{
// 		menu.GET("", fetchMenus(mu))
// 		menu.GET("/category/:id", getMenusByCategory(mu))
// 		menu.POST("", createMenu(mu))
// 		menu.PUT("/:id", updateMenu(mu))
// 		menu.DELETE("/:id", deleteMenu(mu))
// 	}
// }

// PERUBAHAN 1: Tambahkan catUsecase domain.CateringUsecase di parameter ini
func RegisterHandlers(r *gin.Engine, cu domain.CategoryUsecase, mu domain.MenuUsecase, catUsecase domain.CateringUsecase) {
	api := r.Group("/api")

	// Routes Kategori
	cat := api.Group("/categories")
	{
		cat.GET("", fetchCategories(cu))
		cat.POST("", createCategory(cu))
		cat.PUT("/:id", updateCategory(cu))
		cat.DELETE("/:id", deleteCategory(cu))
	}

	// Routes Menu
	menu := api.Group("/menus")
	{
		menu.GET("", fetchMenus(mu))
		menu.GET("/category/:id", getMenusByCategory(mu))
		menu.POST("", createMenu(mu))
		menu.PUT("/:id", updateMenu(mu))
		menu.DELETE("/:id", deleteMenu(mu))
	}

	// Routes Catering
	catering := api.Group("/catering")
	{
		catering.GET("/greeting", getGreeting(catUsecase)) 
		catering.PUT("/greeting", updateGreeting(catUsecase))
		catering.GET("/packages", getCaterings(catUsecase))
		catering.POST("/packages", createCateringPackage(catUsecase))
		catering.GET("/bookings", getBookings(catUsecase))
		catering.POST("/bookings", createBooking(catUsecase))
		catering.PUT("/bookings/:id/approve", approveBooking(catUsecase))
		catering.PUT("/bookings/:id/reject", rejectBooking(catUsecase))
	}
}

// === HANDLERS KATEGORI ===
func fetchCategories(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

// func getMenusByCategory(u domain.MenuUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		menus, err := u.GetByCategoryID(uint(id))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": menus})
// 	}
// }

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

func createCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cat domain.Category
		if err := c.ShouldBindJSON(&cat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Create(&cat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil dibuat", "data": cat})
	}
}

func updateCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var cat domain.Category
		if err := c.ShouldBindJSON(&cat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Update(uint(id), &cat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diupdate"})
	}
}

func deleteCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
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

// func createMenu(u domain.MenuUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var m domain.Menu
// 		if err := c.ShouldBindJSON(&m); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if err := u.Create(&m); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusCreated, gin.H{"message": "Menu berhasil dibuat", "data": m})
// 	}
// }

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

// func updateMenu(u domain.MenuUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		var m domain.Menu
// 		if err := c.ShouldBindJSON(&m); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if err := u.Update(uint(id), &m); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil diupdate"})
// 	}
// }

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

// ====================================================================
// === HANDLERS CATERING (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

func getGreeting(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg, err := u.GetGreeting()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": msg})
	}
}

func updateGreeting(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg domain.GreetingMessage
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.UpdateGreeting(&msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Pesan sambutan berhasil diupdate"})
	}
}

func getBookings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func approveBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.ApproveBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui"})
	}
}

func rejectBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.RejectBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ditolak"})
	}
}

func getCaterings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, _ := u.GetAllCaterings()
		c.JSON(200, gin.H{"data": res})
	}
}

func createCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.CateringImage
		
		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.CateringImage{ImageURL: url})
				}
			}
		}

		catering := domain.Catering{ Name: name, Description: description, Images: images }
		if err := u.CreateCatering(&catering); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Katering ditambahkan"})
	}
}

func createBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b domain.Booking
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := u.CreateBooking(&b); err != nil {
			// Menampilkan error jika jadwal bentrok
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Booking berhasil diajukan"})
	}
}
