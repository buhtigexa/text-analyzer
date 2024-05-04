package main

import (
	"fmt"
	"log"
	"os"
	"semanticAnalysis/models"
)

func main() {
	c := models.NewCountVectorizer()
	m := models.NewTfIdf()
	di := models.NewDataIngestor()

	f, err := os.Open("text")
	if err != nil {
		log.Fatalf("%s", err)
	}

	done := m.Fit(c.Fit(di.Ingest(f)))
	<-done
	wordIndex := make([]string, len(di.Index))

	for k, v := range di.Index {
		wordIndex[v] = k
	}

	//for i := 1; i < len(c.Tf); i++ {
	//	for j := i + 1; j <= len(c.Tf); j++ {
	//		a := c.Vector(i, di.Index)
	//		b := c.Vector(j, di.Index)
	//		fmt.Println("wordIndex:", wordIndex)
	//		fmt.Printf("Coseno (%d,%d) = %.3f\n", i, j, models.CosineSimilarity(a, b))
	//	}
	//}

	fmt.Println(" TF IDF >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	for i := 1; i < len(m.Tfidf); i++ {
		for j := i + 1; j <= len(m.Tfidf); j++ {
			//fmt.Println(" TF IDF >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			//fmt.Println("wordIndex:", wordIndex)
			a := m.Vector(i, di.Index)
			b := m.Vector(j, di.Index)
			fmt.Printf("Coseno (%d,%d) = %.3f\n", i, j, models.CosineSimilarity(a, b))
		}
	}

}
