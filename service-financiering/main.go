package main

import (
	m "Financiering/models"
	"fmt"
)

func main() {
	fmt.Println("shut up golang")
	var FD m.FinancieringsDossier
	FD.VraagBudgetAan(10)
	fmt.Println(FD.BedragAangevraagd)
	return
}
