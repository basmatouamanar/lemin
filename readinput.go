package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(fileName string) (AntFarm, []string, error) {
	var farm AntFarm
	var state string
	var foundStart bool
	var foundEnd bool
	var slice []string

	data, err := os.ReadFile(fileName)
	if err != nil {
		return AntFarm{}, []string{}, fmt.Errorf("failed to open file: %v", err)
	}

	// Split file into lines
	lines := strings.Split(string(data), "\n")
	var file []string

	for i, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		file = append(file, rawLine)

		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				if foundStart {
					return AntFarm{}, []string{}, fmt.Errorf("ERROR: multiple start")
				}
				state = "start"
				foundStart = true
			} else if line == "##end" {
				if foundEnd {
					return AntFarm{}, []string{}, fmt.Errorf("ERROR: multiple end")
				}
				state = "end"
				foundEnd = true
			}
			continue
		}

		parts := strings.Fields(line)

		// Number of ants
		if len(parts) == 1 && i == 0 && !strings.Contains(parts[0], "-") {
			ants, err := strconv.Atoi(parts[0])
			if err != nil || ants < 1 {
				return AntFarm{}, []string{}, fmt.Errorf("invalid number of ants")
			}
			farm.Ants = ants

		// Room
		} else if len(parts) == 3 {
			if strings.HasPrefix(parts[0], "L") {
				return AntFarm{}, []string{}, fmt.Errorf("room name %s cannot start with 'L'", parts[0])
			}

			// Check for duplicate names
			for _, v := range slice {
				if v == parts[0] {
					return AntFarm{}, []string{}, fmt.Errorf("invalid room: duplicate room name %s", v)
				}
			}
			slice = append(slice, parts[0])

			room, err := parseRoom(parts)
			if err != nil {
				return AntFarm{}, []string{}, fmt.Errorf("failed to parse room: %v", err)
			}
			if err := checkRoom(farm.Rooms, room); err != nil {
				return AntFarm{}, []string{}, err
			}

			if state == "start" {
				if farm.Start != (Room{}) {
					return AntFarm{}, []string{}, fmt.Errorf("duplicate start room defined")
				}
				farm.Start, state = room, ""
			} else if state == "end" {
				if farm.End != (Room{}) {
					return AntFarm{}, []string{}, fmt.Errorf("duplicate end room defined")
				}
				farm.End, state = room, ""
			}

			farm.Rooms = append(farm.Rooms, room)

		// Tunnel
		} else if len(parts) == 1 && strings.Contains(parts[0], "-") {
			tunnelParts := strings.Split(parts[0], "-")
			if len(tunnelParts) != 2 || tunnelParts[0] == "" || tunnelParts[1] == "" {
				return AntFarm{}, []string{}, fmt.Errorf("invalid format tunnel")
			}
			if tunnelParts[0] == tunnelParts[1] {
				return AntFarm{}, []string{}, fmt.Errorf("invalid tunnel: from == to")
			}

			tunnel := Tunnel{From: tunnelParts[0], To: tunnelParts[1]}
			if err := checkTunnel(farm.Tunnels, tunnel.From, tunnel.To); err != nil {
				return AntFarm{}, []string{}, err
			}
			farm.Tunnels = append(farm.Tunnels, tunnel)

			// Make sure both rooms exist
			found1, found2 := false, false
			for _, r := range slice {
				if r == tunnel.From {
					found1 = true
				}
				if r == tunnel.To {
					found2 = true
				}
			}
			if !found1 || !found2 {
				return AntFarm{}, []string{}, fmt.Errorf("link points to invalid room")
			}

		}else {
			return AntFarm{}, []string{}, fmt.Errorf("invalid data format")
		}
	}

	// Final validation
	if farm.Ants <= 0 || len(farm.Tunnels) == 0 {
		return AntFarm{}, []string{}, fmt.Errorf("invalid data format")
	}
	if farm.Start.Name == "" || farm.End.Name == "" {
		return AntFarm{}, []string{}, fmt.Errorf("invalid path! Check start and end rooms")
	}

	return farm, file, nil
}

func parseRoom(parts []string) (Room, error) {
	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return Room{}, fmt.Errorf("invalid coordinates")
	}
	return Room{Name: parts[0], X: x, Y: y}, nil
}

func checkRoom(rooms []Room, room Room) error {
	for _, r := range rooms {
		if r.Name == room.Name {
			return fmt.Errorf("room %s is duplicate", r.Name)
		}
		if r.X == room.X && r.Y == room.Y {
			return fmt.Errorf("duplicate coordinates detected for rooms '%s' and '%s'", r.Name, room.Name)
		}
	}
	return nil
}

func checkTunnel(tunnels []Tunnel, from, to string) error {
	for _, t := range tunnels {
		if (t.From == from && t.To == to) || (t.From == to && t.To == from) {
			return fmt.Errorf("duplicate tunnel defined for rooms '%s' and '%s'", from, to)
		}
	}
	return nil
}
