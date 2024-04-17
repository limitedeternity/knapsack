package bounded

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, 0, len(ts))
	for _, t := range ts {
		result = append(result, fn(t))
	}

	return result
}

func Reduce[T, M any](ts []T, fn func(M, T) M, init M) M {
	acc := init
	for _, t := range ts {
		acc = fn(acc, t)
	}

	return acc
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func Zip[T, U any](ts []T, us []U) []Pair[T, U] {
	if len(ts) != len(us) {
		panic("slices have different length")
	}

	pairs := make([]Pair[T, U], 0, len(ts))
	for i := range ts {
		pairs = append(pairs, Pair[T, U]{ts[i], us[i]})
	}

	return pairs
}

func TestBruteForce(t *testing.T) {
	items := []Item{
		{Item: "USTIPS", Weight: int(math.Round(65.67)), Value: int(math.Round(65.67)), Pieces: 481},
		{Item: "FFIN", Weight: int(math.Round(562.39)), Value: int(math.Round(562.39)), Pieces: 9},
		{Item: "US_IT", Weight: int(math.Round(149.11)), Value: int(math.Round(149.11)), Pieces: 69},
		{Item: "UCITS", Weight: int(math.Round(87.64)), Value: int(math.Round(87.64)), Pieces: 433},
		{Item: "GLOB", Weight: int(math.Round(0.94)), Value: int(math.Round(0.94)), Pieces: 20996},
		{Item: "FALL", Weight: int(math.Round(92.36)), Value: int(math.Round(92.36)), Pieces: 89},
		{Item: "DevMarket", Weight: int(math.Round(107.19)), Value: int(math.Round(107.19)), Pieces: 342},
		{Item: "Em", Weight: int(math.Round(87.83)), Value: int(math.Round(87.83)), Pieces: 260},
		{Item: "Gold", Weight: int(math.Round(132.86)), Value: int(math.Round(132.86)), Pieces: 73},
	}

	sol := NewKnapsack[BruteForce](items).Pack(100_000)

	fmt.Println("Taking:")
	for i, q := range sol.Quantities {
		if q > 0 {
			fmt.Printf("+ %s: %d/%d\n", items[i].Item, q, items[i].Pieces)
		}
	}

	fmt.Printf("Total value: %d\n", sol.Value)
	fmt.Printf("Total weight: %d\n", sol.Weight)

	require.LessOrEqual(t, sol.Weight, 100_000)
	require.Equal(t, sol.Value, sol.Weight)

	require.Equal(t, sol.Weight, Reduce(
		Zip(
			Map(items, func(item Item) int { return item.Weight }),
			sol.Quantities,
		),
		func(acc int, pair Pair[int, int]) int {
			itemWeight, itemQuantity := pair.First, pair.Second
			return acc + itemWeight*itemQuantity
		},
		0),
	)
}
