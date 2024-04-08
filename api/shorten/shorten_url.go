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

func shorten(url string) (string, error) {
	var shortURL string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("short_urls"))
		id, _ := bucket.NextSequence()
		shortURL = fmt.Sprintf("https://shorti.com/%d", id)

		return bucket.Put([]byte(shortURL), []byte(url))
	})

	return shortURL, err
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

func ShortenUrl(c *gin.Context) {

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	short_url, err := shorten(payload.Url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"shortened_url": short_url})
}

func GetOriginalUrlFromDb(c *gin.Context) {
	short_url := c.Param("short_url")

	var url string

	var cleaned_short_url string

	if len(short_url) > 0 && short_url[0] == '/' {
		//create slice of string from index 1 to end
		cleaned_short_url = short_url[1:]
	}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("short_urls"))
		url = string(bucket.Get([]byte(cleaned_short_url)))
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func init() {
	if err := initDb(); err != nil {
		panic(err)
	}
}
