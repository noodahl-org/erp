package main

import (
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	api "github.com/noodahl-org/erp/api/conf"
	pg "github.com/noodahl-org/erp/api/clients/postgres"
	"github.com/noodahl-org/erp/api/server"
	"github.com/noodahl-org/erp/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serv server.Server

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	valid := validator.New()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	config := api.ApiConf{
		Port:        port,
		BraveAPIKey: os.Getenv("BRAVE_API_KEY"),
		DbConf: api.DbConf{
			DbPort: dbPort,
			DbUser: os.Getenv("DB_USER"),
			DbPass: os.Getenv("DB_PASS"),
			DbHost: os.Getenv("DB_HOST"),
			DbName: os.Getenv("DB_NAME"),
		},
	}

	if err := valid.Struct(config); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(config.DbConf.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // not necessary in this example but a good idea later on...
	})
	if err != nil {
		log.Fatal(err)
	}

	serv = *server.NewServer(
		server.WithConfig(&config),
		server.WithDB(
			pg.NewPGDB(
				pg.WithGormDB(db),
			),
		),
	)

}

func main() {
	defs := []interface{}{
		models.Equipment{},
		models.User{},
		models.UserEquipment{},
		models.MaintenanceTask{},
		models.UserMaintenanceTask{},
		models.EquipmentComponent{},
	}
	for _, d := range defs {
		serv.MigrateDomainModel(&d)
	}

	serv.Start()
}
