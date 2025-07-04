package models

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

func (b *Budget) UpdateBudget(bedrag float64) {
	b.GebruiktBedrag += bedrag
	b.BeschikbaarBedrag -= bedrag
}
