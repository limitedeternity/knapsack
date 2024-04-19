package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mkideal/cli"
	"gopkg.in/yaml.v3"

	"knapsack/bounded"
	. "knapsack/common"
	"knapsack/unbounded"
)

type argT struct {
	cli.Helper2
	ItemsYaml string `cli:"*items,i" usage:"Yaml file with an array of items"`
	Capacity  int    `cli:"*capacity,c" usage:"Knapsack capacity"`
	Knapsack  string `cli:"knapsack,k" usage:"Knapsack type" dft:"bounded"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	if absPath, err := filepath.Abs(argv.ItemsYaml); err == nil {
		argv.ItemsYaml = absPath
	} else {
		return err
	}

	if _, err := os.Stat(argv.ItemsYaml); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file does not exist: %s", ctx.Color().Red(argv.ItemsYaml))
	}

	if argv.Capacity <= 0 {
		return fmt.Errorf("knapsack capacity must be greater than zero: %d", argv.Capacity)
	}

	if argv.Knapsack != "bounded" && argv.Knapsack != "unbounded" {
		return fmt.Errorf("knapsack type must be one of: bounded, unbounded")
	}

	return nil
}

func indicator(cancel <-chan struct{}, ack chan<- struct{}) {
	t0 := time.Now()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("\rRunning: %.0f seconds", t.Sub(t0).Seconds())
		case <-cancel:
			close(ack)
			return
		}
	}
}

func runIndicator() (cancel chan<- struct{}, ack <-chan struct{}) {
	cancelCh := make(chan struct{})
	ackCh := make(chan struct{})

	go indicator(cancelCh, ackCh)
	return cancelCh, ackCh
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) (err error) {
		argv := ctx.Argv().(*argT)
		cancel, ack := runIndicator()

		var yamlFile []byte
		if yamlFile, err = os.ReadFile(argv.ItemsYaml); err != nil {
			return err
		}

		var sol Solution

		switch argv.Knapsack {
		case "unbounded":
			var items []unbounded.Item
			if err = yaml.Unmarshal(yamlFile, &items); err != nil {
				return err
			}

			sol = NewKnapsack[unbounded.Item, *unbounded.DPSolver](argv.Capacity).Pack(items)

		default:
			var items []bounded.Item
			if err = yaml.Unmarshal(yamlFile, &items); err != nil {
				return err
			}

			sol = NewKnapsack[bounded.Item, *bounded.DPSolver](argv.Capacity).Pack(items)
		}

		close(cancel)
		<-ack

		// TODO: Print solution
		ctx.String("\r%+v\n", sol)

		return nil
	}))
}
