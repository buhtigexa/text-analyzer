package vectorizers

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

func (di *DataIngestor) Ingest(r io.ReadCloser) <-chan interface{} {
	out := make(chan interface{})
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
				if _, ok := di.Index[term]; !ok {
					di.Index[term] = wordIndex
					wordIndex++
				}
			}
			docId++
		}
		if err := scanner.Err(); err != nil {
			log.Printf("reading from file error: %s", err)
		}
	}()
	return out
}
