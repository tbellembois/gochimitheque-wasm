package types

import (
	"syscall/js"
)

var (
	Jq func(args ...interface{}) Jquery
)

type Jquery struct {
	Object js.Value
}

func NewJquery(args ...interface{}) Jquery {

	return Jquery{Object: js.Global().Get("jQuery").New(args...)}

}

func (jq Jquery) Append(i ...interface{}) Jquery {

	jq.Object = jq.Object.Call("append", i...)
	return jq

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
