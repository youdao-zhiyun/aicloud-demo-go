package main

import (
	"demo/apidemo/utils"
	"demo/apidemo/utils/authv4"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"time"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

// 识别音频路径, 例windows路径：PATH = "C:\\youdao\\media.wav"
var path = ""

/*
 * 网易有道智云流式语音翻译服务api调用demo
 * api接口: wss://openapi.youdao.com/stream_speech_trans
 */
func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	// 添加鉴权相关参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 创建websocket连接
	ws, wg := utils.InitConnectionWithParams("wss://openapi.youdao.com/stream_speech_trans", paramsMap)
	// 发送流式数据
	sendData(path, 6400, ws)
	wg.Wait()
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/%E5%AE%9E%E6%97%B6%E8%AF%AD%E9%9F%B3%E7%BF%BB%E8%AF%91/API%E6%96%87%E6%A1%A3/%E5%AE%9E%E6%97%B6%E8%AF%AD%E9%9F%B3%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1/%E5%AE%9E%E6%97%B6%E8%AF%AD%E9%9F%B3%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1-API%E6%96%87%E6%A1%A3.html
	*/
	from := "源语言语种"
	to := "目标语言语种"
	rate := "音频数据采样率, 推荐16000"
	format := "音频格式, 推荐wav"

	return map[string][]string{
		"from":    {from},
		"to":      {to},
		"rate":    {rate},
		"format":  {format},
		"channel": {"1"},
		"version": {"v1"},
	}
}

func sendData(path string, step int, ws *websocket.Conn) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open file error", err)
		return
	}
	defer file.Close()
	bytes := make([]byte, step)
	for {
		d, err := file.Read(bytes[:])
		if err != nil {
			fmt.Println("read file error", err)
		}
		if d > 0 {
			utils.SendBinaryMessage(ws, bytes[0:d])
			time.Sleep(200 * time.Millisecond)
			if d < step {
				break
			}
		} else {
			break
		}
	}
	// 发送结束message
	endMessage := "{\"end\": \"true\"}"
	utils.SendBinaryMessage(ws, []byte(endMessage))
}
