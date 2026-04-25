// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"LayerProxy/models"
	"encoding/json"
	"os"
)

func saveConfigToFile(cfg models.ConfigFile) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("./config/config.json", data, 0644)
}
