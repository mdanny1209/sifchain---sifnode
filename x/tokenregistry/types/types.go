package types

import "strings"

const (
	// ModuleName is the name of the whitelist module
	ModuleName = "tokenregistry"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route
	QuerierRoute = ModuleName

	// RouterKey is the msg router key
	RouterKey = ModuleName
)

func (r *RegistryEntry) Sanitize() {
	r.Path = strings.Trim(r.Path, "/")
}
