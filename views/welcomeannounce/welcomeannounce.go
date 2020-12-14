package welcomeannounce

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
)

func WelcomeAnnounce_listCallback(this js.Value, args []js.Value) interface{} {

	done := func(data js.Value) {

		var (
			welcomeAnnounce WelcomeAnnounce
			err             error
		)

		if err = json.Unmarshal([]byte(data.String()), &welcomeAnnounce); err != nil {
			fmt.Println(err)
		}
		Jq("#welcomeannounce_text").SetHtml(welcomeAnnounce.WelcomeAnnounceText)

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: "get",
		URL:    fmt.Sprintf("%swelcomeannounce", ApplicationProxyPath),
		Done:   done,
		Fail:   fail,
	}.Send()

	Jq("#search").Hide()
	Jq("#actions").Hide()

	return nil

}
func SaveWelcomeAnnounce(this js.Value, args []js.Value) interface{} {

	var (
		welcomeAnnounce WelcomeAnnounce
		dataBytes       []byte
		err             error
	)

	welcomeAnnounce.WelcomeAnnounceText = Jq("#welcomeannounce_text").GetVal().String()
	if dataBytes, err = json.Marshal(welcomeAnnounce); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    fmt.Sprintf("%swelcomeannounce", ApplicationProxyPath),
		Method: "put",
		Data:   dataBytes,
		Done: func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("welcomeannounce_text_modificationsuccess", HTTPHeaderAcceptLanguage))

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
