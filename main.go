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
	// fmt.Println(graphe, "heeee")
	var save [][]string
	
	for {
		path := helpers.Bfs(graphe, StartRoom, EndRoom)
		if path == nil || len(path) == 0 {
			break
		}
		fmt.Println(path)
		save = append(save, path)
		
		for i := 0; i < len(path); i++ {
			for j := 0; j < len(graphe[path[i]]); j++ {
				m := graphe[path[i]][j]
				if i+1 < len(path) && m == path[i+1] {
					graphe[path[i]] = append(graphe[path[i]][:j], graphe[path[i]][j+1:]...)
					j--
					// fmt.Println(graphe[path[i]], "heeeeeeeeeeeeeeeeeee")
				}
			}
		}
		fmt.Println(save)
		path = nil
	}

}
