package sortedmap

import (
	"sort"
)

type Pair struct {
	Key   string
	Value int
}

type DescendingPairList []Pair

func (p DescendingPairList) Len() int           { return len(p) }
func (p DescendingPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DescendingPairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func Sort(m map[string]int) DescendingPairList {
	p := make(DescendingPairList, len(m))

	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	return p
}
