package unbounded

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	. "knapsack/common"
	ft "knapsack/utils/functools"
)

var (
	capacity = len(items)
	items    = []Item{
		{Item: "1m", Weight: 1, Value: 1},
		{Item: "2m", Weight: 2, Value: 5},
		{Item: "3m", Weight: 3, Value: 8},
		{Item: "4m", Weight: 4, Value: 9},
		{Item: "5m", Weight: 5, Value: 10},
		{Item: "6m", Weight: 6, Value: 17},
		{Item: "7m", Weight: 7, Value: 17},
		{Item: "8m", Weight: 8, Value: 20},
	}

	solutions = make(map[string]Solution, 1)
)

func printSolution(sol *Solution) {
	formattedSol := sol.String(
		struct{ ItemNames []string }{
			ItemNames: ft.Reduce(items,
				func(acc []string, val Item) []string {
					return append(acc, val.Item)
				}, nil,
			),
		})

	fmt.Print(formattedSol)
}

func TestDPSolver(t *testing.T) {
	t.Parallel()
	solverName := "DPSolver"

	var (
		sol Solution
		ok  bool
	)

	if sol, ok = solutions[solverName]; !ok {
		sol = NewKnapsack[Item, *DPSolver]().WithCapacity(capacity).Pack(items)
		solutions[solverName] = sol
	}

	printSolution(&sol)

	require.LessOrEqual(t, sol.Weight, capacity)
	require.Equal(t, sol.Value, 22)

	require.Equal(t, sol.Weight, ft.Reduce(
		ft.Zip(
			ft.Map(items, func(item Item) int { return item.Weight }),
			sol.Quantities,
		),
		func(acc int, pair ft.Pair[int, int]) int {
			itemWeight, itemQuantity := pair.First, pair.Second
			return acc + itemWeight*itemQuantity
		},
		0),
	)
}

func TestItem_Unmarshal(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		data := []byte(`
item: 1m
weight: 1
value: 1
`)

		var item Item
		require.NoError(t, yaml.Unmarshal(data, &item))
		require.Equal(t, item, items[0])
	})

	t.Run("IgnorePieces", func(t *testing.T) {
		data := []byte(`
item: 1m
weight: 1
value: 1
pieces: 481
`)

		var item Item
		require.NoError(t, yaml.Unmarshal(data, &item))
		require.Equal(t, item, items[0])
	})
}

func TestItem_Marshal(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		expected := `item: 1m
weight: 1
value: 1
`

		data, err := yaml.Marshal(items[0])
		require.NoError(t, err)
		require.Equal(t, expected, string(data))
	})
}
