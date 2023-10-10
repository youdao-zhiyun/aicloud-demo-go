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
	// 1、请求version服务
	// 添加请求参数
	paramsMap1 := createRequestParams1()
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap1)
	// 请求api服务
	version := utils.DoGet("https://openapi.youdao.com/textEmbedding/queryTextEmbeddingVersion", header, paramsMap1, "application/json")
	// 打印返回结果
	if version != nil {
		fmt.Println("version: " + string(version))
	}

	// 2、请求embedding服务
	// 添加请求参数
	paramsMap2 := createRequestParams2()
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap2)
	// 请求api服务
	embedding := utils.DoPost("https://openapi.youdao.com/textEmbedding/queryTextEmbeddings", header, paramsMap2, "application/json")
	// 打印返回结果
	if embedding != nil {
		fmt.Println("embedding: " + string(embedding))
	}
}

func createRequestParams1() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		注: q参数本身非必传, 但是计算签名时需要用作空字符串处理""
	*/
	q := ""

	return map[string][]string{
		"q": {q},
	}
}

func createRequestParams2() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		注: 包含16个q时效果最佳, 每个q不要超过500个token
	*/
	q1 := "待输入文本1"
	q2 := "待输入文本2"
	q3 := "待输入文本3"
	// q4...

	return map[string][]string{
		"q": {q1, q2, q3},
	}
}
