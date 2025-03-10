package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"fit-tracker/api"
	"fit-tracker/database"
	"fit-tracker/database/mongo"
	"fit-tracker/ingestor"
	ingestor_handler "fit-tracker/ingestor/handler"

	"github.com/gorilla/websocket"
)

const API_VERSION = "v1"

func main() {
	var (
		ctx, stop  = signal.NotifyContext(context.Background(), os.Interrupt)
		httpClient = http.DefaultClient
		wsClient   = websocket.DefaultDialer
	)

	defer stop()

	// database
	mongoClient, err := mongo.New(ctx)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Client.Disconnect(ctx)

	db := database.New(database.WithIngestorRepository(mongoClient))

	// ingestor
	ingestorHandler := ingestor_handler.New(httpClient, wsClient)
	ftIngestor := ingestor.New(ingestorHandler, db)

	go func() {
		if err := ftIngestor.Run(ctx); err != nil {
			log.Fatal("shutting down the server")
		}
	}()

	a := api.New(db)
	a.Run(ctx)
}
