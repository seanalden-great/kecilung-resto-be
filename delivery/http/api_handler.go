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
func RegisterHandlers(r *gin.Engine, cu domain.CategoryUsecase, mu domain.MenuUsecase, catUsecase domain.CateringUsecase, momentUsecase domain.MomentUsecase, articleUsecase domain.ArticleUsecase, contactUsecase domain.ContactUsecase) {
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
		catering.GET("/packages/:id", getCateringPackageByID(catUsecase))
		catering.POST("/packages", createCateringPackage(catUsecase))
		// Tambahkan 2 baris ini di dalam RegisterHandlers -> catering := api.Group("/catering")
		catering.PUT("/packages/:id", updateCateringPackage(catUsecase))
		catering.DELETE("/packages/:id", deleteCateringPackage(catUsecase))
		catering.GET("/bookings", getBookings(catUsecase))
		catering.POST("/bookings", createBooking(catUsecase))
		catering.PUT("/bookings/:id/approve", approveBooking(catUsecase))
		catering.PUT("/bookings/:id/reject", rejectBooking(catUsecase))
	}

	moment := api.Group("/moments")
	{
		moment.GET("/packages", getMoments(momentUsecase))
		moment.GET("/packages/:id", getMomentByID(momentUsecase))
		moment.POST("/packages", createMomentPackage(momentUsecase))
		moment.PUT("/packages/:id", updateMomentPackage(momentUsecase))
		moment.DELETE("/packages/:id", deleteMomentPackage(momentUsecase))

		moment.GET("/bookings", getMomentBookings(momentUsecase))
		moment.POST("/bookings", createMomentBooking(momentUsecase))
		moment.PUT("/bookings/:id/approve", approveMomentBooking(momentUsecase))
		moment.PUT("/bookings/:id/reject", rejectMomentBooking(momentUsecase))
	}

	// === TAMBAHKAN ROUTES ARTICLE DI SINI ===
	article := api.Group("/articles")
	{
		article.GET("", getArticles(articleUsecase))
		article.GET("/:id", getArticleByID(articleUsecase))
		article.POST("", createArticle(articleUsecase))
		article.PUT("/:id", updateArticle(articleUsecase))
		article.DELETE("/:id", deleteArticle(articleUsecase))
	}

	// === TAMBAHKAN ROUTES CONTACT DI SINI ===
	contact := api.Group("/contacts")
	{
		contact.POST("", createContact(contactUsecase))
		contact.GET("", getContacts(contactUsecase)) // Opsional untuk Admin
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

// func createCategory(u domain.CategoryUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var cat domain.Category
// 		if err := c.ShouldBindJSON(&cat); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if err := u.Create(&cat); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil dibuat", "data": cat})
// 	}
// }

// func updateCategory(u domain.CategoryUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		var cat domain.Category
// 		if err := c.ShouldBindJSON(&cat); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if err := u.Update(uint(id), &cat); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diupdate"})
// 	}
// }

func createCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil input teks dari form-data
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var imageURL string

		// 2. Cek apakah ada file gambar yang dilampirkan
		file, err := c.FormFile("image")
		if err == nil {
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		// 3. Susun struct Kategori
		cat := domain.Category{
			Code:        code,
			Name:        name,
			Description: description,
			ImageURL:    imageURL,
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

		// 1. Ambil data kategori lama
		existingCat, err := u.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		// 2. Ambil data baru
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		// Set default gambar ke gambar lama
		imageURL := existingCat.ImageURL

		// 3. Cek apakah ada gambar baru
		file, err := c.FormFile("image")
		if err == nil {
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		cat := domain.Category{
			Code:        code,
			Name:        name,
			Description: description,
			ImageURL:    imageURL,
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

func getCateringPackageByID(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetCateringByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data katering tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
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

		catering := domain.Catering{Name: name, Description: description, Images: images}
		if err := u.CreateCatering(&catering); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Katering ditambahkan"})
	}
}

func updateCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.CateringImage
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

		catering := domain.Catering{Name: name, Description: description, Images: images}
		if err := u.UpdateCatering(uint(id), &catering); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Katering berhasil diupdate"})
	}
}

func deleteCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.DeleteCatering(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Katering berhasil dihapus"})
	}
}

func createMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b domain.MomentBooking
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

func getMomentBookings(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func approveMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.ApproveBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui"})
	}
}

func rejectMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.RejectBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ditolak"})
	}
}

func getMoments(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, _ := u.GetAllMoments()
		c.JSON(200, gin.H{"data": res})
	}
}

func getMomentByID(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetMomentByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data moment tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.MomentImage

		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.MomentImage{ImageURL: url})
				}
			}
		}

		moment := domain.Moment{Name: name, Description: description, Images: images}
		if err := u.CreateMoment(&moment); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Moment ditambahkan"})
	}
}

func updateMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.MomentImage
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.MomentImage{ImageURL: url})
				}
			}
		}

		moment := domain.Moment{Name: name, Description: description, Images: images}
		if err := u.UpdateMoment(uint(id), &moment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Moment berhasil diupdate"})
	}
}

func deleteMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.DeleteMoment(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Moment berhasil dihapus"})
	}
}

// ====================================================================
// === HANDLERS ARTICLE (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

func getArticles(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func getArticleByID(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data artikel tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.ArticleImage

		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.ArticleImage{ImageURL: url})
				}
			}
		}

		article := domain.Article{Code: code, Name: name, Description: description, Images: images}
		if err := u.Create(&article); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Artikel berhasil ditambahkan", "data": article})
	}
}

func updateArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.ArticleImage
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.ArticleImage{ImageURL: url})
				}
			}
		}

		article := domain.Article{Code: code, Name: name, Description: description, Images: images}
		if err := u.Update(uint(id), &article); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil diupdate"})
	}
}

func deleteArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil dihapus"})
	}
}

// ====================================================================
// === HANDLERS CONTACT (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

func createContact(u domain.ContactUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var contactData domain.ContactUs
		if err := c.ShouldBindJSON(&contactData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Create(&contactData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Pesan Anda berhasil dikirim", "data": contactData})
	}
}

func getContacts(u domain.ContactUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}
