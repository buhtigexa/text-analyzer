package models

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type DataIngestor struct {
	Index map[string]int
}

func NewDataIngestor() *DataIngestor {
	return &DataIngestor{
		Index: make(map[string]int),
	}
}

func (m *DataIngestor) Ingest(r io.ReadCloser) <-chan DocTermData {
	out := make(chan DocTermData)
	scanner := bufio.NewScanner(r)
	go func() {
		defer func() {
			r.Close()
			close(out)
		}()
		docId := 1
		wordIndex := 0
		for scanner.Scan() {
			words := strings.Fields(scanner.Text())
			for _, term := range words {
				out <- DocTermData{term: strings.Trim(term, " "), docId: docId, count: len(words)}
				if _, ok := m.Index[term]; !ok {
					m.Index[term] = wordIndex
					wordIndex++
				}
			}
			docId++
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("reading from file error: %s", err)
		}
	}()
	return out
}
