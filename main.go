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
	bet	uint32
}

type Table struct {
	seats		[]*Seat
	board		[]deck.Card
	bigBlind	uint32
	ante		uint32
	totalPot	uint32
}

func newTable() *Table {
	table := Table{}
	table.bigBlind = 200
	table.ante = 200
	table.totalPot = 0

	return &table
}

func newSeat() *Seat {
	seat := Seat{}
	seat.dealer = false
	seat.chips = 10000
	seat.playing = true
	seat.bet = 0

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

func (table *Table) boardAction(deck *deck.Deck) error {
	
	if len(table.board) == 5 {
		return fmt.Errorf("The board is already full")
	} else if len(table.board) >= 3 {
		dealtCard, err := deck.Pop()
		if err != nil {
			return err
		}	
		table.board = append(table.board, dealtCard)
		return nil
	
	} else if len(table.board) == 0 {
		for range 3 {
			dealtCard, err := deck.Pop()
			if err != nil {
				return err
			}
			table.board = append(table.board, dealtCard)
		}
		return nil
	
	} else {
		return fmt.Errorf("There is something wrong with the board.")
	}	
}

func isFlush(cards []deck.Card) bool {
	suitCtr := 0
	for _, suit := range deck.Suits {
		for _, card := range cards {
			if card.Suit == suit {
				suitCtr += 1
			}
			if suitCtr == 5 {
				return true
			}
		}
	}
	return false
}

func containsRank(cards []deck.Card, rank string) (bool, deck.Card) {
	for _, card := range cards {
		if card.Rank == rank {
			return true, card
		}
	}
	return false, deck.Card{}
}

func findStraight(cards []deck.Card) []deck.Card {
	straight := []deck.Card {}
	for _, rank := range deck.Ranks {	
		containsRank, card := containsRank(cards, rank) 
		if containsRank {
			straight = append(straight, card) 
		} else {
			if len(straight) >= 5 {
				return straight
			}
			straight = []deck.Card {}
		}
	}
	return straight
}

func isStraightFlush(cards []deck.Card) bool {
	if isFlush(findStraight(cards)) {
		return true
	}
	return false
}

func isRoyalFlush(cards []deck.Card) bool {
	containsAce, _ := containsRank(findStraight(cards), "A")	
	if (containsAce && isFlush(findStraight(cards))) {
		return true
	}
	return false
}

func findStrongestSet(seat *Seat, board []deck.Card) {
		
}

func main() {
	deckInPlay := deck.NewSingleSuitDeck()//deck.NewDeck()
	//deckInPlay.Shuffle()
	
	table := newTable()

	seat1 := newSeat()
	//seat2 := newSeat()
	//seat3 := newSeat()

	table.Sit(seat1)
	//table.Sit(seat2)
	//table.Sit(seat3)
	
	table.dealCards(deckInPlay)
	
	for i, seat := range table.seats {
		fmt.Printf("Seat %d\n", i)
		for _, card := range seat.hand {
			card.Print()
		}
	}

	for range 3 {
		err := table.boardAction(deckInPlay)
		if err != nil {
			panic(err)
		}
	}
	
//	for _, card := range table.board {
//		card.Print()
//	}
	
	fmt.Printf("Table\n")
	for _, card := range table.board {
		card.Print()
	}
	//fmt.Printf(seat1.hand)
	commonCards := append(table.board, seat1.hand...)
	fmt.Printf("IsFlush: %b", isFlush(commonCards))
	straight := findStraight(commonCards)
	isStraight := false
	if len(straight) >= 5 {
		isStraight = true
	}
	fmt.Printf("IsStraight: %b", isStraight)
	fmt.Printf("IsRoyalFlush: %b", isRoyalFlush(commonCards))


}
