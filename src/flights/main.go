package main

import (
	"flights/controllers"
	"flights/objects"
	"flights/utils"
	"log"

	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDBConnection(cnf utils.DBConfiguration) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cnf.Host, cnf.User, cnf.Password, cnf.Name, cnf.Port)
	db, e := gorm.Open(cnf.Type, dsn)

	if e != nil {
		utils.Logger.Print("DB Connection failed")
		utils.Logger.Print(e)
		panic("DB Connection failed")
	} else {
		utils.Logger.Print("DB Connection Established")
	}

	db.SingularTable(true)
	db.AutoMigrate(&objects.Airport{})
	db.AutoMigrate(&objects.Flight{})

	flight := &objects.Flight{
		Id:           1,
		FlightNumber: "AFL031",
		Datetime:     "2024-10-08 20:00",
		FromAirport:  objects.Airport{Id: 2, Name: "Пулково", City: "Санкт-Петербург", Country: "Россия"},
		ToAirport:    objects.Airport{Id: 1, Name: "Шереметьево", City: "Москва", Country: "Россия"},
		Price:        6000,
	}
	db.FirstOrCreate(flight)
	flight = &objects.Flight{
		Id:           1,
		FlightNumber: "BGH789",
		Datetime:     "2024-13-08 12:45",
		FromAirport:  objects.Airport{Id: 2, Name: "Пулково", City: "Санкт-Петербург", Country: "Россия"},
		ToAirport:    objects.Airport{Id: 1, Name: "Домодедово", City: "Москва", Country: "Россия"},
		Price:        3000,
	}
	db.Create(flight)
	return db
}

func main() {
	rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	db := initDBConnection(utils.Config.DB)
	defer db.Close()
	r := controllers.InitRouter(db)

	utils.Logger.Print("Server started")
	log.Printf("Server is running on http://localhost:%d\n", utils.Config.Port)
	code := controllers.RunRouter(r, utils.Config.Port)

	utils.Logger.Printf("Server ended with code %s", code)
}
