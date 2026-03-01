package main

import (
	"log"

	"github.com/hannanaarif/Social/internal/db"
	"github.com/hannanaarif/Social/internal/env"
	"github.com/hannanaarif/Social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
