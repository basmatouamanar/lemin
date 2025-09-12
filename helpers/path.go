package helpers

import (
	"strings"
)

func Path(lignes []string, rooms *[]*Room) (*Room, *Room) {
	var startRoom *Room
	var endRoom *Room

	for i := 0; i < len(lignes); i++ {
		if lignes[i] == "##start" {
			line := strings.TrimSpace(lignes[i+1])
			s := strings.Split(line, " ")
			c := s[0]
			for j := range *rooms {
				if (*rooms)[j].Name == c {
					startRoom = (*rooms)[j]
					(*rooms)[j].Start = true
				}
			}
		}
		if lignes[i] == "##end" {
			line := strings.TrimSpace(lignes[i+1])
			s := strings.Split(line, " ")
			c := s[0]
			for j := range *rooms {
				if (*rooms)[j].Name == c {
					endRoom = (*rooms)[j]
					(*rooms)[j].End = true
				}
			}
		}
		
	}
	
		return startRoom, endRoom
	
	
}
