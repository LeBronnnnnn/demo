package main

import (
	"UdpFileSender/common"
	"encoding/json"
	"log"
	"math"
	"net"
	"os"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 服务器监听UDP，本地主机的9991端口
	addr, err := net.ResolveUDPAddr("udp", "192.168.11.3:9991")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP port: %v", err)
	}

	for {
		var buf = make([]byte, 1024)
		// 读取客户端请求（读到的字节数、发送方地址、错误）
		n, addr, err := conn.ReadFromUDP(buf[0:])
		// 处理数据
		go handleClient(n, addr, buf, err, conn)
	}
}

func handleClient(n int, addr *net.UDPAddr, buf []byte, err error, conn *net.UDPConn) {
	if err != nil {
		log.Printf("Error reading from UDP connection: %v", err)
		return
	}
	var fileReq common.FileRequest           //创建请求报文（保存请求数据）
	err = json.Unmarshal(buf[0:n], &fileReq) //读到的请求数据，转化为json
	if err != nil {
		log.Printf("Error unmarshaling JSON request: %v", err)
		return
	}

	// 打开文件
	file, err := os.Open("./output.bin")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return
	}
	var fileResp common.FileResponse
	// 处理关于文件大小的请求
	if fileReq.Start == math.MaxInt64 {
		//生成响应报文（文件size）
		fileResp = common.FileResponse{
			Content: nil,
			Start:   fileInfo.Size(),
			End:     fileInfo.Size(),
		}
	} else { // 处理关于文件内容的请求
		file.Seek(int64(fileReq.Start), 0)
		// 相关文件内容保存到buff中
		content := make([]byte, fileReq.End-fileReq.Start)
		_, err := file.Read(content)
		if err != nil {
			log.Printf("Error reading file content from(%d) - to(%d): %v", fileReq.Start, fileReq.End, err)
			return
		}
		// 生成响应报文（文件内容）
		fileResp = common.FileResponse{
			Content: content,
			Start:   fileReq.Start,
			End:     fileReq.End,
		}
	}
	resp, err := json.Marshal(fileResp) // 响应报文转为json
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		return
	}
	_, err = conn.WriteToUDP(resp, addr) // 发送响应
	if err != nil {
		log.Printf("Error writing response to UDP connection: %v", err)
		return
	}
}
