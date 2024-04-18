package common

import (
	"log"
	"reflect"
)

type Solution struct {
	Weight     int
	Value      int
	Quantities []int
}

type Solver interface {
	GetBase() ISolverBase
	Solve() Solution
}

type ISolverBase interface {
	InjectKnapsack(knapsack any)
}

type SolverBase[I any, S Solver] struct {
	Knapsack *Knapsack[I, S]
}

func (s *SolverBase[I, S]) InjectKnapsack(knapsack any) {
	switch v := knapsack.(type) {
	case *Knapsack[I, S]:
		s.Knapsack = v
	default:
		log.Fatalf("Unsupported knapsack type: %T, expected: %T", knapsack, s.Knapsack)
	}
}

type Knapsack[I any, S Solver] struct {
	Limit int
	Items []I
}

func NewKnapsack[I any, S Solver]() *Knapsack[I, S] {
	return &Knapsack[I, S]{}
}

func (k *Knapsack[I, S]) WithCapacity(limit int) *Knapsack[I, S] {
	k.Limit = limit
	return k
}

func (k *Knapsack[I, S]) Pack(items []I) Solution {
	if k.Limit <= 0 {
		log.Fatal("Knapsack capacity must be greater than zero")
	}

	k.Items = items

	solver := reflect.New(reflect.TypeFor[S]().Elem()).Interface().(Solver)
	solver.GetBase().InjectKnapsack(k)
	sol := solver.Solve()

	solver = nil
	return sol
}