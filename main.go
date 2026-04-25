// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"LayerProxy/database"
	"LayerProxy/http"
	"LayerProxy/logger"
	"LayerProxy/models"
	"LayerProxy/proxy"
	"LayerProxy/setup"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var cfg models.ConfigFile

func main() {
	logger.InitLogFile()
	setup.InitFiles()

	logger.Info("LayerProxy 正在启动...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	rootCtx, rootCancel := context.WithCancel(context.Background())

	configFile, err := os.ReadFile("./config/config.json")
	if err != nil {
		logger.Error("无法加载配置文件 ./config/config.json:", err.Error())
		return
	}

	if err := json.Unmarshal(configFile, &cfg); err != nil {
		logger.Error("配置文件格式解析失败:", err.Error())
		return
	}

	database.InitDB()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartAPI(rootCtx, cfg)
	}()

	servers := database.GetServersFromDB()

	if len(servers) == 0 {
		logger.Warning("当前数据库中未保存任何实例！仅启动 Web 服务器。")
	} else {
		if cfg.Wildcard.EnableWildcard {
			wg.Go(func() {
				proxy.StartWildcardServer(rootCtx, cfg.Wildcard.WildcardMainPort, servers)
			})
		} else {
			logger.Info("LayerProxy 以 [Port] 模式启动")
			for i, server := range servers {
				if server.Name == "" {
					return
				}

				wg.Add(1)
				assignedAddr := fmt.Sprintf(":%d", cfg.Port.PortStartAt+i)
				go func(s models.ProxyInstance, addr string) {
					defer wg.Done()
					proxy.StartPortServer(rootCtx, addr, s)
				}(server, assignedAddr)
			}

			logger.Info("正在自检...")
			for _, inst := range servers {
				if inst.Name == "" {
					return
				}

				_, err := net.Dial("tcp", inst.BackendIP)
				if err != nil {
					logger.Warning(fmt.Sprintf("无法连接到 Minecraft 服务器 %s: %s", inst.BackendIP, err.Error()))
				} else {
					logger.Info(fmt.Sprintf("成功连接到 Minecraft 服务器 %s", inst.BackendIP))
				}
			}
		}
	}

	logger.Info("LayerProxy 系统初始化完成，运行中... (Crtl+C 退出)")

	sig := <-sigChan
	logger.Info(fmt.Sprintf("接收到退出信号 [%v]，正在释放资源...", sig))

	rootCancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("所有实例已安全结束。")
	case <-time.After(5 * time.Second):
		logger.Info("强制退出：部分连接未能在规定时间内关闭。")
	}

	logger.Info("LayerProxy 已完全退出。")
	os.Exit(0)
}
