package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const viacepAPI = "http://viacep.com.br/ws/%s/json/"
const hgAPI = "http://api.hgbrasil.com/weather?key=SUA-CHAVE&fields=only_results,temp&city_name=%s"

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"ciaf"`
}

type HGResponse struct {
	Temp float64 `json:"temp"`
}

type WeatherValues struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

func main() {
	http.HandleFunc("/weather", handleWeather)
	http.ListenAndServe(":8080", nil)
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	log.Println("cep:", cep)
	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zip code"))
		return
	}

	city, err := getCityName(cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error getting city name"))
		return
	}
	log.Println("city:", city)
	weather, err := getWeather(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error getting weather: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func getCityName(cep string) (string, error) {
	url := fmt.Sprintf(viacepAPI, cep)
	log.Println("cep url:", url)
	req, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	var data ViaCEPResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", err
	}
	return data.Localidade, nil
}

func getWeather(city string) (WeatherValues, error) {
	url := fmt.Sprintf(hgAPI, city)
	log.Println("weather url:", url)
	req, err := http.Get(url)
	if err != nil {
		return WeatherValues{}, err
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return WeatherValues{}, err
	}
	var data HGResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		return WeatherValues{}, err
	}
	log.Println(string(res))
	return WeatherValues{
		Celsius:    data.Temp,
		Fahrenheit: data.Temp*1.8 + 32,
		Kelvin:     data.Temp + 273.15,
	}, nil
}
