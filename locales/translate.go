//go:build go1.24 && js && wasm

package locales

import (
	"fmt"
	"strings"
	"syscall/js"

	"golang.org/x/text/language"
)

func TranslateWrapper(this js.Value, args []js.Value) interface{} {

	Translate(args[0].String(), args[1].String())
	return nil

}

// Translate translates the s string in the "accept" language.
func Translate(s string, accept string) string {
	if s == "" {
		return ""
	}

	ts, _, e := language.ParseAcceptLanguage(accept)
	if e != nil {
		fmt.Println(e)
		ts = []language.Tag{language.English}
	}

	// The t entries are
	// ordered by the preferred language.
	for _, t := range ts {

		newt := strings.Replace(t.String(), "-", "_", -1)
		js_locale_varname := fmt.Sprintf("locale_%s_%s", newt, s)

		translated := js.Global().Get(js_locale_varname)
		if !translated.IsUndefined() {

			return translated.String()

		} else {

			js_locale_varname = fmt.Sprintf("locale_en_EN_%s", s)
			translated = js.Global().Get(js_locale_varname)
			if !translated.IsUndefined() {
				return translated.String()
			}

		}

	}

	// We should never come here.
	return "translation error for " + s
}
