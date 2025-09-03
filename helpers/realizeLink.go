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
}

func RealizeLink(data []string) {
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
		}
		if strings.Contains(line, "-") && !strings.HasPrefix(line, "#") {
			s := strings.Split(line, "-")
			if len(s) == 2 {
				room1 := rooms[s[0]]
				room2 := rooms[s[1]]
				if room1 != nil && room2 != nil {
					room1.linker = append(room1.linker, room2)
					room2.linker = append(room2.linker, room1)
				}
			}
		}
	}
	/*for _, r := range rooms["1"].linker {
		fmt.Println(r.Name, r.X, r.Y)
	}*/

}
