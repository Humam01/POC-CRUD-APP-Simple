package handlers

import (
	"net/http"

	"POC-CRUD-APP/config"
	"POC-CRUD-APP/models"

	"github.com/gin-gonic/gin"
)

// READ ALL - GET /api/tasks
func GetProducts(c *gin.Context) {
	var product []models.Product                      // [] List/array  menampilkann semua data yang ada di database dengan tipe data Task yang sudah dibuat di models/task.go
	config.DB.Order("created_at desc").Find(&product) //mengakses variabel config dan method DB untuk mengambil semua data dengan model Task, dengan parameter order untuk mengurutkan data berdasarkan field created_at secara menurun (desc) dan  mengambil semua data dengan method Find() dan menyimpan hasilnya ke variabel tasks
	c.JSON(http.StatusOK, product)
}

// READ ONE - GET /api/tasks/:id
func GetProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func SearchProduct(c *gin.Context) {
	var product []models.Product
	queryName := c.Query("name")
	queryID := c.Query("id")

	if queryName == "" && queryID == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Silakan masukkan kriteria pencarian",
			"data":    []models.Product{},
		})
		return
	}

	db := config.DB

	if queryName != "" {
		db = db.Where("name LIKE ?", "%"+queryName+"%")
	}
	if queryID != "" {
		db = db.Where("id = ?", queryID)
	}

	if err := db.Find(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal pencarian"})
		return
	}

	if len(product) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "product tidak ditemukan",
			"data":    []models.Product{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product ditemukan",
		"data":    product,
	})
}

func CreateProduct(c *gin.Context) { //mengakses data HTTP berupa JSON
	var product models.Product                         // task adalah sebuah struct, membuat variabel task dengan mengakses method dari variabel Task dalam folder models/task.go
	if err := c.ShouldBindJSON(&product); err != nil { //membuat variabel err(fungsi error) dengan mengikat/mengubah data JSON ke struct Go format dari data variabel task dan menghasilkan error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //memproses data JSON & struct dengan parameter menampilkan kesalahan dari user dengan mapping sehingga memberikan langsung pesan error
		return
	}

	result := config.DB.Create(&product) //membuat variabel result dan mengakses variabel config dan method DB dari config/database.go
	if result.Error != nil {             //result.Error memberikan informasi lebih result.RowsAffected: Berapa baris yang berhasil masuk ke database?
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()}) //memproses data JSON & struct dengan parameter yang menampilkan kesalahan dari server atau HTTP dengan mapping sehingga membuka koper atau data variabel result dan memberikan pesan error

		return
	}
	models.Products = append(models.Products, product)
	c.JSON(http.StatusCreated, product) //memberikan informasi data baru ditambahkan ke DB
}

// UPDATE - PUT /api/tasks/:id
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id") //mengambil nilai dari jalur path URL contoh:(http/id=4)

	if err := config.DB.First(&product, id).Error; err != nil { //GORM mencari data DB dengan ID yang cocok First(), .Error GORM menghasilkan error errReocrdNotFound jika tidak ditemukan, jika ditemukan maka akan mengupdate data dengan data baru yang diinputkan
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"}) //http.StatusNotFound memberikan informasi bahwa data tidak ditemukan, dengan mapping memberikan pesan error "Task tidak ditemukan"
		return
	}
	var input models.Product //sama seperti var task models.Task tapi variabel/sturct input digunakan untuk menyimpan data baru yang akan diupdate, sedangkan variabel task menyimpan data lama yang akan diupdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	result := config.DB.Model(&product).Updates(models.Product{ //mengakses variabel config dan method DB untuk mengupdate data dengan model Task, dengan parameter task yang merupakan data lama yang akan diupdate, dan data baru yang akan diupdate dengan mengakses struct Task dengan field Title, Description, Completed yang diambil dari variabel input
		Name:  input.Name,
		Price: input.Price,
		Stock: input.Stock,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate product"}) // bisa menggunakan result.Error.Error() untuk memberikan pesan error yang lebih spesifik dari GORM, tetapi jika koneksi terputus, constraint error dan database lock dan ternyata gagal maka server tetap mengirim status ok padahal gagal
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada perubahan data", "data": product}) // memberikan informasi bahwa tidak ada perubahan data yang terjadi, dengan mapping memberikan pesan "Tidak ada perubahan data" dan data yang lama
		return
	}
	models.Products = append(models.Products, product)

	c.JSON(http.StatusOK, product) //memberikan informasi bahwa data berhasil
}

// DELETE - DELETE /api/tasks/:id
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product Tidak Ditemukan"})
		return
	}

	result := config.DB.Delete(&product) //mengakses variabel config dan method DB untuk menghapus data dengan model Task, dengan parameter task yang merupakan data yang akan dihapus
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus product"})
		return
	}
	models.Products = append(models.Products, product)
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil Dihapus"})

}
