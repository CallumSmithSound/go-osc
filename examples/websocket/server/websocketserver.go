package main

import (
	"fmt"
	"log"

	"github.com/CallumSmithSound/go-osc/osc"
)

func main() {

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("*", func(msg *osc.Message, clientAddress string) {
		fmt.Printf("SERVER: FROM CLIENT %v: ", clientAddress)
		osc.PrintMessage(msg)
	})

	server := osc.WebServer{
		Addr:       ":8085",
		Pattern:    "/osc",
		Dispatcher: d,
	}

	log.Println("starting server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
