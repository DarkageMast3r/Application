package models

type Factuur struct {
	factuurID int
	productID int
	bedrag    float64
	betaald   bool
}

func (f Factuur) factuurBetaald() {
	// open recieving the paying, mark it as paid(maybe some checks? idk)
	f.betaald = true
}
