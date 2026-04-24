package proxy

import (
	"LayerProxy/logger"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type ProxyInstance struct {
	Name      string `json:"name"`
	BackendIP string `json:"backend_ip"`
	Subdomain string `json:"subdomain"`
}

func StartWildcardServer(listenAddr string, instances []ProxyInstance) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Error(fmt.Sprintf("Wildcard 主端口 %s 监听失败: %s", listenAddr, err.Error()))
		return
	}

	logger.Info(fmt.Sprintf("LayerProxy [Wildcard] 模式启动 | 监听端口: %s", listenAddr))
	for _, inst := range instances {
		logger.Info(fmt.Sprintf(" -> 路由加载: %s -> %s [%s]", inst.Subdomain, inst.BackendIP, inst.Name))
	}

	logger.Info("正在自检...")

	for _, inst := range instances {
		_, err := net.Dial("tcp", inst.BackendIP)
		if err != nil {
			logger.Warning(fmt.Sprintf("无法连接到 Minecraft 服务器 %s: %s", inst.BackendIP, err.Error()))
		} else {
			logger.Info(fmt.Sprintf("成功连接到 Minecraft 服务器 %s", inst.BackendIP))
		}
	}

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleWildcardRouting(clientConn, instances)
	}
}

func StartPortServer(listenAddr string, inst ProxyInstance) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] 端口 %s 监听失败: %s", inst.Name, listenAddr, err.Error()))
		return
	}

	logger.Info(fmt.Sprintf("LayerProxy [Port] 模式启动 | 实例: %s | 监听: %s -> %s", inst.Name, listenAddr, inst.BackendIP))

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleSimpleForward(clientConn, inst)
	}
}

func handleWildcardRouting(clientConn net.Conn, instances []ProxyInstance) {
	defer clientConn.Close()

	playerName, isLogin, connectedDomain, initialData, err := identifyPlayer(clientConn)
	if err != nil {
		return
	}

	var target *ProxyInstance
	for _, inst := range instances {
		if strings.HasPrefix(connectedDomain, inst.Subdomain+".") || connectedDomain == inst.Subdomain {
			target = &inst
			break
		}
	}

	if target == nil {
		return
	}

	logPlayer(playerName, connectedDomain, target.BackendIP, isLogin, true)

	executeForward(clientConn, target.BackendIP, initialData)

	if isLogin {
		logger.Info(fmt.Sprintf("Proxy: %s (%s) -> Disconnect", playerName, connectedDomain))
	}
}

func handleSimpleForward(clientConn net.Conn, inst ProxyInstance) {
	defer clientConn.Close()

	playerName, isLogin, connectedDomain, initialData, err := identifyPlayer(clientConn)
	if err != nil {
		return
	}

	logPlayer(playerName, connectedDomain, inst.BackendIP, isLogin, false)

	executeForward(clientConn, inst.BackendIP, initialData)

	if isLogin {
		logger.Info(fmt.Sprintf("Proxy: %s -> Disconnect", playerName))
	}
}

func logPlayer(name, domain, target string, isLogin bool, showDomain bool) {
	displayName := name
	if showDomain {
		displayName = fmt.Sprintf("%s (%s)", name, domain)
	}

	if isLogin {
		logger.Info(fmt.Sprintf("Proxy: %s -> %s", displayName, target))
	} else {
		logger.Info(fmt.Sprintf("Proxy: %s -> Refresh (%s)", displayName, target))
	}
}

func executeForward(clientConn net.Conn, backendIP string, initialData *bytes.Buffer) {
	serverConn, err := net.Dial("tcp", backendIP)
	if err != nil {
		logger.Warning(fmt.Sprintf("无法连接到 Minecraft 服务器 %s: %s", backendIP, err.Error()))
		return
	}
	defer serverConn.Close()

	serverConn.Write(initialData.Bytes())

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(serverConn, clientConn)
		serverConn.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(clientConn, serverConn)
		clientConn.Close()
	}()
	wg.Wait()
}

func identifyPlayer(conn net.Conn) (string, bool, string, *bytes.Buffer, error) {
	fullBuffer := new(bytes.Buffer)
	teeReader := io.TeeReader(conn, fullBuffer)

	// --- Handshake ---
	length, _ := readVarInt(teeReader)
	if length == 0 {
		return "", false, "", nil, fmt.Errorf("EOF")
	}
	packetID, _ := readVarInt(teeReader)
	if packetID != 0x00 {
		return "Unknown", false, "", fullBuffer, nil
	}

	_, _ = readVarInt(teeReader) // Protocol Version

	// 提取域名
	addrLen, _ := readVarInt(teeReader)
	addrBuf := make([]byte, addrLen)
	io.ReadFull(teeReader, addrBuf)
	connectedDomain := strings.Trim(string(addrBuf), "\x00")

	// 端口和状态
	portBuf := make([]byte, 2)
	io.ReadFull(teeReader, portBuf)
	nextState, _ := readVarInt(teeReader) // 1: Status, 2: Login

	if nextState == 1 {
		return "Guest", false, connectedDomain, fullBuffer, nil
	}

	// --- Login Start ---
	loginLen, _ := readVarInt(teeReader)
	if loginLen > 0 {
		loginID, _ := readVarInt(teeReader)
		if loginID == 0x00 {
			nameLen, _ := readVarInt(teeReader)
			nameBuf := make([]byte, nameLen)
			io.ReadFull(teeReader, nameBuf)
			return string(nameBuf), true, connectedDomain, fullBuffer, nil
		}
	}
	return "UnknownPlayer", true, connectedDomain, fullBuffer, nil
}

// readVarInt 处理 MC 协议中的变长整数
func readVarInt(r io.Reader) (int, error) {
	var res int
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return 0, err
		}
		res |= int(b[0]&0x7F) << shift
		if b[0]&0x80 == 0 {
			break
		}
		shift += 7
	}
	return res, nil
}
