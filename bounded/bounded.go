package bounded

type Item struct {
	Item   string
	Weight int
	Value  int
	Pieces int
}

type Knapsack[S Solver] struct {
	Items []Item
}

type Solution struct {
	Weight     int
	Value      int
	Quantities []int
}

type Solver interface {
	WithKnapsack(knapsack any) Solver
	Solve(limit int) Solution
}

func NewKnapsack[S Solver](items []Item) *Knapsack[S] {
	return &Knapsack[S]{Items: items}
}

func (k *Knapsack[S]) Pack(limit int) Solution {
	solver := new(S)
	sol := (*solver).WithKnapsack(k).Solve(limit)

	solver = nil
	return sol
}
