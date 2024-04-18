package unbounded

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	. "knapsack/functools"
)

var (
	items = []Item{
		{Item: "1m", Weight: 1, Value: 1},
		{Item: "2m", Weight: 2, Value: 5},
		{Item: "3m", Weight: 3, Value: 8},
		{Item: "4m", Weight: 4, Value: 9},
		{Item: "5m", Weight: 5, Value: 10},
		{Item: "6m", Weight: 6, Value: 17},
		{Item: "7m", Weight: 7, Value: 17},
		{Item: "8m", Weight: 8, Value: 20},
	}
	capacity  = len(items)
	solutions = map[string]Solution{}
)

func printSolution(sol Solution) {
	fmt.Println("Taking:")
	for i, q := range sol.Quantities {
		if q > 0 {
			fmt.Printf("+ %s: %d\n", items[i].Item, q)
		}
	}

	fmt.Printf("Total value: %d\n", sol.Value)
	fmt.Printf("Total weight: %d\n", sol.Weight)
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
	require.Equal(t, sol.Value, 22)

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
