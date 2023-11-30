package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// go run main.go <input> > <output>
func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("giff file")
	}
	bs, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}

	m := map[byte]string{
		'+': "AH",
		'-': "OH",
		'>': "YES",
		'<': "FUCK",
		',': "MORE",
		'.': "YEAH",
		'[': "AHH",
		']': "OOH",
	}
	out := &strings.Builder{}
	for _, b := range bs {
		if tok, ok := m[b]; ok {
			// initial empty space does not matter anyway
			out.WriteString(" ")
			out.WriteString(tok)
		}
	}
	fmt.Println(out.String())
}
