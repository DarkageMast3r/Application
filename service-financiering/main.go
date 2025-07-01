package main

import (
	m "Financiering/models"
	"fmt"
)

func main() {
	fmt.Println("shut up golang")
	var FD m.FinancieringsDossier
	FD.NieuwDossier(1, 2, 3)
	FD.VraagBudgetAan(10)
	fmt.Println(FD.Budget)
	fmt.Println(FD.Budget.BudgetStatus)
}
