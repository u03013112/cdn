package app

import (
	"fmt"
	"os"
)

type Apper interface {
	Main(args []string) error
}

var vppMap = map[string]Apper{}

func Add(name string, v Apper) {
	vppMap[name] = v
}

func Run() error {

	args := os.Args
	if len(args) == 1 {
		fmt.Printf("app:\n")
		for k, _ := range vppMap {
			fmt.Printf("\t%s\n", k)
		}

		return nil
	}
	if v, ok := vppMap[args[1]]; ok {
		return v.Main(args[1:])
	}
	return nil
}
