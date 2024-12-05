package product

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque/models"
)

func FillInProductForm(p Product, id string) {

	jquery.Jq(fmt.Sprintf("#%s #product_id", id)).SetVal(p.ProductID)

	select2Category := select2.NewSelect2(jquery.Jq("select#category"), nil)
	select2Category.Select2Clear()
	if p.CategoryID != nil {
		select2Category.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.Category.CategoryLabel,
				Value:           strconv.Itoa(int(*p.Category.CategoryID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Tags := select2.NewSelect2(jquery.Jq("select#tags"), nil)
	select2Tags.Select2Clear()
	for _, tag := range p.Tags {
		select2Tags.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            tag.TagLabel,
				Value:           strconv.Itoa(tag.TagID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Name := select2.NewSelect2(jquery.Jq("select#name"), nil)
	select2Name.Select2Clear()
	select2Name.Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            p.Name.NameLabel,
			Value:           strconv.Itoa(int(p.Name.NameID)),
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())

	select2Synonyms := select2.NewSelect2(jquery.Jq("select#synonyms"), nil)
	select2Synonyms.Select2Clear()
	for _, synonym := range p.Synonyms {
		select2Synonyms.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            synonym.NameLabel,
				Value:           strconv.Itoa(synonym.NameID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Producer := select2.NewSelect2(jquery.Jq("select#producer"), nil)
	select2Producer.Select2Clear()
	if p.Producer != nil && p.Producer.ProducerID != nil {
		select2Producer.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.Producer.ProducerLabel,
				Value:           strconv.Itoa(int(*p.Producer.ProducerID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2ProducerRef := select2.NewSelect2(jquery.Jq("select#producer_ref"), nil)
	select2ProducerRef.Select2Clear()
	if p.ProducerRef.ProducerRefID != nil {
		select2ProducerRef.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.ProducerRef.ProducerRefLabel,
				Value:           strconv.Itoa(int(*p.ProducerRef.ProducerRefID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2SupplierRef := select2.NewSelect2(jquery.Jq("select#supplier_refs"), nil)
	select2SupplierRef.Select2Clear()
	for _, supplierref := range p.SupplierRefs {

		supplierrefToSupplier[supplierref.SupplierRefLabel] = *supplierref.Supplier.SupplierID

		select2SupplierRef.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            fmt.Sprintf("%s@%s", supplierref.SupplierRefLabel, *supplierref.Supplier.SupplierLabel),
				Value:           strconv.Itoa(supplierref.SupplierRefID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_temperature").SetVal("")
	if p.ProductTemperature != nil {
		jquery.Jq("#product_temperature").SetVal(*p.ProductTemperature)
	}

	select2UnitTemperature := select2.NewSelect2(jquery.Jq("select#unit_temperature"), nil)
	select2UnitTemperature.Select2Clear()
	if p.UnitTemperature.UnitID != nil {
		select2UnitTemperature.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.UnitTemperature.UnitLabel,
				Value:           strconv.Itoa(int(*p.UnitTemperature.UnitID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_molecularweight").SetVal("")
	if p.ProductMolecularWeight != nil {
		jquery.Jq("#product_molecularweight").SetVal(*p.ProductMolecularWeight)
	}

	select2UnitMolecularWeight := select2.NewSelect2(jquery.Jq("select#unit_molecularweight"), nil)
	select2UnitMolecularWeight.Select2Clear()
	if p.UnitMolecularWeight.UnitID != nil {
		select2UnitMolecularWeight.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.UnitMolecularWeight.UnitLabel,
				Value:           strconv.Itoa(int(*p.UnitMolecularWeight.UnitID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2EmpiricalFormula := select2.NewSelect2(jquery.Jq("select#empirical_formula"), nil)
	select2EmpiricalFormula.Select2Clear()
	if p.EmpiricalFormula.EmpiricalFormulaID != nil {
		select2EmpiricalFormula.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.EmpiricalFormula.EmpiricalFormulaLabel,
				Value:           strconv.Itoa(int(*p.EmpiricalFormula.EmpiricalFormulaID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linear_formula"), nil)
	select2LinearFormula.Select2Clear()
	if p.LinearFormula.LinearFormulaID != nil {
		select2LinearFormula.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.LinearFormula.LinearFormulaLabel,
				Value:           strconv.Itoa(int(*p.LinearFormula.LinearFormulaID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Casnumber := select2.NewSelect2(jquery.Jq("select#cas_number"), nil)
	select2Casnumber.Select2Clear()
	// if p.CasNumber.CasNumberID.Valid {
	if p.CasNumber.CasNumberID != nil {
		select2Casnumber.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text: *p.CasNumber.CasNumberLabel,
				// Value:           strconv.Itoa(int(p.CasNumber.CasNumberID.Int64)),
				Value:           strconv.Itoa(int(*p.CasNumber.CasNumberID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Cenumber := select2.NewSelect2(jquery.Jq("select#ce_number"), nil)
	select2Cenumber.Select2Clear()
	if p.CeNumber.CeNumberID != nil {
		select2Cenumber.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.CeNumber.CeNumberLabel,
				Value:           strconv.Itoa(int(*p.CeNumber.CeNumberID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_inchi").SetVal("")
	if p.ProductInchi != nil {
		jquery.Jq("#product_inchi").SetVal(*p.ProductInchi)
	}
	jquery.Jq("#product_inchikey").SetVal("")
	if p.ProductInchikey != nil {
		jquery.Jq("#product_inchikey").SetVal(*p.ProductInchikey)
	}
	jquery.Jq("#product_canonicalsmiles").SetVal("")
	if p.ProductCanonicalSmiles != nil {
		jquery.Jq("#product_canonicalsmiles").SetVal(*p.ProductCanonicalSmiles)
	}

	jquery.Jq("#product_specificity").SetVal("")
	if p.ProductSpecificity != nil {
		jquery.Jq("#product_specificity").SetVal(*p.ProductSpecificity)
	}
	jquery.Jq("#product_msds").SetVal("")
	if p.ProductMSDS != nil {
		jquery.Jq("#product_msds").SetVal(*p.ProductMSDS)
	}
	jquery.Jq("#product_sheet").SetVal("")
	if p.ProductSheet != nil {
		jquery.Jq("#product_sheet").SetVal(*p.ProductSheet)
	}
	jquery.Jq("#product_threedformula").SetVal("")
	if p.ProductThreeDFormula != nil {
		jquery.Jq("#product_threedformula").SetVal(*p.ProductThreeDFormula)
	}

	select2PhysicalState := select2.NewSelect2(jquery.Jq("select#physical_state"), nil)
	select2PhysicalState.Select2Clear()
	if p.PhysicalState.PhysicalStateID != nil {
		select2PhysicalState.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.PhysicalState.PhysicalStateLabel,
				Value:           strconv.Itoa(int(*p.PhysicalState.PhysicalStateID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Coc := select2.NewSelect2(jquery.Jq("select#class_of_compound"), nil)
	select2Coc.Select2Clear()
	for _, coc := range p.ClassOfCompound {
		select2Coc.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            coc.ClassOfCompoundLabel,
				Value:           strconv.Itoa(coc.ClassOfCompoundID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2SignalWord := select2.NewSelect2(jquery.Jq("select#signal_word"), nil)
	select2SignalWord.Select2Clear()
	if p.SignalWord.SignalWordID != nil {
		select2SignalWord.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *p.SignalWord.SignalWordLabel,
				Value:           strconv.Itoa(int(*p.SignalWord.SignalWordID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Symbols := select2.NewSelect2(jquery.Jq("select#symbols"), nil)
	select2Symbols.Select2Clear()
	for _, symbol := range p.Symbols {
		select2Symbols.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            symbol.SymbolLabel,
				Value:           strconv.Itoa(symbol.SymbolID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2HS := select2.NewSelect2(jquery.Jq("select#hazard_statements"), nil)
	select2HS.Select2Clear()
	for _, hs := range p.HazardStatements {
		select2HS.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            hs.HazardStatementLabel,
				Value:           strconv.Itoa(hs.HazardStatementID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionary_statements"), nil)
	select2PS.Select2Clear()
	for _, ps := range p.PrecautionaryStatements {
		select2PS.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            ps.PrecautionaryStatementLabel,
				Value:           strconv.Itoa(ps.PrecautionaryStatementID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_restricted").SetProp("checked", false)
	if p.ProductRestricted {
		jquery.Jq("#product_restricted").SetProp("checked", "checked")
	}
	jquery.Jq("#product_radioactive").SetProp("checked", false)
	if p.ProductRadioactive {
		jquery.Jq("#product_radioactive").SetProp("checked", "checked")
	}

	jquery.Jq("#product_disposalcomment").SetVal("")
	if p.ProductDisposalComment != nil {
		jquery.Jq("#product_disposalcomment").SetVal(*p.ProductDisposalComment)
	}
	jquery.Jq("#product_remark").SetVal("")
	if p.ProductRemark != nil {
		jquery.Jq("#product_remark").SetVal(*p.ProductRemark)
	}

	jquery.Jq("#product_number_per_carton").SetVal("")
	if p.ProductNumberPerCarton != nil && *p.ProductNumberPerCarton > 0 {
		jquery.Jq("#product_number_per_carton").SetVal(*p.ProductNumberPerCarton)
	}
	jquery.Jq("#product_number_per_bag").SetVal("")
	if p.ProductNumberPerBag != nil {
		jquery.Jq("#product_number_per_bag").SetVal(*p.ProductNumberPerBag)
	}

	// Chem/Bio/Consu detection.
	switch p.ProductType {
	case "CONS":
		Consufy()
	case "BIO":
		Biofy()
	default:
		Chemify()
	}

}

func SaveProduct(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#product"), nil).Valid() {
		return nil
	}

	globals.CurrentProduct = Product{Product: &models.Product{}}
	if jquery.Jq("input#product_id").GetVal().Truthy() {
		if globals.CurrentProduct.ProductID, err = strconv.Atoi(jquery.Jq("input#product_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if jquery.Jq("input#showchem:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductType = "chem"
	}
	if jquery.Jq("input#showbio:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductType = "bio"
	}
	if jquery.Jq("input#showconsu:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductType = "cons"
	}

	if jquery.Jq("input#product_temperature").GetVal().Truthy() {
		var productTemperature int
		if productTemperature, err = strconv.Atoi(jquery.Jq("input#product_temperature").GetVal().String()); err != nil {
			return nil
		}
		var ProductTemperaturePointer *int64 = new(int64)
		*ProductTemperaturePointer = int64(productTemperature)
		globals.CurrentProduct.ProductTemperature = ProductTemperaturePointer
	}

	if jquery.Jq("input#product_number_per_carton").GetVal().Truthy() {
		var productNumberPerCarton int
		if productNumberPerCarton, err = strconv.Atoi(jquery.Jq("input#product_number_per_carton").GetVal().String()); err != nil {
			return nil
		}
		globals.CurrentProduct.ProductNumberPerCarton = new(int64)
		*globals.CurrentProduct.ProductNumberPerCarton = int64(productNumberPerCarton)
	} else if jquery.Jq("input#showconsu:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductNumberPerCarton = new(int64)
		*globals.CurrentProduct.ProductNumberPerCarton = -1
	}

	if jquery.Jq("input#product_number_per_bag").GetVal().Truthy() {
		var productNumberPerBag int
		if productNumberPerBag, err = strconv.Atoi(jquery.Jq("input#product_number_per_bag").GetVal().String()); err != nil {
			return nil
		}

		globals.CurrentProduct.ProductNumberPerBag = new(int64)
		*globals.CurrentProduct.ProductNumberPerBag = int64(productNumberPerBag)
	}

	if jquery.Jq("input#product_specificity").GetVal().Truthy() {
		globals.CurrentProduct.ProductSpecificity = new(string)
		*globals.CurrentProduct.ProductSpecificity = jquery.Jq("input#product_specificity").GetVal().String()
	}

	if jquery.Jq("input#product_inchi").GetVal().Truthy() {
		globals.CurrentProduct.ProductInchi = new(string)
		*globals.CurrentProduct.ProductInchi = jquery.Jq("input#product_inchi").GetVal().String()
	}

	if jquery.Jq("input#product_inchikey").GetVal().Truthy() {
		globals.CurrentProduct.ProductInchikey = new(string)
		*globals.CurrentProduct.ProductInchikey = jquery.Jq("input#product_inchikey").GetVal().String()
	}

	if jquery.Jq("input#product_canonicalsmiles").GetVal().Truthy() {
		globals.CurrentProduct.ProductCanonicalSmiles = new(string)
		*globals.CurrentProduct.ProductCanonicalSmiles = jquery.Jq("input#product_canonicalsmiles").GetVal().String()
	}

	if jquery.Jq("input#product_molecularweight").GetVal().Truthy() {
		var productMolecularWeight float64
		if productMolecularWeight, err = strconv.ParseFloat(jquery.Jq("input#product_molecularweight").GetVal().String(), 64); err != nil {
			return nil
		}

		globals.CurrentProduct.ProductMolecularWeight = new(float64)
		*globals.CurrentProduct.ProductMolecularWeight = productMolecularWeight
	}

	if jquery.Jq("input#hidden_product_twodformula_content").Html() != "" {
		globals.CurrentProduct.ProductTwoDFormula = new(string)
		*globals.CurrentProduct.ProductTwoDFormula = jquery.Jq("input#hidden_product_twodformula_content").Html()
	}

	if jquery.Jq("input#product_threedformula").GetVal().Truthy() {
		globals.CurrentProduct.ProductThreeDFormula = new(string)
		*globals.CurrentProduct.ProductThreeDFormula = jquery.Jq("input#product_threedformula").GetVal().String()
	}

	// if jquery.Jq("#hidden_product_molformula_content").GetVal().Truthy() {
	// 	globals.CurrentProduct.ProductMolFormula = sql.NullString{
	// 		String: jquery.Jq("#hidden_product_molformula_content").GetVal().String(),
	// 		Valid:  true,
	// 	}
	// }

	if jquery.Jq("input#product_sheet").GetVal().Truthy() {
		globals.CurrentProduct.ProductSheet = new(string)
		*globals.CurrentProduct.ProductSheet = jquery.Jq("input#product_sheet").GetVal().String()
	}

	if jquery.Jq("input#product_msds").GetVal().Truthy() {
		globals.CurrentProduct.ProductMSDS = new(string)
		*globals.CurrentProduct.ProductMSDS = jquery.Jq("input#product_msds").GetVal().String()
	}

	if jquery.Jq("textarea#product_disposalcomment").GetVal().Truthy() {
		globals.CurrentProduct.ProductDisposalComment = new(string)
		*globals.CurrentProduct.ProductDisposalComment = jquery.Jq("textarea#product_disposalcomment").GetVal().String()
	}

	if jquery.Jq("textarea#product_remark").GetVal().Truthy() {
		globals.CurrentProduct.ProductRemark = new(string)
		*globals.CurrentProduct.ProductRemark = jquery.Jq("textarea#product_remark").GetVal().String()
	}

	if jquery.Jq("input#product_restricted:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductRestricted = true
	}

	if jquery.Jq("input#product_radioactive:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductRadioactive = true
	}

	select2UnitTemperature := select2.NewSelect2(jquery.Jq("select#unit_temperature"), nil)
	if len(select2UnitTemperature.Select2Data()) > 0 {
		select2ItemUnitTemperature := select2UnitTemperature.Select2Data()[0]
		globals.CurrentProduct.UnitTemperature = models.Unit{}
		var unitTemperatureId int
		if unitTemperatureId, err = strconv.Atoi(select2ItemUnitTemperature.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.UnitTemperature.UnitID = new(int64)
		globals.CurrentProduct.UnitTemperature.UnitLabel = new(string)
		*globals.CurrentProduct.UnitTemperature.UnitID = int64(unitTemperatureId)
		*globals.CurrentProduct.UnitTemperature.UnitLabel = select2ItemUnitTemperature.Text
	}

	select2UnitMolecularWeight := select2.NewSelect2(jquery.Jq("select#unit_molecularweight"), nil)
	if len(select2UnitMolecularWeight.Select2Data()) > 0 {
		select2ItemUnitMolecularWeight := select2UnitMolecularWeight.Select2Data()[0]
		globals.CurrentProduct.UnitMolecularWeight = models.Unit{}
		var unitMolecularWeightId int
		if unitMolecularWeightId, err = strconv.Atoi(select2ItemUnitMolecularWeight.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.UnitMolecularWeight.UnitID = new(int64)
		globals.CurrentProduct.UnitMolecularWeight.UnitLabel = new(string)
		*globals.CurrentProduct.UnitMolecularWeight.UnitID = int64(unitMolecularWeightId)
		*globals.CurrentProduct.UnitMolecularWeight.UnitLabel = select2ItemUnitMolecularWeight.Text
	}

	select2CasNumber := select2.NewSelect2(jquery.Jq("select#cas_number"), nil)
	if len(select2CasNumber.Select2Data()) > 0 {
		select2ItemCasNumber := select2CasNumber.Select2Data()[0]
		globals.CurrentProduct.CasNumber = models.CasNumber{}
		var casNumberId = -1

		if select2ItemCasNumber.IDIsDigit() {
			if casNumberId, err = strconv.Atoi(select2ItemCasNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		// globals.CurrentProduct.CasNumber.CasNumberID = sql.NullInt64{
		// 	Int64: int64(casNumberId),
		// 	Valid: true,
		// }
		var casNumberIdPointer *int64 = new(int64)
		*casNumberIdPointer = int64(casNumberId)
		var casNumberLabelPointer *string = new(string)
		*casNumberLabelPointer = select2ItemCasNumber.Text
		globals.CurrentProduct.CasNumber.CasNumberID = casNumberIdPointer
		globals.CurrentProduct.CasNumber.CasNumberLabel = casNumberLabelPointer
	}

	select2CeNumber := select2.NewSelect2(jquery.Jq("select#ce_number"), nil)
	if len(select2CeNumber.Select2Data()) > 0 {
		select2ItemCeNumber := select2CeNumber.Select2Data()[0]
		globals.CurrentProduct.CeNumber = models.CeNumber{}
		var ceNumberId = -1

		if select2ItemCeNumber.IDIsDigit() {
			if ceNumberId, err = strconv.Atoi(select2ItemCeNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.CeNumber.CeNumberID = new(int64)
		globals.CurrentProduct.CeNumber.CeNumberLabel = new(string)
		*globals.CurrentProduct.CeNumber.CeNumberID = int64(ceNumberId)
		*globals.CurrentProduct.CeNumber.CeNumberLabel = select2ItemCeNumber.Text
	}

	select2EmpiricalFormula := select2.NewSelect2(jquery.Jq("select#empirical_formula"), nil)
	if len(select2EmpiricalFormula.Select2Data()) > 0 {
		select2ItemEmpiricalFormula := select2EmpiricalFormula.Select2Data()[0]
		globals.CurrentProduct.EmpiricalFormula = models.EmpiricalFormula{}
		var empiricalFormulaId = -1

		if select2ItemEmpiricalFormula.IDIsDigit() {
			if empiricalFormulaId, err = strconv.Atoi(select2ItemEmpiricalFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		var EmpiricalFormulaIDPointer *int64 = new(int64)
		var EmpiricalFormulaLabelPointer *string = new(string)
		*EmpiricalFormulaIDPointer = int64(empiricalFormulaId)
		*EmpiricalFormulaLabelPointer = select2ItemEmpiricalFormula.Text
		globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID = EmpiricalFormulaIDPointer
		globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaLabel = EmpiricalFormulaLabelPointer
	}

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linear_formula"), nil)
	if len(select2LinearFormula.Select2Data()) > 0 {
		select2ItemLinearFormula := select2LinearFormula.Select2Data()[0]
		globals.CurrentProduct.LinearFormula = models.LinearFormula{}
		var linearFormulaId = -1

		if select2ItemLinearFormula.IDIsDigit() {
			if linearFormulaId, err = strconv.Atoi(select2ItemLinearFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		var LinearFormulaIDPointer *int64 = new(int64)
		var LinearFormulaLabelPointer *string = new(string)
		*LinearFormulaIDPointer = int64(linearFormulaId)
		*LinearFormulaLabelPointer = select2ItemLinearFormula.Text
		globals.CurrentProduct.LinearFormula.LinearFormulaID = LinearFormulaIDPointer
		globals.CurrentProduct.LinearFormula.LinearFormulaLabel = LinearFormulaLabelPointer

	}

	select2Name := select2.NewSelect2(jquery.Jq("select#name"), nil)
	if len(select2Name.Select2Data()) > 0 {
		select2ItemName := select2Name.Select2Data()[0]
		globals.CurrentProduct.Name = models.Name{}
		var nameId = -1

		if select2ItemName.IDIsDigit() {
			if nameId, err = strconv.Atoi(select2ItemName.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.Name.NameID = nameId
		globals.CurrentProduct.Name.NameLabel = select2ItemName.Text
	}

	select2PhysicalState := select2.NewSelect2(jquery.Jq("select#physical_state"), nil)
	if len(select2PhysicalState.Select2Data()) > 0 {
		select2ItemPhysicalState := select2PhysicalState.Select2Data()[0]
		globals.CurrentProduct.PhysicalState = models.PhysicalState{}
		var physicalStateId = -1

		if select2ItemPhysicalState.IDIsDigit() {
			if physicalStateId, err = strconv.Atoi(select2ItemPhysicalState.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.PhysicalState.PhysicalStateID = new(int64)
		globals.CurrentProduct.PhysicalState.PhysicalStateLabel = new(string)
		*globals.CurrentProduct.PhysicalState.PhysicalStateID = int64(physicalStateId)
		*globals.CurrentProduct.PhysicalState.PhysicalStateLabel = select2ItemPhysicalState.Text
	}

	select2SignalWord := select2.NewSelect2(jquery.Jq("select#signal_word"), nil)
	if len(select2SignalWord.Select2Data()) > 0 {
		select2ItemSignalWord := select2SignalWord.Select2Data()[0]
		globals.CurrentProduct.SignalWord = models.SignalWord{}
		var signalWordId int
		if signalWordId, err = strconv.Atoi(select2ItemSignalWord.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.SignalWord.SignalWordID = new(int64)
		globals.CurrentProduct.SignalWord.SignalWordLabel = new(string)
		*globals.CurrentProduct.SignalWord.SignalWordID = int64(signalWordId)
		*globals.CurrentProduct.SignalWord.SignalWordLabel = select2ItemSignalWord.Text
	}

	select2Category := select2.NewSelect2(jquery.Jq("select#category"), nil)
	if len(select2Category.Select2Data()) > 0 {
		select2ItemCategory := select2Category.Select2Data()[0]
		globals.CurrentProduct.Category = models.Category{}
		var categoryId = -1

		if select2ItemCategory.IDIsDigit() {
			if categoryId, err = strconv.Atoi(select2ItemCategory.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.Category.CategoryID = new(int64)
		globals.CurrentProduct.Category.CategoryLabel = new(string)
		*globals.CurrentProduct.Category.CategoryID = int64(categoryId)
		*globals.CurrentProduct.Category.CategoryLabel = select2ItemCategory.Text
	}

	select2ProducerRef := select2.NewSelect2(jquery.Jq("select#producer_ref"), nil)
	if len(select2ProducerRef.Select2Data()) > 0 {
		select2ItemProducerRef := select2ProducerRef.Select2Data()[0]
		globals.CurrentProduct.ProducerRef = models.ProducerRef{}
		var producerrefId = -1

		if select2ItemProducerRef.IDIsDigit() && !(select2ItemProducerRef.Id == select2ItemProducerRef.Text) {
			if producerrefId, err = strconv.Atoi(select2ItemProducerRef.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.ProducerRef.ProducerRefID = new(int64)
		globals.CurrentProduct.ProducerRef.ProducerRefLabel = new(string)
		*globals.CurrentProduct.ProducerRef.ProducerRefID = int64(producerrefId)
		*globals.CurrentProduct.ProducerRef.ProducerRefLabel = select2ItemProducerRef.Text

		var producerId int
		select2ItemProducer := select2.NewSelect2(jquery.Jq("select#producer"), nil).Select2Data()[0]
		if producerId, err = strconv.Atoi(select2ItemProducer.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		producerIdInt64 := int64(producerId)
		globals.CurrentProduct.ProducerRef.Producer = &models.Producer{
			ProducerID: &producerIdInt64,
		}
	}

	select2Coc := select2.NewSelect2(jquery.Jq("select#class_of_compound"), nil)
	for _, select2Item := range select2Coc.Select2Data() {
		classofcompound := models.ClassOfCompound{}
		var classofcompoundID = -1

		if select2Item.IDIsDigit() {
			if classofcompoundID, err = strconv.Atoi(select2Item.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}
		classofcompound.ClassOfCompoundID = classofcompoundID
		classofcompound.ClassOfCompoundLabel = select2Item.Text

		globals.CurrentProduct.ClassOfCompound = append(globals.CurrentProduct.ClassOfCompound, classofcompound)
	}

	select2Synonyms := select2.NewSelect2(jquery.Jq("select#synonyms"), nil)
	for _, select2Item := range select2Synonyms.Select2Data() {
		synonym := models.Name{}
		var nameID = -1

		if select2Item.IDIsDigit() {
			if nameID, err = strconv.Atoi(select2Item.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}
		synonym.NameID = nameID
		synonym.NameLabel = select2Item.Text

		globals.CurrentProduct.Synonyms = append(globals.CurrentProduct.Synonyms, synonym)
	}

	select2Symbols := select2.NewSelect2(jquery.Jq("select#symbols"), nil)
	for _, select2Item := range select2Symbols.Select2Data() {
		symbol := models.Symbol{}
		if symbol.SymbolID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		symbol.SymbolLabel = select2Item.Text

		globals.CurrentProduct.Symbols = append(globals.CurrentProduct.Symbols, symbol)
	}

	select2HS := select2.NewSelect2(jquery.Jq("select#hazard_statements"), nil)
	for _, select2Item := range select2HS.Select2Data() {
		hazardstatement := models.HazardStatement{}
		if hazardstatement.HazardStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		hazardstatement.HazardStatementLabel = select2Item.Text

		globals.CurrentProduct.HazardStatements = append(globals.CurrentProduct.HazardStatements, hazardstatement)
	}

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionary_statements"), nil)
	for _, select2Item := range select2PS.Select2Data() {
		precautionarystatement := models.PrecautionaryStatement{}
		if precautionarystatement.PrecautionaryStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		precautionarystatement.PrecautionaryStatementLabel = select2Item.Text

		globals.CurrentProduct.PrecautionaryStatements = append(globals.CurrentProduct.PrecautionaryStatements, precautionarystatement)
	}

	select2Tags := select2.NewSelect2(jquery.Jq("select#tags"), nil)
	for _, select2Item := range select2Tags.Select2Data() {
		tag := models.Tag{}
		var tagId = -1

		if select2Item.IDIsDigit() {
			if tagId, err = strconv.Atoi(select2Item.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}
		tag.TagID = tagId
		tag.TagLabel = select2Item.Text

		globals.CurrentProduct.Tags = append(globals.CurrentProduct.Tags, tag)
	}

	select2SupplierRefs := select2.NewSelect2(jquery.Jq("select#supplier_refs"), nil)
	for _, select2Item := range select2SupplierRefs.Select2Data() {
		supplierref := models.SupplierRef{}
		var supplierrefId = -1

		if select2Item.IDIsDigit() && !strings.HasPrefix(select2Item.Text, fmt.Sprintf("%s@", select2Item.Id)) {
			if supplierrefId, err = strconv.Atoi(select2Item.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		supplierref.SupplierRefID = supplierrefId
		supplierref.SupplierRefLabel = strings.Split(select2Item.Text, "@")[0]

		supplierId := supplierrefToSupplier[supplierref.SupplierRefLabel]

		supplierref.Supplier = &models.Supplier{
			SupplierID: &supplierId,
		}

		globals.CurrentProduct.SupplierRefs = append(globals.CurrentProduct.SupplierRefs, supplierref)
	}

	if (!jquery.Jq("form#product input#product_id").GetVal().IsUndefined()) && jquery.Jq("form#product input#product_id").GetVal().String() != "" {
		ajaxURL = fmt.Sprintf("%sproducts/%d", ApplicationProxyPath, globals.CurrentProduct.ProductID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sproducts", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	if dataBytes, err = json.Marshal(globals.CurrentProduct); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			var (
				err error
			)

			if err = json.Unmarshal([]byte(data.String()), &CurrentProduct); err != nil {
				jsutils.DisplayGenericErrorMessage()
				fmt.Println(err)
			} else {
				href := fmt.Sprintf("%sv/products", ApplicationProxyPath)
				jsutils.ClearSearch(js.Null(), nil)
				jsutils.LoadContent("div#content", "product", href, Product_SaveCallback, globals.CurrentProduct.ProductID)
			}

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
