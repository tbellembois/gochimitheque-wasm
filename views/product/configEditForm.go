package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func FillInProductForm(p Product, id string) {

	Jq(fmt.Sprintf("#%s #product_id", id)).SetVal(p.ProductID)

	Jq("select#category").Select2Clear()
	if p.CategoryID.Valid {
		Jq("select#category").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.Category.CategoryLabel.String,
				Value:           strconv.Itoa(int(p.Category.CategoryID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#tags").Select2Clear()
	for _, tag := range p.Tags {
		Jq("select#tags").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            tag.TagLabel,
				Value:           strconv.Itoa(tag.TagID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#name").Select2Clear()
	Jq("select#name").Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            p.Name.NameLabel,
			Value:           strconv.Itoa(int(p.Name.NameID)),
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())
	Jq("select#synonyms").Select2Clear()
	for _, synonym := range p.Synonyms {
		Jq("select#synonyms").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            synonym.NameLabel,
				Value:           strconv.Itoa(synonym.NameID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#producer").Select2Clear()
	if p.Producer.ProducerID.Valid {
		Jq("select#producer").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.Producer.ProducerLabel.String,
				Value:           strconv.Itoa(int(p.Producer.ProducerID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#producerref").Select2Clear()
	if p.ProducerRef.ProducerRefID.Valid {
		Jq("select#producerref").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.ProducerRef.ProducerRefLabel.String,
				Value:           strconv.Itoa(int(p.ProducerRef.ProducerRefID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#supplierrefs").Select2Clear()
	for _, supplierref := range p.SupplierRefs {

		supplierrefToSupplier[supplierref.SupplierRefLabel] = supplierref.Supplier.SupplierID.Int64

		Jq("select#supplierrefs").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            fmt.Sprintf("%s@%s", supplierref.SupplierRefLabel, supplierref.Supplier.SupplierLabel.String),
				Value:           strconv.Itoa(supplierref.SupplierRefID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#product_temperature").SetVal("")
	if p.ProductTemperature.Valid {
		Jq("#product_temperature").SetVal(p.ProductTemperature.Int64)
	}
	Jq("select#unit_temperature").Select2Clear()
	if p.UnitTemperature.UnitID.Valid {
		Jq("select#unit_temperature").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.UnitTemperature.UnitLabel.String,
				Value:           strconv.Itoa(int(p.UnitTemperature.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#empiricalformula").Select2Clear()
	if p.EmpiricalFormula.EmpiricalFormulaID.Valid {
		Jq("select#empiricalformula").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.EmpiricalFormula.EmpiricalFormulaLabel.String,
				Value:           strconv.Itoa(int(p.EmpiricalFormula.EmpiricalFormulaID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#linearformula").Select2Clear()
	if p.LinearFormula.LinearFormulaID.Valid {
		Jq("select#linearformula").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.LinearFormula.LinearFormulaLabel.String,
				Value:           strconv.Itoa(int(p.LinearFormula.LinearFormulaID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#casnumber").Select2Clear()
	if p.CasNumber.CasNumberID.Valid {
		Jq("select#casnumber").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.CasNumber.CasNumberLabel.String,
				Value:           strconv.Itoa(int(p.CasNumber.CasNumberID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#cenumber").Select2Clear()
	if p.CeNumber.CeNumberID.Valid {
		Jq("select#cenumber").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.CeNumber.CeNumberLabel.String,
				Value:           strconv.Itoa(int(p.CeNumber.CeNumberID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#product_specificity").SetVal("")
	if p.ProductSpecificity.Valid {
		Jq("#product_specificity").SetVal(p.ProductSpecificity.String)
	}
	Jq("#product_msds").SetVal("")
	if p.ProductMSDS.Valid {
		Jq("#product_msds").SetVal(p.ProductMSDS.String)
	}
	Jq("#product_sheet").SetVal("")
	if p.ProductSheet.Valid {
		Jq("#product_sheet").SetVal(p.ProductSheet.String)
	}

	Jq("#product_threedformula").SetVal("")
	if p.ProductThreeDFormula.Valid {
		Jq("#product_threedformula").SetVal(p.ProductThreeDFormula.String)
	}

	Jq("select#physicalstate").Select2Clear()
	if p.PhysicalState.PhysicalStateID.Valid {
		Jq("select#physicalstate").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.PhysicalState.PhysicalStateLabel.String,
				Value:           strconv.Itoa(int(p.PhysicalState.PhysicalStateID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#classofcompound").Select2Clear()
	for _, coc := range p.ClassOfCompound {
		Jq("select#classofcompound").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            coc.ClassOfCompoundLabel,
				Value:           strconv.Itoa(coc.ClassOfCompoundID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#signalword").Select2Clear()
	if p.SignalWord.SignalWordID.Valid {
		Jq("select#signalword").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            p.SignalWord.SignalWordLabel.String,
				Value:           strconv.Itoa(int(p.SignalWord.SignalWordID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#symbols").Select2Clear()
	for _, symbol := range p.Symbols {
		Jq("select#symbols").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            symbol.SymbolLabel,
				Value:           strconv.Itoa(symbol.SymbolID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#hazardstatements").Select2Clear()
	for _, hs := range p.HazardStatements {
		Jq("select#hazardstatements").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            hs.HazardStatementLabel,
				Value:           strconv.Itoa(hs.HazardStatementID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}
	Jq("select#precautionarystatements").Select2Clear()
	for _, ps := range p.PrecautionaryStatements {
		Jq("select#precautionarystatements").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            ps.PrecautionaryStatementLabel,
				Value:           strconv.Itoa(ps.PrecautionaryStatementID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#product_restricted").SetProp("checked", false)
	if p.ProductRestricted.Valid && p.ProductRestricted.Bool {
		Jq("#product_restricted").SetProp("checked", "checked")
	}
	Jq("#product_radioactive").SetProp("checked", false)
	if p.ProductRadioactive.Valid && p.ProductRadioactive.Bool {
		Jq("#product_radioactive").SetProp("checked", "checked")
	}

	Jq("#product_disposalcomment").SetVal("")
	if p.ProductDisposalComment.Valid {
		Jq("#product_disposalcomment").SetVal(p.ProductDisposalComment.String)
	}
	Jq("#product_remark").SetVal("")
	if p.ProductRemark.Valid {
		Jq("#product_remark").SetVal(p.ProductRemark.String)
	}

	// Chem/Bio detection.
	if Jq("select#producerref").GetVal().Truthy() {
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

	if !(Jq("form#product").Valid()) {
		return nil
	}

	types.CurrentProduct = Product{}
	if Jq("input#product_id").GetVal().Truthy() {
		if types.CurrentProduct.ProductID, err = strconv.Atoi(Jq("input#product_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if Jq("input#product_temperature").GetVal().Truthy() {
		var productTemperature int
		if productTemperature, err = strconv.Atoi(Jq("input#product_temperature").GetVal().String()); err != nil {
			return nil
		}
		types.CurrentProduct.ProductTemperature = sql.NullInt64{
			Int64: int64(productTemperature),
			Valid: true,
		}
	}

	if Jq("input#product_specificity").GetVal().Truthy() {
		types.CurrentProduct.ProductSpecificity = sql.NullString{
			String: Jq("input#product_specificity").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("input#hidden_product_twodformula_content").Html() != "" {
		types.CurrentProduct.ProductTwoDFormula = sql.NullString{
			String: Jq("input#hidden_product_twodformula_content").Html(),
			Valid:  true,
		}
	}

	if Jq("input#product_threedformula").GetVal().Truthy() {
		types.CurrentProduct.ProductThreeDFormula = sql.NullString{
			String: Jq("input#product_threedformula").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("#hidden_product_molformula_content").GetVal().Truthy() {
		types.CurrentProduct.ProductMolFormula = sql.NullString{
			String: Jq("#hidden_product_molformula_content").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("input#product_sheet").GetVal().Truthy() {
		types.CurrentProduct.ProductSheet = sql.NullString{
			String: Jq("input#product_sheet").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("input#product_msds").GetVal().Truthy() {
		types.CurrentProduct.ProductMSDS = sql.NullString{
			String: Jq("input#product_msds").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("textarea#product_disposalcomment").GetVal().Truthy() {
		types.CurrentProduct.ProductDisposalComment = sql.NullString{
			String: Jq("textarea#product_disposalcomment").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("textarea#product_remark").GetVal().Truthy() {
		types.CurrentProduct.ProductRemark = sql.NullString{
			String: Jq("textarea#product_remark").GetVal().String(),
			Valid:  true,
		}
	}

	if Jq("input#product_restricted:checked").Object.Length() > 0 {
		types.CurrentProduct.ProductRestricted = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if Jq("input#product_radioactive:checked").Object.Length() > 0 {
		types.CurrentProduct.ProductRadioactive = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if len(Jq("select#unit_temperature").Select2Data()) > 0 {
		select2ItemUnitTemperature := Jq("select#unit_temperature").Select2Data()[0]
		types.CurrentProduct.UnitTemperature = Unit{}
		var unitTemperatureId int
		if unitTemperatureId, err = strconv.Atoi(select2ItemUnitTemperature.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentProduct.UnitTemperature.UnitID = sql.NullInt64{
			Int64: int64(unitTemperatureId),
			Valid: true,
		}
		types.CurrentProduct.UnitTemperature.UnitLabel = sql.NullString{
			String: select2ItemUnitTemperature.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#casnumber").Select2Data()) > 0 {
		select2ItemCasNumber := Jq("select#casnumber").Select2Data()[0]
		types.CurrentProduct.CasNumber = CasNumber{}
		var casNumberId = -1

		if select2ItemCasNumber.IDIsDigit() {
			if casNumberId, err = strconv.Atoi(select2ItemCasNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.CasNumber.CasNumberID = sql.NullInt64{
			Int64: int64(casNumberId),
			Valid: true,
		}
		types.CurrentProduct.CasNumber.CasNumberLabel = sql.NullString{
			String: select2ItemCasNumber.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#cenumber").Select2Data()) > 0 {
		select2ItemCeNumber := Jq("select#cenumber").Select2Data()[0]
		types.CurrentProduct.CeNumber = CeNumber{}
		var ceNumberId = -1

		if select2ItemCeNumber.IDIsDigit() {
			if ceNumberId, err = strconv.Atoi(select2ItemCeNumber.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.CeNumber.CeNumberID = sql.NullInt64{
			Int64: int64(ceNumberId),
			Valid: true,
		}
		types.CurrentProduct.CeNumber.CeNumberLabel = sql.NullString{
			String: select2ItemCeNumber.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#empiricalformula").Select2Data()) > 0 {
		select2ItemEmpiricalFormula := Jq("select#empiricalformula").Select2Data()[0]
		types.CurrentProduct.EmpiricalFormula = EmpiricalFormula{}
		var empiricalFormulaId = -1

		if select2ItemEmpiricalFormula.IDIsDigit() {
			if empiricalFormulaId, err = strconv.Atoi(select2ItemEmpiricalFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{
			Int64: int64(empiricalFormulaId),
			Valid: true,
		}
		types.CurrentProduct.EmpiricalFormula.EmpiricalFormulaLabel = sql.NullString{
			String: select2ItemEmpiricalFormula.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#linearformula").Select2Data()) > 0 {
		select2ItemLinearFormula := Jq("select#linearformula").Select2Data()[0]
		types.CurrentProduct.LinearFormula = LinearFormula{}
		var linearFormulaId = -1

		if select2ItemLinearFormula.IDIsDigit() {
			if linearFormulaId, err = strconv.Atoi(select2ItemLinearFormula.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.LinearFormula.LinearFormulaID = sql.NullInt64{
			Int64: int64(linearFormulaId),
			Valid: true,
		}
		types.CurrentProduct.LinearFormula.LinearFormulaLabel = sql.NullString{
			String: select2ItemLinearFormula.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#name").Select2Data()) > 0 {
		select2ItemName := Jq("select#name").Select2Data()[0]
		types.CurrentProduct.Name = Name{}
		var nameId = -1

		if select2ItemName.IDIsDigit() {
			if nameId, err = strconv.Atoi(select2ItemName.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.Name.NameID = nameId
		types.CurrentProduct.Name.NameLabel = select2ItemName.Text
	}

	if len(Jq("select#physicalstate").Select2Data()) > 0 {
		select2ItemPhysicalState := Jq("select#physicalstate").Select2Data()[0]
		types.CurrentProduct.PhysicalState = PhysicalState{}
		var physicalStateId = -1

		if select2ItemPhysicalState.IDIsDigit() {
			if physicalStateId, err = strconv.Atoi(select2ItemPhysicalState.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.PhysicalState.PhysicalStateID = sql.NullInt64{
			Int64: int64(physicalStateId),
			Valid: true,
		}
		types.CurrentProduct.PhysicalState.PhysicalStateLabel = sql.NullString{
			String: select2ItemPhysicalState.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#signalword").Select2Data()) > 0 {
		select2ItemSignalWord := Jq("select#signalword").Select2Data()[0]
		types.CurrentProduct.SignalWord = SignalWord{}
		var signalWordId int
		if signalWordId, err = strconv.Atoi(select2ItemSignalWord.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentProduct.SignalWord.SignalWordID = sql.NullInt64{
			Int64: int64(signalWordId),
			Valid: true,
		}
		types.CurrentProduct.SignalWord.SignalWordLabel = sql.NullString{
			String: select2ItemSignalWord.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#category").Select2Data()) > 0 {
		select2ItemCategory := Jq("select#category").Select2Data()[0]
		types.CurrentProduct.Category = Category{}
		var categoryId = -1

		if select2ItemCategory.IDIsDigit() {
			if categoryId, err = strconv.Atoi(select2ItemCategory.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.Category.CategoryID = sql.NullInt64{
			Int64: int64(categoryId),
			Valid: true,
		}
		types.CurrentProduct.Category.CategoryLabel = sql.NullString{
			String: select2ItemCategory.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#producerref").Select2Data()) > 0 {
		select2ItemProducerRef := Jq("select#producerref").Select2Data()[0]
		types.CurrentProduct.ProducerRef = ProducerRef{}
		var producerrefId = -1

		if select2ItemProducerRef.IDIsDigit() && !(select2ItemProducerRef.Id == select2ItemProducerRef.Text) {
			if producerrefId, err = strconv.Atoi(select2ItemProducerRef.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentProduct.ProducerRef.ProducerRefID = sql.NullInt64{
			Int64: int64(producerrefId),
			Valid: true,
		}
		types.CurrentProduct.ProducerRef.ProducerRefLabel = sql.NullString{
			String: select2ItemProducerRef.Text,
			Valid:  true,
		}

		var producerId int
		select2ItemProducer := Jq("select#producer").Select2Data()[0]
		if producerId, err = strconv.Atoi(select2ItemProducer.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentProduct.ProducerRef.Producer = &Producer{
			ProducerID: sql.NullInt64{
				Int64: int64(producerId),
				Valid: true,
			},
		}
	}

	for _, select2Item := range Jq("select#classofcompound").Select2Data() {
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

		types.CurrentProduct.ClassOfCompound = append(types.CurrentProduct.ClassOfCompound, classofcompound)
	}

	for _, select2Item := range Jq("select#synonyms").Select2Data() {
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

		types.CurrentProduct.Synonyms = append(types.CurrentProduct.Synonyms, synonym)
	}

	for _, select2Item := range Jq("select#symbols").Select2Data() {
		symbol := Symbol{}
		if symbol.SymbolID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		symbol.SymbolLabel = select2Item.Text

		types.CurrentProduct.Symbols = append(types.CurrentProduct.Symbols, symbol)
	}

	for _, select2Item := range Jq("select#hazardstatements").Select2Data() {
		hazardstatement := HazardStatement{}
		if hazardstatement.HazardStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		hazardstatement.HazardStatementLabel = select2Item.Text

		types.CurrentProduct.HazardStatements = append(types.CurrentProduct.HazardStatements, hazardstatement)
	}

	for _, select2Item := range Jq("select#precautionarystatements").Select2Data() {
		precautionarystatement := PrecautionaryStatement{}
		if precautionarystatement.PrecautionaryStatementID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		precautionarystatement.PrecautionaryStatementLabel = select2Item.Text

		types.CurrentProduct.PrecautionaryStatements = append(types.CurrentProduct.PrecautionaryStatements, precautionarystatement)
	}

	for _, select2Item := range Jq("select#tags").Select2Data() {
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

		types.CurrentProduct.Tags = append(types.CurrentProduct.Tags, tag)
	}

	for _, select2Item := range Jq("select#supplierrefs").Select2Data() {
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

		types.CurrentProduct.SupplierRefs = append(types.CurrentProduct.SupplierRefs, supplierref)
	}

	if (!Jq("form#product input#product_id").GetVal().IsUndefined()) && Jq("form#product input#product_id").GetVal().String() != "" {
		ajaxURL = fmt.Sprintf("%sproducts/%d", ApplicationProxyPath, types.CurrentProduct.ProductID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sproducts", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	if dataBytes, err = json.Marshal(types.CurrentProduct); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			var (
				err error
			)

			if err = json.Unmarshal([]byte(data.String()), &CurrentProduct); err != nil {
				utils.DisplayGenericErrorMessage()
				fmt.Println(err)
			} else {
				href := fmt.Sprintf("%sv/products", ApplicationProxyPath)
				search.ClearSearch(js.Null(), nil)
				utils.LoadContent("product", href, Product_SaveCallback, types.CurrentProduct.ProductID)
			}

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
