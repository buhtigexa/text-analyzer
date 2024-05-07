package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"os"
	"semanticAnalysis/internal/vectorizers"
	"semanticAnalysis/models"
)

func prettyPrintMatrix(m mat.Matrix) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%3.2f ", m.At(i, j))
		}
		fmt.Println()
	}
}

func main() {
	c := vectorizers.NewCountVectorizer()
	m := vectorizers.NewTfIdfVectorizer()
	di := vectorizers.NewDataIngestor()
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

	tfmatrix, _ := models.ToMatrix(di.Index, c)
	fmt.Println(" Matrix TF IDF .......")
	prettyPrintMatrix(tfmatrix)

	tfIdfMatrix, _ := models.ToMatrix(di.Index, c)
	fmt.Println(" Matrix TF IDF .......")
	prettyPrintMatrix(tfIdfMatrix)

	var svd mat.SVD
	if ok := svd.Factorize(tfmatrix, mat.SVDFull); !ok {
		log.Fatalf("SVD factorization failed")
	}

	var U *mat.Dense = &mat.Dense{}
	svd.UTo(U)
	//svd.VTo(v)

	S := svd.Values(nil)
	prettyPrintMatrix(U)

	Sigma := mat.NewDiagDense(len(S), S)

	// Compute U * Sigma
	USigma := mat.NewDense(5, 5, nil)
	USigma.Product(U, Sigma)

	// Print the resulting matrix
	fa := mat.Formatted(USigma, mat.Prefix(" "), mat.Squeeze())
	fmt.Printf("U * Sigma = %v\n", fa)

	fmt.Println("************************* U  ********************************")

	rows, _ := U.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			u1 := U.RowView(i)
			u2 := U.RowView(j)
			cos, _ := models.CosineSimilarity(u1, u2)
			fmt.Printf(" Cosine  (%d,%d)=%4.3f: \n", i, j, cos)
		}
	}

	fmt.Println("************************* U SIGMA ********************************")

	rows, _ = USigma.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			u1 := USigma.RowView(i)
			u2 := USigma.RowView(j)
			cos, _ := models.CosineSimilarity(u1, u2)
			fmt.Printf(" Cosine US  (%d,%d)=%4.3f: \n", i, j, cos)
		}
	}

	fmt.Println("**************  TF IDF MATRIX *******************************************")

	rows, _ = tfIdfMatrix.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			u1, _ := models.CreateVector(i, di.Index, m)
			u2, _ := models.CreateVector(j, di.Index, m)
			cos, _ := models.CosineSimilarity(u1, u2)
			fmt.Printf(" Cosine US  (%d,%d)=%4.3f: \n", i, j, cos)
		}
	}

}
