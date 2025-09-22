package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
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

func readInput(fileName string) (AntFarm, []string, error) {
	var farm AntFarm
	var state string
	var file []string

	fileHandle, err := os.Open(fileName)
	if err != nil {
		return AntFarm{}, []string{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for i := 0; scanner.Scan(); i++ {
		file = append(file, scanner.Text())
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				state = "start"
			} else if line == "##end" {
				state = "end"
			}
			continue
		}

		parts := strings.Fields(line)
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
// find all possible path from start to end using dfs recursive
func findPaths(tunnels []Tunnel, start, end string, path []string) {
	path = append(path, start)
	if start == end {
		newPath := make([]string, len(path))
		copy(newPath, path)
		Paths = append(Paths, newPath)
		return
	}

	for _, tunnel := range tunnels {
		if start == tunnel.From && !containsRoom(path, tunnel.To) {
			findPaths(tunnels, tunnel.To, end, path)
		} else if start == tunnel.To && !containsRoom(path, tunnel.From) {
			findPaths(tunnels, tunnel.From, end, path)
		}
	}
}

// check if the room is already in the path
func containsRoom(path []string, room string) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}


// sort paths according to their lenght  
func sortPaths(paths *[][]string) {
	sort.Slice(*paths, func(i, j int) bool {
		return len((*paths)[i]) < len((*paths)[j])
	})
}

//find two path that share the same rooms
func pathsOverlap(a, b []string) bool {
	set := make(map[string]bool)
	for _, room := range a[1 : len(a)-1] { // ignore start and end for overlap check
		set[room] = true
	}
	for _, room := range b[1 : len(b)-1] {
		if set[room] {
			return true
		}
	}
	return false
}

//keep only path do not intersect
func getDisjointPaths(paths [][]string) [][]string {
	var result [][]string
	for _, path := range paths {
		conflict := false
		for _, existing := range result {
			if pathsOverlap(existing, path) {
				conflict = true
				break
			}
		}
		if !conflict {
			result = append(result, path)
		}
	}
	return result
}

//select the best sets of path according to lenght and number of ants
func getBestPaths(paths [][]string, antCount int) [][]string {
	sortPaths(&paths)
	disjoint := getDisjointPaths(paths)

	best := [][]string{}
	minTurns := int(^uint(0) >> 1) // Max int

	for i := 1; i <= len(disjoint); i++ {
		set := disjoint[:i]
		turns := calculateTurns(set, antCount)
		if turns < minTurns {
			minTurns = turns
			best = set
		}
	}

	return best
}

//calculate how many turn we need  for all the ants to reash the end
func calculateTurns(paths [][]string, antCount int) int {
	pathLens := make([]int, len(paths))
	for i, path := range paths {
		pathLens[i] = len(path)
	}

	turns := 0
	for {
		total := 0
		for i := 0; i < len(paths); i++ {
			if turns >= pathLens[i]-1 {
				total += turns - pathLens[i] + 1
			}
		}
		if total >= antCount {
			break
		}
		turns++
	}
	return turns
}

// distribute the ants on path fairly
func distributeAnts(totalAnts int, paths [][]string) ([][]string, []int) {
	pathLens := make([]int, len(paths))
	for i, path := range paths {
		pathLens[i] = len(path)
	}

	ants := make([]int, len(paths))
	turns := calculateTurns(paths, totalAnts)

	remaining := totalAnts
	for i := 0; i < len(paths); i++ {
		ants[i] = turns - pathLens[i] + 1
		if ants[i] < 0 {
			ants[i] = 0
		}
		remaining -= ants[i]
	}

	i := 0
	for remaining > 0 {
		ants[i]++
		remaining--
		i = (i + 1) % len(paths)
	}

	// **Keep full paths including start room** for correct movement simulation
	finalPaths := make([][]string, len(paths))
	for i := 0; i < len(paths); i++ {
		finalPaths[i] = paths[i]
	}

	return finalPaths, ants
}


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

