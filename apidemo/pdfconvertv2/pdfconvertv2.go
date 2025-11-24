package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	utils "demo/apidemo/utils"
)

// === 请填写你自己的本地文件地址 ===
const LOCAL_FILE_PATH = "请填写你自己的本地文件地址"

// === 请填写你自己的应用ID和密钥 ===
const APP_KEY = "你的应用ID"
const APP_SECRET = "你的应用密钥"

// === 高级版 API 基础URL ===
const BASE_URL = "https://openapi.youdao.com/file_convert/v2"
const UPLOAD_URL = BASE_URL + "/upload"
const QUERY_URL = BASE_URL + "/query"

type UploadResponse struct {
	Code string `json:"code"`
	Data struct {
		Flownumber string `json:"flownumber"`
	} `json:"data"`
}

type QueryResponse struct {
	Code string `json:"code"`
	Data struct {
		Status       int    `json:"status"`
		StatusString string `json:"statusString"`
		ResultUrl    string `json:"resultUrl"`
	} `json:"data"`
}

// truncate 根据文档规则生成 input 字符串
func truncate(q string) string {
	if len(q) <= 20 {
		return q
	}
	return q[:10] + strconv.Itoa(len(q)) + q[len(q)-10:]
}

// sha256Digest 生成 SHA256 签名
func sha256Digest(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// genSign 生成签名所需参数：salt, curtime, sign
func genSign(qOrFlownumber string) (salt, curtime, sign string) {
	salt = getUUID()
	curtime = strconv.FormatInt(time.Now().Unix(), 10)
	inputStr := truncate(qOrFlownumber)
	signStr := APP_KEY + inputStr + salt + curtime + APP_SECRET
	sign = sha256Digest(signStr)
	return salt, curtime, sign
}

// getUUID 生成UUID
func getUUID() string {
	// 简化版本的UUID生成器
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// uploadPDF 上传PDF文件
func uploadPDF(filePath string, targetType string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filePath)
	}

	// 读取文件并转换为base64
	fileBase64, err := utils.ReadFileAsBase64(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file as base64: %v", err)
	}

	salt, curtime, sign := genSign(fileBase64)

	// 准备表单数据 (使用multipart/form-data)
	// 注意：这里使用map[string]string而不是map[string][]string
	formData := map[string]string{
		"appKey":         APP_KEY,
		"salt":           salt,
		"curtime":        curtime,
		"sign":           sign,
		"signType":       "v3",
		"q":              fileBase64,
		"fileName":       filepath.Base(filePath),
		"fileType":       "pdf",
		"targetFileType": targetType,
	}

	fmt.Println("正在以 multipart/form-data 上传文件...")
	// 发送POST请求 (使用multipart/form-data)
	result := utils.DoPostMultipart(UPLOAD_URL, formData, "application/json")
	if result == nil {
		return "", fmt.Errorf("upload request failed")
	}

	fmt.Printf("响应内容: %s\n", string(result))

	// 解析响应
	var resp UploadResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if resp.Code == "0" && resp.Data.Flownumber != "" {
		flownumber := resp.Data.Flownumber
		fmt.Printf("上传成功，任务流水号: %s\n", flownumber)
		return flownumber, nil
	} else {
		return "", fmt.Errorf("上传失败: %s", string(result))
	}
}

// queryTask 查询任务状态
func queryTask(flownumber string) (*QueryResponse, error) {
	salt, curtime, sign := genSign(flownumber)

	// 准备表单数据 (使用multipart/form-data)
	// 注意：这里使用map[string]string而不是map[string][]string
	formData := map[string]string{
		"appKey":     APP_KEY,
		"salt":       salt,
		"curtime":    curtime,
		"sign":       sign,
		"signType":   "v3",
		"flownumber": flownumber,
	}

	fmt.Printf("查询任务状态: %s\n", flownumber)
	// 发送POST请求 (使用multipart/form-data)
	result := utils.DoPostMultipart(QUERY_URL, formData, "application/json")
	if result == nil {
		return nil, fmt.Errorf("query request failed")
	}

	fmt.Printf("响应内容: %s\n", string(result))

	var resp QueryResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if resp.Code != "0" {
		return nil, fmt.Errorf("查询失败: %s", string(result))
	}
	return &resp, nil
}

// waitForResult 轮询任务状态直到完成或失败
func waitForResult(flownumber string, interval int, timeout int) (string, error) {
	fmt.Println("查询任务进度中...")
	startTime := time.Now()
	for time.Since(startTime).Seconds() < float64(timeout) {
		result, err := queryTask(flownumber)
		if err != nil {
			return "", err
		}

		status := result.Data.Status
		statusStr := result.Data.StatusString
		fmt.Printf("任务状态: %d (%s)\n", status, statusStr)

		if status == 4 {
			url := result.Data.ResultUrl
			fmt.Printf("转换完成，下载地址: %s\n", url)
			return url, nil
		} else if status == -2 {
			return "", fmt.Errorf("转换失败")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return "", fmt.Errorf("任务超时未完成")
}

// downloadResult 下载转换后的结果文件
func downloadResult(resultURL string, savePath string) error {
	fmt.Println("正在下载结果...")
	resp, err := http.Get(resultURL)
	if err != nil {
		return fmt.Errorf("failed to download result: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	fmt.Printf("文件已保存到: %s\n", savePath)
	return nil
}

func main() {
	// === 示例使用 ===
	// 输出的Word文件名
	outputFile := "result.docx"

	flownumber, err := uploadPDF(LOCAL_FILE_PATH, "docx")
	if err != nil {
		fmt.Printf("上传PDF失败: %v\n", err)
		return
	}

	resultURL, err := waitForResult(flownumber, 5, 300)
	if err != nil {
		fmt.Printf("等待结果失败: %v\n", err)
		return
	}

	// 保存文件到本地
	err = downloadResult(resultURL, outputFile)
	if err != nil {
		fmt.Printf("下载结果失败: %v\n", err)
		return
	}
}
