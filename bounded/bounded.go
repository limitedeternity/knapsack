package bounded

type Item struct {
	Item   string
	Weight int
	Value  int
	Pieces int
}

type Knapsack[S Solver] struct {
	Limit int
	Items []Item
}

type Solution struct {
	Weight     int
	Value      int
	Quantities []int
}

type Solver interface {
	WithKnapsack(knapsack any) Solver
	Solve() Solution
}

func NewKnapsack[S Solver](limit int) *Knapsack[S] {
	return &Knapsack[S]{Limit: limit}
}

func (k *Knapsack[S]) Pack(items []Item) Solution {
	k.Items = items

	solver := new(S)
	sol := (*solver).WithKnapsack(k).Solve()

	solver = nil
	return sol
}
