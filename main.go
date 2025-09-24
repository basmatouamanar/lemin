package main

import (
	"fmt"

	"lemin/helpers"
)

var rooms []*helpers.Room

func main() {
	data, err := helpers.ReadFile("sample.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if helpers.ValidateData(data) {
		fmt.Println("✅")
	}else {
		return
	}
	helpers.RealizeLink(data, &rooms)
	StartRoom, EndRoom, ants := helpers.Path(data, &rooms)
	fmt.Println(ants, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaants")
	graphe := helpers.Graphe(rooms)
	// fmt.Println(graphe, "heeee")
	var save [][]string

	for {
		path := helpers.Bfs(graphe, StartRoom, EndRoom)
		if path == nil || len(path) == 0 {
			break
		}

		save = append(save, path)

		for i := 0; i < len(path)-1; i++ {
			current := path[i]
			next := path[i+1]

			newNeighbors := []string{}
			for _, neighbor := range graphe[current] {
				if neighbor != next {
					newNeighbors = append(newNeighbors, neighbor)
				}
			}
			graphe[current] = newNeighbors
		}
	}
	fmt.Println(save)
	helpers.Distribute(save, ants)
	fmt.Println(ants, "heeeeeeeere")
}
