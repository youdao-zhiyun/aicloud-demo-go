package main

import (
	"demo/apidemo/utils"
	"demo/apidemo/utils/authv4"
	"fmt"
	"github.com/gorilla/websocket"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

// 语音合成文本
var text = "语音合成文本"

/*
 * 网易有道智云流式语音合成服务api调用demo
 * api接口: wss://openapi.youdao.com/stream_tts
 */
func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	// 添加鉴权相关参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 创建websocket连接
	ws, wg := utils.InitConnectionWithParams("wss://openapi.youdao.com/stream_tts", paramsMap)
	// 发送流式数据
	sendData(ws)
	wg.Wait()
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
	*/
	langType := "语言类型"
	voice := "发言人"
	rate := "音频数据采样率"
	format := "音频格式"
	volume := "音量, 取值0.1-5"
	speed := "语速, 取值0.5-2"

	return map[string][]string{
		"langType": {langType},
		"voice":    {voice},
		"rate":     {rate},
		"format":   {format},
		"volume":   {volume},
		"speed":    {speed},
	}
}

func sendData(ws *websocket.Conn) {
	utils.SendTextMessage(ws, fmt.Sprintf("{\"text\":\"%s\"}", text))
	// 发送结束message
	endMessage := "{\"end\": \"true\"}"
	utils.SendBinaryMessage(ws, []byte(endMessage))
}
