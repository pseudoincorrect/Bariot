package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() {
	fmt.Println("connecting to the DB service_db...")
	dsn := "host=service_db user=service_db password=service_db dbname=service_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("...connected to DB service_db")

}
