package models

type Factuur struct {
	FactuurID int
	ProductID int
	Bedrag    float64
	Betaald   bool
}

func (f *Factuur) FactuurBetaald() {
	// open recieving the payment, mark it as paid(maybe some checks? idk)
	f.Betaald = true
}
