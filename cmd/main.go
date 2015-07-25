package main

import (
	"fmt"
	"log"

	"github.com/f2prateek/polygot"
	"github.com/google/go-github/github"
	"github.com/tj/docopt"
)

const (
	Usage = `Polygot.

Polygot will show you your Github Activity by language.

Usage:
  polygot <token>
  polygot -h | --help
  polygot --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	Version = "1.0.0"
)

var client *github.Client

func main() {
	arguments, err := docopt.Parse(Usage, nil, true, Version, false)
	check(err)

	token := arguments["<token>"].(string)
	p := polygot.New(token)

	counts, err := p.Counts()
	check(err)

	for lang, count := range counts {
		fmt.Printf("%d\t events for %s\n", count, lang)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
