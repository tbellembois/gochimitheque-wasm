package storage

import (
	"database/sql"
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
)

func Storage_closeQR(this js.Value, args []js.Value) interface{} {

	js.Global().Get("window").Get("qrScanner").Call("destroy")
	js.Global().Get("window").Set("qrScanner", nil)

	jquery.Jq("#video").AddClass("invisible")

	return nil

}

func ScanQRdone(this js.Value, args []js.Value) interface{} {

	qr := args[0].String()
	BSTableQueryFilter.Clean()
	BSTableQueryFilter.Storage = qr

	storageCallbackWrapper := func(args ...interface{}) {
		Storage_listCallback(js.Null(), nil)
	}

	Storage_closeQR(js.Null(), nil)
	jsutils.LoadContent("div#content", "storage", fmt.Sprintf("%sv/storages", ApplicationProxyPath), storageCallbackWrapper)

	return nil

}

func SaveBorrowing(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
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
		globals.CurrentStorage = Storage{}.FromJsJSONValue(row)
		s = &globals.CurrentStorage

	} else {
		// When coming from Storage_operateEventsBorrow (unborrow).

		var storageId int

		s = &Storage{}
		s.Borrowing = &Borrowing{}

		// TODO: do not get value from dom
		if storageId, err = strconv.Atoi(jquery.Jq("input#bstorage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}

		s.StorageID = sql.NullInt64{
			Valid: true,
			Int64: int64(storageId),
		}

	}

	if jquery.Jq("textarea#borrowing_comment").GetVal().Truthy() {

		s.Borrowing.BorrowingComment = sql.NullString{
			Valid:  true,
			String: jquery.Jq("textarea#borrowing_comment").GetVal().String(),
		}

	}

	select2Borrower := select2.NewSelect2(jquery.Jq("select#borrower"), nil)
	if len(select2Borrower.Select2Data()) > 0 {

		select2ItemBorrower := select2Borrower.Select2Data()[0]
		s.Borrowing.Borrower = &Person{}
		if s.Borrowing.Borrower.PersonId, err = strconv.Atoi(select2ItemBorrower.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		s.Borrowing.Borrower.PersonEmail = select2ItemBorrower.Text

	}

	ajaxURL = fmt.Sprintf("%sborrowings", ApplicationProxyPath)
	ajaxMethod = "put"

	if dataBytes, err = json.Marshal(s); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
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

func SaveStorage(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
		storageId           int
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#storage"), nil).Valid() {
		return nil
	}

	globals.CurrentStorage = Storage{}
	globals.CurrentStorage.Product = Product{}

	if globals.CurrentStorage.Product.ProductID, err = strconv.Atoi(jquery.Jq("input#product_id").GetVal().String()); err != nil {
		fmt.Println(err)
		return nil
	}

	if jquery.Jq("input#storage_id").GetVal().Truthy() {
		if storageId, err = strconv.Atoi(jquery.Jq("input#storage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageID = sql.NullInt64{
			Valid: true,
			Int64: int64(storageId),
		}
	}

	if jquery.Jq("input#storage_nbitem").GetVal().Truthy() {
		if globals.CurrentStorage.StorageNbItem, err = strconv.Atoi(jquery.Jq("input#storage_nbitem").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	if jquery.Jq("input#storage_identicalbarecode:checked").Object.Length() > 0 {
		globals.CurrentStorage.StorageIdenticalBarecode = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#storelocation"), nil)
	if len(select2StoreLocation.Select2Data()) > 0 {
		select2ItemStorelocation := select2StoreLocation.Select2Data()[0]
		globals.CurrentStorage.StoreLocation = StoreLocation{}
		var storelocationId int
		if storelocationId, err = strconv.Atoi(select2ItemStorelocation.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StoreLocation.StoreLocationID = sql.NullInt64{
			Int64: int64(storelocationId),
			Valid: true,
		}
		globals.CurrentStorage.StoreLocation.StoreLocationName = sql.NullString{
			String: select2ItemStorelocation.Text,
			Valid:  true,
		}
	}

	if jquery.Jq("input#storage_quantity").GetVal().Truthy() {
		var storageQuantity float64
		if storageQuantity, err = strconv.ParseFloat(jquery.Jq("#storage_quantity").GetVal().String(), 64); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageQuantity = sql.NullFloat64{
			Valid:   true,
			Float64: storageQuantity,
		}
	}

	select2UnitQuantity := select2.NewSelect2(jquery.Jq("select#unit_quantity"), nil)
	if len(select2UnitQuantity.Select2Data()) > 0 {
		select2ItemUnitQuantity := select2UnitQuantity.Select2Data()[0]
		globals.CurrentStorage.UnitQuantity = Unit{}
		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitQuantity.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.UnitQuantity.UnitID = sql.NullInt64{
			Int64: int64(unitId),
			Valid: true,
		}
		globals.CurrentStorage.UnitQuantity.UnitLabel = sql.NullString{
			String: select2ItemUnitQuantity.Text,
			Valid:  true,
		}
	}

	if jquery.Jq("input#storage_concentration").GetVal().Truthy() {
		var storageConcentration int
		if storageConcentration, err = strconv.Atoi(jquery.Jq("#storage_concentration").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageConcentration = sql.NullInt64{
			Valid: true,
			Int64: int64(storageConcentration),
		}
	}

	select2UnitConcentration := select2.NewSelect2(jquery.Jq("select#unit_concentration"), nil)
	if len(select2UnitConcentration.Select2Data()) > 0 {
		select2ItemUnitConcentration := select2UnitConcentration.Select2Data()[0]
		globals.CurrentStorage.UnitConcentration = Unit{}
		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitConcentration.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.UnitConcentration.UnitID = sql.NullInt64{
			Int64: int64(unitId),
			Valid: true,
		}
		globals.CurrentStorage.UnitConcentration.UnitLabel = sql.NullString{
			String: select2ItemUnitConcentration.Text,
			Valid:  true,
		}
	}

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)
	if len(select2Supplier.Select2Data()) > 0 {
		select2ItemSupplier := select2Supplier.Select2Data()[0]
		globals.CurrentStorage.Supplier = Supplier{}
		var supplierId = -1

		if select2ItemSupplier.IDIsDigit() {
			if supplierId, err = strconv.Atoi(select2ItemSupplier.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		globals.CurrentStorage.Supplier.SupplierID = sql.NullInt64{
			Int64: int64(supplierId),
			Valid: true,
		}
		globals.CurrentStorage.Supplier.SupplierLabel = sql.NullString{
			String: select2ItemSupplier.Text,
			Valid:  true,
		}
	}

	if jquery.Jq("input#storage_entrydate").GetVal().Truthy() {
		var storageEntryDate time.Time
		if storageEntryDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_entrydate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageEntryDate = sql.NullTime{
			Valid: true,
			Time:  storageEntryDate,
		}
	}
	if jquery.Jq("input#storage_exitdate").GetVal().Truthy() {
		var storageExitDate time.Time
		if storageExitDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_exitdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageExitDate = sql.NullTime{
			Valid: true,
			Time:  storageExitDate,
		}
	}
	if jquery.Jq("input#storage_openingdate").GetVal().Truthy() {
		var storageOpeningDate time.Time
		if storageOpeningDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_openingdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageOpeningDate = sql.NullTime{
			Valid: true,
			Time:  storageOpeningDate,
		}
	}
	if jquery.Jq("input#storage_expirationdate").GetVal().Truthy() {
		var storageExpirationDate time.Time
		if storageExpirationDate, err = time.Parse("2006-01-02", jquery.Jq("#storage_expirationdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		globals.CurrentStorage.StorageExpirationDate = sql.NullTime{
			Valid: true,
			Time:  storageExpirationDate,
		}
	}

	if jquery.Jq("input#storage_reference").GetVal().Truthy() {
		globals.CurrentStorage.StorageReference = sql.NullString{
			Valid:  true,
			String: jquery.Jq("input#storage_reference").GetVal().String(),
		}
	}
	if jquery.Jq("input#storage_batchnumber").GetVal().Truthy() {
		globals.CurrentStorage.StorageBatchNumber = sql.NullString{
			Valid:  true,
			String: jquery.Jq("input#storage_batchnumber").GetVal().String(),
		}
	}
	if jquery.Jq("input#storage_barecode").GetVal().Truthy() {
		globals.CurrentStorage.StorageBarecode = sql.NullString{
			Valid:  true,
			String: jquery.Jq("input#storage_barecode").GetVal().String(),
		}
	}
	if jquery.Jq("input#storage_comment").GetVal().Truthy() {
		globals.CurrentStorage.StorageComment = sql.NullString{
			Valid:  true,
			String: jquery.Jq("input#storage_comment").GetVal().String(),
		}
	}

	if jquery.Jq("input#storage_todestroy:checked").Object.Length() > 0 {
		globals.CurrentStorage.StorageToDestroy = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if (!jquery.Jq("form#storage input#storage_id").GetVal().IsUndefined()) && jquery.Jq("form#storage input#storage_id").GetVal().String() != "" {
		fmt.Println("update")
		ajaxURL = fmt.Sprintf("%sstorages/%d", ApplicationProxyPath, globals.CurrentStorage.StorageID.Int64)
		ajaxMethod = "put"
	} else {
		fmt.Println("create")
		ajaxURL = fmt.Sprintf("%sstorages", ApplicationProxyPath)
		ajaxMethod = "post"
	}

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

			if err = json.Unmarshal([]byte(data.String()), &CurrentStorages); err != nil {
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

	jquery.Jq(fmt.Sprintf("#%s #storage_id", id)).SetVal(s.StorageID.Int64)

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#storelocation"), nil)
	select2StoreLocation.Select2Clear()
	if s.StoreLocation.StoreLocationID.Valid {
		select2StoreLocation.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.StoreLocation.StoreLocationName.String,
				Value:           strconv.Itoa(int(s.StoreLocation.StoreLocationID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#storage_quantity").SetVal("")
	if s.StorageQuantity.Valid {
		jquery.Jq("#storage_quantity").SetVal(s.StorageQuantity.Float64)
	}

	select2UnitQuantity := select2.NewSelect2(jquery.Jq("select#unit_quantity"), nil)
	select2UnitQuantity.Select2Clear()
	if s.UnitQuantity.UnitID.Valid {
		select2UnitQuantity.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.UnitQuantity.UnitLabel.String,
				Value:           strconv.Itoa(int(s.UnitQuantity.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#storage_concentration").SetVal("")
	if s.StorageConcentration.Valid {
		jquery.Jq("#storage_concentration").SetVal(s.StorageConcentration.Int64)
	}

	select2UnitConcentration := select2.NewSelect2(jquery.Jq("select#unit_concentration"), nil)
	select2UnitConcentration.Select2Clear()
	if s.UnitConcentration.UnitID.Valid {
		select2UnitConcentration.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.UnitConcentration.UnitLabel.String,
				Value:           strconv.Itoa(int(s.UnitConcentration.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)
	select2Supplier.Select2Clear()
	if s.Supplier.SupplierID.Valid {
		select2Supplier.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.Supplier.SupplierLabel.String,
				Value:           strconv.Itoa(int(s.Supplier.SupplierID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	jquery.Jq("#storage_entrydate").SetVal("")
	if s.StorageEntryDate.Valid {
		jquery.Jq("#storage_entrydate").SetVal(s.StorageEntryDate.Time.Format("2006-01-02"))
	}
	jquery.Jq("#storage_exitdate").SetVal("")
	if s.StorageExitDate.Valid {
		jquery.Jq("#storage_exitdate").SetVal(s.StorageExitDate.Time.Format("2006-01-02"))
	}
	jquery.Jq("#storage_openingdate").SetVal("")
	if s.StorageOpeningDate.Valid {
		jquery.Jq("#storage_openingdate").SetVal(s.StorageOpeningDate.Time.Format("2006-01-02"))
	}
	jquery.Jq("#storage_expirationdate").SetVal("")
	if s.StorageExpirationDate.Valid {
		jquery.Jq("#storage_expirationdate").SetVal(s.StorageExpirationDate.Time.Format("2006-01-02"))
	}

	jquery.Jq("#storage_reference").SetVal("")
	if s.StorageReference.Valid {
		jquery.Jq("#storage_reference").SetVal(s.StorageReference.String)
	}
	jquery.Jq("#storage_batchnumber").SetVal("")
	if s.StorageBatchNumber.Valid {
		jquery.Jq("#storage_batchnumber").SetVal(s.StorageBatchNumber.String)
	}
	jquery.Jq("#storage_barecode").SetVal("")
	if s.StorageBarecode.Valid {
		jquery.Jq("#storage_barecode").SetVal(s.StorageBarecode.String)
	}
	jquery.Jq("#storage_comment").SetVal("")
	if s.StorageComment.Valid {
		jquery.Jq("#storage_comment").SetVal(s.StorageComment.String)
	}

	jquery.Jq("#storage_todestroy").SetProp("chacked", false)
	if s.StorageToDestroy.Valid && s.StorageToDestroy.Bool {
		jquery.Jq("#storage_todestroy").SetProp("checked", "checked")
	}

}
