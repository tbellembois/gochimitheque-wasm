package personpass

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
)

func PersonPass_listCallback(this js.Value, args []js.Value) interface{} {

	Jq("#personp").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"person_password": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				EqualTo:  "#person_passwordagain",
				Remote: ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"person_passwordagain": {
				EqualTo: "#person_password",
			},
		},
		Messages: map[string]ValidateMessage{
			"person_password": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
				EqualTo:  locales.Translate("not_same_password", HTTPHeaderAcceptLanguage),
			},
			"person_passwordagain": {
				EqualTo: locales.Translate("not_same_password", HTTPHeaderAcceptLanguage),
			},
		},
	})

	Jq("#search").Hide()
	Jq("#actions").Hide()

	return nil

}

func SavePersonPassword(this js.Value, args []js.Value) interface{} {

	var (
		person    Person
		dataBytes []byte
		err       error
	)

	if !Jq("#personp").Valid() {
		return nil
	}

	person.PersonPassword = Jq("input#person_password").GetVal().String()
	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    fmt.Sprintf("%speoplep", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("person_password_updated_message", HTTPHeaderAcceptLanguage))

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
