package main

import (
	"POC-CRUD-APP/config"
	"POC-CRUD-APP/handlers"
	"POC-CRUD-APP/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() { //func adalah method
	if err := godotenv.Load(".env"); err != nil { //godotenv.load mencari dan membaca file .env
		if err2 := godotenv.Load("../.env"); err2 != nil { //jika tidak ditemukan di folder utama, coba cari di path folder parent../env
			log.Println("File .env tidak ditemukan, menggunakan env vars") //dengan menggunakan teknik callback
		}
	}

	//Koneksi ke database
	config.ConnectDatabase()

	//Auto migrate (buat tabel otomatis)
	config.DB.AutoMigrate(&models.Product{}) //AutoMigrate() otomatis membuat atau memperbarui tabel dalam database berdasarkan struktur  &models.Task{} membuat atau memperbarui tabel yang sesuai dengan struktur model Task yang telah didefinisikan di dalam package models.

	//Setup router
	r := gin.Default() // membuat terminal log untuk mencatat aktivitas server

	//CORS Middleware - untuk mengizinkan akses dari frontend
	r.Use(func(c *gin.Context) { //r.Use() mencegat semua permintaan masuk HTTP dan hanya membaca data header dari permintaan c *gin.Context yang memberikan informasi tentang permintaan HTTP yang sedang diproses
		c.Header("Access-Control-Allow-Origin", "*")                               //Header Access-Control-Allow-Origin memberi tahu browser apakah kode JavaScript di situs lain diizinkan untuk mengakses resource dari server ini, dengan nilai "*" berarti mengizinkan akses dari semua domain.
		c.Header("Access-Control-Allow-Methods", "GET,POST, PUT, DELETE, OPTIONS") //Header Access-Control-Allow-Methods memberi tahu browser metode HTTP mana yang diizinkan untuk permintaan CORS, dengan nilai "GET, POST, PUT, DELETE, OPTIONS" berarti mengizinkan metode HTTP tersebut untuk permintaan CORS.
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")    //c.header membuat peraturan browser Access-Control-Allow-Headers menentukan header mana yang diizinkan dalam permintaan CORS, Content-Type memberikan informasi tentang jenis data yang dikirimkan dalam permintaan HTTP, Authorization memberikan informasi tentang kredensial atau token yang digunakan untuk otentikasi login
		if c.Request.Method == "OPTIONS" {                                         //c.Request.Method  informasi method HTTP  jika method HTTP adalah OPTIONS  maka server akan merespons method HTTP apa saja yang diizinkan untuk permintaan CORS dengan header Access-Control-Allow-Methods
			c.AbortWithStatus(http.StatusNoContent) //mengakhiri proses permintaan dengan c.AbortWithStatus(http.StatusNoContent) dan return untuk keluar dari fungsi middleware
			return
		}
		c.Next() //proses selain method HTTP OPTIONS akan diteruskan ke handler berikutnya dengan c.Next() untuk melanjutkan proses permintaan HTTP
	})

	//Routes
	api := r.Group("/api")
	{
		api.POST("/products", handlers.CreateProduct)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/products/search", handlers.SearchProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)
	}

	//Start server
	port := os.Getenv("PORT")
	if port == "" { //callback jika tidak ditemukan variabel PORT di file .env, maka akan menggunakan port default 8081
		port = "8081"
	}
	log.Printf("Server running on port %s\n", port)
	r.Run(":" + port) //r.run enahan program agar terus berjalan (looping) menunggu request masuk.

}
