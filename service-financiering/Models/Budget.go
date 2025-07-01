package models

// als ik het goed begrijp dan vraag je budget aan om iets te aan te schaffen. Dus nadat datgene aangeschaft is dan moet je opnieuw een budget aanvragen voor iets nieuws? of zit ik nu fout?
type Budget struct {
	maxBedrag         float64
	beschikbaarBedrag float64
	gebruiktBedrag    float64
	BudgetStatus      BudgetStatus
}

func (b *Budget) NieuwBudget(max float64) {
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

func (b *Budget) BudgetGoedgekeurd() {
	b.BudgetStatus = b.BudgetStatus.GetStatus("Goedgekeurd")
}

func (b *Budget) BudgetAfgewezen() {
	b.BudgetStatus = b.BudgetStatus.GetStatus("Afgewezen")
}
