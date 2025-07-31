package repositories

import (
	m "Financiering/Models"
)

func NieuwBudget(b *m.Budget, max float64) error {
	b.MaxBedrag = max
	b.BeschikbaarBedrag = max
	b.GebruiktBedrag = 0
	b.BudgetStatus = "Aangevraagd"
	// return NewBudget(b.MaxBedrag, b.BeschikbaarBedrag, b.GebruiktBedrag, b.BudgetStatus)
	return nil //unfinished method
}

func UpdateBudget(b *m.Budget, bedrag float64) error {
	b.GebruiktBedrag += bedrag
	b.BeschikbaarBedrag -= bedrag
	b.BudgetStatus = "Gefactureerd"
	return ProcessPayment(b.GebruiktBedrag, b.BeschikbaarBedrag, b.BudgetStatus, b.ID)
}
