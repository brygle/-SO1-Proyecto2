package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Caso struct {
	Name string `json:"name"`
	Location string `json:"location"`
	Age int16 `json:"age"`
	InfectedType string `json:"infectedType"`
	State string `json:"state"`
}

func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "este es el index")
}

func addCaso(w http.ResponseWriter, r *http.Request){

	//recuperar datos
	var newCaso Caso
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid case")
	}
	json.Unmarshal(reqbody, &newCaso)
	name:= newCaso.Name
	location := newCaso.Location
	age := newCaso.Age
	it := newCaso.InfectedType
	state := newCaso.State

	var jsonstr = string(`{name:` + name + `,location:` + location +`,age:`+strconv.Itoa(int(age))+`,infectedtype:`+it+`,state:`+state+ `}`)

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

	//body := fmt.Sprintf("%s,%s;%s,%s;%s,%s;%s,%s;%s,%s;" , "name", name,"location", location,"age", strconv.Itoa(int(age)),"infectedType", it,"state", state)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: jsonstr,
		})
	failOnError(err, "Fallo en enviar el mensaje")
	
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/", addCaso).Methods("POST")
	router.HandleFunc("/index", index).Methods("GET")
	fmt.Println("El servidor go a la escucha en puerto 5000")
	http.ListenAndServe(":5000",handlers.CORS(headers, methods, origins)(router))
}

