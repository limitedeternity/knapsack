package unbounded

import (
	. "knapsack/common"
)

type DPSolver struct {
	SolverBase[Item, *DPSolver]
}

func (s *DPSolver) GetBase() ISolverBase {
	return &s.SolverBase
}

func (s *DPSolver) Solve() Solution {
	solutions := make([]Solution, s.Knapsack.Limit+1)

	for i, item := range s.Knapsack.Items {
		for weight := s.Knapsack.Limit; weight >= 0; weight-- {
			if weight >= item.Weight {
				withoutSolution := solutions[weight-item.Weight]
				potentialValue := withoutSolution.Value + item.Value

				if potentialValue > solutions[weight].Value {
					solutions[weight] = Solution{Weight: withoutSolution.Weight + item.Weight, Value: potentialValue}
					solutions[weight].Quantities = make([]int, len(s.Knapsack.Items))

					copy(solutions[weight].Quantities, withoutSolution.Quantities)
					solutions[weight].Quantities[i] += 1
				}
			}
		}
	}

	return solutions[s.Knapsack.Limit]
}
