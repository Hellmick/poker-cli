package main

import (
	"fmt"
	"sort"

	deck "github.com/Hellmick/poker-cli/deck"
)

type Seat struct {
	id      int
	hand    []deck.Card
	dealer  bool
	chips   uint32
	playing bool
	bet     uint32
}

type Table struct {
	seats    []*Seat
	board    []deck.Card
	bigBlind uint32
	ante     uint32
	totalPot uint32
}

func newTable() *Table {
	table := Table{}
	table.bigBlind = 200
	table.ante = 200
	table.totalPot = 0

	return &table
}

func newSeat(id int) *Seat {
	seat := Seat{}
	seat.dealer = false
	seat.chips = 10000
	seat.playing = true
	seat.bet = 0
	seat.id = id

	return &seat
}

func (table *Table) Sit(seat *Seat) {
	table.seats = append(table.seats, seat)
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

func containRank(cards []deck.Card, rank string) (bool, deck.Card) {
	for _, card := range cards {
		if card.Rank == rank {
			return true, card
		}
	}
	return false, deck.Card{}
}

func sortCardsByRank(cards []deck.Card) {
	sort.Slice(cards, func(i, j int) bool {
		return deck.RankIndex(cards[i].Rank) > deck.RankIndex((cards[j].Rank))
	})
}

func findFlush(cards []deck.Card) []deck.Card {
	for _, suit := range deck.Suits {
		suited := []deck.Card{}
		for _, card := range cards {
			if card.Suit == suit {
				suited = append(suited, card)
			}
		}
		if len(suited) >= 5 {
			sortCardsByRank(suited)
			return suited[:5]
		}
	}
	return nil
}

func findStraight(cards []deck.Card) []deck.Card {
	var straight []deck.Card
	sortCardsByRank(cards)
	for _, rank := range deck.Ranks {
		ok, card := containRank(cards, rank)
		if ok {
			straight = append(straight, card)
			if len(straight) == 4 && rank == "5" {
				if ok, ace := containRank(cards, "A"); ok {
					straight = append(straight, ace)
					return straight
				}
			}
			if len(straight) >= 5 {
				return straight[len(straight)-5:]
			}
		} else {
			straight = nil
		}
	}
	return nil
}

func findStraightFlush(cards []deck.Card) []deck.Card {
	for _, suit := range deck.Suits {
		suited := []deck.Card{}
		for _, card := range cards {
			if card.Suit == suit {
				suited = append(suited, card)
			}
		}
		if straight := findStraight(suited); straight != nil {
			return straight
		}
	}
	return nil
}

func isRoyalFlush(cards []deck.Card) bool {
	required := []string{"10", "J", "Q", "K", "A"}
	for _, rank := range required {
		if ok, _ := containRank(cards, rank); !ok {
			return false
		}
	}
	return true
}

func groupByRank(cards []deck.Card) map[string][]deck.Card {
	groups := make(map[string][]deck.Card)
	for _, c := range cards {
		groups[c.Rank] = append(groups[c.Rank], c)
	}
	return groups
}

type rankGroup struct {
	rank  string
	cards []deck.Card
}

func classifyGroups(groups map[string][]deck.Card) []rankGroup {
	var result []rankGroup

	for rank, cards := range groups {
		if len(cards) >= 2 {
			result = append(result, rankGroup{rank, cards})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if len(result[i].cards) != len(result[j].cards) {
			return len(result[i].cards) > len(result[j].cards)
		}
		return deck.RankIndex(result[i].rank) > deck.RankIndex(result[j].rank)
	})

	return result
}

func findQuads(groups []rankGroup) []deck.Card {
	if len(groups) > 0 && len(groups[0].cards) == 4 {
		return groups[0].cards
	}
	return nil
}

func findFullHouse(groups []rankGroup) []deck.Card {
	var trip []deck.Card
	var pair []deck.Card

	for _, group := range groups {
		if len(group.cards) >= 3 && trip == nil {
			trip = group.cards[:3]
		} else if len(group.cards) >= 2 && pair == nil {
			pair = group.cards[:2]
		}
	}

	if trip != nil && pair != nil {
		return append(trip, pair...)
	}

	return nil
}

func findTwoPair(groups []rankGroup) []deck.Card {
	pairs := [][]deck.Card{}

	for _, group := range groups {
		if len(group.cards) >= 2 {
			pairs = append(pairs, group.cards[:2])
		}
		if len(pairs) == 2 {
			return append(pairs[0], pairs[1]...)
		}
	}
	return nil
}

func findTrips(groups []rankGroup) []deck.Card {
	for _, group := range groups {
		if len(group.cards) == 3 {
			return group.cards
		}
	}
	return nil
}

func findPair(groups []rankGroup) []deck.Card {
	for _, group := range groups {
		if len(group.cards) == 2 {
			return group.cards
		}
	}
	return nil
}

func strongestSet(cards []deck.Card) ([]deck.Card, string) {
	groups := classifyGroups(groupByRank(cards))

	if sf := findStraightFlush(cards); sf != nil {
		if isRoyalFlush(sf) {
			return sf, "Royal Flush"
		}
		return sf, "Straight Flush"
	}
	if quads := findQuads(groups); quads != nil {
		return findQuads(groups), "Four of a Kind"
	}
	if fullHouse := findFullHouse(groups); fullHouse != nil {
		return findFullHouse(groups), "Full House"
	}
	if flush := findFlush(cards); flush != nil {
		return flush, "Flush"
	}
	if straight := findStraight(cards); straight != nil {
		return straight, "Straight"
	}
	if trips := findTrips(groups); trips != nil {
		return trips, "Three of a Kind"
	}
	if twoPair := findTwoPair(groups); twoPair != nil {
		return twoPair, "Two Pair"
	}
	if pair := findPair(groups); pair != nil {
		return pair, "Pair"
	}

	sortCardsByRank(cards)
	return []deck.Card{cards[0]}, "High Card"
}

func main() {
	deckInPlay := deck.NewDeck()
	deckInPlay.Shuffle()

	table := newTable()
	table.Sit(newSeat(len(table.seats)))
	table.Sit(newSeat(len(table.seats)))
	table.Sit(newSeat(len(table.seats)))

	table.dealCards(deckInPlay)

	fmt.Println("=== Hole Cards ===")
	for i, seat := range table.seats {
		fmt.Printf("Seat %d: ", i)
		for _, card := range seat.hand {
			fmt.Printf("%s%s ", card.Rank, card.Suit)
		}
		fmt.Println()
	}

	for range 3 {
		err := table.boardAction(deckInPlay)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("\n=== Board ===")
	for _, card := range table.board {
		fmt.Printf("%s%s ", card.Rank, card.Suit)
	}
	fmt.Println()

	for _, seat := range table.seats {
		allCards := append(seat.hand, table.board...)
		strongest, handType := strongestSet(allCards)

		fmt.Printf("\n=== Seat %d Best Hand ===\n", seat.id)
		for _, card := range strongest {
			fmt.Printf("%s%s ", card.Rank, card.Suit)
		}
		fmt.Printf("=> %s\n", handType)
	}

}
