package main

import (
	"flag"
	"fmt"

	"bytes"

	"bitbucket.org/cabify/cbot/trello_markov"
)

var (
	trelloApiKey   = flag.String("trelloKey", "", "Trello API Key")
	trelloApiToken = flag.String("trelloToken", "", "Trello API Token")
	trelloBoardID  = flag.String("trelloBoard", "wlJR9jXV", "Trello Board ID used for corpus")
)

func main() {
	flag.Parse()
	trello := trello_markov.NewTrelloCorpus(*trelloApiKey, *trelloApiToken)
	corpusBufs, err := trello.TextCorpus(*trelloBoardID)
	//corpusBufs, err := trello.TestCorpus()
	if err != nil {
		fmt.Println(err)
	}
	chain := trello_markov.NewChain(3)
	for corpus := range corpusBufs {
		chain.Build(bytes.NewBuffer(corpus))
	}

	fmt.Println(chain.Generate(32))
	fmt.Println(chain.Generate(32))
	fmt.Println(chain.Generate(32))
}
