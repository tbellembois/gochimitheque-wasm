// +build go1.15,js,wasm

package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/about"
	"github.com/tbellembois/gochimitheque-wasm/views/common"
	"github.com/tbellembois/gochimitheque-wasm/views/entity"
	"github.com/tbellembois/gochimitheque-wasm/views/login"
	"github.com/tbellembois/gochimitheque-wasm/views/person"
	"github.com/tbellembois/gochimitheque-wasm/views/personpass"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/views/storelocation"
	"github.com/tbellembois/gochimitheque-wasm/views/welcomeannounce"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque/data"
)

var (
	fullUrl *url.URL
	signal  = make(chan int)
	err     error
)

func keepAlive() {
	for {
		<-signal
	}
}

func init() {

	types.Jq = types.NewJquery

	Win = js.Global().Get("window")
	Doc = js.Global().Get("document")

	fullUrl, err = url.Parse(js.Global().Get("location").Get("href").String())
	if err != nil {
		panic(err)
	}
	URLParameters, err = url.ParseQuery(fullUrl.RawQuery)
	if err != nil {
		panic(err)
	}

	types.BSTableQueryFilter = types.SafeQueryFilter{
		QueryFilter: types.QueryFilter{},
	}

	// TODO: factorize the js and wasm functions.
	CurrentView = "product"

	// TODO: get the variables from Go instead of JS
	ApplicationProxyPath = js.Global().Get("container").Get("ProxyPath").String()
	HTTPHeaderAcceptLanguage = js.Global().Get("container").Get("PersonLanguage").String()
	DisableCache, err = strconv.ParseBool(js.Global().Get("disableCache").String())
	if err != nil {
		panic(err)
	}

	// Initializing the slices of statements for the magic selector.
	var (
		r       *csv.Reader
		records [][]string
	)
	r = csv.NewReader(strings.NewReader(data.PRECAUTIONARYSTATEMENT))
	r.Comma = '\t'
	if records, err = r.ReadAll(); err != nil {
		panic(err)
	}
	// FIXME: we assume here that the id starts by 1 in the DB
	for id, record := range records {
		types.DBPrecautionaryStatements = append(types.DBPrecautionaryStatements,
			types.PrecautionaryStatement{
				PrecautionaryStatementID:        id + 1,
				PrecautionaryStatementLabel:     record[0],
				PrecautionaryStatementReference: record[1],
			})
	}

	r = csv.NewReader(strings.NewReader(data.HAZARDSTATEMENT))
	r.Comma = '\t'
	if records, err = r.ReadAll(); err != nil {
		panic(err)
	}
	// FIXME: we assume here that the id starts by 1 in the DB
	for id, record := range records {
		types.DBHazardStatements = append(types.DBHazardStatements,
			types.HazardStatement{
				HazardStatementID:        id + 1,
				HazardStatementLabel:     record[0],
				HazardStatementReference: record[1],
			})
	}

}

func Test(this js.Value, args []js.Value) interface{} {

	return nil

}

func main() {

	// Common actions for all pages.
	types.Jq("#table").On("load-success.bs.table",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			search := URLParameters.Get("search")
			if search != "" {
				types.Jq("#table").Bootstraptable(nil).ResetSearch(search)
			}
			return nil
		}))

	// Test.
	js.Global().Set("Test", js.FuncOf(Test))

	// Global functions.
	js.Global().Set("Utils_closeEdit", js.FuncOf(utils.CloseEdit))
	js.Global().Set("Utils_message", js.FuncOf(utils.DisplayMessageWrapper))
	js.Global().Set("Utils_translate", js.FuncOf(locales.TranslateWrapper))

	js.Global().Set("Widgets_permission", js.FuncOf(widgets.PermissionWrapper))
	js.Global().Set("Widgets_title", js.FuncOf(widgets.TitleWrapper))

	// Product/Storage common functions.
	js.Global().Set("Common_SwitchProductStorage", js.FuncOf(common.SwitchProductStorageWrapper))
	js.Global().Set("Common_export", js.FuncOf(common.Export))

	// Product bootstraptable functions.
	js.Global().Set("Product_operateEventsBookmark", js.FuncOf(product.OperateEventsBookmark))
	js.Global().Set("Product_operateEventsStore", js.FuncOf(product.OperateEventsStore))
	js.Global().Set("Product_operateEventsStorages", js.FuncOf(product.OperateEventsStorages))
	js.Global().Set("Product_operateEventsOStorages", js.FuncOf(product.OperateEventsOStorages))
	js.Global().Set("Product_operateEventsEdit", js.FuncOf(product.OperateEventsEdit))
	js.Global().Set("Product_operateEventsDelete", js.FuncOf(product.OperateEventsDelete))
	js.Global().Set("Product_getTableData", js.FuncOf(product.GetTableData))
	js.Global().Set("Product_dataQueryParams", js.FuncOf(product.DataQueryParams))
	js.Global().Set("Product_detailFormatter", js.FuncOf(product.DetailFormatter))
	js.Global().Set("Product_empiricalformulaFormatter", js.FuncOf(product.EmpiricalformulaFormatter))
	js.Global().Set("Product_casnumberFormatter", js.FuncOf(product.CasnumberFormatter))
	js.Global().Set("Product_productspecificityFormatter", js.FuncOf(product.Product_productSpecificityFormatter))
	js.Global().Set("Product_productslFormatter", js.FuncOf(product.Product_productSlFormatter))
	js.Global().Set("Product_operateFormatter", js.FuncOf(product.OperateFormatter))
	js.Global().Set("Product_addProducer", js.FuncOf(product.AddProducer))
	js.Global().Set("Product_addSupplier", js.FuncOf(product.AddSupplier))

	js.Global().Set("Product_linearToEmpirical", js.FuncOf(product.LinearToEmpirical))
	js.Global().Set("Product_noCas", js.FuncOf(product.NoCas))
	js.Global().Set("Product_noEmpiricalFormula", js.FuncOf(product.NoEmpiricalFormula))
	js.Global().Set("Product_magic", js.FuncOf(product.Magic))
	js.Global().Set("Product_howToMagicalSelector", js.FuncOf(product.HowToMagicalSelector))

	js.Global().Set("Common_search", js.FuncOf(search.Search))
	js.Global().Set("Common_clearSearch", js.FuncOf(search.ClearSearch))

	js.Global().Set("Product_saveProduct", js.FuncOf(product.SaveProduct))

	// Product page load callbacks.
	js.Global().Set("Product_listBookmark", js.FuncOf(product.Product_listBookmarkCallback))
	js.Global().Set("Product_list", js.FuncOf(product.Product_listCallback))
	js.Global().Set("Product_create", js.FuncOf(product.Product_createCallback))

	// Storage bootstraptable functions.
	js.Global().Set("Storage_getTableData", js.FuncOf(storage.GetTableData))
	js.Global().Set("Storage_dataQueryParams", js.FuncOf(storage.DataQueryParams))
	js.Global().Set("Storage_detailFormatter", js.FuncOf(storage.DetailFormatter))
	js.Global().Set("Storage_productFormatter", js.FuncOf(storage.Storage_productFormatter))
	js.Global().Set("Storage_storelocationFormatter", js.FuncOf(storage.Storage_storelocationFormatter))
	js.Global().Set("Storage_quantityFormatter", js.FuncOf(storage.Storage_quantityFormatter))
	js.Global().Set("Storage_barecodeFormatter", js.FuncOf(storage.Storage_barecodeFormatter))
	js.Global().Set("Storage_operateFormatter", js.FuncOf(storage.Storage_operateFormatter))
	js.Global().Set("Storage_operateEventsRestore", js.FuncOf(storage.Storage_operateEventsRestore))
	js.Global().Set("Storage_operateEventsClone", js.FuncOf(storage.Storage_operateEventsClone))
	js.Global().Set("Storage_operateEventsHistory", js.FuncOf(storage.Storage_operateEventsHistory))
	js.Global().Set("Storage_operateEventsBorrow", js.FuncOf(storage.Storage_operateEventsBorrow))
	js.Global().Set("Storage_operateEventsEdit", js.FuncOf(storage.Storage_operateEventsEdit))
	js.Global().Set("Storage_operateEventsArchive", js.FuncOf(storage.Storage_operateEventsArchive))
	js.Global().Set("Storage_operateEventsDelete", js.FuncOf(storage.Storage_operateEventsDelete))

	js.Global().Set("Storage_saveStorage", js.FuncOf(storage.SaveStorage))
	js.Global().Set("Storage_saveBorrowing", js.FuncOf(storage.SaveBorrowing))

	js.Global().Set("Storage_scanQRdone", js.FuncOf(storage.ScanQRdone))

	// Storage page load callbacks.
	js.Global().Set("Storage_list", js.FuncOf(storage.Storage_listCallback))

	// Entity bootstraptable functions.
	js.Global().Set("Entity_operateEventsStorelocations", js.FuncOf(entity.OperateEventsStorelocations))
	js.Global().Set("Entity_operateEventsMembers", js.FuncOf(entity.OperateEventsMembers))
	js.Global().Set("Entity_operateEventsEdit", js.FuncOf(entity.OperateEventsEdit))
	js.Global().Set("Entity_operateEventsDelete", js.FuncOf(entity.OperateEventsDelete))
	js.Global().Set("Entity_getTableData", js.FuncOf(entity.GetTableData))
	js.Global().Set("Entity_managersFormatter", js.FuncOf(entity.ManagersFormatter))
	js.Global().Set("Entity_operateFormatter", js.FuncOf(entity.OperateFormatter))
	js.Global().Set("Entity_saveEntity", js.FuncOf(entity.SaveEntity))

	// Entity page load callbacks.
	js.Global().Set("Entity_list", js.FuncOf(entity.Entity_listCallback))
	js.Global().Set("Entity_create", js.FuncOf(entity.Entity_createCallBack))

	// StoreLocation bootstraptable functions.
	js.Global().Set("StoreLocation_operateEventsEdit", js.FuncOf(storelocation.OperateEventsEdit))
	js.Global().Set("StoreLocation_operateEventsDelete", js.FuncOf(storelocation.OperateEventsDelete))
	js.Global().Set("StoreLocation_getTableData", js.FuncOf(storelocation.GetTableData))
	js.Global().Set("StoreLocation_operateFormatter", js.FuncOf(storelocation.OperateFormatter))
	js.Global().Set("StoreLocation_colorFormatter", js.FuncOf(storelocation.ColorFormatter))
	js.Global().Set("StoreLocation_canStoreFormatter", js.FuncOf(storelocation.CanStoreFormatter))
	js.Global().Set("StoreLocation_storeLocationFormatter", js.FuncOf(storelocation.StoreLocationFormatter))
	js.Global().Set("StoreLocation_saveStoreLocation", js.FuncOf(storelocation.SaveStoreLocation))
	js.Global().Set("StoreLocation_dataQueryParams", js.FuncOf(storelocation.DataQueryParams))

	// StoreLocation page load callbacks.
	js.Global().Set("StoreLocation_list", js.FuncOf(storelocation.StoreLocation_listCallback))
	js.Global().Set("StoreLocation_create", js.FuncOf(storelocation.StoreLocation_createCallBack))

	// Person bootstraptable functions.
	js.Global().Set("Person_operateEventsEdit", js.FuncOf(person.OperateEventsEdit))
	js.Global().Set("Person_getTableData", js.FuncOf(person.GetTableData))
	js.Global().Set("Person_savePerson", js.FuncOf(person.SavePerson))
	js.Global().Set("Person_dataQueryParams", js.FuncOf(person.DataQueryParams))
	js.Global().Set("Person_operateFormatter", js.FuncOf(person.OperateFormatter))

	// Person page load callbacks.
	js.Global().Set("Person_list", js.FuncOf(person.Person_listCallback))
	js.Global().Set("Person_create", js.FuncOf(person.Person_createCallBack))

	// Welcome announce
	js.Global().Set("WelcomeAnnounce_saveWelcomeAnnounce", js.FuncOf(welcomeannounce.SaveWelcomeAnnounce))
	js.Global().Set("WelcomeAnnounce_list", js.FuncOf(welcomeannounce.WelcomeAnnounce_listCallback))

	// Person password
	js.Global().Set("PersonPass_savePersonPassword", js.FuncOf(personpass.SavePersonPassword))
	js.Global().Set("PersonPass_list", js.FuncOf(personpass.PersonPass_listCallback))

	// Login
	js.Global().Set("Login_getToken", js.FuncOf(login.GetToken))
	js.Global().Set("Login_getCaptcha", js.FuncOf(login.GetCaptcha))
	js.Global().Set("Login_resetPassword", js.FuncOf(login.ResetPassword))
	js.Global().Set("Login_getAnnounce", js.FuncOf(login.GetAnnounce))

	js.Global().Set("Login_list", js.FuncOf(login.Login_listCallback))

	// About
	js.Global().Set("About_list", js.FuncOf(about.About_listCallback))

	types.Jq("#loading").Empty()
	types.Jq("div.container").RemoveClass("invisible")

	utils.DisplaySuccessMessage(locales.Translate("wasm_loaded", HTTPHeaderAcceptLanguage))

	// For the login and home pages, calling the callback
	// function manually.
	// For the other pages they are called when clicking
	// on the menu.
	if ConnectedUserEmail != "" {
		productCallbackWrapper := func(args ...interface{}) {
			product.Product_listCallback(js.Null(), nil)
		}
		utils.LoadContent("product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper)
	} else {
		loginCallbackWrapper := func(args ...interface{}) {
			login.Login_listCallback(js.Null(), nil)
		}
		utils.LoadContent("login", fmt.Sprintf("%slogin", ApplicationProxyPath), loginCallbackWrapper)
	}

	keepAlive()

}
