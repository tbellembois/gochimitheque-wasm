//go:build go1.24 && js && wasm

package themes

import "fmt"

type Icon interface {
	ToString() string
}

type MDIcon interface {
	Icon
}

// Material Design icons
// https://materialdesignicons.com/
type mdiIcon struct {
	Face IconFace
	Size IconSize
}

type IconClass string
type IconFace string
type IconSize string

const (
	MDI IconClass = "mdi"

	MDI_OK               IconFace = "mdi-check"
	MDI_STORELOCATION    IconFace = "mdi-docker"
	MDI_SUBSTORELOCATION IconFace = "mdi-relation-many-to-many"
	MDI_PEOPLE           IconFace = "mdi-account-group"
	MDI_EDIT             IconFace = "mdi-pencil-outline"
	MDI_DELETE           IconFace = "mdi-delete-outline"
	MDI_ARCHIVE          IconFace = "mdi-archive-outline"
	MDI_NONE             IconFace = "mdi-border-none-variant"
	MDI_ERROR            IconFace = "mdi-alert-circle-outline"
	MDI_CHECK            IconFace = "mdi-checkbox-marked-outline"
	MDI_NO_CHECK         IconFace = "mdi-checkbox-blank-outline"
	MDI_INFO             IconFace = "mdi-information-outline"
	MDI_CLOSE            IconFace = "mdi-close-box"
	MDI_BUG              IconFace = "mdi-car-brake-alert"
	MDI_NO_BOOKMARK      IconFace = "mdi-bookmark-outline"
	MDI_BOOKMARK         IconFace = "mdi-bookmark"
	MDI_PRODUCT          IconFace = "mdi-tag"
	MDI_STORAGE          IconFace = "mdi-cube-unfolded"
	MDI_OSTORAGE         IconFace = "mdi-cube-scan"
	MDI_STORE            IconFace = "mdi-forklift"
	MDI_RESTRICTED       IconFace = "mdi-hand"
	MDI_COLOR            IconFace = "mdi-format-color-fill"
	MDI_LINK             IconFace = "mdi-link-variant"
	MDI_RADIOACTIVE      IconFace = "mdi-radioactive"
	MDI_BORROW           IconFace = "mdi-hand-okay"
	MDI_UNBORROW         IconFace = "mdi-hand-pointing-left"
	MDI_CLONE            IconFace = "mdi-content-copy"
	MDI_HISTORY          IconFace = "mdi-history"
	MDI_RESTORE          IconFace = "mdi-undo"
	MDI_SHOW_DELETED     IconFace = "mdi-archive-outline"
	MDI_HIDE_DELETED     IconFace = "mdi-archive-alert-outline"
	MDI_CONFIRM          IconFace = "mdi-checkbox-marked-outline"
	MDI_PRINT            IconFace = "mdi-printer"
	MDI_DOWNLOAD         IconFace = "mdi-cloud-download-outline"
	MDI_TOTALSTOCK       IconFace = "mdi-sigma"
	MDI_REMOVEFILTER     IconFace = "mdi-filter-off"
	MDI_ARROW_RIGHT      IconFace = "mdi-arrow-right-thin"
	MDI_VIEW             IconFace = "mdi-eye-outline"
	MDI_PUBCHEM          IconFace = "mdi-alpha-c-circle-outline"

	MDI_16PX IconSize = "mdi-16px"
	MDI_24PX IconSize = "mdi-24px"
	MDI_36PX IconSize = "mdi-36px"
	MDI_48PX IconSize = "mdi-48px"
)

func NewMdiIcon(face IconFace, size IconSize) Icon {

	if face == "" {
		face = MDI_NONE
	}
	if size == "" {
		size = MDI_24PX
	}

	return mdiIcon{Face: face, Size: size}

}

func (i mdiIcon) ToString() string {
	return fmt.Sprintf("%s %s %s", MDI, i.Face, i.Size)
}

func (i IconFace) ToString() string {
	return fmt.Sprintf("%v", i)
}
