package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger-service/data"
	"net/http"
	"time"
)

const (
	webPort  = "80"
	mongoURL = "mongodb://mongo:27017"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// go app.serve()
	log.Printf("Starting service on port:", webPort)
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = svr.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

//func (app *Config) serve() {
//	svr := &http.Server{
//		Addr:    fmt.Sprintf(":%s", webPort),
//		Handler: app.routes(),
//	}
//
//	err := svr.ListenAndServe()
//	if err != nil {
//		log.Panic()
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}

	log.Printf("Connect to mongo!")
	return c, nil
}
