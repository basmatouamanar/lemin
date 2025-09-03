package main

import (
	"fmt"
	"lemin/helpers"
)

func main() {
	data, err := helpers.ReadFile("sample.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if helpers.ValidateData(data) {
		fmt.Println("✅")
	}
	helpers.RealizeLink(data)
	
}
