package deck

import (
	"fmt"
	"math/rand"
)

type Card struct {
	Rank string
	Suit string
}

func NewCard(rank string, suit string) *Card {
	return &Card{Rank: rank, Suit: suit}
}

func (card *Card) Print() {
	fmt.Printf("%s%s\n", card.Rank, card.Suit)
}

type Deck struct {
	cards []Card
}

var Ranks = [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var Suits = [4]string{"S", "H", "D", "C"}

// the index is numerated from 1, if the function returns 0 the input rank is not in the rank list
func RankIndex(baseRank string) int {
	for i, rank := range Ranks {
		if baseRank == rank {
			return i + 1
		}
	}
	return 0
}

func NewDeck() *Deck {
	deck := Deck{}
	for _, rank := range Ranks {
		for _, suit := range Suits {
			card := NewCard(rank, suit)
			deck.Push(*card)
		}
	}
	return &deck
}

func NewSingleSuitDeck() *Deck {
	deck := Deck{}
	for _, rank := range Ranks {
		card := NewCard(rank, "S")
		deck.Push(*card)
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
	topCard, err := deck.Top()
	if err != nil {
		return Card{}, err
	}
	deck.cards = deck.cards[:len(deck.cards)-1]
	return topCard, nil
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
