// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func HandleCheckUpdate(c *gin.Context) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", utils.RepoPath)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "无法连接到 GitHub"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "GitHub API 返回错误"})
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := utils.Version

	hasUpdate := latestVersion != currentVersion

	c.JSON(http.StatusOK, gin.H{
		"current_version": currentVersion,
		"latest_version":  latestVersion,
		"has_update":      hasUpdate,
		"release_url":     release.HTMLURL,
	})
}
