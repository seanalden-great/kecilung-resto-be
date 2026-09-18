package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
)

func RegisterHandlers(r *gin.Engine, cu domain.CategoryUsecase, mu domain.MenuUsecase) {
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
		menu.POST("", createMenu(mu))
		menu.PUT("/:id", updateMenu(mu))
		menu.DELETE("/:id", deleteMenu(mu))
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

func createMenu(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m domain.Menu
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Create(&m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Menu berhasil dibuat", "data": m})
	}
}

func updateMenu(u domain.MenuUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var m domain.Menu
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Update(uint(id), &m); err != nil {
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