package main

import (
	"fmt"
	deck "github.com/Hellmick/poker-cli/deck"
)

type Seat struct {
	hand	[]deck.Card
	dealer	bool
	chips	uint32
	playing bool
}

type Table struct {
	seats		[]*Seat
	board		[]deck.Card
	bigBlind	uint32
	ante		uint32
}

func newTable() *Table {
	table := Table{}
	table.bigBlind = 200
	table.ante = 200

	return &table
}

func newSeat() *Seat {
	seat := Seat{}
	seat.dealer = false
	seat.chips = 10000
	seat.playing = true
	
	return &seat
}

func (table *Table) Sit(seat *Seat) {
	table.seats = append(table.seats, newSeat()) 
	
	return
}

func (table *Table) dealCards(deck *deck.Deck) {
	for range 2 {	
		for _, seat := range table.seats {
			dealtCard, err := deck.Pop()
			if err != nil {
				panic(err)
			}
			seat.hand = append(seat.hand, dealtCard)			
		}
	}
}	

func main() {
	deck := deck.NewDeck()
	deck.Shuffle()
	
	table := newTable()

	seat1 := newSeat()
	seat2 := newSeat()

	table.Sit(seat1)
	table.Sit(seat2)

	table.dealCards(deck)
		
	for _, seat := range table.seats {
		fmt.Printf("Seat\n")
		for _, card := range seat.hand {
			card.Print()
		}
	}
}
