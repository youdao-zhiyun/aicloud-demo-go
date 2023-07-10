package main

import (
	utils2 "demo/apidemo/utils"
	"demo/apidemo/utils/authv3"
	"fmt"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

// 待识别音频, 例windows路径：PATH = "C:\\youdao\\media.wav"
var path = ""

/*
 * 网易有道智云语音识别服务api调用demo
 * api接口: https://openapi.youdao.com/asrapi
 */
func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	result := utils2.DoPost("https://openapi.youdao.com/asrapi", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Print(string(result))
	}
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/tts/api/dyysb/index.html
	*/
	rate := "采样率 推荐16000"
	langType := "语种"

	// 数据的base64编码
	q, err := utils2.ReadFileAsBase64(path)
	if err != nil {
		return nil
	}
	return map[string][]string{
		"q":        {q},
		"langType": {langType},
		"rate":     {rate},
		"format":   {"wav"},
		"channel":  {"1"},
		"docType":  {"json"},
		"type":     {"1"},
	}
}
