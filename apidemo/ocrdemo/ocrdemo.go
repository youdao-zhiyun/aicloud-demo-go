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
 * 网易有道智云通用OCR服务api调用demo
 * api接口: https://openapi.youdao.com/ocrapi
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
	result := utils2.DoPost("https://openapi.youdao.com/ocrapi", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Print(string(result))
	}
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/%E6%96%87%E5%AD%97%E8%AF%86%E5%88%ABOCR/API%E6%96%87%E6%A1%A3/%E9%80%9A%E7%94%A8%E6%96%87%E5%AD%97%E8%AF%86%E5%88%AB%E6%9C%8D%E5%8A%A1/%E9%80%9A%E7%94%A8%E6%96%87%E5%AD%97%E8%AF%86%E5%88%AB%E6%9C%8D%E5%8A%A1-API%E6%96%87%E6%A1%A3.html
	*/
	langType := "要识别的语言类型"
	detectType := "识别类型"
	angle := "是否进行360角度识别"
	column := "是否按多列识别"
	rotate := "是否需要获得文字旋转角度"
	docType := "json"
	imageType := "1"

	// 数据的base64编码
	img, err := utils2.ReadFileAsBase64(path)
	if err != nil {
		return nil
	}
	return map[string][]string{
		"img":        {img},
		"langType":   {langType},
		"detectType": {detectType},
		"angle":      {angle},
		"column":     {column},
		"rotate":     {rotate},
		"docType":    {docType},
		"imageType":  {imageType},
	}
}
