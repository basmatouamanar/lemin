package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidateData(str []string) bool {
	foundStart := false
	foundEnd := false
	var slice []string
	var slice2 []string

	for i := 0; i < len(str); i++ {
		ligne := strings.Split(str[i], " ")

		// ignorer les commentaires simples
		if strings.HasPrefix(str[i], "#") && !strings.HasPrefix(str[i], "##") {
			continue
		}

		// Cas des commandes ##start et ##end
		if len(ligne) == 1 {
			if ligne[0] == "##start" {
				if foundStart {
					fmt.Println("ERROR: multiple start")
					return false
				}
				foundStart = true
			}
			if ligne[0] == "##end" {
				if foundEnd {
					fmt.Println("ERROR: multiple end")
					return false
				}
				foundEnd = true
			}

			// Cas des liens entre rooms
			if strings.Contains(ligne[0], "-") && !strings.HasPrefix(ligne[0], "#") {
				rooms := strings.Split(ligne[0], "-")
				if len(rooms) != 2 {
					fmt.Println("this is an invalid link")
					return false
				}
				if rooms[0] == rooms[1] {
					fmt.Println("this is invalid path")
					return false
				}
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
					fmt.Println("link points to invalid room")
					return false
				}
			}
		}

		// Cas des rooms (3 arguments)
		if len(ligne) == 3 {
			s1 := ligne[0]
			// Vérifier doublon
			for _, v := range slice {
				if s1 == v {
					fmt.Println("this is a invalid argument")
					return false
				}
			}
			slice = append(slice, s1)

			arg2, err := strconv.Atoi(ligne[1])
			if err != nil {
				fmt.Println("this is a invalid argument")
				return false
			}
			arg3, err := strconv.Atoi(ligne[2])
			if err != nil {
				fmt.Println("this is a invalid argument")
				return false
			}

			s2 := strconv.Itoa(arg2)
			s3 := strconv.Itoa(arg3)

			for _, v := range slice2 {
				if s2+s3 == v {
					fmt.Println("this is a invalid argument")
					return false
				}
			}
			slice2 = append(slice2, s2+s3)
		}
	}

	return true
}
