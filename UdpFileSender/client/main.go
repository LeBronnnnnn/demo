package main

import (
	"UdpFileSender/common"
	"encoding/json"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const FileBlockSize = 1024 //传输数据块大小
const Concurrency = 50     //并发请求的数量

var fileLock sync.Mutex           //互斥锁
var wg sync.WaitGroup             //等待组
var retriedCount = atomic.Int32{} //原子操作的整型，用于记录重试请求次数
var totalWrote = int64(0)         //原子操作的整型，用于记录已经成功写入文件的数据块数量

// 建立连接、发送请求、解析响应
func sendRequest(fileReq common.FileRequest, server string) common.FileResponse {
	conn, err := net.Dial("udp", server) //建立与服务器的UDP连接
	if err != nil {
		log.Fatalf("Error dialing UDP: %v", err)
	}
	defer conn.Close() // 函数执行完后，关闭连接

	req, err := json.Marshal(fileReq) //:= 	声明并初始化  将请求转化为json
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	_, err = conn.Write(req) //只需要返回值的一部分  发送请求
	if err != nil {
		log.Fatalf("Error writing to UDP connection: %v", err)
	}

	var buf [4096]byte
	n, err := conn.Read(buf[0:]) //读取响应到buf
	if err != nil {
		log.Fatalf("Error reading from UDP connection: %v", err)
	}

	var fileResp common.FileResponse
	err = json.Unmarshal(buf[0:n], &fileResp) //将buf里的json，解析成响应格式
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return fileResp

}

//发送特殊请求获取文件大小--根据文件大小进行分块，创建线程--每个线程请求部分文件，写入本地--

// 每隔5s发送一次特殊请求（获取文件大小），直到得到响应
func getFileSize(server string) int64 {
	fileReq := common.FileRequest{Start: math.MaxInt64, End: math.MaxInt64} //特殊请求，获取文件大小
	response := make(chan common.FileResponse, 1)                           //创建线程间用于传递响应文件的容量为1的管道
	for {
		timer := time.NewTimer(5 * time.Second) //5s定时器
		go func() {                             //新线程执行闭包（发送请求到服务器，将响应发送到管道）
			data := sendRequest(fileReq, server)
			response <- data
		}()

		select {
		case <-timer.C: //如果定时器到时，重新发送请求
			continue
		case result := <-response: //如果响应有结果，返回文件大小
			return result.Start
		}
	}
}

// 按块请求内容，写入本地
func requestFileBlock(region, totalBlock int64, file *os.File, server string) {
	//从块开始索引，每次加并发数，直到最大块数
	for blockId := region; blockId < totalBlock; blockId += Concurrency {
		//当前块的文件开始位置和结束位置
		startFile, endFile := blockId*FileBlockSize, (blockId+1)*FileBlockSize
		//向服务器请求文件内容
		data := requestFileContent(startFile, endFile, server)
		//另一个线程写入本地
		go saveFileBlock(startFile, endFile, totalBlock, file, data)
	}
}

// 每隔1s发送一次内容请求（获取start-end文件内容），直到得到响应
func requestFileContent(start, end int64, server string) []byte {
	fileReq := common.FileRequest{Start: start, End: end} //正常获取文件内容请求
	responseChan := make(chan common.FileResponse, 1)     //创建传递响应的管道
	for {
		timer := time.NewTimer(time.Duration(time.Millisecond * 1000)) //等待时间1s
		//另一个线程发送请求，并将响应写入管道
		go func() {
			data := sendRequest(fileReq, server)
			responseChan <- data
		}()
		select { //收到响应，返回内容
		case res := <-responseChan:
			return res.Content
		case <-timer.C: //超过时间，重试次数+1，重试
			retriedCount.Add(1)
			continue
		}
	}
}

// 写入文件
func saveFileBlock(start, end, totalBlock int64, file *os.File, data []byte) {
	fileLock.Lock()
	defer fileLock.Unlock()
	file.Seek(start, 0) //文件偏移量设为start
	file.Write(data)
	totalWrote += 1 //写入块数+1
	if totalWrote == totalBlock {
		wg.Done() //达到最大块数后，通知线程完成
	}
}

// 并根据文件大小，创建线程，每个线程负责一块
func saveFile(savePath string, server string) {
	fileSize := getFileSize(server)                                      //获取文件大小
	file, _ := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, os.ModePerm) //创建写入文件
	totalBlock := fileSize / FileBlockSize                               //计算所需块数
	actualConcurrency := min(totalBlock, Concurrency)                    //实际并发量（所需块数和最大块数）
	wg.Add(1)                                                            //增加一组
	for i := int64(0); i < actualConcurrency; i++ {
		go requestFileBlock(i, totalBlock, file, server) //每个线程负责其中一部分
	}
	wg.Wait() //等待所有线程完成

	if fileSize%FileBlockSize != 0 { //如果还有不到一个块的文件没有请求
		startOffset := fileSize - fileSize%FileBlockSize          //计算最后一个不完整块的开始位置
		data := requestFileContent(startOffset, fileSize, server) //请求、写入
		saveFileBlock(startOffset, fileSize, totalBlock, file, data)
	}

}

func main() {
	saveFile("./output.bin", "192.168.11.3:9991")
}

func min(a int64, b int64) int64 {
	if a > b {
		return b
	} else {
		return a
	}
}
