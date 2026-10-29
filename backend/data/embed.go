// Package data embeds the generated mock leads into the binary.
// Regenerate leads.json with `make mock`.
package data

import _ "embed"

//go:embed leads.json
var LeadsJSON []byte
