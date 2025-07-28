package models

import (
	r "Financiering/Repositories"
)

type Budget struct {
	ID                int
	MaxBedrag         float64
	BeschikbaarBedrag float64
	GebruiktBedrag    float64
	BudgetStatus      string
}

func (b *Budget) NieuwBudget(max float64) error {
	b.MaxBedrag = max
	b.BeschikbaarBedrag = max
	b.GebruiktBedrag = 0
	b.BudgetStatus = "Aangevraagd"
	return r.NewBudget(b.MaxBedrag, b.BeschikbaarBedrag, b.GebruiktBedrag, b.BudgetStatus)
}

func (b *Budget) UpdateBudget(bedrag float64) error {
	b.GebruiktBedrag += bedrag
	b.BeschikbaarBedrag -= bedrag
	b.BudgetStatus = "Gefactureerd"
	return r.ProcessPayment(b.GebruiktBedrag, b.BeschikbaarBedrag, b.BudgetStatus, b.ID)
}
