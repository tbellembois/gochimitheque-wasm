package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"
)

type Storages struct {
	Rows     []*Storage `json:"rows"`
	Total    int        `json:"total"`
	ExportFn string     `json:"exportfn"`
}

type Storage struct {
	StorageID               sql.NullInt64   `json:"storage_id"`
	StorageCreationDate     time.Time       `json:"storage_creationdate"`
	StorageModificationDate time.Time       `json:"storage_modificationdate"`
	StorageEntryDate        sql.NullTime    `json:"storage_entrydate"`
	StorageExitDate         sql.NullTime    `json:"storage_exitdate"`
	StorageOpeningDate      sql.NullTime    `json:"storage_openingdate"`
	StorageExpirationDate   sql.NullTime    `json:"storage_expirationdate"`
	StorageComment          sql.NullString  `json:"storage_comment"`
	StorageReference        sql.NullString  `json:"storage_reference"`
	StorageBatchNumber      sql.NullString  `json:"storage_batchnumber"`
	StorageQuantity         sql.NullFloat64 `json:"storage_quantity"`
	StorageNbItem           int             `json:"storage_nbitem"`
	StorageBarecode         sql.NullString  `json:"storage_barecode"`
	StorageQRCode           []byte          `json:"storage_qrcode"`
	StorageToDestroy        sql.NullBool    `json:"storage_todestroy"`
	StorageArchive          sql.NullBool    `json:"storage_archive"`
	StorageConcentration    sql.NullInt64   `json:"storage_concentration"`
	Person                  `json:"person"`
	Product                 `json:"product"`
	StoreLocation           `json:"storelocation"`
	UnitQuantity            Unit `json:"unit_quantity"`
	UnitConcentration       Unit `json:"unit_concentration"`
	Supplier                `json:"supplier"`
	Storage                 *Storage   `json:"storage"`   // history reference storage
	Borrowing               *Borrowing `json:"borrowing"` // not un db but sqlx requires the "db" entry

	// storage history count
	StorageHC int `json:"storage_hc"` // not in db but sqlx requires the "db" entry
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
