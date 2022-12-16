package main

import (
	"flag"
	"fmt"
)

func main() {
	city := flag.String("city", "Tokyo", "city name")
	flag.Parse()
	fmt.Println("Hello " + *city)
}
