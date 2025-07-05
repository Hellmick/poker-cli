package main

import (
	"fmt"
	
	deck "github.com/Hellmick/poker-cli/card-deck"
)

func main() {
	deck := deck.newDeck()
	deck.Shuffle()
	deck.Print()
}
