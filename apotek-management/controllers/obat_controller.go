package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateObat(c *gin.Context) {
	var obat models.Obat

	if err := c.ShouldBindJSON(&obat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&obat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(obat.Tags) > 0 {
		for _, tag := range obat.Tags {
			var existingTag models.TagObat
			if err := config.DB.First(&existingTag, tag.ID).Error; err == nil {
				config.DB.Model(&obat).Association("Tags").Append(&existingTag)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID: " + err.Error()})
				return
			}
		}
	}

	var obatWithRelations models.Obat
	if err := config.DB.Preload("TipeObat").Preload("Tags").First(&obatWithRelations, obat.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load relations: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obatWithRelations)
}

func CreateBatchObat(c *gin.Context) {
	var obatList []models.Obat

	if err := c.ShouldBindJSON(&obatList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Create(&obatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obatList)
}

func GetAllObat(c *gin.Context) {
	var obatList []models.Obat
	if err := config.DB.Preload("TipeObat").Preload("Tags").Preload("Stok").Find(&obatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, obatList)
}

func GetObatByID(c *gin.Context) {
	id := c.Param("id")
	var obat models.Obat

	// Cari data obat beserta relasinya
	if err := config.DB.
		Preload("TipeObat").
		Preload("Tags").Preload("Stok").
		First(&obat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	// Kirim data obat dengan relasi
	c.JSON(http.StatusOK, obat)
}

func UpdateObat(c *gin.Context) {
	id := c.Param("id")
	var existingObat models.Obat

	// Cari data obat berdasarkan ID
	if err := config.DB.First(&existingObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	var updatedObat struct {
		KodeObat   string           `json:"kode_obat"`
		NamaObat   string           `json:"nama_obat"`
		Deskripsi  string           `json:"deskripsi"`
		HargaObat  uint64           `json:"harga_obat"`
		TipeObatID uint             `json:"id_tipe_obat"`
		Tags       []models.TagObat `json:"tags"`
	}

	if err := c.ShouldBindJSON(&updatedObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Update field di existingObat
	existingObat.KodeObat = updatedObat.KodeObat
	existingObat.NamaObat = updatedObat.NamaObat
	existingObat.Deskripsi = updatedObat.Deskripsi
	existingObat.Harga = updatedObat.HargaObat
	existingObat.TipeObatID = updatedObat.TipeObatID

	// Update data obat
	if err := config.DB.Save(&existingObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Obat: " + err.Error()})
		return
	}

	// Update relasi Tags jika ada
	if len(updatedObat.Tags) > 0 {
		var tagIDs []uint
		for _, tag := range updatedObat.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}

		var tags []models.TagObat
		if err := config.DB.Where("id_tag_obat IN ?", tagIDs).Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags: " + err.Error()})
			return
		}
		if err := config.DB.Model(&existingObat).Association("Tags").Replace(&tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags association: " + err.Error()})
			return
		}
	}

	// Reload data obat dengan relasi
	if err := config.DB.Preload("TipeObat").Preload("Tags").First(&existingObat, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated Obat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingObat)
}

func UpdateBatchObat(c *gin.Context) {
	var obatList []models.Obat

	if err := c.ShouldBindJSON(&obatList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	for _, obat := range obatList {
		if err := config.DB.Save(&obat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, obatList)
}

func DeleteObat(c *gin.Context) {
	id := c.Param("id")
	var obat models.Obat

	// Cari data obat berdasarkan ID
	if err := config.DB.First(&obat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	// Hapus semua tag terkait obat di tabel many-to-many
	if err := config.DB.Model(&obat).Association("Tags").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related tags: " + err.Error()})
		return
	}

	// Hapus stok terkait obat
	if err := config.DB.Where("obat_id = ?", obat.ID).Delete(&models.Stok{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related stock: " + err.Error()})
		return
	}

	// Hapus data obat
	if err := config.DB.Delete(&obat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete obat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat and related data deleted successfully"})
}

func DeleteBatchObat(c *gin.Context) {
	var ids []uint

	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Delete(&models.Obat{}, ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat(s) deleted successfully"})
}
