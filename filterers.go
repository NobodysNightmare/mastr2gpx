package main

type Filterer interface {
	// Decides whether a given generator should be filtered or not.
	// Returns true if the generator should be part of the result,
	// false if it should be rejected from the result.
	Filter(EinheitSolar) bool
}

type PostalCodeFilter struct {
	PostalCode string
}

func (f PostalCodeFilter) Filter(e EinheitSolar) bool {
	return e.PostalCode == f.PostalCode
}

type BoundingBoxFilter struct {
	left   float64
	right  float64
	bottom float64
	top    float64
}

func (f BoundingBoxFilter) Filter(e EinheitSolar) bool {
	return e.Lng >= f.left && e.Lng <= f.right && e.Lat >= f.bottom && e.Lat <= f.top
}
