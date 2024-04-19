package common

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Solution struct {
	Weight     int
	Value      int
	Quantities []int
}

func (s Solution) String(params struct{ ItemNames []string }) string {
	var builder strings.Builder

	builder.WriteString("Taking:\n")
	for i, q := range s.Quantities {
		if q > 0 {
			builder.WriteString(fmt.Sprintf("+ %s: %d\n", params.ItemNames[i], q))
		}
	}

	builder.WriteString(fmt.Sprintf("Total value: %d\n", s.Value))
	builder.WriteString(fmt.Sprintf("Total weight: %d\n", s.Weight))

	return builder.String()
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
	case nil:
		s.Knapsack = nil
	default:
		log.Fatalf("Unsupported knapsack type: %T, expected: %T", knapsack, s.Knapsack)
	}
}

type Knapsack[I any, S Solver] struct {
	Limit int
	Items []I
}

func NewKnapsack[I any, S Solver](limit ...int) *Knapsack[I, S] {
	if len(limit) == 0 {
		return &Knapsack[I, S]{}
	} else {
		return &Knapsack[I, S]{Limit: limit[0]}
	}
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

	if reflect.TypeFor[S]().Kind() != reflect.Ptr {
		log.Fatal("Solver type parameter must be a pointer")
	}

	solver := reflect.New(reflect.TypeFor[S]().Elem()).Interface().(Solver)
	solver.GetBase().InjectKnapsack(k)

	sol := solver.Solve()
	solver.GetBase().InjectKnapsack(nil)

	solver = nil
	return sol
}
