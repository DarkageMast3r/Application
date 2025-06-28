package models

// value object
// should the money be in floats?
type Budget struct {
	maxBedrag         float64
	beschikbaarBedrag float64
	gebruiktBedrag    float64
	BudgetStatus      BudgetStatus
}

func (b Budget) Constructor(max float64) Budget {
	b.maxBedrag = max
	b.beschikbaarBedrag = max
	b.gebruiktBedrag = 0
	b.BudgetStatus = b.BudgetStatus.GetStatus("Aangevraagd")
	return b
}

func (b Budget) UpdateBudget(bedrag float64) {
	b.gebruiktBedrag += bedrag
	b.beschikbaarBedrag -= bedrag
}
