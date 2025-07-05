package deck

import (
	"fmt"
	"math/rand"
)


type Card struct {
	rank	string
	suit	string
}

func newCard(rank string, suit string) *Card {
	card := Card{rank: rank, suit:suit}
	return &card
}

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	deck := Deck{}
	ranks := [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := [4]string{"S","H","D","C"}

	for _, rank := range ranks {
		for _, suit := range suits {
			card := newCard(rank, suit)
			deck.Push(*card)
		}
	}
	return &deck
}

func (deck *Deck) Push(card Card) {
	deck.cards = append(deck.cards, card)
}

func (deck *Deck) Pop() (Card, error) {
	if deck.IsEmpty() {
		return Card{}, fmt.Errorf("Deck is empty")
	}
	deck.cards = deck.cards[:len(deck.cards)-1]
	return deck.Top()
}

func (deck *Deck) IsEmpty() bool {
	if len(deck.cards) == 0 {
		return true
	}
	return false		
}

func (deck *Deck) Top() (Card, error) {
	if deck.IsEmpty() {
		return Card{}, fmt.Errorf("Stack is empty")
	}
	return deck.cards[len(deck.cards)-1], nil
}

func (deck *Deck) Print() {
	for _, card := range deck.cards {
		fmt.Print(card, " ")
	}
	fmt.Println()
}

func (deck *Deck) Shuffle() {
	for i := range deck.cards {
		j := rand.Intn(i + 1)
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i] 
	}
}

