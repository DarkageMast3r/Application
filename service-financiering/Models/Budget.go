package models

import (
	r "Financiering/Repositories"
	"log"
)

type Budget struct {
	ID                int
	MaxBedrag         float64
	BeschikbaarBedrag float64
	GebruiktBedrag    float64
	BudgetStatus      string
}

func (b *Budget) NewBudget() {
	res, err := r.InsertBudget(b.MaxBedrag, b.BeschikbaarBedrag, b.GebruiktBedrag, b.BudgetStatus)
	var lastid int
	if err == nil {
		val, err := res.LastInsertId()
		if err != nil {
			log.Println("lastinsertid: ", err)
		} else {
			lastid = int(val)
			b.ID = lastid
		}
	}
}