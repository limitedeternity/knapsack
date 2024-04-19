package bounded

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	. "knapsack/common"
	. "knapsack/utils/functools"
)

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

func printSolution(sol *Solution) {
	fmt.Println("Taking:")
	for i, q := range sol.Quantities {
		if q > 0 {
			fmt.Printf("+ %s: %d/%d\n", items[i].Item, q, items[i].Pieces)
		}
	}

	fmt.Printf("Total value: %d\n", sol.Value)
	fmt.Printf("Total weight: %d\n", sol.Weight)
}

func TestSimpleSolver(t *testing.T) {
	var (
		sol Solution
		ok  bool
	)

	if sol, ok = solutions["SimpleSolver"]; !ok {
		sol = NewKnapsack[Item, *SimpleSolver]().WithCapacity(capacity).Pack(items)
		solutions["SimpleSolver"] = sol
	}

	printSolution(&sol)

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

func TestDPSolver(t *testing.T) {
	var (
		sol Solution
		ok  bool
	)

	if sol, ok = solutions["DPSolver"]; !ok {
		sol = NewKnapsack[Item, *DPSolver]().WithCapacity(capacity).Pack(items)
		solutions["DPSolver"] = sol
	}

	printSolution(&sol)

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
	t.Run("TestSimpleSolver", TestSimpleSolver)
	t.Run("TestDPSolver", TestDPSolver)

	require.True(t, reflect.DeepEqual(solutions["SimpleSolver"], solutions["DPSolver"]))
}

func TestItem_Unmarshal(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		data := []byte(`
item: FXIP
weight: 66
value: 66
pieces: 481
`)

		var item Item
		require.NoError(t, yaml.Unmarshal(data, &item))
		require.Equal(t, item, items[0])
	})

	t.Run("DefaultPieces", func(t *testing.T) {
		data := []byte(`
item: Tree
weight: 10
value: 1
`)

		var item Item
		require.NoError(t, yaml.Unmarshal(data, &item))
		require.Equal(t, item, Item{Item: "Tree", Weight: 10, Value: 1, Pieces: 1})
	})
}

func TestItem_Marshal(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		expected := `item: FXIP
weight: 66
value: 66
pieces: 481
`

		data, err := yaml.Marshal(items[0])
		require.NoError(t, err)
		require.Equal(t, expected, string(data))
	})
}
