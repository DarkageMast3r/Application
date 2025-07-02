package models

// als ik het goed begrijp dan vraag je budget aan om iets te aan te schaffen. Dus nadat datgene aangeschaft is dan moet je opnieuw een budget aanvragen voor iets nieuws? of zit ik nu fout?
type Budget struct {
	ID                int
	MaxBedrag         float64
	BeschikbaarBedrag float64
	GebruiktBedrag    float64
	BudgetStatus      string
}

func (b *Budget) NieuwBudget(max float64) {
	b.MaxBedrag = max
	b.BeschikbaarBedrag = max
	b.GebruiktBedrag = 0
	b.BudgetStatus = "Aangevraagd"
}

// is deze überhaupt nodig? fucking business logic
func (b *Budget) UpdateBudget(bedrag float64) {
	b.GebruiktBedrag += bedrag
	b.BeschikbaarBedrag -= bedrag
}

func (b *Budget) BudgetGoedgekeurd() {
	b.BudgetStatus = "Goedgekeurd"
}

func (b *Budget) BudgetAfgewezen() {
	b.BudgetStatus = "Afgewezen"
}
