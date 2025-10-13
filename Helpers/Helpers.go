package Helpers

import "lemin/Var"

// RemovePathsLinks removes links between rooms that are part of a valid path.
// This ensures that paths do not overlap.
func RemovePathsLinks() {
	for pathIndex := 0; pathIndex < len(Var.VPaths); pathIndex++ {
		currentPath := Var.VPaths[pathIndex]

		// Iterate through each room in the path (except the last one)
		for roomIndex := 0; roomIndex < len(currentPath)-1; roomIndex++ {
			currentRoom := Var.Rooms[currentPath[roomIndex]]

			// Remove the link from the current room to the next one
			currentRoom.Links = RemoveLink(currentRoom.Links, currentPath[roomIndex+1])

			// Update the modified room in the map
			Var.Rooms[currentPath[roomIndex]] = currentRoom
		}
	}
}


// SaveBeforeInPath stores the previous room in the path for each room.
// This is used during backtracking to avoid revisiting rooms.
func SaveBeforeInPath() {
	// Get the most recently found valid path
	currentPath := Var.VPaths[len(Var.VPaths)-1]

	// For each room in the path (except the first and last),
	// record which room comes directly before it
	for roomIndex := 1; roomIndex < len(currentPath)-1; roomIndex++ {
		roomName := currentPath[roomIndex]
		previousRoomName := currentPath[roomIndex-1]

		room := Var.Rooms[roomName]
		room.BeforeInPath = previousRoomName
		Var.Rooms[roomName] = room
	}
}



// RemoveLink removes a specific link from a room's list of links.
func RemoveLink(roomLinks []string, targetRoom string) []string {
	for i := 0; i < len(roomLinks); i++ {
		if roomLinks[i] == targetRoom {
			// Remove the target room from the list of links
			if i+1 < len(roomLinks) {
				roomLinks = append(roomLinks[:i], roomLinks[i+1:]...)
			} else {
				roomLinks = roomLinks[:i]
			}
		}
	}
	return roomLinks
}


// CopyRoomsMap creates a deep copy of the Rooms map.
// This is used to reset the state of rooms during pathfinding.
func CopyRooms(originalRooms map[string]Var.Room) map[string]Var.Room {
	clonedRooms := make(map[string]Var.Room)

	for roomName, roomData := range originalRooms {
		// Deep copy the Links slice to avoid shared references
		copiedLinks := make([]string, len(roomData.Links))
		copy(copiedLinks, roomData.Links)

		// Create a new Room struct using the copied data
		clonedRooms[roomName] = Var.Room{
			Links:        copiedLinks,
			IsChecked:    roomData.IsChecked,
			BeforeInPath: roomData.BeforeInPath,
		}
	}

	return clonedRooms
}


// ResetIsChecked resets the IsChecked flag for all rooms.
// This prepares the rooms for a new pathfinding iteration.
func ResetIsChecked() {
	for index := range Var.Rooms {
		room := Var.Rooms[index]
		room.IsChecked = false
		Var.Rooms[index] = room
	}
}