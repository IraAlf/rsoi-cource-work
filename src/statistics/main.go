package main

import (
	"log"
	"statistics/controllers"
	"statistics/repository"
	"statistics/utils"

	"math/rand"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	kafka := utils.InitKafka()
	defer kafka.Close()
	go kafka.ConsumeLoop()

	db := utils.InitDBConnection(utils.Config.DB)
	defer db.Close()
	go utils.StatWriteLoop(repository.NewPGStatisticsRep(db))

	r := controllers.InitRouter(db)
	utils.Logger.Print("Server started")
	log.Printf("Server is running on http://localhost:%d\n", utils.Config.Port)
	code := controllers.RunRouter(r, utils.Config.Port)

	utils.Logger.Printf("Server ended with code %s", code)
}
