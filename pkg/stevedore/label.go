package stevedore

import (
	"math"
	"sort"
)

// Label represents custom dimension or user pereference for overriding values
type Label struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
}

// Labels represents collection for label
type Labels []Label

// Len returns the no of items present in the underlying labels
// implements sort interface
func (labels Labels) Len() int {
	return len(labels)
}

// Less returns whether the label's weight is labels[i].Weight < labels[j].Weight
// implements sort interface
func (labels Labels) Less(i, j int) bool {
	return labels[i].Weight < labels[j].Weight
}

// Swap swaps the label[i] into label[j] and vice versa
// implements sort interface
func (labels Labels) Swap(i, j int) {
	labels[i], labels[j] = labels[j], labels[i]
}

// Weights returns the underlying weights from labels
func (labels Labels) Weights() Weights {
	values := map[string]int{}
	sort.Sort(labels)
	for index, label := range labels {
		values[label.Name] = int(math.Pow(2, float64(index)))
	}
	return Weights{values: values}
}
