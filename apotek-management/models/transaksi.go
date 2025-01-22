package models

import (
	"time"
)

type Transaksi struct {
	ID            uint              `json:"id_transaksi" gorm:"primaryKey"`
	KodeTransaksi string            `json:"kode_transaksi" gorm:"type:varchar(20);unique;not null"`
	TotalHarga    int               `json:"total_harga" gorm:"type:int;not null"`
	Status        string            `json:"status" gorm:"type:varchar(50);not null"`
	ObatID        uint              `json:"id_obat" gorm:"not null"` // Foreign key untuk Obat
	Obats         []TransaksiDetail `json:"obats" gorm:"foreignKey:TransaksiID"`
	CreatedAt     time.Time         `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type TransaksiDetail struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TransaksiID uint      `json:"id_transaksi" gorm:"not null"`
	ObatID      uint      `json:"id_obat" gorm:"not null"`
	Jumlah      int       `json:"jumlah" gorm:"type:int;not null"`
	Obat        Obat      `json:"obat" gorm:"foreignKey:ObatID;references:ID"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// BeforeCreate hook untuk menghitung TotalHarga berdasarkan harga obat dari relasi
func (transaksi *Transaksi) BeforeCreate() (err error) {
	totalHarga := 0
	for _, detail := range transaksi.Obats {
		totalHarga += detail.Jumlah * int(detail.Obat.Harga)
	}
	transaksi.TotalHarga = totalHarga
	return nil
}
