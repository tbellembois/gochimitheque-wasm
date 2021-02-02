package localstorage

import (
	"syscall/js"
)

type LocalStorage struct {
	Object js.Value
}

func NewLocalStorage() LocalStorage {

	return LocalStorage{Object: js.Global().Get("window").Get("localStorage")}

}

// SetItem save val into localstorage under key.
func (l LocalStorage) SetItem(key string, val string) {
	l.Object.Call("setItem", key, val)
}

// RemoveItem remove val into localstorage under key.
func (l LocalStorage) RemoveItem(key string) {
	l.Object.Call("removeItem", key)
}

// GetItem return the key item.
func (l LocalStorage) GetItem(key string) string {
	obj := l.Object.Call("getItem", key)
	if !obj.Truthy() {
		return ""
	}
	return obj.String()
}

// Clear the cache
func (l LocalStorage) Clear() {
	l.Object.Call("clear")
}
