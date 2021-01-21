package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func Storage_closeQR(this js.Value, args []js.Value) interface{} {

	Win.Get("qrScanner").Call("destroy")
	Win.Set("qrScanner", nil)

	Jq("#video").AddClass("invisible")

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
	utils.LoadContent("storage", fmt.Sprintf("%sv/storages", ApplicationProxyPath), storageCallbackWrapper)

	return nil

}

func SaveBorrowing(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
		s                   *Storage
		err                 error
	)

	if !(Jq("form#borrowing").Valid()) {
		return nil
	}

	if len(args) > 0 {
		// When clicking on the "save" button
		// of the borrowing modal.

		row := args[2]
		types.CurrentStorage = Storage{}.FromJsJSONValue(row)
		s = &types.CurrentStorage

	} else {
		// When coming from Storage_operateEventsBorrow (unborrow).

		var storageId int

		s = &Storage{}
		s.Borrowing = &Borrowing{}

		// TODO: do not get value from dom
		if storageId, err = strconv.Atoi(Jq("input#bstorage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}

		s.StorageID = sql.NullInt64{
			Valid: true,
			Int64: int64(storageId),
		}

	}

	if Jq("textarea#borrowing_comment").GetVal().Truthy() {

		s.Borrowing.BorrowingComment = sql.NullString{
			Valid:  true,
			String: Jq("textarea#borrowing_comment").GetVal().String(),
		}

	}

	if len(Jq("select#borrower").Select2Data()) > 0 {

		select2ItemBorrower := Jq("select#borrower").Select2Data()[0]
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

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			Jq("#borrow").Object.Call("modal", "hide")

			utils.DisplaySuccessMessage(locales.Translate("storage_borrow_updated", HTTPHeaderAcceptLanguage))
			Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

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

	if !(Jq("form#storage").Valid()) {
		return nil
	}

	types.CurrentStorage = Storage{}
	types.CurrentStorage.Product = Product{}

	if types.CurrentStorage.Product.ProductID, err = strconv.Atoi(Jq("input#product_id").GetVal().String()); err != nil {
		fmt.Println(err)
		return nil
	}

	if Jq("input#storage_id").GetVal().Truthy() {
		if storageId, err = strconv.Atoi(Jq("input#storage_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageID = sql.NullInt64{
			Valid: true,
			Int64: int64(storageId),
		}
	}

	if Jq("input#storage_nbitem").GetVal().Truthy() {
		if types.CurrentStorage.StorageNbItem, err = strconv.Atoi(Jq("input#storage_nbitem").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	if Jq("input#storage_identicalbarecode:checked").Object.Length() > 0 {
		types.CurrentStorage.StorageIdenticalBarecode = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if len(Jq("select#storelocation").Select2Data()) > 0 {
		select2ItemStorelocation := Jq("select#storelocation").Select2Data()[0]
		types.CurrentStorage.StoreLocation = StoreLocation{}
		var storelocationId int
		if storelocationId, err = strconv.Atoi(select2ItemStorelocation.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StoreLocation.StoreLocationID = sql.NullInt64{
			Int64: int64(storelocationId),
			Valid: true,
		}
		types.CurrentStorage.StoreLocation.StoreLocationName = sql.NullString{
			String: select2ItemStorelocation.Text,
			Valid:  true,
		}
	}

	if Jq("input#storage_quantity").GetVal().Truthy() {
		var storageQuantity float64
		if storageQuantity, err = strconv.ParseFloat(Jq("#storage_quantity").GetVal().String(), 64); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageQuantity = sql.NullFloat64{
			Valid:   true,
			Float64: storageQuantity,
		}
	}

	if len(Jq("select#unit_quantity").Select2Data()) > 0 {
		select2ItemUnitQuantity := Jq("select#unit_quantity").Select2Data()[0]
		types.CurrentStorage.UnitQuantity = Unit{}
		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitQuantity.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.UnitQuantity.UnitID = sql.NullInt64{
			Int64: int64(unitId),
			Valid: true,
		}
		types.CurrentStorage.UnitQuantity.UnitLabel = sql.NullString{
			String: select2ItemUnitQuantity.Text,
			Valid:  true,
		}
	}

	if Jq("input#storage_concentration").GetVal().Truthy() {
		var storageConcentration int
		if storageConcentration, err = strconv.Atoi(Jq("#storage_concentration").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageConcentration = sql.NullInt64{
			Valid: true,
			Int64: int64(storageConcentration),
		}
	}

	if len(Jq("select#unit_concentration").Select2Data()) > 0 {
		select2ItemUnitConcentration := Jq("select#unit_concentration").Select2Data()[0]
		types.CurrentStorage.UnitConcentration = Unit{}
		var unitId int
		if unitId, err = strconv.Atoi(select2ItemUnitConcentration.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.UnitConcentration.UnitID = sql.NullInt64{
			Int64: int64(unitId),
			Valid: true,
		}
		types.CurrentStorage.UnitConcentration.UnitLabel = sql.NullString{
			String: select2ItemUnitConcentration.Text,
			Valid:  true,
		}
	}

	if len(Jq("select#supplier").Select2Data()) > 0 {
		select2ItemSupplier := Jq("select#supplier").Select2Data()[0]
		types.CurrentStorage.Supplier = Supplier{}
		var supplierId = -1

		if select2ItemSupplier.IDIsDigit() {
			if supplierId, err = strconv.Atoi(select2ItemSupplier.Id); err != nil {
				fmt.Println(err)
				return nil
			}
		}

		types.CurrentStorage.Supplier.SupplierID = sql.NullInt64{
			Int64: int64(supplierId),
			Valid: true,
		}
		types.CurrentStorage.Supplier.SupplierLabel = sql.NullString{
			String: select2ItemSupplier.Text,
			Valid:  true,
		}
	}

	if Jq("input#storage_entrydate").GetVal().Truthy() {
		var storageEntryDate time.Time
		if storageEntryDate, err = time.Parse("2006-01-02", Jq("#storage_entrydate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageEntryDate = sql.NullTime{
			Valid: true,
			Time:  storageEntryDate,
		}
	}
	if Jq("input#storage_exitdate").GetVal().Truthy() {
		var storageExitDate time.Time
		if storageExitDate, err = time.Parse("2006-01-02", Jq("#storage_exitdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageExitDate = sql.NullTime{
			Valid: true,
			Time:  storageExitDate,
		}
	}
	if Jq("input#storage_openingdate").GetVal().Truthy() {
		var storageOpeningDate time.Time
		if storageOpeningDate, err = time.Parse("2006-01-02", Jq("#storage_openingdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageOpeningDate = sql.NullTime{
			Valid: true,
			Time:  storageOpeningDate,
		}
	}
	if Jq("input#storage_expirationdate").GetVal().Truthy() {
		var storageExpirationDate time.Time
		if storageExpirationDate, err = time.Parse("2006-01-02", Jq("#storage_expirationdate").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		types.CurrentStorage.StorageExpirationDate = sql.NullTime{
			Valid: true,
			Time:  storageExpirationDate,
		}
	}

	if Jq("input#storage_reference").GetVal().Truthy() {
		types.CurrentStorage.StorageReference = sql.NullString{
			Valid:  true,
			String: Jq("input#storage_reference").GetVal().String(),
		}
	}
	if Jq("input#storage_batchnumber").GetVal().Truthy() {
		types.CurrentStorage.StorageBatchNumber = sql.NullString{
			Valid:  true,
			String: Jq("input#storage_batchnumber").GetVal().String(),
		}
	}
	if Jq("input#storage_barecode").GetVal().Truthy() {
		types.CurrentStorage.StorageBarecode = sql.NullString{
			Valid:  true,
			String: Jq("input#storage_barecode").GetVal().String(),
		}
	}
	if Jq("input#storage_comment").GetVal().Truthy() {
		types.CurrentStorage.StorageComment = sql.NullString{
			Valid:  true,
			String: Jq("input#storage_comment").GetVal().String(),
		}
	}

	if Jq("input#storage_todestroy:checked").Object.Length() > 0 {
		types.CurrentStorage.StorageToDestroy = sql.NullBool{
			Bool:  true,
			Valid: true,
		}
	}

	if (!Jq("form#storage input#storage_id").GetVal().IsUndefined()) && Jq("form#storage input#storage_id").GetVal().String() != "" {
		fmt.Println("update")
		ajaxURL = fmt.Sprintf("%sstorages/%d", ApplicationProxyPath, types.CurrentStorage.StorageID.Int64)
		ajaxMethod = "put"
	} else {
		fmt.Println("create")
		ajaxURL = fmt.Sprintf("%sstorages", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	if dataBytes, err = json.Marshal(types.CurrentStorage); err != nil {
		fmt.Println(err)
		return nil
	}

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			var (
				err error
			)

			if err = json.Unmarshal([]byte(data.String()), &CurrentStorage); err != nil {
				utils.DisplayGenericErrorMessage()
				fmt.Println(err)
			} else {
				href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
				search.ClearSearchExceptProduct(js.Null(), nil)
				utils.LoadContent("storage", href, Storage_SaveCallback, int(types.CurrentStorage.StorageID.Int64))
			}

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func FillInStorageForm(s Storage, id string) {

	Jq(fmt.Sprintf("#%s #storage_id", id)).SetVal(s.StorageID.Int64)

	Jq("select#storelocation").Select2Clear()
	if s.StoreLocation.StoreLocationID.Valid {
		Jq("select#storelocation").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.StoreLocation.StoreLocationName.String,
				Value:           strconv.Itoa(int(s.StoreLocation.StoreLocationID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#storage_quantity").SetVal("")
	if s.StorageQuantity.Valid {
		Jq("#storage_quantity").SetVal(s.StorageQuantity.Float64)
	}

	Jq("select#unit_quantity").Select2Clear()
	if s.UnitQuantity.UnitID.Valid {
		Jq("select#unit_quantity").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.UnitQuantity.UnitLabel.String,
				Value:           strconv.Itoa(int(s.UnitQuantity.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#storage_concentration").SetVal("")
	if s.StorageConcentration.Valid {
		Jq("#storage_concentration").SetVal(s.StorageConcentration.Int64)
	}

	Jq("select#unit_concentration").Select2Clear()
	if s.UnitConcentration.UnitID.Valid {
		Jq("select#unit_concentration").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.UnitConcentration.UnitLabel.String,
				Value:           strconv.Itoa(int(s.UnitConcentration.UnitID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("select#supplier").Select2Clear()
	if s.Supplier.SupplierID.Valid {
		Jq("select#supplier").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            s.Supplier.SupplierLabel.String,
				Value:           strconv.Itoa(int(s.Supplier.SupplierID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	Jq("#storage_entrydate").SetVal("")
	if s.StorageEntryDate.Valid {
		Jq("#storage_entrydate").SetVal(s.StorageEntryDate.Time.Format("2006-01-02"))
	}
	Jq("#storage_exitdate").SetVal("")
	if s.StorageExitDate.Valid {
		Jq("#storage_exitdate").SetVal(s.StorageExitDate.Time.Format("2006-01-02"))
	}
	Jq("#storage_openingdate").SetVal("")
	if s.StorageOpeningDate.Valid {
		Jq("#storage_openingdate").SetVal(s.StorageOpeningDate.Time.Format("2006-01-02"))
	}
	Jq("#storage_expirationdate").SetVal("")
	if s.StorageExpirationDate.Valid {
		Jq("#storage_expirationdate").SetVal(s.StorageExpirationDate.Time.Format("2006-01-02"))
	}

	Jq("#storage_reference").SetVal("")
	if s.StorageReference.Valid {
		Jq("#storage_reference").SetVal(s.StorageReference.String)
	}
	Jq("#storage_batchnumber").SetVal("")
	if s.StorageBatchNumber.Valid {
		Jq("#storage_batchnumber").SetVal(s.StorageBatchNumber.String)
	}
	Jq("#storage_barecode").SetVal("")
	if s.StorageBarecode.Valid {
		Jq("#storage_barecode").SetVal(s.StorageBarecode.String)
	}
	Jq("#storage_comment").SetVal("")
	if s.StorageComment.Valid {
		Jq("#storage_comment").SetVal(s.StorageComment.String)
	}

	Jq("#storage_todestroy").SetProp("chacked", false)
	if s.StorageToDestroy.Valid && s.StorageToDestroy.Bool {
		Jq("#storage_todestroy").SetProp("checked", "checked")
	}

}
