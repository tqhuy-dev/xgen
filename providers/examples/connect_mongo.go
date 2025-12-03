package main

import (
	"time"

	"github.com/tqhuy-dev/xgen/providers/mongo_db"
)

func main() {
	mongoClient, err := mongo_db.NewMongoDB(mongo_db.Option{
		Host:        "localhost",
		Port:        27018,
		DB:          "sample",
		User:        "admin",
		Password:    "123456789",
		MinPoolSize: 5,
		MaxPoolSize: 10,
		MaxIdleTime: 5 * time.Second,
		IsAdmin:     true,
	})
	if err != nil {
		panic(err)
	}
	defer mongoClient.Close()
}
