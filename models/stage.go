package models

type Stage interface {
	Fit(inCh <-chan interface{}) <-chan interface{}
	GetData() map[int]map[string]float64
	GetDataI(int) map[string]float64
}
