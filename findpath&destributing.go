package main

import "sort"

// find all possible path from start to end using dfs recursive
const maxDFSPaths = 2000 // safety cap to avoid exponential blowup on huge graphs

func findPaths(tunnels []Tunnel, start, end string, path []string) {
	if len(Paths) >= maxDFSPaths {
		return
	}
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

// generate all disjoint groups (all subsets with no room overlap excluding start/end)
func generateDisjointGroups(paths [][]string) [][][]string {
	// Precompute room sets for each path (excluding start/end)
	pathRooms := make([]map[string]bool, len(paths))
	for i, p := range paths {
		rooms := make(map[string]bool)
		if len(p) > 2 {
			for _, r := range p[1 : len(p)-1] {
				rooms[r] = true
			}
		}
		pathRooms[i] = rooms
	}

	var groups [][][]string
	usedRooms := make(map[string]bool)

	var backtrack func(idx int, current [][]string)
	backtrack = func(idx int, current [][]string) {
		if idx == len(paths) {
			if len(current) > 0 {
				snapshot := make([][]string, len(current))
				copy(snapshot, current)
				groups = append(groups, snapshot)
			}
			return
		}

		// Option 1: skip current path
		backtrack(idx+1, current)

		// Option 2: take current path if no overlap
		canTake := true
		for r := range pathRooms[idx] {
			if usedRooms[r] {
				canTake = false
				break
			}
		}
		if canTake {
			// mark rooms
			for r := range pathRooms[idx] {
				usedRooms[r] = true
			}
			current = append(current, paths[idx])
			backtrack(idx+1, current)
			// unmark rooms
			current = current[:len(current)-1]
			for r := range pathRooms[idx] {
				delete(usedRooms, r)
			}
		}
	}

	backtrack(0, nil)
	return groups
}

//select the best set of paths considering all disjoint groups, ant count and path lengths
func getBestPaths(paths [][]string, antCount int) [][]string {
	if len(paths) == 0 {
		return nil
	}
	sortPaths(&paths)

	// Cap number of paths to consider for combinatorics
	const maxPathsForGrouping = 16
	if len(paths) > maxPathsForGrouping {
		paths = paths[:maxPathsForGrouping]
	}

	groups := generateDisjointGroups(paths)
	if len(groups) == 0 {
		return nil
	}

	best := [][]string{}
	minTurns := int(^uint(0) >> 1)
	bestTotalLen := int(^uint(0) >> 1)

	for _, g := range groups {
		turns := calculateTurns(g, antCount)
		if turns < minTurns {
			minTurns = turns
			best = g
			// compute total length for tiebreaker
			total := 0
			for _, p := range g {
				total += len(p)
			}
			bestTotalLen = total
			continue
		}
		if turns == minTurns {
			// tiebreaker: prefer more paths; if equal, prefer smaller total length
			if len(g) > len(best) {
				best = g
				total := 0
				for _, p := range g {
					total += len(p)
				}
				bestTotalLen = total
			} else if len(g) == len(best) {
				total := 0
				for _, p := range g {
					total += len(p)
				}
				if total < bestTotalLen {
					best = g
					bestTotalLen = total
				}
			}
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
