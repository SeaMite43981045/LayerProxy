// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package setup

import (
	"LayerProxy/logger"
	"fmt"
	"os"
)

func InitFiles() {
	if _, err := os.Stat("./config"); os.IsNotExist(err) {
		os.Mkdir("./config", 0755)
	}
	if _, err := os.Stat("./config/config.json"); os.IsNotExist(err) {
		file, err := os.Create("./config/config.json")

		if err != nil {
			logger.Error(fmt.Sprintf("初始化配置文件失败: %s", err))
		}
		defer file.Close()

		file.WriteString("{\n    \"server\": {\n        \"web_port\": \"23754\",\n        \"key\": \"\"\n    },\n    \"port\": {\n        \"port_start_at\": 23755\n    },\n    \"wildcard\": {\n        \"enable_wildcard\": false,\n        \"wildcard_domain\": \"\",\n        \"wildcard_main_port\": \"23755\"\n    }\n}\n")
	}
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0755)
	}
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}
}
