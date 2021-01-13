package login

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/localStorage"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
)

func GetAnnounce(this js.Value, args []js.Value) interface{} {

	Ajax{
		URL:    fmt.Sprintf("%swelcomeannounce", ApplicationProxyPath),
		Method: "get",
		Done: func(data js.Value) {

			var (
				welcomeAnnounce WelcomeAnnounce
				err             error
			)

			if err = json.Unmarshal([]byte(data.String()), &welcomeAnnounce); err != nil {
				fmt.Println(err)
			}

			Jq("#wannounce").SetHtml(welcomeAnnounce.WelcomeAnnounceHTML)

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func GetToken(this js.Value, args []js.Value) interface{} {

	var (
		dataBytes []byte
		person    *Person
		err       error
	)

	person = &Person{}
	person.PersonEmail = Jq("#person_email").GetVal().String()
	person.PersonPassword = Jq("#person_password").GetVal().String()

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    fmt.Sprintf("%sget-token", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			localStorage.Clear()

			productCallbackWrapper := func(args ...interface{}) {
				product.Product_listCallback(js.Null(), nil)
			}
			utils.LoadContent("product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper)
			//utils.RedirectTo(ApplicationProxyPath)

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func GetCaptcha(this js.Value, args []js.Value) interface{} {

	if !Jq("#authForm").Valid() {
		return nil
	}

	type captcha struct {
		Image string `json:"image"`
		UID   string `json:"uid"`
	}

	var (
		captchaResponse *captcha
		err             error
	)

	utils.DisplaySuccessMessage(locales.Translate("resetpassword_areyourobot", HTTPHeaderAcceptLanguage))

	Ajax{
		URL:    fmt.Sprintf("%scaptcha", ApplicationProxyPath),
		Method: "get",
		Done: func(data js.Value) {

			if err = json.Unmarshal([]byte(data.String()), &captchaResponse); err != nil {
				fmt.Println(err)
			}

			Jq("#captcha_uid").SetVal(captchaResponse.UID)
			Jq("#captcha-img").SetProp("src", fmt.Sprintf("data:image/png;base64,%s", captchaResponse.Image))
			Jq("#captcha-row").RemoveClass("invisible")

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func ResetPassword(this js.Value, args []js.Value) interface{} {

	var (
		dataBytes []byte
		person    *Person
		err       error
	)

	person = &Person{}
	person.PersonEmail = Jq("#person_email").GetVal().String()
	person.CaptchaText = Jq("#captcha_text").GetVal().String()
	person.CaptchaUID = Jq("#captcha_uid").GetVal().String()

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    fmt.Sprintf("%sreset-password", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			utils.DisplaySuccessMessage(fmt.Sprintf(locales.Translate("resetpassword_message_mailsentto", HTTPHeaderAcceptLanguage), person.PersonEmail))

			Jq("#person_email").SetVal("")
			Jq("#captcha_text").SetVal("")
			Jq("#captcha-row").Hide()

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func Login_listCallback(this js.Value, args []js.Value) interface{} {

	Jq("#authForm").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"person_email": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Email:    true,
			},
		},
		Messages: map[string]ValidateMessage{
			"person_email": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
				Email:    locales.Translate("invalid_email", HTTPHeaderAcceptLanguage),
			},
		},
	})

	return nil

}
