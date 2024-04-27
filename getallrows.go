package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"github.com/anteugen/url-shortener/database"
)

type URLInfo struct {
	ID           int     
	ShortURLCode string    
	OriginalURL  string    
	CreationDate string
}

func storeCSV(URLInfo []URLInfo) {
	file, err := os.Create("urls.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"id",
		"short_url_code",
		"original_url",
		"creation_date",
	}

	writer.Write(headers)

	for _, url := range URLInfo {
		record := []string{
			strconv.Itoa(url.ID),
			url.ShortURLCode,
			url.OriginalURL,
			url.CreationDate,
		}

		writer.Write(record)
	}

	defer writer.Flush()
}

func main() {
	db := database.ConnectToDB()
	defer db.Close()	

	var urls []URLInfo

	rows, err := db.Query("SELECT id, short_url_code, original_url, creation_date FROM urls")
	if err != nil {
		log.Fatal("Error fetching data from database:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url URLInfo
		err := rows.Scan(&url.ID, &url.ShortURLCode, &url.OriginalURL, &url.CreationDate)
		if err != nil {
			log.Fatal("Error scanning data from row:", err)
			return
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error during rows iteration:", err)
	}

	storeCSV(urls)

}