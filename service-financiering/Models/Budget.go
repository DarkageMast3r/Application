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
	if err == nil {
		val, err := res.LastInsertId()
		if err != nil {
			log.Println("lastinsertid: ", err)
		} else {
			b.ID = int(val)
		}
	}
}

func (b *Budget) VerwerkBetaling(Bedrag int) {
	b.BeschikbaarBedrag -= float64(Bedrag)
	b.GebruiktBedrag += float64(Bedrag)
	err := r.ProcessPayment(b.GebruiktBedrag, b.BeschikbaarBedrag, b.ID)
	if err != nil {
		log.Println("VerwerkFactuur: ", err)
	}
}

func GetClientBudget(clientID int) Budget {
	var Budget Budget
	Result, err := r.GetBudgetbyClientID(clientID)
	if err != nil {
		log.Println("GetClientBudget: ", err)
		return Budget
	}
	nextable := Result.Next()
	defer Result.Close()
	if nextable == true {
		Result.Scan(
		&Budget.ID,
		&Budget.MaxBedrag,
		&Budget.BeschikbaarBedrag,
		&Budget.GebruiktBedrag,
		&Budget.BudgetStatus,
	)
	}
	return Budget
}