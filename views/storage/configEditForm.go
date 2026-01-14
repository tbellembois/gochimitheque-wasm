//go:build go1.24 && js && wasm

package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque/models"
)

func ScanQRdone(this js.Value, args []js.Value) interface{} {

	qr := args[0].String()
	BSTableQueryFilter.Clean()
	BSTableQueryFilter.Storage = qr

	storageCallbackWrapper := func(args ...interface{}) {
		Storage_listCallback(js.Null(), nil)
	}

	jsutils.CloseQR(js.Null(), nil)
	jsutils.LoadContent("div#content", "storage", fmt.Sprintf("%sv/storages", ApplicationProxyPath), storageCallbackWrapper)

	return nil

}

func SaveBorrowing(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		s                   *Storage
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#borrowing"), nil).Valid() {
		return nil
	}

	if len(args) > 0 {
		// When clicking on the "save" button
		// of the borrowing modal.

		row := args[2]
		globals.CurrentStorage = Storage{Storage: &models.Storage{}}.FromJsJSONValue(row)
		s = &globals.CurrentStorage

	} else {
		// When coming from Storage_operateEventsBorrow (unborrow).

		var storageId int

		s = &Storage{Storage: &models.Storage{}}
		s.Borrowing = &models.Borrowing{}

		// TODO: do not get value from dom
		if storageId, err = strconv.Atoi(jquery.Jq("input#bstorage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}

		var storageIdInt int64 = int64(storageId)
		s.StorageID = &storageIdInt

	}

	if jquery.Jq("textarea#borrowing_comment").GetVal().Truthy() {

		var _comment = jquery.Jq("textarea#borrowing_comment").GetVal().String()
		s.Borrowing.BorrowingComment = &_comment

	}

	select2Borrower := select2.NewSelect2(jquery.Jq("select#borrower"), nil)
	if len(select2Borrower.Select2Data()) > 0 {

		select2ItemBorrower := select2Borrower.Select2Data()[0]
		s.Borrowing.Borrower = models.Person{}

		var _person_id int
		if _person_id, err = strconv.Atoi(select2ItemBorrower.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		var _person_id_64 int64 = int64(_person_id)
		s.Borrowing.Borrower.PersonID = &_person_id_64
		s.Borrowing.Borrower.PersonEmail = select2ItemBorrower.Text

	}

	var borrowing_comment *string
	if s.Borrowing.BorrowingComment != nil {
		borrowing_comment = s.Borrowing.BorrowingComment
	}
	var borrower_id int64
	if s.Borrowing != nil {
		borrower_id = *s.Borrowing.Borrower.PersonID
	}

	ajaxURL = fmt.Sprintf("%sborrows/%d?borrower_id=%d", BackProxyPath, *s.StorageID, borrower_id)
	ajaxMethod = "get"

	if borrowing_comment != nil {
		ajaxURL += fmt.Sprintf("&borrowing_comment=%s", *borrowing_comment)
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Done: func(data js.Value) {

			jquery.Jq("#borrow").Object.Call("modal", "hide")

			jsutils.DisplaySuccessMessage(locales.Translate("storage_borrow_updated", HTTPHeaderAcceptLanguage))
			bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func SaveStorage(this js.Value, args []js.Value) any {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
		storageId           int
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#storage"), nil).Valid() {
		return nil
	}

	globals.CurrentStorage = Storage{Storage: &models.Storage{Person: models.Person{PersonID: &globals.ConnectedUserID}}}
	globals.CurrentStorage.Product = models.Product{ProductType: "chem", Name: &models.Name{}, Person: models.Person{PersonID: &globals.ConnectedUserID}} // Fields not used but required

	var _productID int
	if _productID, err = strconv.Atoi(jquery.Jq("input#product_id").GetVal().String()); err != nil {
		fmt.Println(err)
		return nil
	}
	var _product_id_64 int64 = int64(_productID)
	globals.CurrentStorage.Product.ProductID = &_product_id_64

	if jquery.Jq("input#storage_id").GetVal().Truthy() {
		if storageId, err = strconv.Atoi(jquery.Jq("input#storage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storageIdInt int64 = int64(storageId)
		globals.CurrentStorage.StorageID = &storageIdInt
	}

	// if jquery.Jq("input#storage_nbitem").GetVal().Truthy() {
	// 	if globals.CurrentStorage.StorageNbItem, err = strconv.Atoi(jquery.Jq("input#storage_nbitem").GetVal().String()); err != nil {
	// 		fmt.Println(err)
	// 		return nil
	// 	}
	// }
	// if jquery.Jq("input#storage_identicalbarecode:checked").Object.Length() > 0 {
	// 	globals.CurrentStorage.StorageIdenticalBarecode = true
	// }

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)
	if len(select2StoreLocation.Select2Data()) > 0 {
		select2ItemStorelocation := select2StoreLocation.Select2Data()[0]

		var storelocationId int
		if storelocationId, err = strconv.Atoi(select2ItemStorelocation.Id); err != nil {
			fmt.Println(err)
			return nil
		}

		var _store_location_id_64 int64 = int64(storelocationId)
		globals.CurrentStorage.StoreLocation.StoreLocationID = &_store_location_id_64

		var _store_location_name string = select2ItemStorelocation.Text
		globals.CurrentStorage.StoreLocation.StoreLocationName = _store_location_name
	}

	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", jquery.Jq("#storage_quantity").GetVal().String()))
	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", jquery.Jq("#storage_concentration").GetVal().String()))

	if jquery.Jq("input#storage_quantity").GetVal().Truthy() {
		var storageQuantity float64
		if storageQuantity, err = strconv.ParseFloat(jquery.Jq("input#storage_quantity").GetVal().String(), 64); err != nil {
			fmt.Println(err)
			return nil
		}
		var storageQuantityFloat float64 = float64(storageQuantity)
		globals.CurrentStorage.StorageQuantity = &storageQuantityFloat
	}

	select2UnitQuantity := select2.NewSelect2(jquery.Jq("select#unit_quantity"), nil)
	if len(select2UnitQuantity.Select2Data()) > 0 {
		select2ItemUnitQuantity := select2UnitQuantity.Select2Data()[0]
		label := ""
		id := int64(0)
		globals.CurrentStorage.UnitQuantity = &models.Unit{
			UnitID:    &id,
			UnitLabel: &label,
		}

		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitQuantity.Id); err != nil {
			fmt.Println(err)
			return nil
		}

		*globals.CurrentStorage.UnitQuantity.UnitID = int64(unitId)
		*globals.CurrentStorage.UnitQuantity.UnitLabel = select2ItemUnitQuantity.Text
	}

	if jquery.Jq("input#storage_concentration").GetVal().Truthy() {
		var storageConcentration int
		if storageConcentration, err = strconv.Atoi(jquery.Jq("input#storage_concentration").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storageConcentrationInt int64 = int64(storageConcentration)
		globals.CurrentStorage.StorageConcentration = &storageConcentrationInt
	}

	select2UnitConcentration := select2.NewSelect2(jquery.Jq("select#unit_concentration"), nil)
	if len(select2UnitConcentration.Select2Data()) > 0 {
		select2ItemUnitConcentration := select2UnitConcentration.Select2Data()[0]

		label := ""
		id := int64(0)
		globals.CurrentStorage.UnitConcentration = &models.Unit{
			UnitID:    &id,
			UnitLabel: &label,
		}

		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitConcentration.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		*globals.CurrentStorage.UnitConcentration.UnitID = int64(unitId)
		*globals.CurrentStorage.UnitConcentration.UnitLabel = select2ItemUnitConcentration.Text
	}

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)
	if len(select2Supplier.Select2Data()) > 0 {
		select2ItemSupplier := select2Supplier.Select2Data()[0]

		var supplierId int

		if select2ItemSupplier.IDIsDigit() {
			if supplierId, err = strconv.Atoi(select2ItemSupplier.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		supplierIdInt64 := int64(supplierId)
		globals.CurrentStorage.Supplier = &models.Supplier{}
		globals.CurrentStorage.Supplier.SupplierID = &supplierIdInt64
		globals.CurrentStorage.Supplier.SupplierLabel = &select2ItemSupplier.Text
	}

	if jquery.Jq("input#storage_entry_date").GetVal().Truthy() {
		var storageEntryDate time.Time
		if storageEntryDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_entry_date").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_entry_date time.Time = storageEntryDate
		globals.CurrentStorage.StorageEntryDate = &storage_entry_date
	}
	if jquery.Jq("input#storage_exit_date").GetVal().Truthy() {
		var storageExitDate time.Time
		if storageExitDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_exit_date").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_exit_date time.Time = storageExitDate
		globals.CurrentStorage.StorageExitDate = &storage_exit_date
	}
	if jquery.Jq("input#storage_opening_date").GetVal().Truthy() {
		var storageOpeningDate time.Time
		if storageOpeningDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_opening_date").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_opening_date time.Time = storageOpeningDate
		globals.CurrentStorage.StorageOpeningDate = &storage_opening_date
	}
	if jquery.Jq("input#storage_expiration_date").GetVal().Truthy() {
		var storageExpirationDate time.Time
		if storageExpirationDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_expiration_date").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_expiration_date time.Time = storageExpirationDate
		globals.CurrentStorage.StorageExpirationDate = &storage_expiration_date
	}

	if jquery.Jq("input#storage_reference").GetVal().Truthy() {
		var storage_reference string = jquery.Jq("input#storage_reference").GetVal().String()
		globals.CurrentStorage.StorageReference = &storage_reference
	}
	if jquery.Jq("input#storage_batch_number").GetVal().Truthy() {
		var storage_batch_number string = jquery.Jq("input#storage_batch_number").GetVal().String()
		globals.CurrentStorage.StorageBatchNumber = &storage_batch_number
	}
	if jquery.Jq("input#storage_barecode").GetVal().Truthy() {
		var storage_barecode string = jquery.Jq("input#storage_barecode").GetVal().String()
		globals.CurrentStorage.StorageBarecode = &storage_barecode
	}
	if jquery.Jq("input#storage_comment").GetVal().Truthy() {
		var storage_comment string = jquery.Jq("input#storage_comment").GetVal().String()
		globals.CurrentStorage.StorageComment = &storage_comment
	}

	if jquery.Jq("input#storage_number_of_bag").GetVal().Truthy() {
		var StorageNumberOfBag int
		if StorageNumberOfBag, err = strconv.Atoi(jquery.Jq("input#storage_number_of_bag").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_number_of_bag int64 = int64(StorageNumberOfBag)
		globals.CurrentStorage.StorageNumberOfBag = &storage_number_of_bag
	}
	if jquery.Jq("input#storage_number_of_carton").GetVal().Truthy() {
		var StorageNumberOfCarton int
		if StorageNumberOfCarton, err = strconv.Atoi(jquery.Jq("input#storage_number_of_carton").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		var storage_number_of_carton int64 = int64(StorageNumberOfCarton)
		globals.CurrentStorage.StorageNumberOfCarton = &storage_number_of_carton
	}

	if jquery.Jq("input#storage_to_destroy:checked").Object.Length() > 0 {
		globals.CurrentStorage.StorageToDestroy = true
	}

	var nb_items int
	if jquery.Jq("input#storage_nbitem").GetVal().Truthy() {
		if nb_items, err = strconv.Atoi(jquery.Jq("input#storage_nbitem").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}

	var identical_barecode bool
	if jquery.Jq("input#storage_identicalbarecode:checked").Object.Length() > 0 {
		identical_barecode = true
	}

	if (!jquery.Jq("form#storage input#storage_id").GetVal().IsUndefined()) && jquery.Jq("form#storage input#storage_id").GetVal().String() != "" {
		ajaxURL = fmt.Sprintf("%sstorages/%d", BackProxyPath, *globals.CurrentStorage.StorageID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sstorages?nb_items=%d&identical_barecode=%t", BackProxyPath, nb_items, identical_barecode)
		ajaxMethod = "post"
	}

	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", *globals.CurrentStorage.Storage))

	if dataBytes, err = json.Marshal(globals.CurrentStorage); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			var (
				err error
			)

			if err = json.Unmarshal([]byte(data.String()), &CurrentStorageIds); err != nil {
				jsutils.DisplayGenericErrorMessage()
				fmt.Println(err)
			} else {
				href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
				jsutils.ClearSearchExceptProduct(js.Null(), nil)
				jsutils.LoadContent("div#content", "storage", href, Storage_SaveCallback)
			}

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func FillInStorageForm(s Storage, id string) {

	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", s.Storage))
	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", s.Storage.StorageID))
	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", id))

	if s.StorageID != nil {
		jquery.Jq(fmt.Sprintf("#%s #storage_id", id)).SetVal(*s.StorageID)
	}

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)
	select2StoreLocation.Select2Clear()
	if s.StoreLocation.StoreLocationID != nil {
		select2StoreLocation.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.StoreLocation.StoreLocationName,
				Value:           strconv.Itoa(int(*s.StoreLocation.StoreLocationID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("input#storage_quantity").SetVal("")
	if s.StorageQuantity != nil {
		jquery.Jq("input#storage_quantity").SetVal(*s.StorageQuantity)
	}

	select2UnitQuantity := select2.NewSelect2(jquery.Jq("select#unit_quantity"), nil)
	select2UnitQuantity.Select2Clear()
	if s.UnitQuantity != nil && s.UnitQuantity.UnitID != nil {
		select2UnitQuantity.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *s.UnitQuantity.UnitLabel,
				Value:           strconv.Itoa(int(*s.UnitQuantity.UnitID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("input#storage_concentration").SetVal("")
	if s.StorageConcentration != nil {
		jquery.Jq("input#storage_concentration").SetVal(*s.StorageConcentration)
	}

	select2UnitConcentration := select2.NewSelect2(jquery.Jq("select#unit_concentration"), nil)
	select2UnitConcentration.Select2Clear()
	if s.UnitConcentration != nil && s.UnitConcentration.UnitID != nil {
		select2UnitConcentration.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *s.UnitConcentration.UnitLabel,
				Value:           strconv.Itoa(int(*s.UnitConcentration.UnitID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)
	select2Supplier.Select2Clear()
	if s.Supplier != nil && s.Supplier.SupplierID != nil {
		select2Supplier.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            *s.Supplier.SupplierLabel,
				Value:           strconv.Itoa(int(*s.Supplier.SupplierID)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("input#storage_number_of_bag").SetVal("")
	if s.StorageNumberOfBag != nil {
		jquery.Jq("input#storage_number_of_bag").SetVal(*s.StorageNumberOfBag)
	}
	jquery.Jq("input#storage_number_of_carton").SetVal("")
	if s.StorageNumberOfCarton != nil {
		jquery.Jq("input#storage_number_of_carton").SetVal(*s.StorageNumberOfCarton)
	}

	jquery.Jq("input#storage_entry_date").SetVal("")
	if s.StorageEntryDate != nil {
		jquery.Jq("input#storage_entry_date").SetVal(s.StorageEntryDate.Format("2006-01-02"))
	}
	jquery.Jq("input#storage_exit_date").SetVal("")
	if s.StorageExitDate != nil {
		jquery.Jq("input#storage_exit_date").SetVal(s.StorageExitDate.Format("2006-01-02"))
	}
	jquery.Jq("#storage_opening_date").SetVal("")
	if s.StorageOpeningDate != nil {
		jquery.Jq("#storage_opening_date").SetVal(s.StorageOpeningDate.Format("2006-01-02"))
	}
	jquery.Jq("input#storage_expiration_date").SetVal("")
	if s.StorageExpirationDate != nil {
		jquery.Jq("input#storage_expiration_date").SetVal(s.StorageExpirationDate.Format("2006-01-02"))
	}

	jquery.Jq("input#storage_reference").SetVal("")
	if s.StorageReference != nil {
		jquery.Jq("input#storage_reference").SetVal(*s.StorageReference)
	}
	jquery.Jq("input#storage_batch_number").SetVal("")
	if s.StorageBatchNumber != nil {
		jquery.Jq("input#storage_batch_number").SetVal(*s.StorageBatchNumber)
	}
	jquery.Jq("input#storage_barecode").SetVal("")
	if s.StorageBarecode != nil {
		jquery.Jq("input#storage_barecode").SetVal(*s.StorageBarecode)
	}
	jquery.Jq("input#storage_comment").SetVal("")
	if s.StorageComment != nil {
		jquery.Jq("input#storage_comment").SetVal(*s.StorageComment)
	}

	jquery.Jq("input#storage_todestroy").SetProp("checked", false)
	if s.StorageToDestroy {
		jquery.Jq("input#storage_todestroy").SetProp("checked", "checked")
	}

}
