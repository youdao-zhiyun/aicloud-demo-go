package main

import (
	"demo/apidemo/utils"
	"demo/apidemo/utils/authv3"
	"encoding/json"
	"fmt"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

// 文档路径, 例windows路径：PATH = "C:\\youdao\\doc.pdf"
var path = ""

func main() {
	kbName := "知识库名称"
	kbId := "知识库id"
	url := "上传url地址"
	fileId := "文档id"
	// NOTE: 1、文档管理
	// 1.1 创建知识库
	createKB(kbName)
	// 1.2 删除知识库
	deleteKB(kbId)
	// 1.3 上传文档
	uploadFile(kbId, path)
	// 1.4 上传url
	uploadUrl(kbId, url)
	// 1.5 删除文档
	deleteFile(kbId, fileId)
	// 1.6 知识库列表
	kbList()
	// 1.7 知识库文档列表
	fileList(kbId)

	// NOTE 2、文档问答
	q := "提问问题"
	// 2.1 问答
	chat(kbId, q)
	// 2.2 问答(流式)
	chatStream(kbId, q)
}

// 创建知识库接口: /create_kb
func createKB(knName string) {
	params := map[string][]string{
		"q": {knName},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string `json:"appKey"`
		Curtime  string `json:"curtime"`
		Q        string `json:"q"`
		Salt     string `json:"salt"`
		Sign     string `json:"sign"`
		SignType string `json:"signType"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/create_kb", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 删除知识库接口: /delete_kb
func deleteKB(knId string) {
	params := map[string][]string{
		"q": {knId},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string `json:"appKey"`
		Curtime  string `json:"curtime"`
		Q        string `json:"q"`
		Salt     string `json:"salt"`
		Sign     string `json:"sign"`
		SignType string `json:"signType"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/delete_kb", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 上传文档接口: /upload_file
func uploadFile(knId string, path string) {
	params := map[string][]string{
		"q": {knId},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	header := map[string][]string{
		"Content-Type": {"multipart/form-data"},
	}
	result := utils.DoPostWithFile("https://openapi.youdao.com/q_anything/paas/upload_file", header, params, "file", path, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 上传url接口: /upload_url
func uploadUrl(knId string, url string) {
	params := map[string][]string{
		"q":   {knId},
		"url": {url},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string `json:"appKey"`
		Curtime  string `json:"curtime"`
		Q        string `json:"q"`
		Salt     string `json:"salt"`
		Sign     string `json:"sign"`
		SignType string `json:"signType"`
		Url      string `json:"url"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
		Url:      params["url"][0],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/upload_url", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 删除文档接口: /delete_file
func deleteFile(knId string, fileId string) {
	params := map[string][]string{
		"q":       {knId},
		"fileIds": {fileId},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string   `json:"appKey"`
		Curtime  string   `json:"curtime"`
		Q        string   `json:"q"`
		Salt     string   `json:"salt"`
		Sign     string   `json:"sign"`
		SignType string   `json:"signType"`
		FileIds  []string `json:"fileIds"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
		FileIds:  params["fileIds"],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/delete_file", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 知识库列表接口: /kb_list
func kbList() {
	params := map[string][]string{
		"q": {""},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string `json:"appKey"`
		Curtime  string `json:"curtime"`
		Salt     string `json:"salt"`
		Sign     string `json:"sign"`
		SignType string `json:"signType"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/kb_list", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 查询知识库文档接口: /file_list
func fileList(kbId string) {
	params := map[string][]string{
		"q": {kbId},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string `json:"appKey"`
		Curtime  string `json:"curtime"`
		Q        string `json:"q"`
		Salt     string `json:"salt"`
		Sign     string `json:"sign"`
		SignType string `json:"signType"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/file_list", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 问答接口: /file_list
func chat(kbId string, q string) {
	params := map[string][]string{
		"kbIds": {kbId},
		"q":     {q},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string   `json:"appKey"`
		Curtime  string   `json:"curtime"`
		Q        string   `json:"q"`
		Salt     string   `json:"salt"`
		Sign     string   `json:"sign"`
		SignType string   `json:"signType"`
		KbIds    []string `json:"kbIds"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
		KbIds:    params["kbIds"],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/chat", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}

// 问答接口(流式): /file_list
func chatStream(kbId string, q string) {
	params := map[string][]string{
		"kbIds": {kbId},
		"q":     {q},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, params)
	paramsStruct := struct {
		AppKey   string   `json:"appKey"`
		Curtime  string   `json:"curtime"`
		Q        string   `json:"q"`
		Salt     string   `json:"salt"`
		Sign     string   `json:"sign"`
		SignType string   `json:"signType"`
		KbIds    []string `json:"kbIds"`
	}{
		AppKey:   appKey,
		Curtime:  params["curtime"][0],
		Q:        params["q"][0],
		Salt:     params["salt"][0],
		Sign:     params["sign"][0],
		SignType: params["signType"][0],
		KbIds:    params["kbIds"],
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}
	paramsJson, err := json.Marshal(paramsStruct)
	if err != nil {
		fmt.Println("parse params error")
	}
	println("请求参数: " + string(paramsJson))
	result := utils.DoPostWithJson("https://openapi.youdao.com/q_anything/paas/chat_stream", header, paramsJson, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
	}
}
