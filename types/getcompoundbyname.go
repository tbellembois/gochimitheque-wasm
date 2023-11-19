package types

type PropValue struct {
	Ival   *int     `json:"ival"`
	Fval   *float64 `json:"fval"`
	Binary *string  `json:"binary"`
	Sval   *string  `json:"sval"`
}

type PropURN struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type Prop struct {
	URN   PropURN   `json:"urn"`
	Value PropValue `json:"value"`
}

type ID struct {
	ID CID `json:"id"`
}

type CID struct {
	CID int `json:"cid"`
}

type PCCompound struct {
	ID    ID     `json:"id"`
	Props []Prop `json:"props"`
}

type Compounds struct {
	PCCompounds []PCCompound `json:"PC_Compounds"`
}
