package handlers

import (
	"net/http"

	"POC-CRUD-APP/config"
	"POC-CRUD-APP/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) { //mengakses data HTTP berupa JSON
	var task models.Task                            // task adalah sebuah struct, membuat variabel task dengan mengakses method dari variabel Task dalam folder models/task.go
	if err := c.ShouldBindJSON(&task); err != nil { //membuat variabel err(fungsi error) dengan mengikat/mengubah data JSON ke struct Go format dari data variabel task dan menghasilkan error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //memproses data JSON & struct dengan parameter menampilkan kesalahan dari user dengan mapping sehingga memberikan langsung pesan error
		return
	}

	result := config.DB.Create(&task) //membuat variabel result dan mengakses variabel config dan method DB dari config/database.go
	if result.Error != nil {          //result.Error memberikan informasi lebih result.RowsAffected: Berapa baris yang berhasil masuk ke database?
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()}) //memproses data JSON & struct dengan parameter yang menampilkan kesalahan dari server atau HTTP dengan mapping sehingga membuka koper atau data variabel result dan memberikan pesan error
		return
	}

	c.JSON(http.StatusCreated, task) //memberikan informasi data baru ditambahkan ke DB
}

// READ ALL - GET /api/tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task                         // [] List/array  menampilkann semua data yang ada di database dengan tipe data Task yang sudah dibuat di models/task.go
	config.DB.Order("created_at desc").Find(&tasks) //mengakses variabel config dan method DB untuk mengambil semua data dengan model Task, dengan parameter order untuk mengurutkan data berdasarkan field created_at secara menurun (desc) dan  mengambil semua data dengan method Find() dan menyimpan hasilnya ke variabel tasks
	c.JSON(http.StatusOK, tasks)
}

// READ ONE - GET /api/tasks/:id
func GetTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UPDATE - PUT /api/tasks/:id
func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id") //mengambil nilai dari jalur path URL contoh:(http/id=4)

	if err := config.DB.First(&task, id).Error; err != nil { //GORM mencari data DB dengan ID yang cocok First(), .Error GORM menghasilkan error errReocrdNotFound jika tidak ditemukan, jika ditemukan maka akan mengupdate data dengan data baru yang diinputkan
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"}) //http.StatusNotFound memberikan informasi bahwa data tidak ditemukan, dengan mapping memberikan pesan error "Task tidak ditemukan"
		return
	}
	var input models.Task //sama seperti var task models.Task tapi variabel/sturct input digunakan untuk menyimpan data baru yang akan diupdate, sedangkan variabel task menyimpan data lama yang akan diupdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	result := config.DB.Model(&task).Updates(models.Task{ //mengakses variabel config dan method DB untuk mengupdate data dengan model Task, dengan parameter task yang merupakan data lama yang akan diupdate, dan data baru yang akan diupdate dengan mengakses struct Task dengan field Title, Description, Completed yang diambil dari variabel input
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate task"}) // bisa menggunakan result.Error.Error() untuk memberikan pesan error yang lebih spesifik dari GORM, tetapi jika koneksi terputus, constraint error dan database lock dan ternyata gagal maka server tetap mengirim status ok padahal gagal
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada perubahan data", "data": task}) // memberikan informasi bahwa tidak ada perubahan data yang terjadi, dengan mapping memberikan pesan "Tidak ada perubahan data" dan data yang lama
		return
	}

	c.JSON(http.StatusOK, task) //memberikan informasi bahwa data berhasil
}

// DELETE - DELETE /api/tasks/:id
func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task Tidak Ditemukan"})
		return
	}

	result := config.DB.Delete(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil Dihapus"})

}
