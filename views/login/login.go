package login

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/views/menu"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque/models"
)

func GetAnnounce(this js.Value, args []js.Value) interface{} {

	ajax.Ajax{
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

			jquery.Jq("#wannounce").SetHtml(welcomeAnnounce.WelcomeAnnounceHTML)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

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

	person = &Person{Person: &models.Person{}}
	person.PersonEmail = jquery.Jq("#person_email").GetVal().String()
	person.PersonPassword = jquery.Jq("#person_password").GetVal().String()

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%sget-token", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			AfterLogin_listCallback(js.Null(), nil)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func ScanQRdone(this js.Value, args []js.Value) interface{} {

	qr := args[0].String()

	var (
		person    *Person
		dataBytes []byte
		err       error
	)

	person = &Person{Person: &models.Person{}}
	person.QRCode = []byte(qr)

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%sget-token", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			AfterLogin_listCallback(js.Null(), nil)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	jsutils.CloseQR(js.Null(), nil)

	return nil

}

func GetCaptcha(this js.Value, args []js.Value) interface{} {

	if !validate.NewValidate(jquery.Jq("#authform"), nil).Valid() {
		return nil
	}

	var (
		err    error
		isLDAP bool
	)

	email := jquery.Jq("#person_email").GetVal().String()

	ajax.Ajax{
		URL:    fmt.Sprintf("%speople/isldap/%s", ApplicationProxyPath, email),
		Method: "get",
		Done: func(data js.Value) {

			if err = json.Unmarshal([]byte(data.String()), &isLDAP); err != nil {
				fmt.Println(err)
			}

			if isLDAP {
				jsutils.DisplayErrorMessage(locales.Translate("resetpassword_alert_ldap", HTTPHeaderAcceptLanguage))
			} else {
				getCaptcha()
			}

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func getCaptcha() {

	type captcha struct {
		Image string `json:"image"`
		UID   string `json:"uid"`
	}

	var (
		captchaResponse *captcha
		err             error
	)

	jsutils.DisplaySuccessMessage(locales.Translate("resetpassword_areyourobot", HTTPHeaderAcceptLanguage))

	ajax.Ajax{
		URL:    fmt.Sprintf("%scaptcha", ApplicationProxyPath),
		Method: "get",
		Done: func(data js.Value) {

			if err = json.Unmarshal([]byte(data.String()), &captchaResponse); err != nil {
				fmt.Println(err)
			}

			jquery.Jq("#captcha_uid").SetVal(captchaResponse.UID)
			jquery.Jq("#captcha-img").SetProp("src", fmt.Sprintf("data:image/png;base64,%s", captchaResponse.Image))
			jquery.Jq("#captcha-row").RemoveClass("invisible")

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

}

func ResetPassword(this js.Value, args []js.Value) interface{} {

	var (
		dataBytes []byte
		person    *Person
		err       error
	)

	person = &Person{Person: &models.Person{}}
	person.PersonEmail = jquery.Jq("#person_email").GetVal().String()
	person.CaptchaText = jquery.Jq("#captcha_text").GetVal().String()
	person.CaptchaUID = jquery.Jq("#captcha_uid").GetVal().String()

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%sreset-password", ApplicationProxyPath),
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			jsutils.DisplaySuccessMessage(fmt.Sprintf(locales.Translate("resetpassword_message_mailsentto", HTTPHeaderAcceptLanguage), person.PersonEmail))

			jquery.Jq("#person_email").SetVal("")
			jquery.Jq("#captcha_text").SetVal("")
			jquery.Jq("#captcha-row").Hide()

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func Login_listCallback(this js.Value, args []js.Value) interface{} {

	validate.NewValidate(jquery.Jq("#authform"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"person_email": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Email:    true,
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"person_email": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
				Email:    locales.Translate("invalid_email", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	jquery.Jq("#authform input").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		event := args[0]
		if !event.Get("which").IsUndefined() && event.Get("which").Int() == 13 {

			event.Call("preventDefault")
			GetToken(js.Null(), nil)

		}

		return nil

	}))

	return nil

}

func AfterLogin_listCallback(this js.Value, args []js.Value) interface{} {

	productCallbackWrapper := func(args ...interface{}) {
		product.Product_listCallback(js.Null(), nil)
	}

	jsutils.LoadContent("div#menu", "menu", fmt.Sprintf("%smenu", ApplicationProxyPath), menu.ShowIfAuthorizedMenuItems)
	jsutils.LoadContent("div#searchbar", "product", fmt.Sprintf("%ssearch", ApplicationProxyPath), search.Search_listCallback)
	jsutils.LoadContent("div#content", "product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper)

	// Can not read Email and ID from container as those values
	// are not set before login.
	cookie := js.Global().Get("document").Get("cookie").String()

	regexEmail := regexp.MustCompile(`email=(\S+)`)
	matchEmail := regexEmail.FindStringSubmatch(cookie)[1]

	regexId := regexp.MustCompile(`id=(\S+)`)
	matchId := regexId.FindStringSubmatch(cookie)[1]

	matchEmail = strings.TrimRight(matchEmail, ";")
	matchId = strings.TrimRight(matchId, ";")

	ConnectedUserEmail = matchEmail
	var err error
	if ConnectedUserID, err = strconv.Atoi(matchId); err != nil {
		panic(err)
	}

	jquery.Jq("#logged").SetHtml(ConnectedUserEmail)

	return nil

}
