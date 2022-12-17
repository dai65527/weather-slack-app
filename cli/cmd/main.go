package main

import (
	"flag"
	"fmt"
	"log"
	"weather-slack-app/weather"
)

func main() {
	city := flag.String("city", "Tokyo", "city name")
	flag.Parse()
	res, err := weather.GetWeather(*city)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
