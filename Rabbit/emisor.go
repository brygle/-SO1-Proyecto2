package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"strconv"
)

func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo al conectar con rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //nombre de la cola a la que nos queremos susbcribir
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Fallo al abrir el canal")

	name:="op"
	location := "guate"
	it := "comunitario"
	state := "sintomatico"

	for i:= 0 ; i<10 ; i++ {
		body := fmt.Sprintf("%s,%s;%s,%s;%s,%s;%s,%s;%s,%s;" , "name", name+strconv.Itoa(i),"location", location,"age", strconv.Itoa(i),"infectedType", it,"state", state)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body: []byte(body),
			})
		log.Printf("Sent %s", body)
		failOnError(err, "Fallo en enviar el mensaje")
	}
}

