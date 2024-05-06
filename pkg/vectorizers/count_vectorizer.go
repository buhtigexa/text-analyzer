package vectorizers

import (
	"log"
	"semanticAnalysis/models"
)

type CountVectorizer struct {
	data map[int]map[string]float64
	models.Vectorizer
}

func NewCountVectorizer() *CountVectorizer {
	return &CountVectorizer{
		data: make(map[int]map[string]float64),
	}
}

func (c *CountVectorizer) GetData() map[int]map[string]float64 {
	return c.data
}

func (c *CountVectorizer) GetDataI(i int) map[string]float64 {
	return c.data[i]
}

func (c *CountVectorizer) Fit(inCh <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for d := range inCh {
			if in, ok := d.(DocTermData); ok {
				if _, ok := c.data[in.docId]; !ok {
					c.data[in.docId] = map[string]float64{}
				}
				c.data[in.docId][in.term]++
				out <- DocTermData{term: in.term, docId: in.docId, count: in.count, freq: c.data[in.docId][in.term] / float64(in.count)}
			} else {
				log.Printf("received invalid data: %v", d)
			}
		}
	}()
	return out
}
