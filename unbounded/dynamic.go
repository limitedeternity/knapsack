package unbounded

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

func (s Dynamic) Solve() Solution {
	solutions := make([]Solution, s.knapsack.Limit+1)

	for i, item := range s.knapsack.Items {
		for weight := s.knapsack.Limit; weight >= 0; weight-- {
			if weight >= item.Weight {
				withoutSolution := solutions[weight-item.Weight]
				potentialValue := withoutSolution.Value + item.Value

				if potentialValue > solutions[weight].Value {
					solutions[weight] = Solution{Weight: withoutSolution.Weight + item.Weight, Value: potentialValue}
					solutions[weight].Quantities = make([]int, len(s.knapsack.Items))

					copy(solutions[weight].Quantities, withoutSolution.Quantities)
					solutions[weight].Quantities[i]++
				}
			}
		}
	}

	return solutions[s.knapsack.Limit]
}
