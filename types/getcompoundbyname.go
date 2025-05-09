//go:build go1.24 && js && wasm

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

type Section struct {
	TOCHeading  *string        `json:"TOCHeading"`
	TOCID       *int           `json:"TOCID"`
	Description string         `json:"Description"`
	URL         string         `json:"URL"`
	Section     *[]Section     `json:"Section"`
	Information *[]Information `json:"Information"`
}

type Information struct {
	ReferenceNumber int      `json:"ReferenceNumber"`
	Name            string   `json:"Name"`
	Description     string   `json:"Description"`
	Reference       []string `json:"Reference"`
	LicenseNote     []string `json:"LicenseNote"`
	LicenseURL      []string `json:"LicenseURL"`
	Value           Value    `json:"Value"`
}

type Value struct {
	Number               []float64           `json:"Number"`
	DateISO8601          []string            `json:"DateISO8601"`
	Boolean              []bool              `json:"Boolean"`
	Binary               []string            `json:"Binary"`
	BinaryToStore        []string            `json:"BinaryToStore"`
	ExternalDataURL      []string            `json:"ExternalDataURL"`
	ExternalTableName    string              `json:"ExternalTableName"`
	Unit                 string              `json:"Unit"`
	MimeType             string              `json:"MimeType"`
	ExternalTableNumRows int                 `json:"ExternalTableNumRows"`
	StringWithMarkup     *[]StringWithMarkup `json:"StringWithMarkup"`
}

type Markup struct {
	Start  float64 `json:"Start"`
	Length float64 `json:"Length"`
	URL    string  `json:"URL"`
	Type   string  `json:"Type"`
	Extra  string  `json:"Extra"`
}

type StringWithMarkup struct {
	String string   `json:"String"`
	Markup []Markup `json:"Markup"`
}

type Record struct {
	Record RecordContent `json:"Record"`
}

type RecordContent struct {
	RecordType        string        `json:"RecordType"`
	RecordNumber      int           `json:"RecordNumber"`
	RecordAccession   string        `json:"RecordAccession"`
	RecordTitle       string        `json:"RecordTitle"`
	RecordExternalURL string        `json:"RecordExternalURL"`
	Section           []Section     `json:"Section"`
	Information       []Information `json:"Information"`
}

type PCCompound struct {
	ID     ID     `json:"id"`
	Props  []Prop `json:"props"`
	Record Record `json:"record"`
}

type Compounds struct {
	PCCompounds []PCCompound `json:"PC_Compounds"`
	Record      Record       `json:"record"`
	Base64Png   string       `json:"base64_png"`
}
