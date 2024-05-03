package tests

import (
	"semanticAnalysis/models"
	"testing"
)

var corpus = [2][]byte{
	[]byte("the quick brown fox jumped over the lazy dog"),
	[]byte("hey diddle diddle, the cat and the fiddle"),
}

func TestAddModel(t *testing.T) {
	tf := models.NewTFModel()
	for _, v := range corpus {
		tf.Ingest(v)
	}

}
