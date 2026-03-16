
package main

import (
	"log"

	"github.com/hannanaarif/Social/internal/db"
	"github.com/hannanaarif/Social/internal/store"
	"github.com/hannanaarif/Social/internal/env"
)

func main(){

	addr:=env.GetString("DB_ADDR","postgres://admin:adminpassword@localhost:5432/social?sslmode=disable")

	conn,err :=db.New(addr,10,10,"10m")
	if err!=nil{
		log.Fatal("error while creating db connection",err)
	}
	defer conn.Close()
	store:=store.NewStorage(conn)
	db.Seed(store,conn)
	
}	