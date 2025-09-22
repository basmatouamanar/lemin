package main

import "sort"

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
	copy(finalPaths, paths)

	return finalPaths, ants
}

