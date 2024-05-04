package models

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

type CountVectorizer struct {
	Tf map[int]map[string]float64
}

func NewCountVectorizer() *CountVectorizer {
	return &CountVectorizer{
		Tf: make(map[int]map[string]float64),
	}
}

func (c *CountVectorizer) Fit(inCh <-chan DocTermData) <-chan DocTermData {
	out := make(chan DocTermData)
	go func() {
		defer close(out)
		for in := range inCh {
			if _, ok := c.Tf[in.docId]; !ok {
				c.Tf[in.docId] = map[string]float64{}
			}
			c.Tf[in.docId][in.term]++
			out <- DocTermData{term: in.term, docId: in.docId, count: in.count, freq: c.Tf[in.docId][in.term] / float64(in.count)}
		}
	}()
	return out
}

func (c *CountVectorizer) Vector(id int, index map[string]int) mat.Vector {
	data := make([]float64, len(index))
	for w, tf := range c.Tf[id] {
		i := index[w]
		data[i] = tf
	}
	matrix := mat.NewVecDense(len(index), data)
	fmt.Printf("%v\n", matrix)
	return matrix
}
