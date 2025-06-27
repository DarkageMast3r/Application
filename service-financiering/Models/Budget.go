package models

// value object
// should the money be in floats?
type Budget struct {
	maxBedrag         float32
	beschikbaarBedrag float32
	gebruiktBedrag    float32
}

func (b Budget) constructor(max float32, beschikbaar float32, gebruikt float32) {
	b.maxBedrag = max
	b.beschikbaarBedrag = beschikbaar
	b.gebruiktBedrag = gebruikt
}

func (b Budget) updateBudget(bedrag float32) {
	b.gebruiktBedrag += bedrag
	b.beschikbaarBedrag -= bedrag
}
