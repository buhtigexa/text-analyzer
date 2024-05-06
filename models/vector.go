package models

import (
	"errors"
	"gonum.org/v1/gonum/mat"
)

func CosineSimilarity(a, b mat.Vector) (float64, error) {
	if a.Len() != b.Len() {
		return 0, errors.New("vectors have different lengths")
	}
	dot := mat.Dot(a, b)
	na := mat.Norm(a, 2)
	nb := mat.Norm(b, 2)
	return dot / (na * nb), nil
}

func ToMatrix(index map[string]int, c Stage) (mat.Matrix, error) {
	numDocs := len(c.GetData())
	numTerms := len(index)
	data := make([]float64, numDocs*numTerms)
	j := 0

	for docID := 1; docID <= numDocs; docID++ {
		v, err := CreateVector(docID, index, c)
		if err != nil {
			return nil, err
		}
		for i := 0; i < v.Len(); i++ {
			data[j] = v.AtVec(i)
			j++
		}
	}

	return mat.NewDense(numDocs, numTerms, data), nil
}

func CreateVector(id int, index map[string]int, c Stage) (mat.Vector, error) {
	data := make([]float64, len(index))
	for w, tf := range c.GetDataI(id) {
		if idx, ok := index[w]; ok {
			data[idx] = tf
		} else {
			return nil, errors.New("index for term not found")
		}
	}
	return mat.NewVecDense(len(index), data), nil
}
