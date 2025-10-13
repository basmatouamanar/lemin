package Display

import (
	"fmt"
	"lemin/Var"
	"strconv"
	"strings"
)

// Printing prints the input data and simulates the movement of ants along the shortest path.
func Printing(antsPerPath []int, totalTurns int, shortestPathIndex int, originalData string) {
	// Print the original input data first
	fmt.Println(originalData + "\n")

	// Result stores the movement sequence for each turn
	var movementLog []string
	antCounter := 1

	// Simulate ant movement along the selected shortest path
	for pathIndex, numAnts := range antsPerPath {
		for antNum := 0; antNum < numAnts; antNum++ {
			// Iterate through rooms in the path for this ant
			for stepIndex, roomName := range Var.AllVPaths[shortestPathIndex][pathIndex][1:] {
				targetIndex := stepIndex + antNum

				// Ensure the result slice has enough entries
				if targetIndex >= len(movementLog) {
					movementLog = append(movementLog, "")
				}

				// Add a space if there is already movement recorded for this turn
				if movementLog[targetIndex] != "" {
					movementLog[targetIndex] += " "
				}

				// Append the movement string "L<ant_number>-<room_name>"
				movementLog[targetIndex] += "L" + strconv.Itoa(antCounter) + "-" + roomName
			}
			antCounter++
		}
	}

	// Print the simulated movements turn by turn
	fmt.Print(strings.Join(movementLog, "\n"))
	fmt.Println()
}

// OrderAnts distributes the ants across the chosen set of paths.
// It assigns ants in a way that balances path length with the number of ants already on each path.
// The goal is to minimize the total number of turns required.
// Returns the number of turns and the distribution of ants per path.
func DistributeAnts(pathIndex int) (int, []int) {
	// Get all valid paths for this index
	paths := Var.AllVPaths[pathIndex]
	totalAnts := Var.AntsNumber

	// Prepare an array to track how many ants go on each path
	antsOnPath := make([]int, len(paths))

	currentPath := 0
	shortestPathLength := 0

	for totalAnts > 0 {
		// Compute total length = path length + ants already on it
		shortestPathLength = len(paths[currentPath]) + antsOnPath[currentPath]

		// Check if the next path is currently shorter
		if len(paths) > 1 &&
			currentPath+1 < len(paths) &&
			len(paths[currentPath+1])+antsOnPath[currentPath+1] < shortestPathLength {
			shortestPathLength = len(paths[currentPath+1]) + antsOnPath[currentPath+1]
			currentPath++
		}

		// Assign one ant to the selected path
		antsOnPath[currentPath]++
		totalAnts--

		// Loop back to the first path if at the end
		if currentPath == len(paths)-1 {
			currentPath = 0
		}
	}

	return shortestPathLength - 1, antsOnPath
}
