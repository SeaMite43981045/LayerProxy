// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"net/http"
	"runtime"
	"time"

	"LayerProxy/models"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func HandleSystemInfo(c *gin.Context) {
	info := models.SystemInfo{}

	// OS info
	info.OSName = runtime.GOOS
	hostInfo, err := host.Info()
	if err == nil {
		info.OSVersion = hostInfo.PlatformVersion
		info.Uptime = hostInfo.Uptime
	}

	// CPU info
	cpuInfos, err := cpu.Info()
	if err == nil && len(cpuInfos) > 0 {
		info.CPUModel = cpuInfos[0].ModelName
		info.CPUCores = int(cpuInfos[0].Cores)
	}
	cpuCounts, err := cpu.Counts(true)
	if err == nil {
		info.CPUThreads = cpuCounts
	}
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercents) > 0 {
		info.CPUUsage = cpuPercents[0]
	}

	// Memory info
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryTotal = memInfo.Total
		info.MemoryUsed = memInfo.Used
		info.MemoryFree = memInfo.Free
	}

	c.JSON(http.StatusOK, info)
}
