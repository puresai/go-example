package main

import (
	"demo10/models"
)

func main() {
	client := models.GetClient()
	client.Start()
	client.Near()
	defer client.Close()
}


