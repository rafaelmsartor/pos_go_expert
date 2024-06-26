package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const brasilAPI = "https://brasilapi.com.br/api/cep/v1/%s"
const viacepAPI = "http://viacep.com.br/ws/%s/json/"

func main() {
	// read the CEP from the command line
	if len(os.Args) < 2 {
		fmt.Println("Please provide a CEP as an argument")
		return
	}
	cep := os.Args[1]

	first_channel := make(chan string)
	second_channel := make(chan string)
	go makeRequest(brasilAPI, cep, first_channel)
	go makeRequest(viacepAPI, cep, second_channel)

	timeout := time.After(1 * time.Second)
	// wait for the responses
	select {
	case res := <-first_channel:
		fmt.Println("Received the response from BrasilAPI first")
		fmt.Println(res)
	case res := <-second_channel:
		fmt.Println("Received the response from ViaCEP first")
		fmt.Println(res)
	case <-timeout:
		fmt.Println("Timeout reached")
	}
}

func makeRequest(api, cep string, ch chan string) {
	req, err := http.Get(fmt.Sprintf(api, cep))
	if err != nil {
		fmt.Println("Error making the request to ", api, err)
		return
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the response from ", api, err)
		return
	}
	ch <- string(res)
}
