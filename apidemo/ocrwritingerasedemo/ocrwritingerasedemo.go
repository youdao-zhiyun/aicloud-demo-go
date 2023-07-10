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

// 待识别图片路径, 例windows路径：PATH = "C:\\youdao\\media.jpg"
var path = ""

/*
 * 网易有道智云手写体擦除服务api调用demo
 * api接口: https://openapi.youdao.com/ocr_writing_erase
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
	result := utils2.DoPost("https://openapi.youdao.com/ocr_writing_erase", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Print(string(result))
	}
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/ocr/api/sxocr/index.html
	*/

	// 是否进行360角度识别 0：不识别，1：识别。
	angle := "0"

	// 数据的base64编码
	q, err := utils2.ReadFileAsBase64(path)
	if err != nil {
		return nil
	}
	return map[string][]string{
		"q":     {q},
		"angle": {angle},
	}
}
