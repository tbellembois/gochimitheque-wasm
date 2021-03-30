package product

import (
	"database/sql"
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
)

func FillInProductForm(p Product, id string) {

	jquery.Jq(fmt.Sprintf("#%s #product_id", id)).SetVal(p.ProductID)

	select2Category := select2.NewSelect2(jquery.Jq("select#category"), nil)
	select2Category.Select2Clear()
	if p.CategoryID.Valid {
		select2Category.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.Category.CategoryLabel.String,
				Value:           strconv.Itoa(int(p.Category.CategoryID.Int64)),
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
	if p.Producer.ProducerID.Valid {
		select2Producer.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.Producer.ProducerLabel.String,
				Value:           strconv.Itoa(int(p.Producer.ProducerID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2ProducerRef := select2.NewSelect2(jquery.Jq("select#producerref"), nil)
	select2ProducerRef.Select2Clear()
	if p.ProducerRef.ProducerRefID.Valid {
		select2ProducerRef.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.ProducerRef.ProducerRefLabel.String,
				Value:           strconv.Itoa(int(p.ProducerRef.ProducerRefID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2SupplierRef := select2.NewSelect2(jquery.Jq("select#supplierrefs"), nil)
	select2SupplierRef.Select2Clear()
	for _, supplierref := range p.SupplierRefs {

		supplierrefToSupplier[supplierref.SupplierRefLabel] = supplierref.Supplier.SupplierID.Int64

		select2SupplierRef.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            fmt.Sprintf("%s@%s", supplierref.SupplierRefLabel, supplierref.Supplier.SupplierLabel.String),
				Value:           strconv.Itoa(supplierref.SupplierRefID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_temperature").SetVal("")
	if p.ProductTemperature.Valid {
		jquery.Jq("#product_temperature").SetVal(p.ProductTemperature.Int64)
	}

	select2UnitTemperature := select2.NewSelect2(jquery.Jq("select#unit_temperature"), nil)
	select2UnitTemperature.Select2Clear()
	if p.UnitTemperature.UnitID.Valid {
		select2UnitTemperature.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.UnitTemperature.UnitLabel.String,
				Value:           strconv.Itoa(int(p.UnitTemperature.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2EmpiricalFormula := select2.NewSelect2(jquery.Jq("select#empiricalformula"), nil)
	select2EmpiricalFormula.Select2Clear()
	if p.EmpiricalFormula.EmpiricalFormulaID.Valid {
		select2EmpiricalFormula.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.EmpiricalFormula.EmpiricalFormulaLabel.String,
				Value:           strconv.Itoa(int(p.EmpiricalFormula.EmpiricalFormulaID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linearformula"), nil)
	select2LinearFormula.Select2Clear()
	if p.LinearFormula.LinearFormulaID.Valid {
		select2LinearFormula.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.LinearFormula.LinearFormulaLabel.String,
				Value:           strconv.Itoa(int(p.LinearFormula.LinearFormulaID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Casnumber := select2.NewSelect2(jquery.Jq("select#casnumber"), nil)
	select2Casnumber.Select2Clear()
	if p.CasNumber.CasNumberID.Valid {
		select2Casnumber.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.CasNumber.CasNumberLabel.String,
				Value:           strconv.Itoa(int(p.CasNumber.CasNumberID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Cenumber := select2.NewSelect2(jquery.Jq("select#cenumber"), nil)
	select2Cenumber.Select2Clear()
	if p.CeNumber.CeNumberID.Valid {
		select2Cenumber.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.CeNumber.CeNumberLabel.String,
				Value:           strconv.Itoa(int(p.CeNumber.CeNumberID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#product_specificity").SetVal("")
	if p.ProductSpecificity.Valid {
		jquery.Jq("#product_specificity").SetVal(p.ProductSpecificity.String)
	}
	jquery.Jq("#product_msds").SetVal("")
	if p.ProductMSDS.Valid {
		jquery.Jq("#product_msds").SetVal(p.ProductMSDS.String)
	}
	jquery.Jq("#product_sheet").SetVal("")
	if p.ProductSheet.Valid {
		jquery.Jq("#product_sheet").SetVal(p.ProductSheet.String)
	}
	jquery.Jq("#product_threedformula").SetVal("")
	if p.ProductThreeDFormula.Valid {
		jquery.Jq("#product_threedformula").SetVal(p.ProductThreeDFormula.String)
	}

	select2PhysicalState := select2.NewSelect2(jquery.Jq("select#physicalstate"), nil)
	select2PhysicalState.Select2Clear()
	if p.PhysicalState.PhysicalStateID.Valid {
		select2PhysicalState.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.PhysicalState.PhysicalStateLabel.String,
				Value:           strconv.Itoa(int(p.PhysicalState.PhysicalStateID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Coc := select2.NewSelect2(jquery.Jq("select#classofcompound"), nil)
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

	select2SignalWord := select2.NewSelect2(jquery.Jq("select#signalword"), nil)
	select2SignalWord.Select2Clear()
	if p.SignalWord.SignalWordID.Valid {
		select2SignalWord.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.SignalWord.SignalWordLabel.String,
				Value:           strconv.Itoa(int(p.SignalWord.SignalWordID.Int64)),
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

	select2HS := select2.NewSelect2(jquery.Jq("select#hazardstatements"), nil)
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

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionarystatements"), nil)
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
	if p.ProductRestricted.Valid && p.ProductRestricted.Bool {
		jquery.Jq("#product_restricted").SetProp("checked", "checked")
	}
	jquery.Jq("#product_radioactive").SetProp("checked", false)
	if p.ProductRadioactive.Valid && p.ProductRadioactive.Bool {
		jquery.Jq("#product_radioactive").SetProp("checked", "checked")
	}

	jquery.Jq("#product_disposalcomment").SetVal("")
	if p.ProductDisposalComment.Valid {
		jquery.Jq("#product_disposalcomment").SetVal(p.ProductDisposalComment.String)
	}
	jquery.Jq("#product_remark").SetVal("")
	if p.ProductRemark.Valid {
		jquery.Jq("#product_remark").SetVal(p.ProductRemark.String)
	}

	jquery.Jq("#product_number_per_carton").SetVal("")
	if p.ProductNumberPerCarton.Valid && p.ProductNumberPerCarton.Int64 > 0 {
		jquery.Jq("#product_number_per_carton").SetVal(p.ProductNumberPerCarton.Int64)
	}
	jquery.Jq("#product_number_per_bag").SetVal("")
	if p.ProductNumberPerBag.Valid {
		jquery.Jq("#product_number_per_bag").SetVal(p.ProductNumberPerBag.Int64)
	}

	// Chem/Bio/Consu detection.
	if p.ProductNumberPerCarton.Valid {
		Consufy()
	} else if p.ProducerRef.ProducerRefID.Valid {
		Biofy()
	} else {
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

	globals.CurrentProduct = Product{}
	if jquery.Jq("input#product_id").GetVal().Truthy() {
		if globals.CurrentProduct.ProductID, err = strconv.Atoi(jquery.Jq("input#product_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if jquery.Jq("input#product_temperature").GetVal().Truthy() {
		var productTemperature int
		if productTemperature, err = strconv.Atoi(jquery.Jq("input#product_temperature").GetVal().String()); err != nil {
			return nil
		}
		globals.CurrentProduct.ProductTemperature = sql.NullInt64{
			Int64: int64(productTemperature),
			Valid: true,
		}
	}

	if jquery.Jq("input#product_number_per_carton").GetVal().Truthy() {
		var productNumberPerCarton int
		if productNumberPerCarton, err = strconv.Atoi(jquery.Jq("input#product_number_per_carton").GetVal().String()); err != nil {
			return nil
		}
		globals.CurrentProduct.ProductNumberPerCarton = sql.NullInt64{
			Int64: int64(productNumberPerCarton),
			Valid: true,
		}
	} else if jquery.Jq("input#showconsu:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductNumberPerCarton = sql.NullInt64{
			Int64: -1,
			Valid: true,
		}
	}

	if jquery.Jq("input#product_number_per_bag").GetVal().Truthy() {
		var productNumberPerBag int
		if productNumberPerBag, err = strconv.Atoi(jquery.Jq("input#product_number_per_bag").GetVal().String()); err != nil {
			return nil
		}
		globals.CurrentProduct.ProductNumberPerBag = sql.NullInt64{
			Int64: int64(productNumberPerBag),
			Valid: true,
		}
	}

	if jquery.Jq("input#product_specificity").GetVal().Truthy() {
		globals.CurrentProduct.ProductSpecificity = sql.NullString{
			String: jquery.Jq("input#product_specificity").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("input#hidden_product_twodformula_content").Html() != "" {
		globals.CurrentProduct.ProductTwoDFormula = sql.NullString{
			String: jquery.Jq("input#hidden_product_twodformula_content").Html(),
			Valid:  true,
		}
	}

	if jquery.Jq("input#product_threedformula").GetVal().Truthy() {
		globals.CurrentProduct.ProductThreeDFormula = sql.NullString{
			String: jquery.Jq("input#product_threedformula").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("#hidden_product_molformula_content").GetVal().Truthy() {
		globals.CurrentProduct.ProductMolFormula = sql.NullString{
			String: jquery.Jq("#hidden_product_molformula_content").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("input#product_sheet").GetVal().Truthy() {
		globals.CurrentProduct.ProductSheet = sql.NullString{
			String: jquery.Jq("input#product_sheet").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("input#product_msds").GetVal().Truthy() {
		globals.CurrentProduct.ProductMSDS = sql.NullString{
			String: jquery.Jq("input#product_msds").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("textarea#product_disposalcomment").GetVal().Truthy() {
		globals.CurrentProduct.ProductDisposalComment = sql.NullString{
			String: jquery.Jq("textarea#product_disposalcomment").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("textarea#product_remark").GetVal().Truthy() {
		globals.CurrentProduct.ProductRemark = sql.NullString{
			String: jquery.Jq("textarea#product_remark").GetVal().String(),
			Valid:  true,
		}
	}

	if jquery.Jq("input#product_restricted:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductRestricted = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if jquery.Jq("input#product_radioactive:checked").Object.Length() > 0 {
		globals.CurrentProduct.ProductRadioactive = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	select2UnitTemperature := select2.NewSelect2(jquery.Jq("select#unit_temperature"), nil)
	if len(select2UnitTemperature.Select2Data()) > 0 {
		select2ItemUnitTemperature := select2UnitTemperature.Select2Data()[0]
		globals.CurrentProduct.UnitTemperature = Unit{}
		var unitTemperatureId int
		if unitTemperatureId, err = strconv.Atoi(select2ItemUnitTemperature.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.UnitTemperature.UnitID = sql.NullInt64{
			Int64: int64(unitTemperatureId),
			Valid: true,
		}
		globals.CurrentProduct.UnitTemperature.UnitLabel = sql.NullString{
			String: select2ItemUnitTemperature.Text,
			Valid:  true,
		}
	}

	select2CasNumber := select2.NewSelect2(jquery.Jq("select#casnumber"), nil)
	if len(select2CasNumber.Select2Data()) > 0 {
		select2ItemCasNumber := select2CasNumber.Select2Data()[0]
		globals.CurrentProduct.CasNumber = CasNumber{}
		var casNumberId = -1

		if select2ItemCasNumber.IDIsDigit() {
			if casNumberId, err = strconv.Atoi(select2ItemCasNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.CasNumber.CasNumberID = sql.NullInt64{
			Int64: int64(casNumberId),
			Valid: true,
		}
		globals.CurrentProduct.CasNumber.CasNumberLabel = sql.NullString{
			String: select2ItemCasNumber.Text,
			Valid:  true,
		}
	}

	select2CeNumber := select2.NewSelect2(jquery.Jq("select#cenumber"), nil)
	if len(select2CeNumber.Select2Data()) > 0 {
		select2ItemCeNumber := select2CeNumber.Select2Data()[0]
		globals.CurrentProduct.CeNumber = CeNumber{}
		var ceNumberId = -1

		if select2ItemCeNumber.IDIsDigit() {
			if ceNumberId, err = strconv.Atoi(select2ItemCeNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.CeNumber.CeNumberID = sql.NullInt64{
			Int64: int64(ceNumberId),
			Valid: true,
		}
		globals.CurrentProduct.CeNumber.CeNumberLabel = sql.NullString{
			String: select2ItemCeNumber.Text,
			Valid:  true,
		}
	}

	select2EmpiricalFormula := select2.NewSelect2(jquery.Jq("select#empiricalformula"), nil)
	if len(select2EmpiricalFormula.Select2Data()) > 0 {
		select2ItemEmpiricalFormula := select2EmpiricalFormula.Select2Data()[0]
		globals.CurrentProduct.EmpiricalFormula = EmpiricalFormula{}
		var empiricalFormulaId = -1

		if select2ItemEmpiricalFormula.IDIsDigit() {
			if empiricalFormulaId, err = strconv.Atoi(select2ItemEmpiricalFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{
			Int64: int64(empiricalFormulaId),
			Valid: true,
		}
		globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaLabel = sql.NullString{
			String: select2ItemEmpiricalFormula.Text,
			Valid:  true,
		}
	}

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linearformula"), nil)
	if len(select2LinearFormula.Select2Data()) > 0 {
		select2ItemLinearFormula := select2LinearFormula.Select2Data()[0]
		globals.CurrentProduct.LinearFormula = LinearFormula{}
		var linearFormulaId = -1

		if select2ItemLinearFormula.IDIsDigit() {
			if linearFormulaId, err = strconv.Atoi(select2ItemLinearFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.LinearFormula.LinearFormulaID = sql.NullInt64{
			Int64: int64(linearFormulaId),
			Valid: true,
		}
		globals.CurrentProduct.LinearFormula.LinearFormulaLabel = sql.NullString{
			String: select2ItemLinearFormula.Text,
			Valid:  true,
		}
	}

	select2Name := select2.NewSelect2(jquery.Jq("select#name"), nil)
	if len(select2Name.Select2Data()) > 0 {
		select2ItemName := select2Name.Select2Data()[0]
		globals.CurrentProduct.Name = Name{}
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

	select2PhysicalState := select2.NewSelect2(jquery.Jq("select#physicalstate"), nil)
	if len(select2PhysicalState.Select2Data()) > 0 {
		select2ItemPhysicalState := select2PhysicalState.Select2Data()[0]
		globals.CurrentProduct.PhysicalState = PhysicalState{}
		var physicalStateId = -1

		if select2ItemPhysicalState.IDIsDigit() {
			if physicalStateId, err = strconv.Atoi(select2ItemPhysicalState.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.PhysicalState.PhysicalStateID = sql.NullInt64{
			Int64: int64(physicalStateId),
			Valid: true,
		}
		globals.CurrentProduct.PhysicalState.PhysicalStateLabel = sql.NullString{
			String: select2ItemPhysicalState.Text,
			Valid:  true,
		}
	}

	select2SignalWord := select2.NewSelect2(jquery.Jq("select#signalword"), nil)
	if len(select2SignalWord.Select2Data()) > 0 {
		select2ItemSignalWord := select2SignalWord.Select2Data()[0]
		globals.CurrentProduct.SignalWord = SignalWord{}
		var signalWordId int
		if signalWordId, err = strconv.Atoi(select2ItemSignalWord.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.SignalWord.SignalWordID = sql.NullInt64{
			Int64: int64(signalWordId),
			Valid: true,
		}
		globals.CurrentProduct.SignalWord.SignalWordLabel = sql.NullString{
			String: select2ItemSignalWord.Text,
			Valid:  true,
		}
	}

	select2Category := select2.NewSelect2(jquery.Jq("select#category"), nil)
	if len(select2Category.Select2Data()) > 0 {
		select2ItemCategory := select2Category.Select2Data()[0]
		globals.CurrentProduct.Category = Category{}
		var categoryId = -1

		if select2ItemCategory.IDIsDigit() {
			if categoryId, err = strconv.Atoi(select2ItemCategory.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.Category.CategoryID = sql.NullInt64{
			Int64: int64(categoryId),
			Valid: true,
		}
		globals.CurrentProduct.Category.CategoryLabel = sql.NullString{
			String: select2ItemCategory.Text,
			Valid:  true,
		}
	}

	select2ProducerRef := select2.NewSelect2(jquery.Jq("select#producerref"), nil)
	if len(select2ProducerRef.Select2Data()) > 0 {
		select2ItemProducerRef := select2ProducerRef.Select2Data()[0]
		globals.CurrentProduct.ProducerRef = ProducerRef{}
		var producerrefId = -1

		if select2ItemProducerRef.IDIsDigit() && !(select2ItemProducerRef.Id == select2ItemProducerRef.Text) {
			if producerrefId, err = strconv.Atoi(select2ItemProducerRef.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentProduct.ProducerRef.ProducerRefID = sql.NullInt64{
			Int64: int64(producerrefId),
			Valid: true,
		}
		globals.CurrentProduct.ProducerRef.ProducerRefLabel = sql.NullString{
			String: select2ItemProducerRef.Text,
			Valid:  true,
		}

		var producerId int
		select2ItemProducer := select2.NewSelect2(jquery.Jq("select#producer"), nil).Select2Data()[0]
		if producerId, err = strconv.Atoi(select2ItemProducer.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentProduct.ProducerRef.Producer = &Producer{
			ProducerID: sql.NullInt64{
				Int64: int64(producerId),
				Valid: true,
			},
		}
	}

	select2Coc := select2.NewSelect2(jquery.Jq("select#classofcompound"), nil)
	for _, select2Item := range select2Coc.Select2Data() {
		classofcompound := ClassOfCompound{}
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
		synonym := Name{}
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
		symbol := Symbol{}
		if symbol.SymbolID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		symbol.SymbolLabel = select2Item.Text

		globals.CurrentProduct.Symbols = append(globals.CurrentProduct.Symbols, symbol)
	}

	select2HS := select2.NewSelect2(jquery.Jq("select#hazardstatements"), nil)
	for _, select2Item := range select2HS.Select2Data() {
		hazardstatement := HazardStatement{}
		if hazardstatement.HazardStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		hazardstatement.HazardStatementLabel = select2Item.Text

		globals.CurrentProduct.HazardStatements = append(globals.CurrentProduct.HazardStatements, hazardstatement)
	}

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionarystatements"), nil)
	for _, select2Item := range select2PS.Select2Data() {
		precautionarystatement := PrecautionaryStatement{}
		if precautionarystatement.PrecautionaryStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		precautionarystatement.PrecautionaryStatementLabel = select2Item.Text

		globals.CurrentProduct.PrecautionaryStatements = append(globals.CurrentProduct.PrecautionaryStatements, precautionarystatement)
	}

	select2Tags := select2.NewSelect2(jquery.Jq("select#tags"), nil)
	for _, select2Item := range select2Tags.Select2Data() {
		tag := Tag{}
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

	select2SupplierRefs := select2.NewSelect2(jquery.Jq("select#supplierrefs"), nil)
	for _, select2Item := range select2SupplierRefs.Select2Data() {
		supplierref := SupplierRef{}
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

		supplierref.Supplier = &Supplier{
			SupplierID: sql.NullInt64{
				Valid: true,
				Int64: supplierId,
			},
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
