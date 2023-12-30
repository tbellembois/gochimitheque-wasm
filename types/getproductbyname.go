package types

type PubchemProduct struct {
	Name                *string   `json:"name"`
	Inchi               *string   `json:"inchi"`
	InchiKey            *string   `json:"inchi_key"`
	CanonicalSmiles     *string   `json:"canonical_smiles"`
	MolecularFormula    *string   `json:"molecular_formula"`
	Cas                 *string   `json:"cas"`
	Ec                  *string   `json:"ec"`
	MolecularWeight     *string   `json:"molecular_weight"`
	MolecularWeightUnit *string   `json:"molecular_weight_unit"`
	Synonyms            *[]string `json:"synonyms"`
	Symbols             *[]string `json:"symbols"`
	Signal              *[]string `json:"signal"`
	Hs                  *[]string `json:"hs"`
	Ps                  *[]string `json:"ps"`
	Twodpicture         *string   `json:"twodpicture"` // base64 encoded png
}
