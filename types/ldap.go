//go:build go1.24 && js && wasm

package types

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/go-ldap/ldap/v3"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type LDAPSearchResults struct {
	NbResults int              `json:"NbResults"`
	R         LDAPSearchResult `json:"R"`
}

type LDAPSearchResult struct {
	Entries []*LDAPEntry `json:"Entries"`
}

type LDAPEntry struct {
	*ldap.Entry
}

func (l LDAPSearchResults) GetRowConcreteTypeName() string {

	return ""

}

func (l LDAPSearchResults) IsExactMatch() bool {

	return false

}

func (l LDAPSearchResults) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		ldapSearchResults LDAPSearchResults
		err               error
	)

	jsLDAPSearchResultsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsLDAPSearchResultsString), &ldapSearchResults); err != nil {
		fmt.Println(err)
	}

	return ldapSearchResults

}

func (l LDAPSearchResults) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, l.NbResults)

	if l.NbResults > 0 {
		for i, entry := range l.R.Entries {
			select2ItemAble[i] = entry
		}
	}

	return select2ItemAble

}

func (l LDAPSearchResults) GetTotal() int {

	return l.NbResults

}

func (l *LDAPEntry) MarshalJSON() ([]byte, error) {
	type Copy LDAPEntry
	return json.Marshal(&struct {
		Id   int64  `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   l.GetSelect2Id(),
		Text: l.GetSelect2Text(),

		Copy: (*Copy)(l),
	})
}

func (l LDAPEntry) GetTotal() int {

	return 0

}

func (l LDAPEntry) GetSelect2Id() int64 {

	h := sha1.New()
	h.Write([]byte(l.Entry.DN))
	bs := h.Sum(nil)

	return int64(binary.BigEndian.Uint64(bs))

}

func (l LDAPEntry) GetSelect2Text() string {

	if l.Entry != nil {
		// return l.Entry.GetAttributeValue("cn")
		return l.Entry.DN
	}

	return ""

}

func (l LDAPEntry) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		ldapEntry LDAPEntry
		err       error
	)

	jsLDAPEntryString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsLDAPEntryString), &ldapEntry); err != nil {
		fmt.Println(err)
	}

	return ldapEntry

}
