package models

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"os"
)

type TfIdfModel struct {
	df    map[string]map[int]bool
	tf    map[int]map[string]float64
	idf   map[string]float64
	tfidf map[int]map[string]float64
}

func NewTfIdf() *TfIdfModel {
	return &TfIdfModel{
		df:    make(map[string]map[int]bool),
		tf:    make(map[int]map[string]float64),
		idf:   make(map[string]float64),
		tfidf: make(map[int]map[string]float64),
	}
}

func Build(fname string) {
	c := CountVectorizer{}
	m := NewTfIdf()
	i := DataIngestor{}

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("%s", err)
	}

	done := m.Fit(c.Fit(i.Ingest(f)))
	<-done
}

type DocTermData struct {
	term  string
	docId int
	count int
	freq  float64
}

func (m *TfIdfModel) Fit(inCh <-chan DocTermData) <-chan struct{} {
	done := make(chan struct{})
	maxDocs := 0
	go func() {
		defer func() {
			for docId, wfs := range m.tf {
				for w, _ := range wfs {
					if _, exists := m.tfidf[docId]; !exists {
						m.tfidf[docId] = map[string]float64{}
					}
					fmt.Printf("TF*IDF >>>: m.tfidf[%d][%s] = m.tf[%d][%s] * m.idf[%s] =%.3f\n", docId, w, docId, w, w, m.tf[docId][w]*m.idf[w])
					m.tfidf[docId][w] = m.tf[docId][w] * m.idf[w]
				}
			}
			done <- struct{}{}
		}()
		for data := range inCh {
			if maxDocs < data.docId {
				maxDocs = data.docId
			}

			if _, exists := m.df[data.term]; !exists {
				m.df[data.term] = make(map[int]bool)
			}
			m.df[data.term][data.docId] = true
			fmt.Printf(">>>  m.df[%s][%d]=%t  and len(%s)=%d\n", data.term, data.docId, m.df[data.term][data.docId], data.term, len(m.df[data.term]))
			if _, exists := m.tf[data.docId]; !exists {
				m.tf[data.docId] = make(map[string]float64)
			}
			m.tf[data.docId][data.term] = data.freq

			for w, docs := range m.df {
				dtf := len(docs)
				m.idf[w] = math.Log10(float64(maxDocs) / float64(dtf))
				fmt.Printf("m.idf[%s] = math.Log(float64(%d / %d)) =%.3f \n", w, maxDocs, dtf, m.idf[w])
			}
		}
	}()
	return done
}

func (c *TfIdfModel) Vector(id int, index map[string]int) *mat.Dense {
	data := make([]float64, len(index))
	for w, tf := range c.tfidf[id] {
		i := index[w]
		data[i] = tf
	}
	matrix := mat.NewDense(1, len(c.tf[id]), data)
	return matrix
}
