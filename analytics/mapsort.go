package analytics

import (
	"sort"
)

type (
	Pair struct {
		Key   string
		Value int
	}

	PairList []Pair
)

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//rankByValue adapted from post by Andrew Gerrand - https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
func rankByValue(mip map[string]int) PairList {
	pl := make(PairList, len(mip))
	i := 0
	for k, v := range mip {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
