package product

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

func ValidateProductCeNumberBeforeSend(this js.Value, args []js.Value) interface{} {

	settings := args[1]

	id := "-1"
	pid := jquery.Jq("input#product_id")

	if pid.Object.Length() > 0 {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/product/%s/cenumber/", ApplicationProxyPath, id))

	return nil

}

func ValidateProductCeNumberData(this js.Value, args []js.Value) interface{} {

	return select2.NewSelect2(jquery.Jq("select#cenumber"), nil).Select2Data()[0].Text

}

func ValidateProductCasNumberBeforeSend(this js.Value, args []js.Value) interface{} {

	settings := args[1]

	id := "-1"
	pid := jquery.Jq("input#product_id")

	if pid.Object.Length() > 0 {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/product/%s/casnumber/", ApplicationProxyPath, id))

	return nil

}

func ValidateProductEmpiricalFormulaBeforeSend(this js.Value, args []js.Value) interface{} {

	settings := args[1]

	id := "-1"
	pid := jquery.Jq("input#product_id")

	if pid.Object.Length() > 0 {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/product/%s/empiricalformula/", ApplicationProxyPath, id))

	return nil

}

func ValidateProductCasNumberData1(this js.Value, args []js.Value) interface{} {

	return select2.NewSelect2(jquery.Jq("select#casnumber"), nil).Select2Data()[0].Text

}

func ValidateProductCasNumberData2(this js.Value, args []js.Value) interface{} {

	return jquery.Jq("#product_specificity").GetVal().String()

}

func ValidateProductEmpiricalFormulaData(this js.Value, args []js.Value) interface{} {

	return select2.NewSelect2(jquery.Jq("select#empiricalformula"), nil).Select2Data()[0].Text

}
