package controllers

import (
	"main/db_client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RFID_to_JAN struct {
	ApiKey   string `json:"api_key" binding:"required"`
	RFID     string `json:"rfid" binding:"required"`
	JanCode1 string `json:"jancode_1"`
	JanCode2 string `json:"jancode_2"`
	Status   string `json:"status"`
}

func GetJanCode(c *gin.Context) {
	var reqBody RFID_to_JAN
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isSuccess": false,
			"message":   "Invalid request format(not JSON) or missing (rfid, api_key) information",
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
	jancode_1, jan_code_2, status, data_exist, _ := db_client.ConvertFromRFID(db, reqBody.RFID)

	data := gin.H{
		"rfid":      reqBody.RFID,
		"jancode_1": jancode_1,
		"jancode_2": jan_code_2,
		"status":    status,
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
