package controllers

import (
	"main/db_client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateLog struct {
	ApiKey     string `json:"api_key" binding:"required"`
	StoreCode  string `json:"store_code"  binding:"required"`
	RFID       string `json:"rfid"  binding:"required"`
	CreateDate string `json:"create_date"`
	Mode       string `json:"mode"`
}

func PostDataToLogTable(c *gin.Context) {
	var reqBody UpdateLog
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isSuccess": false,
			"message":   "Invalid request format(not JSON) or missing (api_key, store_code, rfid) information",
		})
		return
	}

	if err := VerifyApiKey(reqBody.ApiKey); !err {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"isSuccess": false,
			"message":   "Invalid API key",
		})
		return
	}

	db, _ := db_client.DbConnection()
	defer db.Close()
	if err := db_client.UpdateLog(db, reqBody.CreateDate, reqBody.StoreCode, reqBody.RFID, reqBody.Mode); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"isSuccess":        false,
			"message":          "Data fields (store_code or rfid) existed or missing",
			"database_message": err,
		})
		return
	}
	data := gin.H{
		"create_date": reqBody.CreateDate,
		"store_code":  reqBody.StoreCode,
		"rfid":        reqBody.RFID,
		"mode":        reqBody.Mode,
	}

	c.JSON(http.StatusCreated, gin.H{
		"isSuccess": true,
		"data":      data,
	})
}
