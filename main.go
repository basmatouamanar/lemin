package main

import (
	"fmt"
	"os"
	"sort"

	"lemin/Display"
	"lemin/Find"
	"lemin/Helpers"
	"lemin/Parsing"
	"lemin/Var"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("USAGE: go run . \"data.txt\"")
		return // Exit if no input file is provided
	}

	// Read the input file
	dataBytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	// Parse the input data and initialize global variables
	err = Parsing.ParseData(string(dataBytes))

	Var.OriginalRooms = Helpers.CopyRooms(Var.Rooms)
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	Find.FindValidPaths()
	Var.AllVPaths = append(Var.AllVPaths, Var.VPaths)

	// Sort valid paths by length (shortest first)
	sort.Slice(Var.VPaths, func(i, j int) bool {
		return len(Var.VPaths[i]) < len(Var.VPaths[j])
	})

	// Assign ants to the initial path and calculate the required turns
	selectedPathIndex := 0
	minTurns, orderedAnts := Display.DistributeAnts(0)

	// Evaluate all other paths to find the path with the minimum number of turns
	for pathIndex := 1; pathIndex < len(Var.AllVPaths); pathIndex++ {
		turns, antsPerPath := Display.DistributeAnts(pathIndex)
		if turns < minTurns {
			orderedAnts = antsPerPath
			minTurns = turns
			selectedPathIndex = pathIndex
		}
	}

	// Print the results, including input data and ant movements
	Display.Printing(orderedAnts, minTurns, selectedPathIndex, string(dataBytes))
	
}
