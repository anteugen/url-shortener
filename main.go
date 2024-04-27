package main

import (
	"database/sql"
	"fmt"
	"log"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"

	"github.com/anteugen/url-shortener/database"
)

func generateShortURL(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func insertRow(db *sql.DB, originalURL, shortURL string) {
	query := `INSERT INTO urls (original_url, short_url_code) VALUES ($1, $2)`
	_, err := db.Exec(query, originalURL, shortURL)
	if err != nil {
		log.Fatal("Error inserting into table: ", err)
	}
}

func main() {
	db := database.ConnectToDB()
	defer db.Close()

	shortCode := generateShortURL(6)
	fmt.Println(shortCode)

	router := gin.Default()

	router.POST("/shorten", func(c *gin.Context) {
		var urlRequest struct {
			URL string `json:"url"`
		}

		if err := c.ShouldBindJSON(&urlRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortURL := generateShortURL(6)

		c.JSON(http.StatusOK, gin.H{"original_url": urlRequest.URL, "short_url": shortURL})
		insertRow(db, urlRequest.URL, shortURL)
	})

	router.GET("/r/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		
		var originalURL string
		err := db.QueryRow("SELECT original_url FROM urls WHERE short_url_code = $1", shortCode).Scan(&originalURL)
		if err != nil {
			log.Fatal("Error fetching original url from database:", err)
			return
		}

		c.Redirect(http.StatusMovedPermanently, originalURL)
	})

	router.Run(":8080")
}
