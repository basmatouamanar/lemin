package main

import (
	"fmt"
	"lemin/helpers"

)


var (
	rooms []*helpers.Room
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
	helpers.RealizeLink(data, &rooms)
	StartRoom, EndRoom := helpers.Path(data, &rooms)
	graphe := helpers.Graphe(rooms)
	helpers.Bfs(graphe, StartRoom, EndRoom)


}
