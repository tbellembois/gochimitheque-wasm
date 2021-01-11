package storage

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "image/png"
	"net/url"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"honnef.co/go/js/dom/v2"
)

func Storage_operateEventsRestore(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	Jq(fmt.Sprintf("button#restore%d", types.CurrentStorage.StorageID.Int64)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sstorages/%d/r", ApplicationProxyPath, types.CurrentStorage.StorageID.Int64)
		method := "put"

		done := func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("storage_restored_message", HTTPHeaderAcceptLanguage))
			BSTableQueryFilter.Lock()
			BSTableQueryFilter.QueryFilter.StorageArchive = false
			BSTableQueryFilter.QueryFilter.Storage = strconv.Itoa(int(types.CurrentStorage.StorageID.Int64))
			BSTableQueryFilter.QueryFilter.StorageFilterLabel = fmt.Sprintf("#%d", CurrentStorage.StorageID.Int64)
			Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

		}
		fail := func(data js.Value) {

			utils.DisplayGenericErrorMessage()

		}

		Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

		return nil

	}))

	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_CONFIRM, ""),
		Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
	})
	Jq(fmt.Sprintf("button#restore%d", types.CurrentStorage.StorageID.Int64)).SetHtml("")
	Jq(fmt.Sprintf("button#restore%d", types.CurrentStorage.StorageID.Int64)).Append(buttonTitle.OuterHTML())

	return nil

}

func Storage_operateEventsClone(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	types.CurrentStorage.StorageID = sql.NullInt64{
		Valid: false,
		Int64: 0,
	}

	href := fmt.Sprintf("%svc/storages", ApplicationProxyPath)
	utils.LoadContent("storage", href, Storage_createCallback, types.CurrentStorage)

	return nil

}

func Storage_operateEventsHistory(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.StorageHistory = true
	BSTableQueryFilter.QueryFilter.Storage = strconv.Itoa(int(types.CurrentStorage.StorageID.Int64))
	BSTableQueryFilter.QueryFilter.StorageFilterLabel = fmt.Sprintf("#%d", CurrentStorage.StorageID.Int64)
	Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

	return nil

}

func Storage_operateEventsBorrow(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	Jq("input#bstorage_id").SetVal(types.CurrentStorage.StorageID.Int64)

	if types.CurrentStorage.Borrowing.BorrowingID.Valid {

		// The storage has a borrowing.

		// Unborrow.
		SaveBorrowing(this, args)

	} else {

		// The storage does not have a borrowing.
		// Displaying the modal.
		Jq("#borrow").Object.Call("modal", "show")

		// Selecting the connected user.
		Jq("select#borrower").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            ConnectedUserEmail,
				Value:           strconv.Itoa(ConnectedUserID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())

	}

	return nil

}

func Storage_operateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	href := fmt.Sprintf("%svc/storages", ApplicationProxyPath)
	utils.LoadContent("storage", href, Storage_createCallback, types.CurrentStorage)

	return nil

}

func Storage_operateEventsArchive(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	Jq(fmt.Sprintf("button#archive%d", types.CurrentStorage.StorageID.Int64)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sstorages/%d/a", ApplicationProxyPath, types.CurrentStorage.StorageID.Int64)
		method := "delete"

		done := func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("storage_trashed_message", HTTPHeaderAcceptLanguage))
			Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

		}
		fail := func(data js.Value) {

			utils.DisplayGenericErrorMessage()

		}

		Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

		return nil

	}))

	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_CONFIRM, ""),
		Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
	})
	Jq(fmt.Sprintf("button#archive%d", types.CurrentStorage.StorageID.Int64)).SetHtml("")
	Jq(fmt.Sprintf("button#archive%d", types.CurrentStorage.StorageID.Int64)).Append(buttonTitle.OuterHTML())

	return nil

}

func Storage_operateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	Jq(fmt.Sprintf("button#delete%d", types.CurrentStorage.StorageID.Int64)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sstorages/%d", ApplicationProxyPath, types.CurrentStorage.StorageID.Int64)
		method := "delete"

		done := func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("storage_deleted_message", HTTPHeaderAcceptLanguage))
			BSTableQueryFilter.Lock()
			BSTableQueryFilter.QueryFilter.Storage = ""
			BSTableQueryFilter.QueryFilter.StorageFilterLabel = ""
			Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

		}
		fail := func(data js.Value) {

			utils.DisplayGenericErrorMessage()

		}

		Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

		return nil

	}))

	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_CONFIRM, ""),
		Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
	})
	Jq(fmt.Sprintf("button#delete%d", types.CurrentStorage.StorageID.Int64)).SetHtml("")
	Jq(fmt.Sprintf("button#delete%d", types.CurrentStorage.StorageID.Int64)).Append(buttonTitle.OuterHTML())

	return nil

}

func ShowIfAuthorizedActionButtons(this js.Value, args []js.Value) interface{} {

	// Iterating other the button with the class "borrow"
	// (we could choose "clone" or "delete")
	// to retrieve once the storage id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
	for _, button := range buttons {
		if button.Class().Contains("borrow") || button.Class().Contains("restore") {
			storageId := button.GetAttribute("sid")

			Jq("#history" + storageId).FadeIn()

			utils.HasPermission("storages", storageId, "put", func() {
				Jq("#edit" + storageId).FadeIn()
				Jq("#clone" + storageId).FadeIn()
				Jq("#borrow" + storageId).FadeIn()
			}, func() {
			})

			utils.HasPermission("storages", storageId, "delete", func() {
				Jq("#delete" + storageId).FadeIn()
				Jq("#archive" + storageId).FadeIn()
				Jq("#restore" + storageId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}

func Storage_productFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	for _, chunk := range utils.Chunks(types.CurrentStorage.Product.Name.NameLabel, 40) {
		d.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: chunk,
			}),
		)
		d.AppendChild(
			widgets.NewBr(widgets.BrAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
			}),
		)
	}

	return d.OuterHTML()

}

func Storage_storelocationFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	iconColor := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Attributes: map[string]string{
				"style": fmt.Sprintf("color: %s", types.CurrentStorage.StoreLocation.StoreLocationColor.String),
			},
		},
		Icon: themes.NewMdiIcon(themes.MDI_COLOR, themes.MDI_24PX),
	})

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: types.CurrentStorage.StoreLocation.StoreLocationFullPath,
	})

	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})
	d.AppendChild(iconColor)
	d.AppendChild(spanLabel)

	return d.OuterHTML()

}

func Storage_quantityFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	return fmt.Sprintf("%v %s", types.CurrentStorage.StorageQuantity.Float64, types.CurrentStorage.UnitQuantity.UnitLabel.String)

}

func Storage_barecodeFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	for _, chunk := range utils.Chunks(types.CurrentStorage.StorageBarecode.String, 10) {
		d.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: chunk,
			}),
		)
		d.AppendChild(
			widgets.NewBr(widgets.BrAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
			}),
		)
	}

	return d.OuterHTML()

}

func Storage_operateFormatter(this js.Value, args []js.Value) interface{} {

	var (
		buttonClone    string
		buttonRestore  string
		buttonDelete   string
		buttonArchive  string
		buttonBorrow   string
		buttonEdit     string
		buttonHistory  string
		iconBorrowing  themes.IconFace
		borrowingTitle string
	)

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	if types.CurrentStorage.Borrowing == nil || !types.CurrentStorage.Borrowing.BorrowingID.Valid {
		iconBorrowing = themes.MDI_BORROW
		borrowingTitle = locales.Translate("storage_borrow", HTTPHeaderAcceptLanguage)
	} else {
		iconBorrowing = themes.MDI_UNBORROW
		borrowingTitle = locales.Translate("storage_unborrow", HTTPHeaderAcceptLanguage)
	}

	if types.CurrentStorage.StorageArchive.Valid && types.CurrentStorage.StorageArchive.Bool {

		// This is an archive.
		buttonClone = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "clone" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"clone"},
					Visible: false,
				},
				Title: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_CLONE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

		buttonRestore = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:         "restore" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes:    []string{"restore"},
					Visible:    false,
					Attributes: map[string]string{"sid": strconv.Itoa(int(types.CurrentStorage.StorageID.Int64))},
				},
				Title: locales.Translate("storage_restore", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("storage_restore", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_RESTORE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

		buttonDelete = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "delete" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"storagedelete"},
					Visible: false,
				},
				Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("delete", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

	} else if types.CurrentStorage.Storage.StorageID.Valid {

		// This is an history.
		buttonClone = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "clone" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"clone"},
					Visible: false,
				},
				Title: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_CLONE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

	} else {

		buttonEdit = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "edit" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"storageedit"},
					Visible: false,
				},
				Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("edit", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
		buttonClone = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "clone" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"clone"},
					Visible: false,
				},
				Title: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("storage_clone", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_CLONE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
		buttonArchive = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "archive" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"archive"},
					Visible: false,
				},
				Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("archive", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_ARCHIVE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
		buttonBorrow = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:         "borrow" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes:    []string{"borrow"},
					Visible:    false,
					Attributes: map[string]string{"sid": strconv.Itoa(int(types.CurrentStorage.StorageID.Int64))},
				},
				Title: borrowingTitle,
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Attributes: map[string]string{
						"data-target": "#borrow",
					},
					Classes: []string{"iconlabel"},
				},
				Text: borrowingTitle,
				Icon: themes.NewMdiIcon(iconBorrowing, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

	}

	if types.CurrentStorage.StorageHC != 0 {

		buttonHistory = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "history" + strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
					Classes: []string{"history"},
					Visible: false,
				},
				Title: locales.Translate("storage_history", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("storage_history", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_HISTORY, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()

	}

	spanId := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Text: strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
	}).OuterHTML()

	return buttonClone + buttonRestore + buttonDelete + buttonArchive + buttonBorrow + buttonEdit + buttonHistory + spanId

}

func DetailFormatter(this js.Value, args []js.Value) interface{} {

	var (
		qrCode string
	)

	row := args[1]
	types.CurrentStorage = Storage{}.FromJsJSONValue(row)

	if len(types.CurrentStorage.StorageQRCode) == 0 {
		qrCode = `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`
	} else {
		qrCode = base64.StdEncoding.EncodeToString(types.CurrentStorage.StorageQRCode)
	}

	//
	// Qrcode, to destroy and ID.
	//
	rowQrcodeAndID := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Qrcode.
	colQrcode := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	colQrcode.AppendChild(widgets.NewImg(widgets.ImgAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Height: "128px",
		Width:  "128px",
		Src:    fmt.Sprintf("data:image/png;base64,%s", qrCode),
	}))
	// To destroy.
	colToDestroy := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	if types.CurrentStorage.StorageToDestroy.Valid && types.CurrentStorage.StorageToDestroy.Bool {
		colToDestroy.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{".text-danger", "mr-sm-2"},
				},
				Text: locales.Translate("storage_todestroy_title", HTTPHeaderAcceptLanguage),
			}))
	}
	// ID.
	colID := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: "id",
		}))
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: strconv.Itoa(int(types.CurrentStorage.StorageID.Int64)),
		}))

	rowQrcodeAndID.AppendChild(colQrcode)
	rowQrcodeAndID.AppendChild(colToDestroy)
	rowQrcodeAndID.AppendChild(colID)

	//
	// Print QRCode.
	//
	rowPrintQrcode := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	colPrint := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12"},
		},
	})

	var icon widgets.Widget
	icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: locales.Translate("storage_print_qrcode", HTTPHeaderAcceptLanguage),
		Icon: themes.NewMdiIcon(themes.MDI_PRINT, themes.MDI_24PX),
	})

	colPrint.AppendChild(widgets.NewLink(widgets.LinkAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Title:   locales.Translate("storage_print_qrcode", HTTPHeaderAcceptLanguage),
		Onclick: "PrintImage('" + fmt.Sprintf("data:image/png;base64,%s", qrCode) + "', '" + fmt.Sprintf("%s - %s", types.CurrentStorage.Product.Name.NameLabel, types.CurrentStorage.StoreLocation.StoreLocationFullPath) + "'); return false;",
		Href:    "#",
		Label:   icon,
	}))

	rowPrintQrcode.AppendChild(colPrint)

	//
	// Product and store location.
	//
	rowProductAndStorelocation := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Product.
	colProduct := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colProduct.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "mr-sm-2"},
		},
		Text: locales.Translate("storage_product_table_header", HTTPHeaderAcceptLanguage),
	}))
	colProduct.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: types.CurrentStorage.Product.Name.NameLabel,
		}))
	// Store location.
	colStorelocation := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colStorelocation.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "mr-sm-2"},
		},
		Text: locales.Translate("storage_storelocation_table_header", HTTPHeaderAcceptLanguage),
	}))
	colStorelocation.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: types.CurrentStorage.StoreLocation.StoreLocationFullPath,
		}))

	rowProductAndStorelocation.AppendChild(colProduct)
	rowProductAndStorelocation.AppendChild(colStorelocation)

	//
	// Quantity and barecode.
	//
	rowQuantityandBarecode := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Quantity.
	colQuantity := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colQuantity.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "mr-sm-2"},
		},
		Text: locales.Translate("storage_quantity_title", HTTPHeaderAcceptLanguage),
	}))
	colQuantity.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: fmt.Sprintf("%v %s", types.CurrentStorage.StorageQuantity.Float64, types.CurrentStorage.UnitQuantity.UnitLabel.String),
		}))
	// Barecode.
	colBarecode := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colBarecode.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "mr-sm-2"},
		},
		Text: locales.Translate("storage_barecode_title", HTTPHeaderAcceptLanguage),
	}))
	colBarecode.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: types.CurrentStorage.StorageBarecode.String,
		}))

	rowQuantityandBarecode.AppendChild(colQuantity)
	rowQuantityandBarecode.AppendChild(colBarecode)

	//
	// Concentration and batch number
	//
	rowConcentrationAndBatchnumber := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Concentration.
	colConcentration := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if types.CurrentStorage.StorageConcentration.Valid {
		colConcentration.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_concentration_title", HTTPHeaderAcceptLanguage),
		}))
		colConcentration.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: fmt.Sprintf("%d %s", types.CurrentStorage.StorageConcentration.Int64, types.CurrentStorage.UnitConcentration.UnitLabel.String),
			}))
	}
	// Batch number.
	colBatchnumber := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if types.CurrentStorage.StorageBatchNumber.Valid && types.CurrentStorage.StorageBatchNumber.String != "" {
		colBatchnumber.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_batchnumber_title", HTTPHeaderAcceptLanguage),
		}))
		colBatchnumber.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageBatchNumber.String,
			}))
	}

	rowConcentrationAndBatchnumber.AppendChild(colConcentration)
	rowConcentrationAndBatchnumber.AppendChild(colBatchnumber)

	//
	// Supplier and reference.
	//
	rowSupplierAndReference := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Supplier.
	colSupplier := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if types.CurrentStorage.Supplier.SupplierID.Valid {
		colSupplier.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("supplier_label_title", HTTPHeaderAcceptLanguage),
		}))
		colSupplier.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.Supplier.SupplierLabel.String,
			}))
	}
	// Reference.
	colReference := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if types.CurrentStorage.StorageReference.Valid {
		colReference.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_reference_title", HTTPHeaderAcceptLanguage),
		}))
		colReference.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageReference.String,
			}))
	}

	rowSupplierAndReference.AppendChild(colSupplier)
	rowSupplierAndReference.AppendChild(colReference)

	//
	// Dates.
	//
	rowDates := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Entry date.
	colEntryDate := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-3"},
		},
	})
	if types.CurrentStorage.StorageEntryDate.Valid {
		colEntryDate.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_entrydate_title", HTTPHeaderAcceptLanguage),
		}))
		colEntryDate.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageEntryDate.Time.String(),
			}))
	}
	// Exit date.
	colExitDate := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-3"},
		},
	})
	if types.CurrentStorage.StorageExitDate.Valid {
		colExitDate.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_exitdate_title", HTTPHeaderAcceptLanguage),
		}))
		colExitDate.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageExitDate.Time.String(),
			}))
	}
	// Opening date.
	colOpeningDate := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-3"},
		},
	})
	if types.CurrentStorage.StorageOpeningDate.Valid {
		colOpeningDate.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_openingdate_title", HTTPHeaderAcceptLanguage),
		}))
		colOpeningDate.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageOpeningDate.Time.String(),
			}))
	}
	// Expiration date.
	colExpirationDate := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-3"},
		},
	})
	if types.CurrentStorage.StorageExpirationDate.Valid {
		colExpirationDate.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_expirationdate_title", HTTPHeaderAcceptLanguage),
		}))
		colExpirationDate.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageExpirationDate.Time.String(),
			}))
	}

	rowDates.AppendChild(colEntryDate)
	rowDates.AppendChild(colExitDate)
	rowDates.AppendChild(colOpeningDate)
	rowDates.AppendChild(colExpirationDate)

	//
	// Comment
	//
	rowComment := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Comment.
	colComment := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12"},
		},
	})
	if types.CurrentStorage.StorageComment.Valid && types.CurrentStorage.StorageComment.String != "" {
		colComment.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("storage_comment_title", HTTPHeaderAcceptLanguage),
		}))
		colComment.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: types.CurrentStorage.StorageComment.String,
			}))
	}

	rowComment.AppendChild(colComment)

	//
	// Creation date and person
	//
	rowCreationDateAndPerson := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Creation date.
	colCreationDate := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colCreationDate.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "mr-sm-2"},
		},
		Text: locales.Translate("created", HTTPHeaderAcceptLanguage),
	}))
	colCreationDate.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: types.CurrentStorage.StorageCreationDate.String(),
		}))
	// Person.
	colPerson := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	colPerson.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"blockquote-footer", "mr-sm-2"},
		},
		Text: types.CurrentStorage.Person.PersonEmail,
	}))

	rowCreationDateAndPerson.AppendChild(colCreationDate)
	rowCreationDateAndPerson.AppendChild(colPerson)

	return rowQrcodeAndID.OuterHTML() +
		rowPrintQrcode.OuterHTML() +
		rowProductAndStorelocation.OuterHTML() +
		rowQuantityandBarecode.OuterHTML() +
		rowConcentrationAndBatchnumber.OuterHTML() +
		rowSupplierAndReference.OuterHTML() +
		rowDates.OuterHTML() +
		rowComment.OuterHTML() +
		rowCreationDateAndPerson.OuterHTML()

}

func DataQueryParams(this js.Value, args []js.Value) interface{} {

	params := args[0]

	queryFilter := QueryFilterFromJsJSONValue(params)

	// Storage_SaveCallback storage id.
	queryFilter.Product = BSTableQueryFilter.Product
	queryFilter.ProductFilterLabel = BSTableQueryFilter.ProductFilterLabel
	queryFilter.Storage = BSTableQueryFilter.Storage
	queryFilter.StorageFilterLabel = BSTableQueryFilter.StorageFilterLabel
	queryFilter.StorageArchive = BSTableQueryFilter.StorageArchive
	queryFilter.StorageHistory = BSTableQueryFilter.StorageHistory
	queryFilter.Export = BSTableQueryFilter.Export
	BSTableQueryFilter.Export = false
	BSTableQueryFilter.Unlock()

	if Jq("select#s_storelocation").Select2IsInitialized() {
		i := Jq("select#s_storelocation").Select2Data()
		if len(i) > 0 {
			queryFilter.StoreLocation = i[0].Id
			queryFilter.StoreLocationFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_name").Select2IsInitialized() {
		i := Jq("select#s_name").Select2Data()
		if len(i) > 0 {
			queryFilter.Name = i[0].Id
			queryFilter.NameFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_casnumber").Select2IsInitialized() {
		i := Jq("select#s_casnumber").Select2Data()
		if len(i) > 0 {
			queryFilter.CasNumber = i[0].Id
			queryFilter.CasNumberFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_empiricalformula").Select2IsInitialized() {
		i := Jq("select#s_empiricalformula").Select2Data()
		if len(i) > 0 {
			queryFilter.EmpiricalFormula = i[0].Id
			queryFilter.EmpiricalFormulaFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_signalword").Select2IsInitialized() {
		i := Jq("select#s_signalword").Select2Data()
		if len(i) > 0 {
			queryFilter.SignalWord = i[0].Id
		}
	}
	if Jq("select#s_hazardstatements").Select2IsInitialized() {
		i := Jq("select#s_hazardstatements").Select2Data()
		if len(i) > 0 {
			for _, hs := range i {
				queryFilter.HazardStatements = append(queryFilter.HazardStatements, hs.Id)
			}
		}
	}
	if Jq("select#s_precautionarystatements").Select2IsInitialized() {
		i := Jq("select#s_precautionarystatements").Select2Data()
		if len(i) > 0 {
			for _, ps := range i {
				queryFilter.PrecautionaryStatements = append(queryFilter.PrecautionaryStatements, ps.Id)
			}
		}
	}
	if Jq("select#s_symbols").Select2IsInitialized() {
		i := Jq("select#s_symbols").Select2Data()
		if len(i) > 0 {
			for _, s := range i {
				queryFilter.Symbols = append(queryFilter.Symbols, s.Id)
			}
		}
	}

	if Jq("#s_storage_barecode").GetVal().Truthy() {
		queryFilter.StorageBarecode = Jq("#s_storage_barecode").GetVal().String()
	}
	if Jq("#s_custom_name_part_of").GetVal().Truthy() {
		queryFilter.CustomNamePartOf = Jq("#s_custom_name_part_of").GetVal().String()
	}
	if Jq("#s_casnumber_cmr:checked").Object.Length() > 0 {
		queryFilter.CasNumberCMR = true
	}
	if Jq("#s_borrowing:checked").Object.Length() > 0 {
		queryFilter.Borrowing = true
	}
	if Jq("#s_storage_to_destroy:checked").Object.Length() > 0 {
		queryFilter.StorageToDestroy = true
	}

	queryFilter.DisplayFilter()

	return queryFilter.ToJsValue()

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "storages"}
		u.RawQuery = params.Data.ToRawQuery()

		ajax := Ajax{
			URL:    u.String(),
			Method: "get",
			Done: func(data js.Value) {

				var (
					storages Storages
					err      error
				)
				if err = json.Unmarshal([]byte(data.String()), &storages); err != nil {
					fmt.Println(err)
				}

				if storages.GetExportFn() != "" {
					utils.DisplaySuccessMessage(locales.Translate("export_done", HTTPHeaderAcceptLanguage))

					var icon widgets.Widget
					icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Text: locales.Translate("download_export", HTTPHeaderAcceptLanguage),
						Icon: themes.NewMdiIcon(themes.MDI_DOWNLOAD, themes.MDI_24PX),
					})

					downloadLink := widgets.NewLink(widgets.LinkAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
							Classes: []string{"iconlabel"},
						},
						Onclick: "$('#export').hide()",
						Title:   locales.Translate("download_export", HTTPHeaderAcceptLanguage),
						Href:    fmt.Sprintf("%sdownload/%s", ApplicationProxyPath, storages.GetExportFn()),
						Label:   icon,
					})

					Jq("#export-body").SetHtml("")
					Jq("#export-body").Append(downloadLink.OuterHTML())
					Jq("#export").Show()

				} else if storages.GetTotal() != 0 {
					row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))
				} else {
					// FIXME: Does not work
					//Jq("#Product_table").Bootstraptable(nil).RemoveAll()
					utils.DisplayErrorMessage(locales.Translate("no_result", HTTPHeaderAcceptLanguage))
				}

			},
			Fail: func(jqXHR js.Value) {

				utils.DisplayGenericErrorMessage()

			},
		}

		ajax.Send()

	}()

	return nil

}
