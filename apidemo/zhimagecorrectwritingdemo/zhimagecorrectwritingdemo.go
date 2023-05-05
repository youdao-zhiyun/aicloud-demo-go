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

// 作文图片路径, 例windows路径：path = "C:\\youdao\\media.jpg"
var path = ""

func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	result := utils2.DoPost("https://openapi.youdao.com/correct_writing_cn_image", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Print(string(result))
	}
}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/%E4%BD%9C%E6%96%87%E6%89%B9%E6%94%B9/API%E6%96%87%E6%A1%A3/%E4%B8%AD%E6%96%87%E4%BD%9C%E6%96%87%E6%89%B9%E6%94%B9%EF%BC%88%E5%9B%BE%E5%83%8F%E8%AF%86%E5%88%AB%EF%BC%89/%E4%B8%AD%E6%96%87%E4%BD%9C%E6%96%87%E6%89%B9%E6%94%B9%EF%BC%88%E5%9B%BE%E5%83%8F%E8%AF%86%E5%88%AB%EF%BC%89-API%E6%96%87%E6%A1%A3.html
	*/
	q := "正文文本"
	grade := "作文等级"
	title := "作文标题"
	requirement := "题目要求"

	// 数据的base64编码
	q, err := utils2.ReadFileAsBase64(path)
	if err != nil {
		return nil
	}
	return map[string][]string{
		"q":           {q},
		"grade":       {grade},
		"title":       {title},
		"requirement": {requirement},
	}
}
