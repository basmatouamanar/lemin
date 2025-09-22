package main

import (
	"fmt"
	"strings"
)

//print the mouvement of ants step by step
func printAnts(paths [][]string, pathCounts []int) {
	type antState struct {
		id    int
		path  []string
		index int // Position in path, -delay means waiting to start
	}

	var ants []antState
	antID := 1

	// Assign ants to paths with delays so they start one by one on the path
	for i, count := range pathCounts {
		for j := 0; j < count; j++ {
			ants = append(ants, antState{id: antID, path: paths[i], index: -j})
			antID++
		}
	}

	for {
		moved := false
		output := []string{}

		for i := range ants {
			ants[i].index++
			if ants[i].index >= 0 && ants[i].index < len(ants[i].path) {
				moved = true
				output = append(output, fmt.Sprintf("L%d-%s", ants[i].id, ants[i].path[ants[i].index]))
			}
		}

		if !moved {
			break
		}

		fmt.Println(strings.Join(output, " "))
	}
}

