package main

import (
	_ "gowebserver/boot"
)
import "gowebserver/app"

// @title        MPU Application platform services Swagger API
// @version      2023.03.08
// @description  MPU Application platform services
// @license.name Apache 2.0
// @schemes      http https
// @BasePath     /mpuaps/v1
// @accept 		 json
// @produce 	 json
// @author       LiYongJun
// @date         2023.03.08
func main() {
	app.Run()
}
