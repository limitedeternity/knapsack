package bounded

import "log"

type Dynamic struct {
	knapsack *Knapsack[Dynamic]
}

func (s Dynamic) WithKnapsack(knapsack any) Solver {
	switch v := knapsack.(type) {
	case *Knapsack[Dynamic]:
		s.knapsack = v
	default:
		log.Fatalf("Unsupported knapsack type: %T", v)
	}

	return s
}

func (s Dynamic) Solve(limit int) Solution {
	solutions := make([]Solution, limit+1)
	for i := 0; i < limit+1; i++ {
		solutions[i] = Solution{}
	}

	for i, item := range s.knapsack.Items {
		for weight := limit; weight >= 0; weight-- {
			if weight >= item.Weight {
				quantity := min(item.Pieces, weight/item.Weight)

				for q := 1; q <= quantity; q++ {
					reduced := solutions[weight-q*item.Weight]
					potentialValue := reduced.Value + q*item.Value

					if potentialValue > solutions[weight].Value {
						solutions[weight] = Solution{Weight: reduced.Weight + q*item.Weight, Value: potentialValue}
						solutions[weight].Quantities = make([]int, len(s.knapsack.Items))

						copy(solutions[weight].Quantities, reduced.Quantities)
						solutions[weight].Quantities[i] += q
					}
				}
			}
		}
	}

	return solutions[limit]
}
