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
	ws, wg := utils.InitConnectionWithParams("wss://openapi.youdao.com/stream_asropenapi", paramsMap)
	// 发送流式数据
	sendData(path, 6400, ws)
	wg.Wait()
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		https://ai.youdao.com/DOCSIRMA/html/tts/api/ssyysb/index.html
	*/
	langType := "语种"
	rate := "音频数据采样率, 推荐16000"

	return map[string][]string{
		"langType":{from},
		"rate":    {rate},
		"format":  {"wav"},
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
