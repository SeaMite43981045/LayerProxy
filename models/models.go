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

type SystemInfo struct {
	CPUModel    string  `json:"cpu_model"`
	CPUCores    int     `json:"cpu_cores"`
	CPUThreads  int     `json:"cpu_threads"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryFree  uint64  `json:"memory_free"`
	OSName      string  `json:"os_name"`
	OSVersion   string  `json:"os_version"`
	Uptime      uint64  `json:"uptime"`
}

type UserPreferences struct {
	Language string `json:"language"`
	Theme    string `json:"theme"`
}
