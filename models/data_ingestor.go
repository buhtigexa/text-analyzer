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
		make(map[string]int),
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
		for scanner.Scan() {
			words := strings.Fields(scanner.Text())
			for _, term := range words {
				out <- DocTermData{term: strings.Trim(term, " "), docId: docId, count: len(words)}
				m.Index[term]++
			}
			docId++
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("reading from file error: %s", err)
		}
	}()
	return out
}
