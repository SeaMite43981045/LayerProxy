package main

import (
	"LayerProxy/logger"
	"LayerProxy/proxy"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

type ConfigFile struct {
	Port struct {
		PortStartAt int `json:"port_start_at"`
	} `json:"port"`
	Wildcard struct {
		EnableWildcard   bool   `json:"enable_wildcard"`
		WildcardDomain   string `json:"wildcard_domain"`
		WildcardMainPort string `json:"wildcard_main_port"`
	} `json:"wildcard"`
}

func main() {
	logger.InitLogFile()
	logger.Info("LayerProxy 正在启动...")

	configFile, err := os.ReadFile("./config/config.json")
	if err != nil {
		logger.Error("无法加载配置文件 ./config/config.json:", err.Error())
		return
	}

	var cfg ConfigFile
	if err := json.Unmarshal(configFile, &cfg); err != nil {
		logger.Error("配置文件格式解析失败:", err.Error())
		return
	}

	proxy.InitDB()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		proxy.StartAPI()
	}()

	servers := proxy.GetServersFromDB()

	if len(servers) == 0 {
		logger.Warning("当前数据库中未保存任何实例！仅启动 API 服务器。")
	} else {
		if cfg.Wildcard.EnableWildcard {
			wg.Add(1)
			go func() {
				defer wg.Done()
				proxy.StartWildcardServer(cfg.Wildcard.WildcardMainPort, servers)
			}()
		} else {
			logger.Info("LayerProxy 以 [Port] 模式启动")
			for i, server := range servers {
				wg.Add(1)
				assignedAddr := fmt.Sprintf(":%d", cfg.Port.PortStartAt+i)
				go func(s proxy.ProxyInstance, addr string) {
					defer wg.Done()
					proxy.StartPortServer(addr, s)
				}(server, assignedAddr)
			}

			logger.Info("正在自检...")
			for _, inst := range servers {
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

	wg.Wait()
}
