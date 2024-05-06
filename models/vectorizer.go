package models

import "gonum.org/v1/gonum/mat"

type Vectorizer interface {
	ToMatrix(index map[string]int, c Stage) mat.Matrix
}
