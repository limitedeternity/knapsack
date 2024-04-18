package bounded

import (
	"fmt"
	"math"
	"reflect"
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

var (
	items = []Item{
		{Item: "FXIP", Weight: int(math.Round(65.67)), Value: int(math.Round(65.67)), Pieces: 481},
		{Item: "FXKZ", Weight: int(math.Round(562.39)), Value: int(math.Round(562.39)), Pieces: 9},
		{Item: "FXIM", Weight: int(math.Round(149.11)), Value: int(math.Round(149.11)), Pieces: 69},
		{Item: "FXUS", Weight: int(math.Round(87.64)), Value: int(math.Round(87.64)), Pieces: 433},
		{Item: "FXRW", Weight: int(math.Round(0.94)), Value: int(math.Round(0.94)), Pieces: 20996},
		{Item: "FXFA", Weight: int(math.Round(92.36)), Value: int(math.Round(92.36)), Pieces: 89},
		{Item: "FXDM", Weight: int(math.Round(107.19)), Value: int(math.Round(107.19)), Pieces: 342},
		{Item: "FXEM", Weight: int(math.Round(87.83)), Value: int(math.Round(87.83)), Pieces: 260},
		{Item: "FXGD", Weight: int(math.Round(132.86)), Value: int(math.Round(132.86)), Pieces: 73},
	}
	capacity  = 100_000
	solutions = map[string]Solution{}
)

func printSolution(sol Solution) {
	fmt.Println("Taking:")
	for i, q := range sol.Quantities {
		if q > 0 {
			fmt.Printf("+ %s: %d/%d\n", items[i].Item, q, items[i].Pieces)
		}
	}

	fmt.Printf("Total value: %d\n", sol.Value)
	fmt.Printf("Total weight: %d\n", sol.Weight)
}

func TestBruteForce(t *testing.T) {
	var (
		sol Solution
		ok  bool
	)

	if sol, ok = solutions["BruteForce"]; !ok {
		sol = NewKnapsack[BruteForce](capacity).Pack(items)
		solutions["BruteForce"] = sol
	}

	printSolution(sol)

	require.LessOrEqual(t, sol.Weight, capacity)
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

func TestDynamic(t *testing.T) {
	var (
		sol Solution
		ok  bool
	)

	if sol, ok = solutions["Dynamic"]; !ok {
		sol = NewKnapsack[Dynamic](capacity).Pack(items)
		solutions["Dynamic"] = sol
	}

	printSolution(sol)

	require.LessOrEqual(t, sol.Weight, capacity)
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

func TestSameResults(t *testing.T) {
	t.Run("TestBruteForce", TestBruteForce)
	t.Run("TestDynamic", TestDynamic)

	require.True(t, reflect.DeepEqual(solutions["BruteForce"], solutions["Dynamic"]))
}
