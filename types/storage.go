package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
)

type Storages struct {
	Rows     []*Storage `json:"rows"`
	Total    int        `json:"total"`
	ExportFn string     `json:"exportfn"`
}

type Storage struct {
	*models.Storage
}

type Borrowing struct {
	BorrowingID      sql.NullInt64  `json:"borrowing_id"`
	BorrowingComment sql.NullString `json:"borrowing_comment"`
	Person           *Person        `json:"person"` // logged person
	//Storage          `json:"storage"`
	Borrower *Person `json:"borrower"` // logged person
}

func (s Storage) FromJsJSONValue(jsvalue js.Value) Storage {

	var (
		storage Storage
		err     error
	)

	jsEntityString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsEntityString), &storage); err != nil {
		fmt.Println(err)
	}

	return storage

}

func (s Storages) GetTotal() int {

	return s.Total

}

func (s Storages) GetExportFn() string {

	return s.ExportFn

}
