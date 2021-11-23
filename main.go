package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name   string `json: "name"`
	Author string `json: "author"`
}

var db *gorm.DB

func main() {
	// Connect to DB
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// ตอนสร้าง db ยังไม่มีtableต้องmigrate db ก่อน และถ้า schema เปลี่ยนก็จะเปลี่ยนให้ด้วย
	db.AutoMigrate(&Book{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Routes
	r.POST("/books", NewBook)
	r.GET("/books", ListBook)
	r.GET("/books/:id", GetBook)
	r.PUT("/books/:id", PutBook)

	r.Run("0.0.0.0:3030")
}

func NewBook(c *gin.Context) {
	var book Book
	if err := c.Bind(&book); err != nil { //gin จะทำการbind json ให้เหมือนกับstructของเรา ถ้ามีปัญหาก็จะเข้าในifแล้วพ่น error มาในรูปแบบjson
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := db.Create(&book)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"result": &book,
	})
}

func ListBook(c *gin.Context) {
	var books []Book
	result := db.Find(&books)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	id := c.Param("id")

	n, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var book Book
	result := db.First(&book, n)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, book)

}

func PutBook(c *gin.Context) {
	var book Book
	id := c.Param("id")
	//แปลง string > int
	n, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//หาว่ามีข้อมูลอยู่ในdbไหม
	data := db.First(&book, n)
	if err := data.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.Bind(&book); err != nil { //gin จะทำการbind json ให้เหมือนกับstructของเรา ถ้ามีปัญหาก็จะเข้าในifแล้วพ่น error มาในรูปแบบjson
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := db.Save(&book)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": &book,
	})

}
