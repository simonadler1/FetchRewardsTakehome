package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessReceipt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var requestBody Receipt
	requestBody.Retailer = "test retailer"
	requestBody.Total = "200"
	requestBody.PurchaseDate = "2022-08-01"
	requestBody.PurchaseTime = "14:30"

	requestBodyBytes, _ := json.Marshal(requestBody)

	router := setupRouter()
	request, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(requestBodyBytes))
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetPoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	receipt := &Receipt{
		ID:       "testid",
		Retailer: "test retailer",
		Total:    "200",
		PurchaseDate: "2022-08-01",
		PurchaseTime: "14:30",
	}

	receiptStore[receipt.ID] = *receipt

	router := setupRouter()
	request, _ := http.NewRequest(http.MethodGet, "/receipts/" + receipt.ID + "/points", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestProcessReceipt_Target(t *testing.T) {
	gin.SetMode(gin.TestMode)

	requestBody := Receipt{
		Retailer:     " Target ",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []struct {
			ShortDescription string `json:"shortDescription"`
			Price            string `json:"price"`
		}{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "Klarbrunn 12-PK 12 FL OZ", Price: "12.00"},
		},
		Total: "35.35",
	}

	requestBodyBytes, _ := json.Marshal(requestBody)

	router := setupRouter()
	request, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(requestBodyBytes))
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var body map[string]string
	json.Unmarshal(responseRecorder.Body.Bytes(), &body)

	request2, _ := http.NewRequest(http.MethodGet, "/receipts/"+body["id"]+"/points", nil)
	responseRecorder2 := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder2, request2)
	assert.Equal(t, http.StatusOK, responseRecorder2.Code)

	var body2 PointsResponse
	json.Unmarshal(responseRecorder2.Body.Bytes(), &body2)
	assert.Equal(t, 28, body2.Points)
}

func TestProcessReceipt_MMCornerMarket(t *testing.T) {
    gin.SetMode(gin.TestMode)

    requestBody := Receipt{
    Retailer:     "M&M Corner Market  ",
    PurchaseDate: "2022-03-20",
    PurchaseTime: "14:33",
    Items: []struct {
        ShortDescription string `json:"shortDescription"`
        Price            string `json:"price"`
    }{
        {ShortDescription: "Gatorade", Price: "2.25"},
        {ShortDescription: "Gatorade", Price: "2.25"},
        {ShortDescription: "Gatorade", Price: "2.25"},
        {ShortDescription: "Gatorade", Price: "2.25"},
    },
    Total: "9.00",
}


    requestBodyBytes, _ := json.Marshal(requestBody)

    router := setupRouter()
    request, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(requestBodyBytes))
    responseRecorder := httptest.NewRecorder()

    router.ServeHTTP(responseRecorder, request)

    assert.Equal(t, http.StatusOK, responseRecorder.Code)

    var body map[string]string
    json.Unmarshal(responseRecorder.Body.Bytes(), &body)

    request2, _ := http.NewRequest(http.MethodGet, "/receipts/"+body["id"]+"/points", nil)
    responseRecorder2 := httptest.NewRecorder()

    router.ServeHTTP(responseRecorder2, request2)
    assert.Equal(t, http.StatusOK, responseRecorder2.Code)

    var body2 map[string]int
    json.Unmarshal(responseRecorder2.Body.Bytes(), &body2)
    assert.Equal(t, 109, body2["points"])
}