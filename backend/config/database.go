package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf( //membuat string koneksi dengan format yang sesuai untuk postgres, menggunakan fmt.Sprintf untuk menyisipkan nilai dari environment variabel ke dalam string koneksi
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", //format connection string untuk postgres, ssl
		os.Getenv("DB_HOST"), //mengambil environment variabel dari sistem .env
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) //membuka koneksi ke database dengan menggunakan gorm dan postgres &gorm.Config{} untuk konfigurasi tambahan dan pointer ke struct gorm.Config
	if err != nil {
		log.Fatal("Gagal konek ke database", err)
	}

	log.Println("Berhasil konek ke database")
	DB = db //koneksi database yang tadi dibuat disimpan ke variabel global:
}

// analogi seperti kartu kosong var DB*gorm.db hanya berisikan alamat gudang, sedangkan ConnectDatabase() adalah proses mengisi kartu tersebut dengan barang (koneksi database)
// gorm.Open() seperti bos atau manajer yang membuka pintu gudang dan memastikan kartu akses siap untuk digunakan, jika ada masalah saat membuka pintu (koneksi), maka log.Fatal akan memberikan pesan kesalahan dan menghentikan program.
// Jika berhasil, maka log.Println akan memberikan pesan bahwa koneksi berhasil dibuat.
