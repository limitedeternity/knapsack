package bounded

import (
	"github.com/creasty/defaults"
)

type Item struct {
	Item   string `yaml:"item"`
	Weight int    `yaml:"weight"`
	Value  int    `yaml:"value"`
	Pieces int    `default:"1" yaml:"pieces"`
}

func (i *Item) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(i); err != nil {
		return err
	}

	type plain Item
	if err := unmarshal((*plain)(i)); err != nil {
		return err
	}

	return nil
}
