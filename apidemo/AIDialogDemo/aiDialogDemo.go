package main

import (
	"bufio"
	utils2 "demo/apidemo/utils"
	"demo/apidemo/utils/authv3"
	"encoding/json"
	"fmt"
	"os"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

type RequestDTO struct {
	TaskId    string                 `json:"taskId"`
	UserLevel string                 `json:"userLevel"`
	Scene     map[string]interface{} `json:"scene"`
	History   []HistoryItem          `json:"history"`
}

type HistoryItem struct {
	Speaker string `json:"speaker"`
	Content string `json:"content"`
	Voice   string `json:"voice"`
}

var globalRequestDTO RequestDTO

/*
 * 网易有道智云AI口语老师服务api调用demo
 * api接口: https://openapi.youdao.com
 */
func main() {

	// 1、 获取默认场景接口，并使用第一个场景关键词
	fmt.Println("request get default topics: ")
	topic := getDefaultTopicHelper()
	if topic == "" {
		fmt.Println("get default topics error")
		return
	}
	fmt.Println()

	// 2、 请求生成场景接口
	fmt.Println("request generate topic: ")
	generateTopicResult := generateTopicHelper(topic)
	if generateTopicResult == nil {
		fmt.Println("generate topic error")
		return
	}
	fmt.Println()

	// 3.1、记录场景参数
	globalRequestDTO.TaskId = generateTopicResult["taskId"].(string)
	globalRequestDTO.Scene = generateTopicResult["scene"].(map[string]interface{})
	/*
	 * note: 将下列变量替换为需要请求的参数
	 * 取值参考文档: https://ai.youdao.com/DOCSIRMA/html/aigc/api/AIkyls/index.html
	 */
	globalRequestDTO.UserLevel = "0"

	// 3.2 请求生成对话接口
	sysWords := ""
	userWords := ""
	fmt.Println("request generate dialog: ")
	dialogIndex := 0
	for dialogIndex < 3 {
		sysWords, globalRequestDTO.History = generateDialogHelper(globalRequestDTO, sysWords, userWords)
		if sysWords == "" {
			fmt.Println("generate dialog error")
			return
		}
		fmt.Println("系统：", sysWords)
		fmt.Println("用户：")
		reader := bufio.NewReader(os.Stdin)
		userWords, _ = reader.ReadString('\n')
		dialogIndex++
		fmt.Println()
	}

	// 4、 请求生成推荐接口
	fmt.Println("request generate recommendation: ")
	generateRecommendationHelper(globalRequestDTO)
	fmt.Println()

	// 5、 请求生成报告接口
	fmt.Println("request generate report: ")
	generateReportHelper(globalRequestDTO)
}

func generateReportHelper(globalRequestDTO RequestDTO) string {
	// 添加鉴权相关参数
	paramsMap := authv3.AddAuthParamsWithQ(appKey, appSecret, globalRequestDTO.TaskId)
	// 添加音频参数示例
	voice, err := utils2.ReadFileAsBase64("音频路径")
	if err != nil {
		return ""
	}
	globalRequestDTO.History[0].Voice = voice
	// 添加请求参数
	paramsMap["taskId"] = globalRequestDTO.TaskId
	paramsMap["userLevel"] = globalRequestDTO.UserLevel
	paramsMap["scene"] = globalRequestDTO.Scene
	paramsMap["history"] = globalRequestDTO.History

	params, err := json.Marshal(paramsMap)
	if err != nil {
		fmt.Println("generate report params error")
		return ""
	}
	fmt.Println(string(params))
	fmt.Println()

	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/json;charset=utf-8"},
	}
	result := utils2.DoPostWithJson("http://openapi.youdao.com/ai_dialog/generate_report", header, params, "application/json")
	if result == nil {
		return ""
	}
	fmt.Println(string(result))
	return string(result)
}

func generateRecommendationHelper(requestDto RequestDTO) string {
	// 添加鉴权相关参数
	paramsMap := authv3.AddAuthParamsWithQ(appKey, appSecret, requestDto.TaskId)
	// 添加请求参数
	/*
	 * note: 将下列变量替换为需要请求的参数
	 * 取值参考文档: https://ai.youdao.com/DOCSIRMA/html/aigc/api/AIkyls/index.html
	 */
	indexArr := []string{"1"}
	count := "1"

	paramsMap["taskId"] = requestDto.TaskId
	paramsMap["userLevel"] = requestDto.UserLevel
	paramsMap["scene"] = requestDto.Scene
	paramsMap["history"] = requestDto.History
	paramsMap["indexArr"] = indexArr
	paramsMap["count"] = count

	params, err := json.Marshal(paramsMap)
	if err != nil {
		fmt.Println("generate recommendation params error")
		return ""
	}
	fmt.Println(string(params))
	fmt.Println()
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/json;charset=utf-8"},
	}
	result := utils2.DoPostWithJson("http://openapi.youdao.com/ai_dialog/generate_recommendation", header, params, "application/json")
	if result == nil {
		return ""
	}
	fmt.Println(string(result))
	return string(result)
}

func generateDialogHelper(requestDto RequestDTO, sysWords string, userWords string) (string, []HistoryItem) {
	// 添加鉴权相关参数
	paramsMap := authv3.AddAuthParamsWithQ(appKey, appSecret, requestDto.TaskId)
	// 添加请求参数
	paramsMap["taskId"] = requestDto.TaskId
	paramsMap["userLevel"] = requestDto.UserLevel
	paramsMap["scene"] = requestDto.Scene
	history := make([]HistoryItem, 0)
	if sysWords != "" || userWords != "" {
		var sysItem HistoryItem
		sysItem.Speaker = "System"
		sysItem.Content = sysWords

		var userItem HistoryItem
		userItem.Speaker = "User"
		userItem.Content = userWords

		history = append(requestDto.History, sysItem, userItem)
	}
	paramsMap["history"] = history
	params, err := json.Marshal(paramsMap)
	if err != nil {
		fmt.Println("generate dialog params error")
		return "", nil
	}
	fmt.Println(string(params))
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/json;charset=utf-8"},
	}
	result := utils2.DoPostWithJson("http://openapi.youdao.com/ai_dialog/generate_dialog", header, params, "application/json")
	if result == nil {
		return "", nil
	}
	fmt.Println(string(result))
	fmt.Println()
	var resultData map[string]interface{}
	if err := json.Unmarshal(result, &resultData); err != nil {
		fmt.Println("generage dialog result parse error")
		return "", nil
	}
	code := resultData["code"]
	if code != "0" {
		return "", nil
	}
	resultArr, ok := resultData["data"].(map[string]interface{})["resultArr"].([]interface{})
	if !ok || len(resultArr) == 0 {
		return "", nil
	}
	return resultArr[0].(map[string]interface{})["result"].([]interface{})[0].(string), history
}

func generateTopicHelper(topic string) map[string]interface{} {
	// 添加鉴权相关参数
	paramsMap := authv3.AddAuthParamsWithQ(appKey, appSecret, topic)
	// 添加请求参数
	paramsMap["topic"] = topic
	params, err := json.Marshal(paramsMap)
	if err != nil {
		fmt.Println("generate topic params error")
		return nil
	}
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/json;charset=utf-8"},
	}
	result := utils2.DoPostWithJson("http://openapi.youdao.com/ai_dialog/generate_topic", header, params, "application/json")
	if result == nil {
		return nil
	}
	fmt.Println(string(result))
	fmt.Println()
	var resultData map[string]interface{}
	if err := json.Unmarshal(result, &resultData); err != nil {
		fmt.Println("generate topic result parse error")
		return nil
	}
	code := resultData["code"]
	if code != "0" {
		return nil
	}
	return resultData["data"].(map[string]interface{})
}

func getDefaultTopicHelper() string {
	q := "topics"
	// 添加鉴权相关参数
	paramsMap := authv3.AddAuthParamsWithQ(appKey, appSecret, q)
	// 添加请求参数
	paramsMap["q"] = q
	params, err := json.Marshal(paramsMap)
	if err != nil {
		fmt.Println("get default topic params error")
		return ""
	}
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/json;charset=utf-8"},
	}
	result := utils2.DoPostWithJson("http://openapi.youdao.com/ai_dialog/get_default_topic", header, params, "application/json")
	if result == nil {
		return ""
	}
	fmt.Println(string(result))
	fmt.Println()
	var resultData map[string]interface{}
	if err := json.Unmarshal(result, &resultData); err != nil {
		fmt.Println("get default topic result parse error")
		return ""
	}
	code := resultData["code"]
	if code != "0" {
		return ""
	}
	topics := resultData["data"].(map[string]interface{})["topicList"].([]interface{})[0].(map[string]interface{})["topics"].([]interface{})
	firstTopic := topics[0].(map[string]interface{})["enName"].(string)
	return firstTopic
}
