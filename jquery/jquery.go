//go:build go1.24 && js && wasm

package jquery

import (
	"syscall/js"
)

var (
	Jq func(args ...interface{}) Jquery
)

type Jquery struct {
	Object js.Value
}

func init() {

	Jq = NewJquery

}

func NewJquery(args ...interface{}) Jquery {

	return Jquery{Object: js.Global().Get("jQuery").New(args...)}

}

func (jq Jquery) Append(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("append", i...)
	return jq

}

func (jq Jquery) Html() string {

	return jq.Object.Call("html").String()

}

func (jq Jquery) SetHtml(i interface{}) Jquery {

	switch i.(type) {
	case func(int, string) string, string:
	default:
		print("SetHtml Argument should be 'string' or 'func(int, string) string'")
	}

	jq.Object = jq.Object.Call("html", i)
	return jq

}

func (jq Jquery) Prop(i ...interface{}) interface{} {

	return jq.Object.Call("prop", i...)

}

func (jq Jquery) SetProp(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("prop", i...)
	return jq

}

func (jq Jquery) RemoveProp(property string) Jquery {

	jq.Object = jq.Object.Call("removeProp", property)
	return jq

}

func (jq Jquery) SetVal(i interface{}) Jquery {

	jq.Object.Call("val", i)
	return jq

}

func (jq Jquery) SetAttr(i ...interface{}) Jquery {

	jq.Object.Call("attr", i...)
	return jq

}

func (jq Jquery) GetAttr(i interface{}) js.Value {

	return jq.Object.Call("attr", i)

}

func (jq Jquery) GetVal() js.Value {

	return jq.Object.Call("val")

}

func (jq Jquery) Show() Jquery {

	jq.Object = jq.Object.Call("collapse", "show")
	return jq

}

func (jq Jquery) Hide() Jquery {

	jq.Object = jq.Object.Call("collapse", "hide")
	return jq

}

func (jq Jquery) FadeIn() Jquery {

	jq.Object = jq.Object.Call("fadeIn")
	return jq

}

func (jq Jquery) FadeOut() Jquery {

	jq.Object = jq.Object.Call("fadeOut")
	return jq

}

func (jq Jquery) On(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("on", i...)
	return jq

}

func (jq Jquery) Find(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("find", i...)
	return jq

}

func (jq Jquery) Remove(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("remove", i...)
	return jq

}

func (jq Jquery) Empty() Jquery {

	jq.Object = jq.Object.Call("empty")
	return jq

}

func (jq Jquery) RemoveClass(property string) Jquery {

	jq.Object = jq.Object.Call("removeClass", property)
	return jq

}

func (jq Jquery) AddClass(property string) Jquery {

	jq.Object = jq.Object.Call("addClass", property)
	return jq

}

func (jq Jquery) HasClass(class string) bool {

	return jq.Object.Call("hasClass", class).Bool()

}

func (jq Jquery) SetVisible() Jquery {

	jq = jq.RemoveClass("invisible")
	jq = jq.AddClass("visible")
	return jq

}

func (jq Jquery) SetInvisible() Jquery {

	jq = jq.RemoveClass("visible")
	jq = jq.AddClass("invisible")
	return jq

}

func (jq Jquery) Not(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("not", i...)
	return jq

}

func (jq Jquery) Is(i ...interface{}) bool {

	return jq.Object.Call("is", i...).Bool()

}
