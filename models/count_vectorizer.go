package models

import "gonum.org/v1/gonum/mat"

type CountVectorizer struct {
	tf map[int]map[string]float64
}

func NewCountVectorizer() *CountVectorizer {
	return &CountVectorizer{
		tf: make(map[int]map[string]float64),
	}
}

func (c *CountVectorizer) Fit(inCh <-chan DocTermData) <-chan DocTermData {
	out := make(chan DocTermData)
	go func() {
		defer close(out)
		for in := range inCh {
			if _, ok := c.tf[in.docId]; !ok {
				c.tf[in.docId] = map[string]float64{}
			}
			c.tf[in.docId][in.term]++
			out <- DocTermData{term: in.term, docId: in.docId, count: in.count, freq: c.tf[in.docId][in.term] / float64(in.count)}
		}
	}()
	return out
}

func (c *CountVectorizer) Vector(id int, index map[string]int) *mat.Dense {
	data := make([]float64, len(index))
	for w, tf := range c.tf[id] {
		i := index[w]
		data[i] = tf
	}
	matrix := mat.NewDense(1, len(c.tf[id]), data)
	return matrix
}
