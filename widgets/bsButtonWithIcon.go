//go:build go1.24 && js && wasm

package widgets

import (
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

type bsButtonWithIcon struct {
	button
}

func NewBSButtonWithIcon(buttonAttrs ButtonAttributes, iconAttrs IconAttributes, buttonStyles []themes.BSClass) *bsButtonWithIcon {

	bi := &bsButtonWithIcon{}

	for _, style := range buttonStyles {
		buttonAttrs.Classes = append(buttonAttrs.Classes, style.ToString())
	}

	b := NewButton(buttonAttrs)
	i := NewIcon(iconAttrs)

	b.AppendChild(i)
	bi.button.HTMLElement = b

	return bi

}
