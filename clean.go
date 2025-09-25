package main

import (
	"fmt"
	"strconv"
	"strings"
)

type nml struct {
	Name string
	X    int
	Y    int
}

func Clean(s []string) []string {
	ss := []string{}
	for i := 0; i < len(s); i++ {
		if !strings.HasPrefix(s[i], "##") && strings.HasPrefix(s[i], "#") {
			continue
		}
		ss = append(ss, s[i])
	}
	nml, err := strconv.Atoi(ss[0])
	if err != nil {
		fmt.Println("nr of ants is not valid")
	}
	fmt.Println("nbr ants",nml)
	// ///
	isstar := false
	isend := false
	for i := 0; i < len(ss); i++ {
		if isstar && ss[i] == "##start" {
			fmt.Println("errrr start")
			return []string{}
		}
		if ss[i] == "##start" {
			isstar = true
		}


		if isend && ss[i] == "##end" {
			fmt.Println("errrr end")
			return []string{}
		}
		if ss[i] == "##end" {
			isend = true
		}

	}
	if !isstar{
		fmt.Println("er starts")
			return []string{}
	}
	if !isend{
		fmt.Println("er ends")
			return []string{}
	}
	
	return ss
}
