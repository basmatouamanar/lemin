package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type lbyout struct {
	Name string
}
type trgan struct {
	From, To string
}
type farm struct {
	nmel  int
	room  map[string]*lbyout
	start string
	end   string
	trig  []trgan
	jiran map[string][]string
}
type massar []string

type move struct {
	antID int
	room  string
	turn  int
}

func NewAntFarm() *farm {
	return &farm{
		room:  make(map[string]*lbyout),
		jiran: make(map[string][]string),
		trig:  make([]trgan, 0),
	}
}

func main() {
	furma := NewAntFarm()
	err := furma.parseinput()
	if err != nil {
		log.Fatal("can't parse the input")
	}

	fmt.Println(furma.nmel)
	fmt.Println(furma.jiran)
	fmt.Println(furma.trig)
	fmt.Println(furma.room)
	fmt.Println(furma.end)

	furma.simulate(furma.findPaths())
}

func (firma *farm) parseinput() error {
	file, err := os.Open("test.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	read := bufio.NewScanner(file)
	if !read.Scan() {
		return fmt.Errorf("can't read number of ants")
	}

	NumOfAnts, err := strconv.Atoi(read.Text())
	if err != nil {
		return fmt.Errorf("invalid number of ants: %v", err)
	}
	firma.nmel = NumOfAnts

	istart := false
	isend := false

	for read.Scan() {
		line := strings.TrimSpace(read.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				istart = true
				continue
			}
			if line == "##end" {
				isend = true
				continue
			}
			continue
		}
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return fmt.Errorf("invalid tunnel format: %s", line)
			}
			from, to := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			firma.trig = append(firma.trig, trgan{From: from, To: to})
			firma.jiran[from] = append(firma.jiran[from], to)
			firma.jiran[to] = append(firma.jiran[to], from)
		} else {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return fmt.Errorf("invalid room format: %s", line)
			}
			name := parts[0]
			lbit := &lbyout{Name: name}
			firma.room[name] = lbit

			if istart {
				firma.start = name
				istart = false
			}
			if isend {
				firma.end = name
				isend = false
			}
		}
	}

	if firma.start == "" || firma.end == "" {
		return fmt.Errorf("start or end room not defined")
	}

	return nil
}

/**************************************************/
func (f *farm) findPaths() [][]string {
	var paths [][]string
	queue := [][]string{{f.start}} //

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		last := path[len(path)-1]

		if last == f.end {
			paths = append(paths, path)
			continue
		}

		for _, neigh := range f.jiran[last] {
			if !contains(path, neigh) {
				newPath := append([]string{}, path...)
				newPath = append(newPath, neigh)
				queue = append(queue, newPath)
			}
		}
	}

	return paths
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

/****************************************/

func (f *farm) simulate(paths [][]string) {
	// Sort paths by length (shortest first)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Track the position of each ant
	type ant struct {
		id    int
		path  []string
		index int // Current position in the path
	}
	var ants []ant

	// Assign ants to paths dynamically
	antID := 1
	for f.nmel > 0 {
		for _, path := range paths {
			if f.nmel == 0 {
				break
			}
			ants = append(ants, ant{id: antID, path: path, index: 0})
			antID++
			f.nmel--
		}
	}

	// Simulate movement turn by turn
	turn := 0
	for len(ants) > 0 {
		var nextAnts []ant
		fmt.Printf("Turn %d:\n", turn)

		for _, a := range ants {
			// Move the ant to the next room
			a.index++
			if a.index < len(a.path) {
				fmt.Printf("Ant %d -> %s\n", a.id, a.path[a.index])
				nextAnts = append(nextAnts, a)
			}
		}

		ants = nextAnts
		turn++
	}
}
