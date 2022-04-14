package personpass

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque/models"
)

func PersonPass_listCallback(this js.Value, args []js.Value) interface{} {

	validate.NewValidate(jquery.Jq("#personp"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"person_password": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				EqualTo:  "#person_passwordagain",
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"person_passwordagain": {
				EqualTo: "#person_password",
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"person_password": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
				EqualTo:  locales.Translate("not_same_password", HTTPHeaderAcceptLanguage),
			},
			"person_passwordagain": {
				EqualTo: locales.Translate("not_same_password", HTTPHeaderAcceptLanguage),
			},
		},
	})
	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

	return nil

}

func SavePersonPassword(this js.Value, args []js.Value) interface{} {

	var (
		person    Person
		dataBytes []byte
		err       error
	)

	person.Person = &models.Person{}

	if !validate.NewValidate(jquery.Jq("#personp"), nil).Valid() {
		return nil
	}

	person.PersonPassword = jquery.Jq("input#person_password").GetVal().String()
	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%speoplep", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("person_password_updated_message", HTTPHeaderAcceptLanguage))

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
