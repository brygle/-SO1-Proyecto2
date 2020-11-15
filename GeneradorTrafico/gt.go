package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"log"
	"net/http"
	"bytes"
)

var urlBC string
var cgr string
var casos string
var ruta string

type Caso struct {
	Name string `json:"name"`
	Location string `json:"location"`
	Age int16 `json:"age"`
	InfectedType string `json:"infectedType"`
	State string `json:"state"`
}

var arrCasos []Caso 

func main() {
	fmt.Println("\n***** ***** *****  ** SISTEMAS OPERATIVOS 1 *  ***** ***** *****")
	fmt.Println("***** ***** ***** ***** ** PROYECTO 2 ** ***** ***** ***** *****\n")
	
	for true {	
		fmt.Println("Ingrese operación\n")
		fmt.Println("1. Ingresar URL de balanceador de carga")
		fmt.Println("2. Ingresar cantidad de gorutinas")
		fmt.Println("3. Ingresar cantidad de casos a enviar")
		fmt.Println("4. Ingresar ruta de archivo a cargar")
		fmt.Println("5. Enviar datos")
		fmt.Println("6. Salir\n") 
		fmt.Printf("SOPES1:~$ ")
		var operacion string
		rdr := bufio.NewReader(os.Stdin)
		operacion, _ = rdr.ReadString('\n')
		operacion = strings.Replace(operacion, "\n", "", -1)
		fmt.Println()

		if operacion == "6" {
			break
		} else if operacion == "1" {
			ingresarURL()
		} else if operacion == "2" {
			ingresarCantidadGoRutinas()
		} else if operacion == "3" {
			ingresarCantidadCasos()
		} else if operacion == "4" {
			ingresarRutaArchivo()
		} else if operacion == "5" {
			enviarDatos()
		} else {
			fmt.Println("***** OPERACION INCORRECTA *****")
		}
	}
	fmt.Println("Finalizando...")
}

func ingresarURL(){
	fmt.Printf("Ingrese URL: ")
	rdr := bufio.NewReader(os.Stdin)
	urlBC, _ = rdr.ReadString('\n')
	urlBC = strings.Replace(urlBC, "\n", "", -1)
	fmt.Println()
}

func ingresarCantidadGoRutinas(){
	fmt.Printf("Ingrese Cantidad Gorrutinas: ")
	rdr := bufio.NewReader(os.Stdin)
	cgr, _ = rdr.ReadString('\n')
	cgr = strings.Replace(cgr, "\n", "", -1)
	fmt.Println()
}

func ingresarCantidadCasos(){
	fmt.Printf("Ingrese Cantidad Casos: ")
	rdr := bufio.NewReader(os.Stdin)
	casos, _ = rdr.ReadString('\n')
	casos = strings.Replace(casos, "\n", "", -1)
	fmt.Println()
}
func ingresarRutaArchivo(){
	fmt.Printf("Ingrese ruta archivo: ")
	rdr := bufio.NewReader(os.Stdin)
	ruta, _ = rdr.ReadString('\n')
	ruta = strings.Replace(ruta, "\n", "", -1)

	file_data, err := ioutil.ReadFile(ruta)

	if err != nil {
		fmt.Println("***** ERROR AL LEER EL ARCHIVO *****")
	}
	json.Unmarshal([]byte(file_data), &arrCasos )
	fmt.Println()
}
func enviarDatos(){
	fmt.Println("URL: " , urlBC)
	fmt.Println("Gorrutinas: " , cgr)
	fmt.Println("Casos: " , casos)
	fmt.Println("Ruta Archivo: " , ruta)
	fmt.Println("Casos", arrCasos)
	fmt.Println()

	c := make(chan int)

	intGR , _ := strconv.Atoi(cgr)
	intCasos , _ := strconv.Atoi(casos)
	longitudArreglo := len(arrCasos)

	nCaso := 0

	fin := 0
	if(intGR <= intCasos && intGR <= longitudArreglo){
		fin = intGR
	} else if(intCasos <= intGR && intCasos <= longitudArreglo){
		fin = intCasos
	} else {
		fin = longitudArreglo
	}
	fmt.Println("El fin va a ser", fin)
	for i := 0 ; i < fin ; i++ {
		//canal := <-c
		go postDatos(i, c , nCaso)
		nCaso++;
	}

	for i := 0 ; i < fin ; i++ {
		canal := <-c
		fmt.Println("Finalizada la gorrutina numero " , canal)
	}

	close(c)

}

func postDatos(id int, c chan int, nCaso int){
	fmt.Println("Enviando la gorrutina " , id)
	clienteHttp := &http.Client{}

	var caso Caso
	caso = (arrCasos[nCaso])
	datosjson, err := json.Marshal(caso)
	if err != nil {
		log.Fatalf("Error codificando caso como JSON: %v", err)
	}
	
	direccion := "http://" + urlBC
	peticion, err := http.NewRequest("POST", direccion, bytes.NewBuffer(datosjson))
	if err != nil {
		// Maneja el error de acuerdo a tu situación
		log.Fatalf("Error creando petición: %v", err)
	}
	
	peticion.Header.Add("Content-Type", "application/json")
	fmt.Println("llego aqui 1", direccion)
	respuesta, err := clienteHttp.Do(peticion)
	
	if err != nil {
		// Maneja el error de acuerdo a tu situación
		log.Fatalf("Error haciendo petición: %v", err)
	}
	defer respuesta.Body.Close()

	c <- id
}