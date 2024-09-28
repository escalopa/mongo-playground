package main

import (
	"context"
	"fmt"

	"github.com/escalopa/mongo-playground/server"
	"github.com/escalopa/mongo-playground/storage"
)

var (
	addr = ":8080"

	password = "example"
	dsn      = fmt.Sprintf("mongodb://root:%s@localhost:27017", password)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongodb, err := storage.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	srv := server.New(mongodb)
	if err = srv.Run(addr); err != nil {
		panic(err)
	}
}
