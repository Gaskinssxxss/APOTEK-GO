package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateStok menambahkan transaksi stok (MASUK/KELUAR)
func CreateStok(c *gin.Context) {
	var stok models.Stok
	if err := c.ShouldBindJSON(&stok); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Menyimpan transaksi stok
	if err := config.DB.Create(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Memuat relasi Obat, TipeObat, dan Tags
	var stokWithRelations models.Stok
	if err := config.DB.First(&stokWithRelations, stok.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load relations: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stokWithRelations)
}

// GetAllStok mengambil seluruh data transaksi stok
func GetAllStok(c *gin.Context) {
	var stokList []models.Stok
	if err := config.DB.Find(&stokList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stokList)
}

// GetStokByID mengambil transaksi stok berdasarkan ID
func GetStokByID(c *gin.Context) {
	id := c.Param("id")
	var stok models.Stok
	if err := config.DB.First(&stok, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok not found"})
		return
	}

	c.JSON(http.StatusOK, stok)
}

// UpdateStok mengupdate data transaksi stok berdasarkan ID
func UpdateStok(c *gin.Context) {
	id := c.Param("id")
	var stok models.Stok
	if err := config.DB.First(&stok, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok not found"})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&stok); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data stok
	if err := config.DB.Save(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stok)
}

// DeleteStok menghapus transaksi stok berdasarkan ID
func DeleteStok(c *gin.Context) {
	id := c.Param("id")
	var stok models.Stok
	if err := config.DB.First(&stok, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok not found"})
		return
	}

	// Hapus transaksi stok
	if err := config.DB.Delete(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stok deleted successfully"})
}
