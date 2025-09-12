package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

type Room struct {
	Name   string
	X, Y   int
	linker []*Room
	Start  bool
    End    bool
}

func RealizeLink(data []string, Rooms *[]*Room) {
	rooms := make(map[string]*Room)
	for _, line := range data {
		line = strings.TrimSpace(line)
		c := strings.Split(line, " ")
		if len(c) == 3 {
			x, err1 := strconv.Atoi(c[1])
			y, err2 := strconv.Atoi(c[2])
			if err1 != nil || err2 != nil {
				fmt.Println("somthing is wrong")
				return
			}
			room := &Room{
				Name: c[0],
				X:    x,
				Y:    y,
			}
			rooms[room.Name] = room

			*Rooms = append(*Rooms, room)
		}
		if strings.Contains(line, "-") && !strings.HasPrefix(line, "#") {
			s := strings.Split(line, "-")
			if len(s) == 2 {
				room1 := rooms[s[0]]
				room2 := rooms[s[1]]
				for i := range *Rooms {
					if (*Rooms)[i].Name == s[0] {
						(*Rooms)[i].linker = append((*Rooms)[i].linker, room2)
					}
					if (*Rooms)[i].Name == s[1] {
						(*Rooms)[i].linker = append((*Rooms)[i].linker, room1)
					}
				}
			}	
		}
		
	}
	for _, r := range *Rooms {
	fmt.Print(r.Name, ": ")
	for _, l := range r.linker {
		fmt.Print(l.Name, " ")
	}
	fmt.Println()
}


}
// (*Rooms)[i].linker = original.linker
