//go:build go1.24 && js && wasm

package types

import "github.com/tbellembois/gochimitheque/models"

type WelcomeAnnounce struct {
	*models.WelcomeAnnounce
}
