package vectorizers

import (
	"fmt"
	"math"
	"semanticAnalysis/models"
)

type TfIdfVectorizer struct {
	df   map[string]map[int]bool
	tf   map[int]map[string]float64
	idf  map[string]float64
	data map[int]map[string]float64
	models.Vectorizer
}

func NewTfIdfVectorizer() *TfIdfVectorizer {
	return &TfIdfVectorizer{
		df:   make(map[string]map[int]bool),
		tf:   make(map[int]map[string]float64),
		idf:  make(map[string]float64),
		data: make(map[int]map[string]float64),
	}
}

type DocTermData struct {
	term  string
	docId int
	count int
	freq  float64
}

func (c *TfIdfVectorizer) GetData() map[int]map[string]float64 {
	return c.data
}

func (c *TfIdfVectorizer) GetDataI(i int) map[string]float64 {
	return c.data[i]
}
func (c *TfIdfVectorizer) updateIDFs(maxDocs int) {
	for term, docs := range c.df {
		docFrequency := len(docs)
		c.idf[term] = math.Log(float64(maxDocs+1)/float64(docFrequency+1)) + 1
	}
}

func (c *TfIdfVectorizer) Fit(inCh <-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	maxDocs := 0

	go func() {
		defer close(done)
		for in := range inCh {
			if data, ok := in.(DocTermData); ok {
				if maxDocs < data.docId {
					maxDocs = data.docId
				}
				if _, exists := c.df[data.term]; !exists {
					c.df[data.term] = make(map[int]bool)
				}
				c.df[data.term][data.docId] = true

				if _, exists := c.tf[data.docId]; !exists {
					c.tf[data.docId] = make(map[string]float64)
				}
				c.tf[data.docId][data.term] = data.freq
			} else {
				fmt.Println("Received data of unexpected type.")
			}
		}
		c.updateIDFs(maxDocs)
		for docId, terms := range c.tf {
			if _, exists := c.data[docId]; !exists {
				c.data[docId] = make(map[string]float64)
			}
			for term := range terms {
				c.data[docId][term] = c.tf[docId][term] * c.idf[term]
			}
		}
	}()

	return done
}
