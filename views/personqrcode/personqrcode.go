package personqrcode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func PersonQRCode_listCallback(this js.Value, args []js.Value) interface{} {

	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

	go func() {

		url := fmt.Sprintf("%speople/%d", globals.ApplicationProxyPath, globals.ConnectedUserID)
		method := "get"

		done := func(data js.Value) {

			var (
				person Person
				err    error
			)

			if err = json.Unmarshal([]byte(data.String()), &person); err != nil {
				jsutils.DisplayGenericErrorMessage()
				return
			}

			qrcode := base64.StdEncoding.EncodeToString(person.QRCode)
			qrcodeImg := widgets.NewImg(widgets.ImgAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Attributes: map[string]string{
						"style": "border: 1px solid black;",
					},
				},
				Height: "128px",
				Width:  "128px",
				Src:    fmt.Sprintf("data:image/png;base64,%s", qrcode),
			})
			jquery.Jq("#qrcode").Empty()
			jquery.Jq("#qrcode").Append(qrcodeImg.OuterHTML())

		}
		fail := func(data js.Value) {

			jsutils.DisplayGenericErrorMessage()

		}

		ajax.Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

	}()

	return nil

}
