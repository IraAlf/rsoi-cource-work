package main

import (
	"gateway/controllers"
	"gateway/utils"
	"log"

	"math/rand"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	r := controllers.InitRouter()
	utils.Logger.Print("Server started")
	log.Printf("Server is running on http://localhost:%d\n", utils.Config.Port)
	code := controllers.RunRouter(r, utils.Config.Port)

	utils.Logger.Printf("Server ended with code %s", code)
}
