package main

import (
	"fmt"
	"os"
	"strings"
)

var Paths [][]string

type Room struct {
	Name string
	X    int
	Y    int
}

type Tunnel struct {
	From string
	To   string
}

type AntFarm struct {
	Rooms   []Room
	Tunnels []Tunnel
	Start   Room
	End     Room
	Ants    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	antFarm, file, err := readInput(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	findPaths(antFarm.Tunnels, antFarm.Start.Name, antFarm.End.Name, []string{})
	if len(Paths) == 0 {
		fmt.Println("ERROR: invalid path")
		return
	}

	fmt.Println(strings.Join(file, "\n") + "\n")

	bestPaths := getBestPaths(Paths, antFarm.Ants)
	if len(bestPaths) == 0 {
		fmt.Println("ERROR: no usable paths")
		return
	}

	finalPaths, pathCounts := distributeAnts(antFarm.Ants, bestPaths)
	printAnts(finalPaths, pathCounts)
}
