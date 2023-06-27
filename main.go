package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"unicode"

	"strconv"

	"math"

	"strings"

	"github.com/gin-gonic/gin"
)

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []item `json:"items"`
}

var receipts = []receipt{}

// func getReceipts(context *gin.Context) {
// 	context.IndentedJSON(http.StatusOK, receipts)
// }

func addReceipt(context *gin.Context) {
	newUUID := uuid.New()
	var newReceipt receipt
	newReceipt.ID = newUUID.String()
	if err := context.BindJSON(&newReceipt); err != nil {
		return
	}

	receipts = append(receipts, newReceipt)
	context.IndentedJSON(http.StatusCreated, gin.H{"id": newUUID})
}

func getReceipt(context *gin.Context) {
	id := context.Param("id")
	points, err := calculatePoints(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Receipt Not Found"})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"points": points})

}

func calculatePoints(id string) (int, error) {
	points := 0
	for _, r := range receipts {
		if r.ID == id {
			// One point for every alphanumeric character in the retailer name.
			for _, char := range r.Retailer {
				if unicode.IsLetter(char) || unicode.IsNumber(char) {
					points++
				}
			}

			// 50 points if the total is a round dollar amount with no cents
			totalArray := strings.Split(r.Total, ".")
			if totalArray[1] == "00" {
				points = points + 50
			}

			// 25 points if the total is a multiple of 0.25.
			totalFloat, err := strconv.ParseFloat(r.Total, 64)
			if err != nil {
				log.Fatal(err)
			}
			multiple := totalFloat / 0.25

			if multiple == float64(int(multiple)) {
				points = points + 25
			}

			// 5 points for every two items on the receipt.
			totalItems := len(r.Items) / 2
			points = points + 5*int(float64(totalItems))

			// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
			for _, item := range r.Items {
				desc := len(strings.TrimSpace(item.ShortDescription))
				if desc%3 == 0 {
					priceFloat, _ := strconv.ParseFloat(item.Price, 64)
					points = points + int(math.Ceil(priceFloat*0.2))
				}
			}

			// 6 points if the day in the purchase date is odd.
			date := strings.Split(r.PurchaseDate, "-")
			day, _ := strconv.ParseInt(date[len(date)-1], 0, 64)
			if day%2 == 1 {
				points = points + 6
			}

			//10 points if the time of purchase is after 2:00pm and before 4:00pm.
			timeArray := strings.Split(r.PurchaseTime, ":")
			hour, err := strconv.ParseInt(timeArray[0], 0, 64)
			if err != nil {
				log.Fatal(err)
			}
			min, err := strconv.ParseInt(timeArray[1], 0, 64)
			if err != nil {
				log.Fatal(err)
			}

			if hour == 15 || (hour == 14 && min > 0) {

				points = points + 10
			}

			return points, nil
		}
	}
	return -1, errors.New("receipt not found")
}

func main() {

	router := gin.Default()
	//router.GET("/receipts/getReceipts", getReceipts)
	router.GET("/receipts/:id/points", getReceipt)
	router.POST("/receipts/process", addReceipt)
	router.Run(":9090")

}
