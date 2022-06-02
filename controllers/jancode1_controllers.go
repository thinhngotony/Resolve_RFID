package controllers

import (
	"main/db_client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JAN_to_RFID struct {
	ApiKey   string `json:"api_key" binding:"required"`
	JanCode1 string `json:"jancode_1" binding:"required"`
	RFID     string `json:"rfid"`
	JanCode2 string `json:"jancode_2"`
	Status   string `json:"status"`
}

func GetRFIDfromJAN1(c *gin.Context) {
	var reqBody JAN_to_RFID
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isSuccess": false,
			"message":   "Invalid request format(not JSON) or missing (jancode_1, api_key) information",
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
	rfid_list, data_exist, _ := db_client.ConvertFromJan1(db, reqBody.JanCode1)

	data := gin.H{
		"rfid_list": rfid_list,
		"jancode_1": reqBody.JanCode1,
	}

	if data_exist {
		c.JSON(http.StatusOK, gin.H{
			"isSuccess": true,
			"data":      data,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"isSuccess": false,
			"message":   "Data not exist in database",
		})
	}

}
