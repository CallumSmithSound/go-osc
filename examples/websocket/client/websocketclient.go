package main

import (
	"log"

	"github.com/CallumSmithSound/go-osc/osc"
)

func main() {
	client := osc.NewWebClient("localhost:8085", "/osc")
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	msg := osc.NewMessage("/hehas/wdwd", "hello World", int32(2323))
	err = client.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent msg %v\n", msg.String())
}
