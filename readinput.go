package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(fileName string) (AntFarm, []string, error) {
	var farm AntFarm
	var state string
	var file []string
	var foundStart bool
	var foundEnd bool
	var slice []string

	fileHandle, err := os.Open(fileName)
	if err != nil {
		return AntFarm{}, []string{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for i := 0; scanner.Scan(); i++ {
		file = append(file, scanner.Text())
		line := strings.TrimSpace(scanner.Text())
		if line == "##start" {
			if foundStart {
				return AntFarm{}, []string{}, fmt.Errorf("ERROR: multiple start")
			}
			foundStart = true
		}
		if line == "##end" {
			if foundEnd {
				return AntFarm{}, []string{}, fmt.Errorf("ERROR: multiple end")
			}
			foundEnd = true
		}
		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				state = "start"
			} else if line == "##end" {
				state = "end"
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			s1 := parts[0]
			// Vérifier doublon
			for _, v := range slice {
				if s1 == v {
					fmt.Println("this is a invalid argument")
					return AntFarm{}, []string{}, fmt.Errorf("invalid room")
				}
			}
			slice = append(slice, s1)

		}
		if len(parts) == 1 && i == 0 && !strings.Contains(parts[0], "-") {
			ants, err := strconv.Atoi(parts[0])
			if err != nil || ants < 1 {
				return AntFarm{}, []string{}, fmt.Errorf("invalid number of ants")
			}
			farm.Ants = ants
		} else if len(parts) == 3 {
			if strings.HasPrefix(parts[0], "L") {
				return AntFarm{}, []string{}, fmt.Errorf("room name %s cannot start with 'L'", parts[0])
			}

			room, err := parseRoom(parts)
			if err != nil {
				return AntFarm{}, []string{}, fmt.Errorf("failed to parse room: %v", parts[0])
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
		} else if len(parts) == 1 && strings.Contains(parts[0], "-") {
			tunnelParts := strings.Split(line, "-")
			if len(tunnelParts) == 2 {
				if tunnelParts[0] == "" || tunnelParts[1] == "" {
					return AntFarm{}, []string{}, fmt.Errorf("invalid format tunnel")
				}
				if tunnelParts[0] == tunnelParts[1] {
					return AntFarm{}, []string{}, fmt.Errorf("invalid format from == to")
				}

				tunnel := Tunnel{From: tunnelParts[0], To: tunnelParts[1]}
				if err := checkTunnel(farm.Tunnels, tunnel.From, tunnel.To); err != nil {
					return AntFarm{}, []string{}, err
				}
				farm.Tunnels = append(farm.Tunnels, tunnel)
			} else {
				return AntFarm{}, []string{}, fmt.Errorf("invalid data format")
			}
			rooms := strings.Split(parts[0], "-")
			found1, found2 := false, false
				for _, r := range slice {
					if r == rooms[0] {
						found1 = true
					}
					if r == rooms[1] {
						found2 = true
					}
				}
				if !found1 || !found2 {
					return AntFarm{}, []string{}, fmt.Errorf("link points to invalid room")
				}
		} else {
			return AntFarm{}, []string{}, fmt.Errorf("invalid data format")
		}
	}

	if err = scanner.Err(); err != nil {
		return AntFarm{}, []string{}, err
	}

	if farm.Ants <= 0 || len(farm.Tunnels) <= 0 {
		return AntFarm{}, []string{}, fmt.Errorf("invalid data format")
	}

	if farm.Start.Name == "" || farm.End.Name == "" {
		return AntFarm{}, []string{}, fmt.Errorf("invalid Path! Check Start and End Rooms")
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
