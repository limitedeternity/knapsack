package bounded

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	. "knapsack/common"
	ft "knapsack/utils/functools"
	kl "knapsack/utils/keylock"
)

var (
	capacity = 100_000
	items    = []Item{
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

	solutions = struct {
		Store map[string]Solution
		Lock  *kl.KeyLock
	}{Store: make(map[string]Solution, 2), Lock: kl.NewKeyLock()}

	solverStarted = make(chan string, 200)
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

func TestSimpleSolver(t *testing.T) {
	t.Parallel()
	solverName := "SimpleSolver"

	var (
		sol Solution
		ok  bool
	)

	canceled, unlock := solutions.Lock.LockKeys([]string{solverName}, nil)
	solverStarted <- solverName

	require.False(t, canceled)

	func() {
		defer unlock()

		if sol, ok = solutions.Store[solverName]; !ok {
			sol = NewKnapsack[Item, *SimpleSolver]().WithCapacity(capacity).Pack(items)
			solutions.Store[solverName] = sol
		}
	}()

	printSolution(&sol)

	require.LessOrEqual(t, sol.Weight, capacity)
	require.Equal(t, sol.Value, sol.Weight)

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

func TestDPSolver(t *testing.T) {
	t.Parallel()
	solverName := "DPSolver"

	var (
		sol Solution
		ok  bool
	)

	canceled, unlock := solutions.Lock.LockKeys([]string{solverName}, nil)
	solverStarted <- solverName

	require.False(t, canceled)

	func() {
		defer unlock()

		if sol, ok = solutions.Store[solverName]; !ok {
			sol = NewKnapsack[Item, *DPSolver]().WithCapacity(capacity).Pack(items)
			solutions.Store[solverName] = sol
		}
	}()

	printSolution(&sol)

	require.LessOrEqual(t, sol.Weight, capacity)
	require.Equal(t, sol.Value, sol.Weight)

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

func TestSameResults(t *testing.T) {
	t.Parallel()
	t.Run("TestSimpleSolver", TestSimpleSolver)
	t.Run("TestDPSolver", TestDPSolver)

	expectedSolvers := map[string]bool{"SimpleSolver": false, "DPSolver": false}
	for !ft.Reduce(maps.Values(expectedSolvers), func(acc bool, val bool) bool { return acc && val }, true) {
		solverName := <-solverStarted

		if seen, expected := expectedSolvers[solverName]; expected && !seen {
			expectedSolvers[solverName] = true
		} else {
			solverStarted <- solverName
		}
	}

	canceled, unlock := solutions.Lock.LockKeys(maps.Keys(expectedSolvers), nil)
	require.False(t, canceled)

	defer unlock()

	require.NotZero(t, solutions.Store["SimpleSolver"])
	require.NotZero(t, solutions.Store["DPSolver"])
	require.Empty(t, cmp.Diff(solutions.Store["SimpleSolver"], solutions.Store["DPSolver"]))
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
