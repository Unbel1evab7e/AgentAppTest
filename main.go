package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	mainUrl, err := url.Parse(viper.GetString("geoIPFY.url"))
	if err != nil {
		log.Fatal(err.Error())
	}

	token := viper.GetString("geoIPFY.token")
	params := url.Values{}
	params.Add("apiKey", token)
	mainUrl.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, mainUrl.String(), strings.NewReader(""))
	if err != nil {
		log.Fatal(err.Error())
	}

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err.Error())
	}

	lat := result["location"].(map[string]interface{})["lat"]
	lng := result["location"].(map[string]interface{})["lng"]

	mainUrl, err = url.Parse(viper.GetString("weatherapi.url"))
	if err != nil {
		log.Fatal(err.Error())
	}

	token = viper.GetString("weatherapi.token")
	params = url.Values{}
	params.Add("q", fmt.Sprintf("%v,%v", lat, lng))
	params.Add("key", token)
	mainUrl.RawQuery = params.Encode()

	req, err = http.NewRequest(http.MethodGet, mainUrl.String(), strings.NewReader(""))
	if err != nil {
		log.Fatal(err.Error())
	}

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	response := WeatherApiReponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err.Error())
	}

	output, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(output))
}
