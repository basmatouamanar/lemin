package Parsing

import (
	"fmt"
	"lemin/Var"
	"strconv"
	"strings"
)


// ParsingData parses the input data and initializes global variables.
// It validates the number of ants, rooms, and links.
func ParseData(input string) error {
	var err error

	// Track used room names and coordinates to prevent duplicates
	roomNames := make(map[string]bool)
	roomCoordinates := make(map[string]bool)

	// Split the input data into lines
	lines := strings.Split(input, "\n")

	// Parse number of ants
	Var.AntsNumber, err = strconv.Atoi(lines[0])
	if err != nil || Var.AntsNumber <= 0 {
		return fmt.Errorf("ERROR: invalid number of ants")
	}

	// Process each line in the input data
	for lineIndex := 1; lineIndex < len(lines); lineIndex++ {
		line := lines[lineIndex]

		// Room definition (format: name x y)
		if fields := strings.Split(line, " "); len(fields) == 3 {
			if err := checkIsValid(roomNames, roomCoordinates, fields, lineIndex); err != nil {
				return err
			}

		// ##start directive
		} else if line == "##start" {
			if lineIndex+1 >= len(lines) {
				return fmt.Errorf("ERROR: start-flag trailing at end of file at line %d", lineIndex+1)
			}
			if Var.Start != "" {
				return fmt.Errorf("ERROR: multiple ##start flags found at line %d", lineIndex+1)
			}

			nextFields := strings.Split(lines[lineIndex+1], " ")
			if err := checkIsValid(roomNames, roomCoordinates, nextFields, lineIndex+1); err != nil {
				return err
			}

			Var.Start = nextFields[0]
			lineIndex++ // Skip next line (processed as start room)

		// ##end directive
		} else if line == "##end" {
			if lineIndex+1 >= len(lines) {
				return fmt.Errorf("ERROR: end-flag trailing at end of file at line %d", lineIndex+1)
			}
			if Var.End != "" {
				return fmt.Errorf("ERROR: multiple ##end flags found at line %d", lineIndex+1)
			}

			nextFields := strings.Split(lines[lineIndex+1], " ")
			if len(nextFields) != 3 {
				return fmt.Errorf("ERROR: invalid room initialization after ##end flag at line %d", lineIndex+1)
			}

			if err := checkIsValid(roomNames, roomCoordinates, nextFields, lineIndex+1); err != nil {
				return err
			}

			Var.End = nextFields[0]
			roomCoordinates[strings.Join(nextFields[1:], " ")] = true
			lineIndex++ // Skip next line (processed as end room)

		// Comment line — ignore
		} else if strings.HasPrefix(line, "#") {
			continue

		// Link definition (format: room1-room2)
		} else if linkParts := strings.Split(line, "-"); len(linkParts) == 2 {
			if !roomNames[linkParts[0]] || !roomNames[linkParts[1]] {
				return fmt.Errorf("ERROR: linking uninitialized room at line %d", lineIndex+1)
			}
			if err := fillRoomData(line); err != nil {
				return fmt.Errorf("ERROR: invalid link format at line %d", lineIndex+1)
			}

		// Empty line — ignore
		} else if line == "" {
			continue

		// Invalid format
		} else {
			return fmt.Errorf("ERROR: invalid data format at line %d", lineIndex+1)
		}
	}

	// Final validation
	if Var.Start == "" || Var.End == "" {
		return fmt.Errorf("ERROR: missing start or end room")
	}
	if Var.Start == Var.End {
		return fmt.Errorf("ERROR: start and end rooms cannot be the same")
	}

	return nil
}



// checkIsValid validates room initialization.
// It ensures that room names and coordinates are unique and valid.
func checkIsValid(roomNames map[string]bool, roomCordinations map[string]bool, spaceSplit []string, i int) error {
	var err error
	
	if roomNames[spaceSplit[0]] {
		return fmt.Errorf("ERROR: Duplicated room initialazation %d", i+1)
	}
	roomNames[spaceSplit[0]] = true

	if _, err = strconv.Atoi(spaceSplit[1]); err != nil {
		return fmt.Errorf("ERROR: invalid coordinates for room at line %d", i+1)
	}
	
	if _, err = strconv.Atoi(spaceSplit[2]); err != nil {
		return fmt.Errorf("ERROR: invalid coordinates for room at line %d", i+1)
	}

	if roomCordinations[strings.Join(spaceSplit[1:], " ")] {
		return fmt.Errorf("ERROR: duplicate coordinates for room at line %d", i+1)
	}

	roomCordinations[strings.Join(spaceSplit[1:], " ")] = true
	return nil
}

//  stores links between rooms in the global Rooms map.
// It ensures that links are valid and not duplicated.
func fillRoomData(str string) error {
    split := strings.Split(str, "-")

    if split[0] == split[1] {
        return fmt.Errorf("room %s can't link to itself", split[0])
    }

    for _, link := range Var.Rooms[split[0]].Links {
        if link == split[1] {
            return fmt.Errorf("duplicate link between %s and %s", split[0], split[1])
        }
    }

    if (split[0] == Var.Start && split[1] == Var.End) || (split[0] == Var.End && split[1] == Var.Start) {
        Var.VPaths = [][]string{{Var.Start, Var.End}}
        return nil
    }

    roomA := Var.Rooms[split[0]]
    roomA.Links = append(roomA.Links, split[1])
    Var.Rooms[split[0]] = roomA

    roomB := Var.Rooms[split[1]]
    roomB.Links = append(roomB.Links, split[0])
    Var.Rooms[split[1]] = roomB

    return nil
}