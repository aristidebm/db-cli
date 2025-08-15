package main

import (
	"log"

	"example.com/db/internal/client"
	"example.com/db/internal/config"
)

func main() {
	config := &config.Config{}
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	name := "sakila"
	source := config.Sources[name]
	c, err := client.NewClient(source.URL)
	if err != nil {
		log.Fatal(err)
	}
	c.SourceConfig = source
	if driver, ok := config.Drivers[c.Driver]; ok {
		c.DriverConfig = driver
	}
}
