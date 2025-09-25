package main

import (
	"fmt"
	"os"
	"strings"
)

func main(){
	if len(os.Args)!=2{
		fmt.Println("1 argument")
		return
	}
	content,err:=os.ReadFile(os.Args[1])
	if err!=nil{
		fmt.Println("err in reading file")
	}
	s:=string(content)
	ss:=strings.Split(s, "\n")
	
	z:=Clean(ss)
	check(z)
}