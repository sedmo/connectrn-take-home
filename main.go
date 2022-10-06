package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/nfnt/resize"

	"github.com/gin-gonic/gin"
)

const (
	maxImageByteSize = 256 * 256
)

type Account struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	CreatedOn   int    `json:"created_on"`
}

type PostAccountResponse struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	DayofTheWeek string `json:"day_of_the_week"`
	RFC3339EST   string `json:"rfc_3339_est"`
}

func CreatePostAccountResponse(account Account) (r PostAccountResponse) {
	r.Name = account.Name
	r.UserID = account.UserID
	r.DayofTheWeek = GetDayOfTheWeek(account.DateOfBirth)
	r.RFC3339EST = getRFC3339EST(account.CreatedOn)
	return r
}

func GetDayOfTheWeek(DateOfBirth string) string {
	// conver to Time object
	if len(DateOfBirth) == 0 {
		return ""
	}
	t, err := time.Parse("2006-01-02", DateOfBirth)
	failOnError(err, "failed to parse dob")
	weekDay := t.Weekday().String()
	return weekDay
}

func getRFC3339EST(CreatedOn int) string {
	est, err := time.LoadLocation("America/New_York")
	failOnError(err, "fail on load location")
	t := time.Unix(int64(CreatedOn), 0).In(est).Format("2006-01-02T15:04:05-07:00")

	return t
}

func failOnError(e error, s string) {
	if e != nil {

		log.Fatalf("%v: %v\n", s, e.Error())
	}
}

func HandleAccounts(c *gin.Context) {
	var accounts []Account
	var postAccounts []PostAccountResponse
	if err := c.BindJSON(&accounts); err != nil {
		return
	}

	for _, account := range accounts {
		postAccounts = append(postAccounts, CreatePostAccountResponse(account))
	}
	c.IndentedJSON(http.StatusOK, postAccounts)
}

func HandleConvertJPEGToPNG(c *gin.Context) {

	file1, _, err := c.Request.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// get Asepct Ratio before conversion
	imgConfig, _, err := image.DecodeConfig(file1)
	failOnError(err, "failed decode config")
	width := imgConfig.Width
	height := imgConfig.Height

	// new reader for each decoder
	file2, _, err := c.Request.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
	jpegImage, err := jpeg.Decode(file2)
	if err != nil {
		c.Status(http.StatusUnsupportedMediaType)
		log.Println(err)
		return
	}

	if width > 256 || height > 256 {

		// resize to 256 and maintain the existing aspect ratio
		jpegImage = resize.Resize(256, 0, jpegImage, resize.Lanczos3)
	}

	// after resizing, convert to png

	buf := new(bytes.Buffer)
	if err = png.Encode(buf, jpegImage); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "image/png", buf, map[string]string{})

}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/accounts", HandleAccounts)
	router.POST("/convert", HandleConvertJPEGToPNG)
	return router
}

func CalculateAspectRatio(h int, w int) string {
	_gcd := gcd(h, w)

	return fmt.Sprintf("%v:%v", w/_gcd, h/_gcd)

}

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func main() {
	router := SetUpRouter()

	router.Run()

}
