package Find

import (
	"fmt"
	"os"

	"lemin/Helpers"
	"lemin/Var"
)

// FindValidPaths runs multiple BFS iterations to find all valid non-overlapping paths.
// It manages backtracking when conflicts are detected between new and old paths.
// The function removes conflicting links and keeps exploring until no more paths exist.
// If no valid paths are found, the program exits with an error.
func FindValidPaths() {
	removedLinks := make(map[string][]string)

	for {
		foundNewPath := BFS()
		if !foundNewPath {
			break
		}

		Helpers.SaveBeforeInPath()

		isBacktrack, fromRoom, toRoom := CheckBT()

		if isBacktrack {
			// Save the removed link so we can rebuild later
			removedLinks[fromRoom] = append(removedLinks[fromRoom], toRoom)

			// Store current valid paths except the last one (the conflicting one)
			Var.AllVPaths = append(Var.AllVPaths, Var.VPaths[:len(Var.VPaths)-1])

			// Reset paths and restore original room graph
			Var.VPaths = [][]string{}
			Var.Rooms = Helpers.CopyRooms(Var.OriginalRooms)

			// Remove backtracking links from the graph
			for src, targets := range removedLinks {
				for _, target := range targets {
					srcRoom := Var.Rooms[src]
					srcRoom.Links = Helpers.RemoveLink(srcRoom.Links, target)
					Var.Rooms[src] = srcRoom

					targetRoom := Var.Rooms[target]
					targetRoom.Links = Helpers.RemoveLink(targetRoom.Links, src)
					Var.Rooms[target] = targetRoom
				}
			}

		} else if Var.AntsNumber == len(Var.VPaths) {
			return
		} else {
			Helpers.RemovePathsLinks()
		}
	}

	if len(Var.VPaths) == 0 {
		fmt.Println("ERROR: No valid paths found!")
		os.Exit(0)
	}
}

// CheckBT inspects the last discovered path to see if it overlaps backwards
// with any previously found paths. It compares links in reverse order to detect conflicts.
// Returns a flag (true/false) and the nodes involved if a backtracking conflict is found.
func CheckBT() (bool, string, string) {
	// The most recently found path
	currentPath := Var.VPaths[len(Var.VPaths)-1]

	// All previously found valid paths
	previousPaths := Var.VPaths[:len(Var.VPaths)-1]

	// Map to store reversed connections from previous paths
	reversedConnections := make(map[string]string)

	// Build reversed connection mapping
	for i := len(previousPaths) - 1; i >= 0; i-- {
		for j := len(previousPaths[i]) - 2; j >= 1; j-- {
			reversedConnections[previousPaths[i][j]] = previousPaths[i][j-1]
		}
	}

	// Check if the current path conflicts with reversed links
	for i := 1; i < len(currentPath)-1; i++ {
		if reversedConnections[currentPath[i]] == currentPath[i+1] {
			// Backtracking conflict found
			return true, currentPath[i], currentPath[i+1]
		}
	}

	// No conflict detected
	return false, "", ""
}

// BFS performs a Breadth-First Search to explore paths from the start room to the end room.
// It expands all possible routes level by level until the end is found.
// Whenever a path to the end is discovered, it is saved as a valid path.
// The function avoids revisiting nodes and cleans up paths that are blocked.
// Returns true if at least one valid path is found, otherwise false.
func BFS() bool {
	start := Var.Start
	end := Var.End

	startRoom := Var.Rooms[start]
	startRoom.IsChecked = true
	Var.Rooms[start] = startRoom

	reverseStepUsed := false
	activePaths := [][]string{{start}}

	for len(activePaths) != 0 {
		// If the start room has no links, there are no paths to explore
		if len(Var.Rooms[start].Links) == 0 {
			return false
		}

		for i := 0; i < len(activePaths); i++ {
			validConnections := 0
			currentRoomName := activePaths[i][len(activePaths[i])-1]

			// Handle "reverse" traversal (backtracking)
			if Var.Rooms[currentRoomName].BeforeInPath != "" && !reverseStepUsed {
				previousRoom := Var.Rooms[currentRoomName].BeforeInPath
				validConnections++

				room := Var.Rooms[currentRoomName]
				room.IsChecked = false
				Var.Rooms[currentRoomName] = room

				activePaths[i] = append(activePaths[i], previousRoom)
				reverseStepUsed = true
				continue
			}

			// Explore all links (neighbors) of the current room
			for linkIndex, neighborName := range Var.Rooms[currentRoomName].Links {

				// ✅ Found a path that reaches the end room
				if neighborName == end {
					if validConnections == 0 {
						activePaths[i] = append(activePaths[i], neighborName)
					} else {
						activePaths = append(activePaths, append(activePaths[i][:len(activePaths[i])-1], neighborName))
					}
					Var.VPaths = append(Var.VPaths, activePaths[i])
					Helpers.ResetIsChecked()
					return true
				}

				// 🧭 Explore an unvisited neighbor
				if !Var.Rooms[neighborName].IsChecked {
					room := Var.Rooms[neighborName]
					room.IsChecked = true
					Var.Rooms[neighborName] = room

					validConnections++

					if validConnections == 1 {
						activePaths[i] = append(activePaths[i], neighborName)
					} else {
						newPath := make([]string, len(activePaths[i]))
						copy(newPath, activePaths[i])
						activePaths = append(activePaths, append(newPath[:len(newPath)-1], neighborName))
					}

					//  Dead-end: no valid links left
				} else if linkIndex == len(Var.Rooms[currentRoomName].Links)-1 && validConnections == 0 {
					if i+1 < len(activePaths) {
						activePaths = append(activePaths[:i], activePaths[i+1:]...)
					} else {
						activePaths = activePaths[:i]
					}
				}
			}
		}
	}

	// No valid path found
	return false
}
