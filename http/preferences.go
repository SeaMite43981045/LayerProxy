// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/database"
	"LayerProxy/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetPreferences(c *gin.Context) {
	var prefs models.UserPreferences
	err := database.DB.QueryRow("SELECT language, theme FROM preferences WHERE id = 1").Scan(&prefs.Language, &prefs.Theme)
	if err != nil {
		prefs.Language = "zh"
		prefs.Theme = "dark"
	}
	c.JSON(http.StatusOK, prefs)
}

func HandleUpdatePreferences(c *gin.Context) {
	var prefs models.UserPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if prefs.Language != "zh" && prefs.Language != "en" {
		prefs.Language = "zh"
	}
	if prefs.Theme != "dark" && prefs.Theme != "light" {
		prefs.Theme = "dark"
	}

	_, err := database.DB.Exec(
		"INSERT INTO preferences (id, language, theme) VALUES (1, ?, ?) ON CONFLICT(id) DO UPDATE SET language=excluded.language, theme=excluded.theme",
		prefs.Language, prefs.Theme,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "偏好设置已保存"})
}
