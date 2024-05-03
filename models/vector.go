package models

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func cosine(a, b mat.Vector) float64 {
	if a.Len() != b.Len() {
		panic("Vectors have different lengths")
	}

	dot := 0.0
	normA, normB := 0.0, 0.0
	for i := 0; i < a.Len(); i++ {
		ai, bi := a.AtVec(i), b.AtVec(i)
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}
	return dot / math.Sqrt(normA) * math.Sqrt(normB)
}
