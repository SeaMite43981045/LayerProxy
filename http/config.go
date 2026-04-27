// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigUpdateRequest struct {
	WebPort          string `json:"web_port"`
	PortStartAt      int    `json:"port_start_at"`
	WildcardDomain   string `json:"wildcard_domain"`
	WildcardMainPort string `json:"wildcard_main_port"`
}

func HandleGetConfig(c *gin.Context) {
	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置未加载"})
		return
	}
	config := cfg.(models.ConfigFile)

	c.JSON(http.StatusOK, gin.H{
		"web_port":            config.Server.WebPort,
		"port_start_at":       config.Port.PortStartAt,
		"wildcard_domain":     config.Wildcard.WildcardDomain,
		"wildcard_main_port":  config.Wildcard.WildcardMainPort,
		"enable_wildcard":     config.Wildcard.EnableWildcard,
	})
}

func HandleUpdateConfig(c *gin.Context) {
	var req ConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置未加载"})
		return
	}
	config := cfg.(models.ConfigFile)

	config.Server.WebPort = req.WebPort
	config.Port.PortStartAt = req.PortStartAt
	config.Wildcard.WildcardDomain = req.WildcardDomain
	config.Wildcard.WildcardMainPort = req.WildcardMainPort

	if err := saveConfigToFile(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "配置已保存，部分设置需重启后生效"})
}
