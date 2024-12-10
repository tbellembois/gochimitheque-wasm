//go:build go1.21 && js && wasm

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/localstorage"
	"github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/about"
	"github.com/tbellembois/gochimitheque-wasm/views/common"
	"github.com/tbellembois/gochimitheque-wasm/views/entity"
	"github.com/tbellembois/gochimitheque-wasm/views/login"
	"github.com/tbellembois/gochimitheque-wasm/views/menu"
	"github.com/tbellembois/gochimitheque-wasm/views/person"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/views/storelocation"
	"github.com/tbellembois/gochimitheque-wasm/views/welcomeannounce"
	"github.com/tbellembois/gochimitheque-wasm/widgets"

	"github.com/tbellembois/gochimitheque/models"
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

// Container is a struct passed to the view.
type Container struct {
	PersonEmail    string `json:"PersonEmail"`
	PersonLanguage string `json:"PersonLanguage"`
	PersonID       int    `json:"PersonID"`
	AppURL         string `json:"AppURL"`
	AppPath        string `json:"AppPath"`
	BuildID        string `json:"BuildID"`
	DisableCache   bool   `json:"DisableCache"`
}

func init() {

	fullUrl, err = url.Parse(js.Global().Get("location").Get("href").String())
	if err != nil {
		panic(err)
	}
	URLParameters, err = url.ParseQuery(fullUrl.RawQuery)
	if err != nil {
		panic(err)
	}

	globals.LocalStorage = localstorage.NewLocalStorage()

	globals.BSTableQueryFilter = ajax.SafeQueryFilter{
		QueryFilter: ajax.QueryFilter{},
	}

	// TODO: factorize the js and wasm functions.
	CurrentView = "product"

	var c Container
	cString := js.Global().Get("JSON").Call("stringify", js.Global().Get("c")).String()
	if err = json.Unmarshal([]byte(cString), &c); err != nil {
		panic(err)
	}

	ApplicationProxyPath = c.AppPath
	HTTPHeaderAcceptLanguage = c.PersonLanguage
	DisableCache = c.DisableCache

	// Initializing the slices of statements for the magic selector.
	var (
		r       *csv.Reader
		records [][]string
	)
	r = csv.NewReader(strings.NewReader(globals.PRECAUTIONARYSTATEMENT))
	r.Comma = '\t'
	if records, err = r.ReadAll(); err != nil {
		panic(err)
	}
	// FIXME: we assume here that the id starts by 1 in the DB
	for id, record := range records {
		globals.DBPrecautionaryStatements = append(globals.DBPrecautionaryStatements,
			types.PrecautionaryStatement{PrecautionaryStatement: &models.PrecautionaryStatement{
				PrecautionaryStatementID:        id + 1,
				PrecautionaryStatementLabel:     record[0],
				PrecautionaryStatementReference: record[1],
			}})
	}

	r = csv.NewReader(strings.NewReader(globals.HAZARDSTATEMENT))
	r.Comma = '\t'
	if records, err = r.ReadAll(); err != nil {
		panic(err)
	}
	// FIXME: we assume here that the id starts by 1 in the DB
	for id, record := range records {
		globals.DBHazardStatements = append(globals.DBHazardStatements,
			types.HazardStatement{HazardStatement: &models.HazardStatement{
				HazardStatementID:        id + 1,
				HazardStatementLabel:     record[0],
				HazardStatementReference: record[1],
			}})
	}

	// Initializing map of symbol images.
	globals.SymbolImages = make(map[string]string)
	globals.SymbolImages["GHS01"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAInSURBVFiFzdi9b45RGMfxz11KNalS0napgVaE0HhLxEtpVEIQxKAbiUhMJCoisTyriYTBLCZGsRCrP8JkkLB4iVlyDL1EU9refZ7raZ3kJPfrub75/a7zWpVSpJSqaoBSGintlVJarzQKJWojo81sqDS4XKUSlcu3LwkuFyoRLh8qCa49UAlw6VDoRUercOlK4T7WtapcNlQHHmfYmgm1CVO4MOPZBow31V4SVDde4i22xrMeXEfVVK5myB4gh/AQVwPqSljbEe/XYVfd9lOgIvAt3AnAS1iBczgV168wVTdOClSAPcMwzmIg4EbRP+u7behZKF6r9q3BTTzFC1wLO49iD/owHioex2nswGpsnC9uU1BYhUE8R8EH3As1DuIYtmAnDsT9SZwPJScxMp8o9RKRtQHSFUk8jBHcxpPIr95QqC+svIxHGKiVDrM4VqpRSik/qqoaxTecwSe8CUWO4Dve4W6o9xFf8Bl9VVV1RgfoDLXfl1J+LhR0bp+nVRjGZoxhLw7jRNhzIwAKXmMCD/AVDVxsRq3ayY/1GEK/6RF+u+k5cTAUGJoxVk1ionaPnjf568HtD6h9GJunY3RjN7qahfobrEYP9Xv0brUuaoCt+VO7oeYGaydcS5N4u+BSlj3ZcKkLxSy4tiytW4Vr62ak2SBLsn1bbLAl3fDWDbosRwQLBV/WQ5W5IP6LY6h/w6VA5YAl2jez1lrBLlhKaaiqP9cJ5Rf+De5Q3HyidwAAAABJRU5ErkJggg=="
	globals.SymbolImages["GHS02"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC"
	globals.SymbolImages["GHS03"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAJhSURBVFiFzdjNq41RFMfxz/FSXK4iUV4GUjJRMpG8lQnJUF0mTAy8xkApimOEMDA2oSgD5R9gcE0obxO6JkooBroJhRttg7Nu94lzzz37OdvLrt16zrPXs9d3/9Ze+5zzNFJKirRGowlSahaZL6XUe6eZSNGbJeYsDVUMrqxSBZUrn75CcMWgsARTSsEVUwq7sLWUcmXS1wK7iOvd+pcDa5++qVgU1/fxKq4n9wrXa/qW4Bb6MIKE2diPmb3A1VYq7MqAuRQ2YTdeY6CXtNZVaj4uYG0FaLR/D3sc0+vC1d3okwJgsAI0iB+Vz5dxe1TdXLg6UHPCvg2AT2E34VobBaflzD8+2AQPRYqu4kUEPh1KzcKOuPck7CMcQF92nOyVtCquqsg8PI5C2IyHWBFjn8NuzM5Mdu7ZGcGO4k1U5EgF9CNO4QuG4t6x7ALLrhY2RLB9uBMAJ7Ea63A+CuMVlobvidzqzz9fmFtR5jvWtPHZHj4Xww5MNO+vHJNktpTSezxAP26klO618bkZah4JRe/mxslOZSiyLZQ43MHnTPicy1Wr1uavBH6Hsx3Gr+ADZudC1TouKoFv4CX624wtwDBO1oH6HSwDDsvxTetrZ2Hl/jKtg3UYs+pAtQfLg1uldcqPhH2qVanPsL4u1Phg3ayIPdiLg3hu7IAdwqEY24vFtRZdew/wtQLTqW+ptYc7gnWYLPbS8i76jFyo7sBqTFri+T86eS/P/dmV/5W/b7nB/uof3m6D/pNXBBMF/6cvVcaD+C9eQ7WHKwJVBqxg+qp9SvYvy3YtpaZGY+y6QPsJlPiFVobY9AkAAAAASUVORK5CYII="
	globals.SymbolImages["GHS04"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAFtSURBVFiFzdixLgRBGADgb4UoBYXKA3gAHYU3EKWH0FKuAk8gGm+hvU4oJAqlaDVqiUSCUbjEJe5ys7v/7NnkT66Ynf/bmZ3bf6ZKKQm5qqoGKdUh/aWUugd1Ig2jjugzGhWGix2pwJGLn74gXCwqEBePCsKVQQXgyqE63lcW1eH+8qiW/fSDatFff6iG/faLatB//6jMPKEoLGEX53jEHTaxj0sMcvN1QmEBWzjGLT6QxsQX7nGUO3KNUdjAAa7wOgEyGjdYazqt+auEQzxnQEbjBdtt3rn5nCq3qqp1nGJuStM3XGMwjIc0fKrGV9YKYQXv/o7Ip58X/AQ7WIxaofnLl7Mh5gkX2MNyqb+NrEYjuNXOkMx8jRr3hRoP6wPX6pNUGtfpI14KF1L2RONCC8UoXJHSuiuu6GakbZJetm9Nk/W64c1NOpMjgmnJZ3qoMgnxL46hxuNCUDGwwOkbjawKNqParFXV7++A6xtDLLIHRMAuWAAAAABJRU5ErkJggg=="
	globals.SymbolImages["GHS05"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI9SURBVFiFzdhPiI1RGMfxz2FmgYVo1CyYhVJISuPPDmFGNwtW2EzKn4TEELKaa2NNdiyUlWxtUGYptiytrRSNIk3GsbinceO+7vvOPe8dp56677k9v+d7nnPOczonxBhlaSE0QYzNLHoxxt6NZiQma+bQzA2VDS5vpjJmLv/0ZYLLC5URLj9UJrh6oDLA1QfVo1+9UD341w+1QJ2exLAG+3AKq9DAGUxgJ1YsFG6gy9k3lb5uibEZQhjHNnxFxEe8w1ucwxfcxxw24GAIYTkCBmOMTSFIulNCKDxbO4N1gEq/P6SsfMdw6hvFDB5hGS5iLYZSRr/hE6bRAikDV2X6sC4B7MZQ13XSAlufII5iBGOllknFNTWGq6W3PJO4ngYygRu4WSoJJTO1N62TXWnqtmK0BNgR3E4DmsQdjJRJRtlMNXAX53EBz7G6BNgAjmEH9hT6dIhfagtjCV5gC67hUgmoA3iNl3iDe7iSQAe7wRWWixBCAyfauj5jo9bOPBRCeNz239kY40yb72lsx5MEth9PY4zvQwgr8aMo7nwrTCWXtWpVGRv+I1vH8SxlaRxLsblKIQ9J6K/aFXiITV1H1mrTMcbZ9o4QwmGc1CrGr/BTa83N4kGMca5T3PmattAjI4uVKhf9hqtUYPsFV6YS9OJcF9S/weqAq6CXVSynTi2iOfxrFe/Fr96R9+X6VjVYXy+8ZYMuyhNBt+CL+qhSBPFfPEN1hssClQcs4/S1W/GFt0r7fVcsvMBWbb8AgnCJLinP5ycAAAAASUVORK5CYII="
	globals.SymbolImages["GHS06"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAK6SURBVFiFzdhPiFdVFMDxzxsnpTIFNQpJZjEoVEwQyoRYRP82JUQySAYuZITEDBeCuvPnInXhzGAaRTMtDNyMBElUi9m4ECVGI6FVMOuilUG4COO0eBfmj2/m997vd0e7cOG9y7vnfN8555577i0iQpZWFC0Q0coiLyK677SCSL2VQ2ZuqGxweS2V0XL53ZcJLi9URrj8UJnglgcqA1zXUHgc7+EsvsQxbOwWrluoAXyDHXgkjfXhC7zbDVw3UIO4gFWVgjmIjzqF6xTqhWSVNxcVzAq0cKgTuI4CF6cxgr4lhTOG850siMZQSeEr+B29NcBONZVfDba0pXrmPH+Ll7EX7+OxNP56GtuBy1iXxot5P9Lu5xtADWICn6AHn2M3Av9gc/puOo19jXNp7FlM4kJtfQ3cdx39eAPncS4p3IPXUrD34xkcwUZ8hW24hJXJtUO1Flhtn/M9Pk6JdBc+wx8pLfRgP+5gLV7ET7iRIA+keXtxpk74FFEu6RNl2eikRSrQoijeSop/w3Cyzq94CbeTpQaSu99JMXUN2zGOX7ATP0TE3xUK5nH0VEFUtYiYiojpiPgrIkZxFE/jJp7HFA4rF8R3CXImIvZFxPWIuBsRk5VQiyhsnpVL1w0pk+wVPKeMvQ8T7EVlDI5hSzt51a4sFd1nyoUuLYriOJ5Srrgn8LPSnffQi1X4E5twFa9iBquxRunaf/FjREzNEVytt9YKKS00jrH0vi8J2zDnm0ct2AmU8TehjL+VuIUPmqWL9nCFsqxZj5G27pmdt1a52W/FcGcJtj3cIEbxZF2wNO9tfIrVtdNTk4DM0rvaxJcLLkvZkxsua6GYC25ZSutu4Zb1MNKpkgdyfGuq7IEeeOsqfShXBO2UP9RLlcUg/hfXUNVwWaDygGV039zeW10+NmwRLUUx+5yh/QdzLVcJBJ5ddQAAAABJRU5ErkJggg=="
	globals.SymbolImages["GHS07"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAF/SURBVFiFzdghT8NAGMbx/y1g0HMLQWDxJDgCkk+AQhE8BlkDGL4JwczyAXBkCY4vQAhBgiI7xJpwK1fW9nm6ccklNe/7/u7eW9c2xBixjBAKAGIsLPlijPqEIkIsZ+HI6UbZcN6dMu6cv30mnBdlxPlRJlw/KAOuP5QYJycHBsARMAYmwCUwVHHyioELIFbmvbpzchuAmwzsSW2rfDaA0wzsTs4rrwz2M7BruRPyWYBRBnYin10FVcIC8FGB7TWJ/fPXrqAS3KQCGzaOr6kro0rYbYJ6bxufq+/5w4WrBPbQJUfVMZAfgWfjuea6+zC1cht4Y7Zjx55WGg5/idsEDhyoeZiIA7aAXQfqN6wjDjgHvoApi+76Det0CsrAXpi/j+0oqHpYSxzwmKCmwEjNK213AjsEXoFP4MyyWPUsJLh1YMOBagbrkNQR32tyJa7flS/l9a1tsaW+8DYtupJPBIuKr/SjSh3iX3yGyuMsKA/M2L50rpmeNgtC+Lk2jG/Rx4o589viKwAAAABJRU5ErkJggg=="
	globals.SymbolImages["GHS08"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAKJSURBVFiFzdixix1FGADw38QDUbQQqyNewiEYCxNN4ECjkuIaFYLNQQpTCBLkzD8gJsVLOqPGJiBYeGB1SeFhYXEhmCIpTgIpRDAEosWhckXU6AknCUyKN49bn+/tze6OnAsDuzA789v5vvnevg0xRkWOEHogxl6R8WKM3Ru9SEytV2LM0qhiuLIrVXDlyoevEK4sqiCuGApTeBo7SuA6o7AX3yKmdh3PdsV1RU3hzwpq0NbxZBdc1/AtjEAN2lddcq41KsG+qYH9jYfb4lqjEuxWDSxissl4ebCMQbBSg/qp7bjjYbk3c7oGttjpoduiEuz5GthbndKkLaqCu5Qgd/FdOv8FD3XaWF1QCfZawqzhWDo/mnv/2FLUEfUA3kuYL/FyOn87GzZm/rbhm8JH+LmSU6t4pHL9A87gmTa4Hf9+px19hBAmQgjzIYQVnMIGHqt0+QKvVq6ncRAzIYSlEMLZEMKu3PmyQ4kTldWI+lX/MD7GZ5jHJH7ESbyJs/qbYnDP79idF8rc7Tu6mK7hHexJfQLewE7jfxWO5Cd/zvblypiJ/sBzQ32vjun7T1hWudjqCdiH20OTrOJz7MJTqd8reAk3RqCu4fGsCDWqLezHr2mSdcxgFss4nkK4iPcxh8v6myTiazyandONCx8H8Ck+wM3KapxP+Ta43tDfqXM4hweb1MzGha+CuzMUpoWU+MPh+yR3g+XD6nEvpMQfAJaxNIT6EKEpKg9Wj3vR5jv/b/r1a4B6t81KNYPV4w7hL3yfcu+e6ivPf/pnZGvcLC7iAl7vimoOq8dN44kSqHawnEm35RPBVpNv60eVcYj/xWeo0bgiqDKwguGrtgkljhh7Qtg8L3DcB497IINNg8B2AAAAAElFTkSuQmCC"
	globals.SymbolImages["GHS09"] = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI8SURBVFiFzdi7a1VBEMDhb0EjYhFfaCEICopBi4haWUiqqGgQW1OKaBpNIQEVvREi2lhaWCiC/4CSImgv2IgQELSSdKKlMUKEtTgr3jxuPI+9iQMD5+y9u/NjdmbO7IYYoywSQgvE2MqyXoyxudKKxKStHGvmhsoGl9dTGT2Xf/syweWFygiXHyoTXHegMsDVgprjLk5hS4wRenLD1fYUhjGans9jF/rRlwOuMhS2Yz924F0CO4EBjOTyXK0YwUVcwAcEHMJr9OaKudqLYBA/cQx7cT+NB1xrCld5cvLOGA7gqSIRricv7sED9Df1XB1PbcVOXMZnzOA7pvEcg3UTqjNY1QDlKmLSsRX+14NNlXamLlQyeDJBfenw+xmcxQ3cxunSCVYpINmYatUfHUlg3xaN9+MmjqR5uxPkaNltrZbChcFYUi8tmjuOqbIxt25pT7uifMThtvcBPMS84kvwKY2fizE+hhDCBkXGzuBZaUsN4qtPUVQj3uOJIhmGFSXlOO4loEmsr5KhtYIfvRjCUTxS1K8reKnIwAkLt/UHhqqUjUblIkFuxi18TRCv8KsNah5v8UbqRqqVi5JwOKhI/ym8wKzlg3+xzmGyrJ2QjC2U4ox4J72NS2fFEMI27CsdwEtlNsY43Wn9BVIlILNoo494t+CytD254bI2irngutJaN4Xr6mGkrpFVOb5VNbaqB96yRtfkiuBfxtf0UqUTxH9xDbU8XBaoPGAZt69dq3awy0uMLSH8fc4gvwFyuYuihNiCxwAAAABJRU5ErkJggg=="

}

func Test(this js.Value, args []js.Value) interface{} {

	return nil

}

func main() {

	// Common actions for all logged pages.
	jquery.Jq("#table").On("load-success.bs.table",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			search := URLParameters.Get("search")
			if search != "" {
				bstable.NewBootstraptable(jquery.Jq("#table"), nil).ResetSearch(search)
			}
			return nil
		}))

	// Test.
	js.Global().Set("Test", js.FuncOf(Test))

	// Global functions.
	js.Global().Set("Utils_closeEdit", js.FuncOf(jsutils.CloseEdit))
	js.Global().Set("Utils_message", js.FuncOf(jsutils.DisplayMessageWrapper))
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
	js.Global().Set("Product_operateEventsSelect", js.FuncOf(product.OperateEventsSelect))
	js.Global().Set("Product_operateEventsTotalStock", js.FuncOf(product.OperateEventsTotalStock))
	js.Global().Set("Product_getTableData", js.FuncOf(product.GetTableData))
	js.Global().Set("Product_dataQueryParams", js.FuncOf(product.DataQueryParams))
	js.Global().Set("Product_detailFormatter", js.FuncOf(product.DetailFormatter))
	js.Global().Set("Product_empirical_formulaFormatter", js.FuncOf(product.EmpiricalformulaFormatter))
	js.Global().Set("Product_twod_formulaFormatter", js.FuncOf(product.TwodformulaFormatter))
	js.Global().Set("Product_nameFormatter", js.FuncOf(product.NameFormatter))
	js.Global().Set("Product_cas_numberFormatter", js.FuncOf(product.CasnumberFormatter))
	js.Global().Set("Product_product_specificityFormatter", js.FuncOf(product.Product_productSpecificityFormatter))
	js.Global().Set("Product_productslFormatter", js.FuncOf(product.Product_productSlFormatter))
	js.Global().Set("Product_operateFormatter", js.FuncOf(product.OperateFormatter))
	js.Global().Set("Product_addProducer", js.FuncOf(product.AddProducer))
	js.Global().Set("Product_addSupplier", js.FuncOf(product.AddSupplier))

	js.Global().Set("Product_linearToEmpirical", js.FuncOf(product.LinearToEmpirical))
	js.Global().Set("Product_noCas", js.FuncOf(product.NoCas))
	js.Global().Set("Product_noEmpiricalFormula", js.FuncOf(product.NoEmpiricalFormula))
	js.Global().Set("Product_magic", js.FuncOf(product.Magic))
	js.Global().Set("Product_howToMagicalSelector", js.FuncOf(product.HowToMagicalSelector))

	js.Global().Set("Product_pubchemSearch", js.FuncOf(product.PubchemSearch))
	js.Global().Set("Product_pubchemGetCompoundByName", js.FuncOf(product.PubchemGetCompoundByName))
	js.Global().Set("Product_pubchemGetProductByName", js.FuncOf(product.PubchemGetProductByName))
	js.Global().Set("Product_pubchemCreateProduct", js.FuncOf(product.PubchemCreateProduct))
	js.Global().Set("Product_pubchemUpdateProduct", js.FuncOf(product.PubchemUpdateProduct))

	js.Global().Set("Common_search", js.FuncOf(jsutils.Search))
	js.Global().Set("Common_clearSearch", js.FuncOf(jsutils.ClearSearch))

	js.Global().Set("Product_saveProduct", js.FuncOf(product.SaveProduct))

	// Product page load callbacks.
	js.Global().Set("Product_listBookmark", js.FuncOf(product.Product_listBookmarkCallback))
	js.Global().Set("Product_list", js.FuncOf(product.Product_listCallback))
	js.Global().Set("Product_create", js.FuncOf(product.ProductCreateCallbackWrapper))

	// Storage bootstraptable functions.
	js.Global().Set("Storage_getTableData", js.FuncOf(storage.GetTableData))
	js.Global().Set("Storage_dataQueryParams", js.FuncOf(storage.DataQueryParams))
	js.Global().Set("Storage_detailFormatter", js.FuncOf(storage.DetailFormatter))
	js.Global().Set("Storage_productFormatter", js.FuncOf(storage.Storage_productFormatter))
	js.Global().Set("Storage_batchnumberFormatter", js.FuncOf(storage.Storage_batchnumberFormatter))
	js.Global().Set("Storage_modificationdateFormatter", js.FuncOf(storage.Storage_modificationdateFormatter))
	js.Global().Set("Storage_store_locationFormatter", js.FuncOf(storage.Storage_storelocationFormatter))
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
	js.Global().Set("StoreLocation_fullPathFormatter", js.FuncOf(storelocation.StoreLocationFullPathFormatter))
	js.Global().Set("StoreLocation_saveStoreLocation", js.FuncOf(storelocation.SaveStoreLocation))
	js.Global().Set("StoreLocation_dataQueryParams", js.FuncOf(storelocation.DataQueryParams))

	// StoreLocation page load callbacks.
	js.Global().Set("StoreLocation_list", js.FuncOf(storelocation.StoreLocation_listCallback))
	js.Global().Set("StoreLocation_create", js.FuncOf(storelocation.StoreLocation_createCallBack))

	// Person bootstraptable functions.
	js.Global().Set("Person_operateEventsEdit", js.FuncOf(person.OperateEventsEdit))
	js.Global().Set("Person_operateEventsDelete", js.FuncOf(person.OperateEventsDelete))
	js.Global().Set("Person_getTableData", js.FuncOf(person.GetTableData))
	js.Global().Set("Person_savePerson", js.FuncOf(person.SavePerson))
	js.Global().Set("Person_dataQueryParams", js.FuncOf(person.DataQueryParams))
	js.Global().Set("Person_operateFormatter", js.FuncOf(person.OperateFormatter))
	js.Global().Set("Person_selectAllEntity", js.FuncOf(person.SelectAllEntity))

	// Person page load callbacks.
	js.Global().Set("Person_list", js.FuncOf(person.Person_listCallback))
	js.Global().Set("Person_create", js.FuncOf(person.Person_createCallBack))

	// Welcome announce
	js.Global().Set("WelcomeAnnounce_saveWelcomeAnnounce", js.FuncOf(welcomeannounce.SaveWelcomeAnnounce))
	js.Global().Set("WelcomeAnnounce_list", js.FuncOf(welcomeannounce.WelcomeAnnounce_listCallback))

	// Login
	js.Global().Set("Login_getAnnounce", js.FuncOf(login.GetAnnounce))

	// About
	js.Global().Set("About_list", js.FuncOf(about.About_listCallback))

	// Menu
	js.Global().Set("Menu_loadContent", js.FuncOf(menu.LoadContentWrapper))

	jquery.Jq("#loading").Remove()
	jquery.Jq("div.container").RemoveClass("invisible")

	// Startup messages
	jsutils.DisplaySuccessMessage(locales.Translate("wasm_loaded", HTTPHeaderAcceptLanguage))
	message := URLParameters.Get("message")
	if message != "" {
		jsutils.DisplaySuccessMessage(message)
	}

	// Load login page.
	loginCallbackWrapper := func(args ...interface{}) {
		login.Login_listCallback(js.Null(), nil)
	}
	jsutils.LoadContent("div#content", "login", fmt.Sprintf("%slogin", ApplicationProxyPath), loginCallbackWrapper)

	fmt.Println("wasm loaded")

	keepAlive()

}
