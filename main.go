package main

import (
	"fmt"
	"main/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	f1 := "server.cfg"
	x := utils.LoadDatabase(f1)
	// fmt.Printf("%+v\n", x)

	fmt.Printf(x.Dbname)

	// gin.SetMode(gin.ReleaseMode)
	// r := gin.Default()
	// r.POST("/api/v1/rfid_to_jan", controllers.GetJanCode)
	// r.POST("/api/v1/jan1_to_rfid", controllers.GetRFIDfromJAN1)
	// r.POST("/api/v1/jan2_to_rfid", controllers.GetRFIDfromJAN2)
	// r.POST("/api/v1/update_log", controllers.PostDataToLogTable)
	// if err := r.Run("0.0.0.0:8026"); err != nil {
	// 	panic(err.Error())
	// }

}
