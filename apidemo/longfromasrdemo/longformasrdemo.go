package main

import (
	utils2 "demo/apidemo/utils"
	"demo/apidemo/utils/authv4"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

// 待转写的音频路径, 例windows路径：PATH = "C:\\youdao\\media.wav"
var path = ""

// 单个分片的大小
var sliceSize = 10485760

type CommonResult struct {
	ErrorCode string `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    string `json:"result"`
}

type GetProcessResult struct {
	ErrorCode string          `json:"errorCode"`
	Msg       string          `json:"msg"`
	Result    []ResultContent `json:"result"`
}

type ResultContent struct {
	Status string `json:"status"`
	TaskId string `json:"taskId"`
}

/*
 * 网易有道智云长语音转写服务api调用demo
 */
func main() {
	// 1、 预处理接口
	prepareResult := prepareHelper()
	if prepareResult == "" {
		fmt.Println("prepare failed")
		return
	}

	// 1.1、 解析预处理返回结果
	var prepareResultJson CommonResult
	err := json.Unmarshal([]byte(prepareResult), &prepareResultJson)
	if err != nil {
		fmt.Println("prepareResult 解析失败")
		return
	}
	prepareErrorCode := prepareResultJson.ErrorCode
	if prepareErrorCode != "0" {
		fmt.Println(prepareResult)
		fmt.Println("prepare failed")
		return
	}
	taskId := prepareResultJson.Result

	// 2、 对文件进行分片并上传
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open file filed: ", err)
		return
	}
	defer file.Close()
	buffer := make([]byte, sliceSize)

	for i := 1; ; i++ {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("file reader failed: ", err)
			return
		}

		sliceFilePath := "D:\\longFormAsrDemo-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(i) + ".wav"
		chunkFile, err := os.Create(sliceFilePath)
		if err != nil {
			fmt.Println("create file error: ", err)
			return
		}

		_, err = chunkFile.Write(buffer[:n])
		if err != nil {
			fmt.Println("write temp file error: ", err)
			return
		}

		/*
			note: 将下列变量替换为需要请求的参数:临时文件的保存地址及格式
		*/
		uploadResult := uploadHelper(taskId, i, sliceFilePath)
		chunkFile.Close()
		err = os.Remove(sliceFilePath)
		if err != nil {
			fmt.Print(err)
		}
		if uploadResult == "" {
			fmt.Println("upload failed")
			return
		}

		var uploadResultJson CommonResult
		err = json.Unmarshal([]byte(uploadResult), &uploadResultJson)
		if err != nil {
			fmt.Println("uploadResult 解析失败")
			return
		}
		uploadErrorCode := uploadResultJson.ErrorCode
		if uploadErrorCode != "0" {
			fmt.Println(uploadResult)
			fmt.Println("upload failed")
			return
		}
		fmt.Println("分片" + strconv.Itoa(i) + "上传成功")
	}

	// 3、 合并文件
	mergeResult := mergeHelper(taskId)
	var mergeResultJson CommonResult
	err = json.Unmarshal([]byte(mergeResult), &mergeResultJson)
	if err != nil {
		fmt.Println("mergeResult 解析失败")
		return
	}
	mergeErrorCode := mergeResultJson.ErrorCode
	if mergeErrorCode != "0" {
		fmt.Println(mergeErrorCode)
		fmt.Println("merge failed")
		return
	}
	fmt.Println("合并成功")

	// 4、查看转写进度
	for {
		fmt.Println("sleep a while...")
		time.Sleep(20 * time.Second)
		getProcessResult := getProcessHelper(taskId)
		var getProcessResultJson GetProcessResult
		err = json.Unmarshal([]byte(getProcessResult), &getProcessResultJson)
		if err != nil {
			fmt.Println("getProcessResult 解析失败")
			continue
		}
		getProcessErrorCode := getProcessResultJson.ErrorCode
		if getProcessErrorCode != "0" {
			fmt.Println(getProcessResult)
			fmt.Println("get process failed")
			continue
		}
		firstElement := getProcessResultJson.Result[0]
		taskStatus := firstElement.Status
		if taskStatus == "9" {
			fmt.Println("任务完成！")
			break
		}
	}

	// 5、 获取结果
	getResultResult := getResultHelper(taskId)
	fmt.Println("转写结果：", getResultResult)
}

func prepareHelper() string {
	// 添加相关参数
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(nil)
		return ""
	}

	fileName := fi.Name()
	fileSize := fi.Size()

	arr := strings.Split(fileName, ".")
	name := arr[0]
	format := arr[1]
	sliceNum := int(math.Ceil(float64(fileSize) / float64(sliceSize)))

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档：https://ai.youdao.com/DOCSIRMA/html/tts/api/cyyzx/index.html
	*/
	langType := "待转写的音频文件的语种"

	paramsMap := map[string][]string{
		"name":     {name},
		"langType": {langType},
		"format":   {format},
		"sliceNum": {strconv.Itoa(sliceNum)},
		"fileSize": {strconv.Itoa(int(fileSize))},
		"type":     {"1"},
	}

	// 添加鉴权相关的参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	result := utils2.DoPost("https://openapi.youdao.com/api/audio/prepare", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
		return string(result)
	}
	return ""
}

func uploadHelper(taskId string, sliceId int, tempFilePath string) string {
	// 添加相关参数
	paramsMap := map[string][]string{
		"q":       {taskId},
		"sliceId": {strconv.Itoa(sliceId)},
		"type":    {"1"},
	}

	// 添加鉴权相关的参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"multipart/form-data"},
	}
	result := utils2.DoPostWithFile("https://openapi.youdao.com/api/audio/upload", header, paramsMap, "file", tempFilePath, "application/json")
	// 打印返回结果
	if result != nil {
		return string(result)
	}
	return ""
}

func mergeHelper(taskId string) string {
	// 添加相关参数
	paramsMap := map[string][]string{
		"q": {taskId},
	}
	// 添加鉴权相关的参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	result := utils2.DoPost("https://openapi.youdao.com/api/audio/merge", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		return string(result)
	}
	return ""
}

func getProcessHelper(taskId string) string {
	// 添加相关参数
	paramsMap := map[string][]string{
		"q": {taskId},
	}
	// 添加鉴权相关的参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	result := utils2.DoPost("https://openapi.youdao.com/api/audio/get_progress", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
		return string(result)
	}
	return ""
}

func getResultHelper(taskId string) string {
	// 添加相关参数
	paramsMap := map[string][]string{
		"q": {taskId},
	}
	// 添加鉴权相关的参数
	authv4.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	result := utils2.DoPost("https://openapi.youdao.com/api/audio/get_result", header, paramsMap, "application/json")
	// 打印返回结果
	if result != nil {
		fmt.Println(string(result))
		return string(result)
	}
	return ""
}
