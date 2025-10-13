package main

import (
	"fmt"
	"lem-in/Find"
	"lem-in/Var"
	"lem-in/Helpers"
	"lem-in/Parsing"
	"os"
	"sort"
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


	Var.OriginalRooms = Helpers.CopyRoomsMap(Var.Rooms)
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	Find.FindValidPaths()
	Var.AllValidPaths = append(Var.AllValidPaths, Var.ValidPaths)

	// Sort valid paths by length (shortest first)
	sort.Slice(Var.ValidPaths, func(i, j int) bool {
		return len(Var.ValidPaths[i]) < len(Var.ValidPaths[j])
	})

	// Assign ants to the initial path and calculate the required turns
	selectedPathIndex := 0
	minTurns, orderedAnts := Find.DistributeAnts(0)

// Evaluate all other paths to find the path with the minimum number of turns
	for pathIndex := 1; pathIndex < len(Var.AllValidPaths); pathIndex++ {
		turns, antsPerPath := Find.DistributeAnts(pathIndex)
		if turns < minTurns {
			orderedAnts = antsPerPath
			minTurns = turns
			selectedPathIndex = pathIndex
		}
	}

// Print the results, including input data and ant movements
	Parsing.Printing(orderedAnts, minTurns, selectedPathIndex, string(dataBytes))
	fmt.Println()

}