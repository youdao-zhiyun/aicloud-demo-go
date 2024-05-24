package main

import (
	"demo/apidemo/utils"
	"demo/apidemo/utils/authv3"
	"fmt"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	events := utils.DoPostBySSE("https://openapi.youdao.com/llm_trans", header, paramsMap)
	for event := range events {
		// 处理接收到的事件
		fmt.Println(event)
	}
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/trans/api/dmxfy/index.html
	*/
	i := "待翻译文本"
	from := "源语种"
	to := "目标语种"

	return map[string][]string{
		"q":    {i},
		"i":    {i},
		"from": {from},
		"to":   {to},
	}
}
