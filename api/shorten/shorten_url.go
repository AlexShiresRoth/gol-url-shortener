package shorten_url

import (
	"fmt"

	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

var db *bolt.DB

type Payload struct {
	Url string `json:"url"`
}

func initDb() error {
	var err error
	db, err = bolt.Open("short_urls.db", 0600, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("short_urls"))
		return err
	})

}

func shorten(url string) (string, uint64, error) {
	var shortURL string
	var urlId uint64

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("short_urls"))
		id, _ := bucket.NextSequence()
		urlId = uint64(id)
		// Convert this to an env var
		shortURL = fmt.Sprintf("https://shorti.com/%d", id)

		return bucket.Put([]byte(fmt.Sprintf("%d", id)), []byte(url)) // Store original URL with ID as the key
	})

	return shortURL, urlId, err
}

func ShortenUrl(c *gin.Context) {

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	short_url, id, err := shorten(payload.Url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"shortened_url": short_url, "id": id})
}

func GetOriginalUrlFromDb(c *gin.Context) {
	id := c.Param("id")

	var shortURL string
	var originalURL string

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("short_urls"))
		originalURL = string(bucket.Get([]byte(id)))        // Retrieve original URL using ID as the key
		shortURL = fmt.Sprintf("https://shorti.com/%s", id) // Construct short URL based on ID
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"original_url": originalURL, "short_url": shortURL})

}

func init() {
	if err := initDb(); err != nil {
		panic(err)
	}
}
