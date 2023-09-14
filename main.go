package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	Total        string `json:"total"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items"`
}
// init storage
var receiptStore = make(map[string]Receipt)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)
	return router
}

func processReceipt(c *gin.Context) {
	var receipt Receipt

	if err := c.BindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	 _, err1 := time.Parse("2006-01-02", receipt.PurchaseDate)
    if err1 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Date"})
        return
    }

    _, err2 := time.Parse("15:04", receipt.PurchaseTime)
    if err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Time"})
        return
    }

	receipt.ID = uuid.New().String()
	
	receiptStore[receipt.ID] = receipt

	c.JSON(http.StatusOK, gin.H{"id": receipt.ID})
}

func getPoints(c *gin.Context) {
	receiptID := c.Param("id")

	receipt, exists := receiptStore[receiptID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt does not exist"})
		return
	}

	points := calculatePoints(&receipt)

	c.JSON(http.StatusOK, PointsResponse{Points: points})
}
func calculatePoints(receipt *Receipt) int {
    var points int

    retailer := strings.TrimSpace(receipt.Retailer)
    retailerLength := 0
    for _, c := range retailer {
        if unicode.IsLetter(c) || unicode.IsNumber(c) {
            retailerLength++
        }
    }
    points += retailerLength

    total, err := strconv.ParseFloat(receipt.Total, 64)
    if err == nil {
        if total == math.Floor(total) && total >= 1 {
            points += 50
        }

        if math.Mod(total*100, 25) == 0 {
            points += 25
        }
    }

    pairCount := 0
    for _, item := range receipt.Items {
        description := strings.TrimSpace(item.ShortDescription)
        itemDescriptionLength := len(description)
        if itemDescriptionLength%3 == 0 {
            price, err := strconv.ParseFloat(strings.TrimSpace(item.Price), 64)
            if err == nil {
                points += int(math.Ceil(price * 0.2))
            }
        }
        pairCount++
    }

    points += 5 * (pairCount / 2)

    parsedDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
    if err == nil && parsedDate.Day()%2 == 1 {
        points += 6
    }

    parsedTime, err := time.Parse("15:04", receipt.PurchaseTime)
    if err == nil && parsedTime.Hour() >= 14 && parsedTime.Hour() < 16 {
        points += 10
    }

    return points
}
type PointsResponse struct {
	Points int `json:"points"`
}