package controllers

import (
	"main/db_client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JAN2_to_RFID struct {
	ApiKey   string `json:"api_key" binding:"required"`
	JanCode2 string `json:"jancode_2" binding:"required"`
	JanCode1 string `json:"jancode_1"`
	RFID     string `json:"rfid"`
	Status   string `json:"status"`
}

func GetRFIDfromJAN2(c *gin.Context) {
	var reqBody JAN2_to_RFID
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isSuccess": false,
			"message":   "Invalid request format(not JSON) or missing (jancode_2, api_key) information",
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
	rfid_list, data_exist, _ := db_client.ConvertFromJan2(db, reqBody.JanCode2)

	data := gin.H{
		"rfid_list": rfid_list,
		"jancode_2": reqBody.JanCode2,
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
