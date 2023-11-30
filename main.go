package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type token int

const (
	BABY token = iota // # BABY (not implemented)
	AH                // + AH
	OH                // - OH
	YES               // > YES
	FUCK              // < FUCK
	MORE              // , MORE
	YEAH              // . YEAH
	AHH               // [ AHH
	OOH               // ] OOH
)

type parser struct {
	tm     map[string]token
	almost map[string]token
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("giff file")
	}
	f, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tm := map[string]token{
		"AH":   AH,
		"OH":   OH,
		"YES":  YES,
		"FUCK": FUCK,
		"MORE": MORE,
		"YEAH": YEAH,
		"AHH":  AHH,
		"OOH":  OOH,
	}
	p := &parser{
		tm:     tm,
		almost: make(map[string]token),
	}

	var tokens []token
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		tokens = append(tokens, p.tokenize(s.Text()))
	}

	var cells [30000]byte
	var pos int
	for i := 0; i < len(tokens); i++ {
		switch tok := tokens[i]; tok {
		case AH:
			cells[pos] += 1
		case OH:
			cells[pos] -= 1
		case MORE:
			r := bufio.NewReader(os.Stdin)
			char, err := r.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			cells[pos] = char
		case YEAH:
			fmt.Print(string(cells[pos]))
		case AHH:
			if cells[pos] == 0 {
				skip := 1
				for skip != 0 {
					i = i + 1
					if tokens[i] == AHH {
						skip += 1
					}
					if tokens[i] == OOH {
						skip -= 1
					}
				}
			}
		case OOH:
			if cells[pos] != 0 {
				skip := 1
				for skip != 0 {
					i = i - 1
					if tokens[i] == OOH {
						skip += 1
					}
					if tokens[i] == AHH {
						skip -= 1
					}
				}
			}
		case FUCK:
			pos -= 1
		case YES:
			pos += 1
		}
	}
}

func (p *parser) tokenize(src string) token {
	src = strings.ToUpper(src)
	if tok, ok := p.tm[src]; ok {
		return tok
	}

	var tok token
	var min int
	for cmp, t := range p.tm {
		dist := leven(src, cmp)
		if min == 0 || dist < min {
			min = dist
			tok = t
		}
	}
	p.almost[src] = tok
	return tok
}

// from https://people.cs.pitt.edu/~kirk/cs1501/Pruhs/Fall2006/Assignments/editdistance/Levenshtein%20Distance.htm
// not handling eq since we only come here if non default token
func leven(s, t string) int {
	n, m := len(s), len(t)
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
	}
	for i := 0; i <= n; i++ {
		d[i][0] = i
	}
	for j := 0; j <= m; j++ {
		d[0][j] = j
	}

	for i := 1; i <= n; i++ {
		si := s[i-1]
		for j := 1; j <= m; j++ {
			tj := t[j-1]
			cost := 0
			if si != tj {
				cost = 1
			}
			d[i][j] = min(
				d[i-1][j]+1,      // ins
				d[i][j-1]+1,      // del
				d[i-1][j-1]+cost, // sub
			)
		}
	}
	return d[n][m]
}

// might want to rename this given min is now identifier
func min(a, b, c int) int {
	res := a
	if b < res {
		res = b
	}
	if c < res {
		res = c
	}
	return res
}
