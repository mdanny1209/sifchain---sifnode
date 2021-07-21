package utils

import (
	"encoding/json"
	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
	"io/ioutil"
	"path/filepath"
)

func ParseDenoms(dir string) (whitelisttypes.DenomWhitelist, error) {
	var denoms whitelisttypes.DenomWhitelist
	file, err := filepath.Abs(dir)
	if err != nil {
		return denoms, err
	}
	o, err := ioutil.ReadFile(file)
	if err != nil {
		return denoms, err
	}

	err = json.Unmarshal(o, &denoms)
	if err != nil {
		return denoms, err
	}
	return denoms, err
}
