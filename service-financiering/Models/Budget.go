package models

// value object
// should the money be in floats?
type Budget struct {
	maxBedrag         float64
	beschikbaarBedrag float64
	gebruiktBedrag    float64
	BudgetStatus      BudgetStatus
}

func (b *Budget) BugdetConstructor(max float64) {
	b.maxBedrag = max
	b.beschikbaarBedrag = max
	b.gebruiktBedrag = 0
	b.BudgetStatus = b.BudgetStatus.GetStatus("Aangevraagd")
}

// is deze überhaupt nodig? fucking business logic
func (b *Budget) UpdateBudget(bedrag float64) {
	b.gebruiktBedrag += bedrag
	b.beschikbaarBedrag -= bedrag
}

func (b *Budget) BudgetAfgewezen() {
	b.BudgetStatus = b.BudgetStatus.GetStatus("Afgewezen")
}
