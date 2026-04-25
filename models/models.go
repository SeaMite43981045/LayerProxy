// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

type ProxyInstance struct {
	Name      string `json:"name"`
	BackendIP string `json:"backend_ip"`
	Subdomain string `json:"subdomain"`
}

type ConfigFile struct {
	Server struct {
		WebPort string `json:"web_port"`
		Key     string `json:"key"`
	} `json:"server"`
	Port struct {
		PortStartAt int `json:"port_start_at"`
	} `json:"port"`
	Wildcard struct {
		EnableWildcard   bool   `json:"enable_wildcard"`
		WildcardDomain   string `json:"wildcard_domain"`
		WildcardMainPort string `json:"wildcard_main_port"`
	} `json:"wildcard"`
}
