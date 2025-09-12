package helpers

import "fmt"

func Bfs(graphe map[string][]string, start, end *Room) {
	queue := []string{start.Name}
	visited := map[string]bool{start.Name: true}
	parent := map[string]string{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end.Name {
			break
		}

		for _, neighbor := range graphe[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	// Reconstruire le chemin
	path := []string{}
	for node := end.Name; node != ""; node = parent[node] {
		path = append([]string{node}, path...)
		if node == start.Name {
			break
		}
	}
	fmt.Println(path,"here")
}

