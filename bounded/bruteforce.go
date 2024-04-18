package bounded

import "log"

type BruteForce struct {
	knapsack *Knapsack[BruteForce]
	cache    map[key]Solution
}

type key struct {
	weight   int
	position int
}

func (s BruteForce) WithKnapsack(knapsack any) Solver {
	switch v := knapsack.(type) {
	case *Knapsack[BruteForce]:
		s.knapsack = v
	default:
		log.Fatalf("Unsupported knapsack type: %T", v)
	}

	return s
}

func (s BruteForce) Solve() Solution {
	s.cache = make(map[key]Solution)
	sol := s.impl(s.knapsack.Limit, len(s.knapsack.Items)-1)

	s.cache = nil
	return sol
}

func (s BruteForce) impl(limit int, position int) Solution {
	if position < 0 || limit <= 0 {
		return Solution{}
	}

	key_ := key{limit, position}
	if sol, ok := s.cache[key_]; ok {
		return sol
	}

	bestI, best := 0, Solution{}
	for i := 0; i*s.knapsack.Items[position].Weight <= limit && i <= s.knapsack.Items[position].Pieces; i++ {
		sol := s.impl(limit-i*s.knapsack.Items[position].Weight, position-1)
		sol.Value += i * s.knapsack.Items[position].Value

		if sol.Value > best.Value {
			bestI, best = i, sol
		}
	}

	if bestI > 0 {
		old := best.Quantities
		best.Quantities = make([]int, len(s.knapsack.Items))
		copy(best.Quantities, old)

		best.Quantities[position] = bestI
		best.Weight += bestI * s.knapsack.Items[position].Weight
	}

	s.cache[key_] = best
	return best
}
