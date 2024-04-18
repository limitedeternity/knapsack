package bounded

import (
	. "knapsack/common"
)

type SimpleSolver struct {
	SolverBase[Item, *SimpleSolver]
	cache map[key]Solution
}

type key struct {
	weight   int
	position int
}

func (s *SimpleSolver) GetBase() ISolverBase {
	return &s.SolverBase
}

func (s *SimpleSolver) Solve() Solution {
	s.cache = make(map[key]Solution)
	sol := s.impl(s.Knapsack.Limit, len(s.Knapsack.Items)-1)

	s.cache = nil
	return sol
}

func (s *SimpleSolver) impl(limit int, position int) Solution {
	if position < 0 || limit <= 0 {
		return Solution{}
	}

	key_ := key{limit, position}
	if sol, ok := s.cache[key_]; ok {
		return sol
	}

	bestQ, best := 0, Solution{}
	for q := 0; q*s.Knapsack.Items[position].Weight <= limit && q <= s.Knapsack.Items[position].Pieces; q++ {
		sol := s.impl(limit-q*s.Knapsack.Items[position].Weight, position-1)
		sol.Value += q * s.Knapsack.Items[position].Value

		if sol.Value > best.Value {
			bestQ, best = q, sol
		}
	}

	if bestQ > 0 {
		old := best.Quantities
		best.Quantities = make([]int, len(s.Knapsack.Items))
		copy(best.Quantities, old)

		best.Quantities[position] = bestQ
		best.Weight += bestQ * s.Knapsack.Items[position].Weight
	}

	s.cache[key_] = best
	return best
}
