package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Room represents a room in the ant farm
type Room struct {
	Name string
	X int
	 Y int
}

// Edge represents a tunnel between two rooms
type Edge struct {
	From string
	 To string
}

// AntFarm represents the entire ant farm structure
type AntFarm struct {
	NumAnts   int
	Rooms     map[string]*Room
	Start     string
	End       string
	Adjacency map[string][]string
	Edges     []Edge
}

// Path represents a sequence of rooms from start to end
type Path []string

// NewAntFarm creates a new ant farm instance
func NewAntFarm() *AntFarm {
	return &AntFarm{
		Rooms:     make(map[string]*Room),
		Adjacency: make(map[string][]string),
		Edges:     make([]Edge, 0),
	}
}

// ParseInput reads and parses the input format
func (af *AntFarm) ParseInput() error {
	scanner := bufio.NewScanner(os.Stdin)
	
	// Read number of ants
	if !scanner.Scan() {
		return fmt.Errorf("failed to read number of ants")
	}
	
	numAnts, err := strconv.Atoi(scanner.Text())
	if err != nil || numAnts <= 0 {
		return fmt.Errorf("invalid number of ants: %s", scanner.Text())
	}
	af.NumAnts = numAnts
	
	// Parse rooms and tunnels
	isStart := false
	isEnd := false
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			// Check for special commands
			if line == "##start" {
				isStart = true
				continue
			}
			if line == "##end" {
				isEnd = true
				continue
			}
			continue
		}
		
		// Check if it's a tunnel (contains '-')
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return fmt.Errorf("invalid tunnel format: %s", line)
			}
			
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			
			// Add edge
			af.Edges = append(af.Edges, Edge{From: from, To: to})
			af.Adjacency[from] = append(af.Adjacency[from], to)
			af.Adjacency[to] = append(af.Adjacency[to], from)
		} else {
			// It's a room
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return fmt.Errorf("invalid room format: %s", line)
			}
			
			name := parts[0]
			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			
			if err1 != nil || err2 != nil {
				return fmt.Errorf("invalid coordinates in room: %s", line)
			}
			
			room := &Room{Name: name, X: x, Y: y}
			af.Rooms[name] = room
			
			// Set start/end rooms
			if isStart {
				af.Start = name
				isStart = false
			}
			if isEnd {
				af.End = name
				isEnd = false
			}
		}
	}
	
	return scanner.Err()
}