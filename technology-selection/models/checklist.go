package models

type CheckableOption struct {
	Selected    bool   `schema:"selected"`
	Name        string `schema:"name"`
	Description string
}

type Checklist struct {
	Options []CheckableOption `schema:"options"`
	Prefix  string
}
