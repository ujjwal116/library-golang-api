package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DB *gorm.DB
}

var DBC *DBConfig

func init() {
	fmt.Println("DbConfig init is running")
	fmt.Println("App configured")
	db, err := setpAndConnectDbConfig()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	} else {
		DBC = &DBConfig{DB: db}
	}
}

func setpAndConnectDbConfig() (*gorm.DB, error) {
	log.Println("Configuring Library DB")
	viper.SetEnvPrefix("DATABASE")
	viper.AutomaticEnv()
	host := viper.GetString("HOST")
	port := viper.GetInt("PORT")
	user := viper.GetString("USER")
	password := viper.GetString("PASSWORD")
	name := viper.GetString("NAME")
	connectionString := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	log.Println("Connecting with DB")
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	return db, err
}