package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/global"
	"regexp"
	"strconv"

	//"github.com/jassue/jassue-gin/app/controllers/app"
	"github.com/jassue/jassue-gin/app/models"
	//"github.com/jassue/jassue-gin/global"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type quarkService struct {
}

var cookie = ""

var QuarkService = new(quarkService)

// 分享资源
func (quarkService *quarkService) ShareItem(fids []string, name string) (url string) {
	cookie := DictService.GetValueByDict("quark", "cookie")
	success, message := share(fids, name, 1, 1, "", cookie)
	if success {
		fmt.Println("拿到的url", message)
		return message
	} else {
		return ""
	}

}

// 分享资源并入库
func (quarkService *quarkService) SaveResouceByUrl(fids []string, name string, data []models.ShareItem, categoryId uint) (url string) {
	url1 := QuarkService.ShareItem(fids, name)
	db := global.App.DB

	var ids []string
	if url1 != "" {
		if len(data) > 0 {
			for i := range data {
				res := &data[i]
				ids = append(ids, res.ID)
				//var tmp models.ResourceItem
				var items []models.NetDiskItem

				// Create a new NetDiskItem
				newItem := models.NetDiskItem{
					Type: 2,
					Url:  url1,
				}

				// Append the new item to the slice
				items = append(items, newItem)

				// Convert the slice to JSON
				jsonData, err := json.MarshalIndent(items, "", "  ")
				if err != nil {
					fmt.Println("err", err.Error())
					continue

				}
				var saveItem = models.ResourceItem{Views: 0, Title: res.Name, DiskItems: string(jsonData), CategoryId: categoryId, Status: 1, CoverImg: ""}
				err1 := db.Create(&saveItem)

				if err1 != nil {
					fmt.Println("err1", err1)
					continue
				}

			}

		}

	}

	return "success"
}

// 查看目录下文件
func (quarkService *quarkService) GetDirInfo(fid string, page int, pageSize int) (res response.DirResponse) {
	client := &QuarkDriveClient{
		Headers: map[string]string{
			// Fill in necessary headers
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Accept":          "application/json, text/plain, */*",
			"Content-Type":    "application/json",
			"Sec-Fetch-Site":  "same-site",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Dest":  "empty",
			"Referer":         "https://pan.quark.cn/",
			"Accept-Encoding": "gzip, deflate, br",
			"Accept-Language": "zh-CN,zh;q=0.9",
		},
	}

	// Call GetDirByFid with desired parameters
	dirResponse, err := client.GetDirByFid(fid, pageSize, page)
	if err != nil {
		fmt.Println("Error:", err)
		return response.DirResponse{Data: []models.FileInfo{}, Code: 500, Status: 500, Total: 0}

	} else {
		return *dirResponse
		//fmt.Println("Directory Response:", dirResponse)
	}

}

type WaitshareItem struct {
	Name string `json:"name"`
	Fid  string `json:"fid"`
}

type ShareResponse struct {
	Status int `json:"status"`
	Code   int `json:"code"`
	Data   struct {
		TaskID  string `json:"task_id"`
		ShareID string `json:"share_id"`
	} `json:"data"`
}

type TaskResponse struct {
	Status int `json:"status"`
	Code   int `json:"code"`
	Data   struct {
		ShareID string `json:"share_id"`
	} `json:"data"`
}

var headers = map[string]string{
	"sec-ch-ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
	"accept":             "application/json, text/plain, */*",
	"content-type":       "application/json",
	"sec-ch-ua-mobile":   "?0",
	"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"sec-ch-ua-platform": `"Windows"`,
	"origin":             "https://pan.quark.cn",
	"sec-fetch-site":     "same-site",
	"sec-fetch-mode":     "cors",
	"sec-fetch-dest":     "empty",
	"referer":            "https://pan.quark.cn/",
	"accept-encoding":    "gzip, deflate, br",
	"accept-language":    "zh-CN,zh;q=0.9",
}

func addCookieHeader(headers map[string]string, cookie string) {
	if _, exists := headers["cookie"]; !exists {
		headers["cookie"] = cookie
	}
}

func createRequest(url, method string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

func executeRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Read the decompressed response body
	return ioutil.ReadAll(reader)
}

func share(fidList []string, title string, expiredType, urlType int, passcode, cookie string) (bool, string) {
	baseURL := "https://drive-pc.quark.cn/1/clouddrive/share"
	queryParams := fmt.Sprintf("?pr=ucpro&fr=pc&uc_param_str=&__dt=%d&__t=%d", rand.Intn(900)+100, time.Now().UnixMilli())

	postData := map[string]interface{}{
		"fid_list":     fidList,
		"title":        title,
		"expired_type": expiredType,
		"url_type":     urlType,
	}
	if passcode != "" {
		postData["passcode"] = passcode
	}
	body, _ := json.Marshal(postData)
	requestJson := ""
	jsonData, err := json.Marshal(postData)
	if err == nil {
		requestJson = string(jsonData)
	}
	addCookieHeader(headers, cookie)
	req, err := createRequest(baseURL+queryParams, "POST", body)
	if err != nil {
		return false, "Request creation failed"
	}

	responseData, err := executeRequest(req)
	if err != nil {
		return false, "Request execution failed"
	}
	fmt.Println(string(responseData))
	var res ShareResponse
	if err = json.Unmarshal(responseData, &res); err != nil || res.Status != 200 || res.Code != 0 {
		return false, "Failed to get task_id"
	}

	taskID := res.Data.TaskID
	for i := 0; i < 4; i++ {
		taskURL := fmt.Sprintf("https://drive-pc.quark.cn/1/clouddrive/task?pr=ucpro&fr=pc&uc_param_str=&task_id=%s&retry_index=%d", taskID, i)
		req, _ = createRequest(taskURL, "GET", nil)
		responseData, _ = executeRequest(req)
		svc := &searchItemService{} // 创建指针类型实例
		var taskRes TaskResponse
		if err := json.Unmarshal(responseData, &taskRes); err == nil && taskRes.Status == 200 && taskRes.Code == 0 {
			if shareID := taskRes.Data.ShareID; shareID != "" {
				if err == nil {
					svc.CreateShareLog(requestJson, string(responseData))
				}

				return getShareURL(shareID)
			}
		}
		time.Sleep(1 * time.Second)
	}
	svc := &searchItemService{} // 创建指针类型实例
	svc.CreateShareLog(requestJson, "分享成功，但获取分享url失败")
	return false, "Failed to retrieve share_id"
}

type Response struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Title    string `json:"title"`
		ShareURL string `json:"share_url"`
	} `json:"data"`
}

func getShareURL(shareID string) (bool, string) {
	shareURL := "https://drive-pc.quark.cn/1/clouddrive/share/password?pr=ucpro&fr=pc&uc_param_str="
	sharePostData := map[string]string{"share_id": shareID}
	body, _ := json.Marshal(sharePostData)

	req, err := createRequest(shareURL, "POST", body)
	if err != nil {
		return false, "Failed to create share URL request"
	}

	responseData, err := executeRequest(req)
	if err != nil {
		return false, "Failed to execute share URL request"
	}
	fmt.Println(string(responseData))
	var shareRes ShareResponse
	if err = json.Unmarshal(responseData, &shareRes); err == nil && shareRes.Code == 0 && shareRes.Status == 200 {
		var resp Response
		if err := json.Unmarshal([]byte(responseData), &resp); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return false, "Failed parse json"
		}

		// Access the fields
		fmt.Println("Title:", resp.Data.Title)
		fmt.Println("Share URL:", resp.Data.ShareURL)
		return true, resp.Data.ShareURL
	}
	return false, "Failed to get share URL"
}

type QuarkDriveClient struct {
	Headers map[string]string
}

func addGetCookieHeader(req *http.Request, cookie string) {
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
}

type Metadata struct {
	Size  int    `json:"_size"`
	ReqID string `json:"req_id"`
	Page  int    `json:"_page"`
	Count int    `json:"_count"`
	Total int    `json:"_total"`
}

type ListResponse struct {
	Status    int       `json:"status"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp int64     `json:"timestamp"`
	Data      DataField `json:"data"`
	Metadata  Metadata  `json:"metadata"`
}

// DataField represents the "data" field in the JSON
type DataField struct {
	LastViewList   []models.FileInfo `json:"last_view_list"`
	RecentFileList []models.FileInfo `json:"recent_file_list"`
	List           []models.FileInfo `json:"list"`
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	var err error

	// Check for gzip encoding
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func (client *QuarkDriveClient) GetDirByFid(fid string, size int, page int) (*response.DirResponse, error) {
	baseURL := "https://drive-pc.quark.cn/1/clouddrive/file/sort"

	// Prepare query parameters
	params := url.Values{}
	params.Set("pr", "ucpro")
	params.Set("fr", "pc")
	params.Set("uc_param_str", "")
	params.Set("pdir_fid", fmt.Sprintf("%s", fid))
	params.Set("_page", fmt.Sprintf("%d", page))
	params.Set("_size", fmt.Sprintf("%d", size))
	params.Set("_fetch_total", "1")
	params.Set("_fetch_sub_dirs", "0")
	params.Set("_sort", "file_type:asc,file_name:asc")

	// Build the full URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Set headers and make the GET request
	cookie = DictService.GetValueByDict("quark", "cookie")
	req, err := http.NewRequest("GET", fullURL, nil)
	addGetCookieHeader(req, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	for key, value := range client.Headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := readResponseBody(resp)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, fmt.Errorf("Error reading response body:", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response
	var res ListResponse
	err2 := json.Unmarshal(body, &res)
	if err2 != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, fmt.Errorf("Error parsing JSON:: %v", err)
	}

	// Output parsed data
	return &response.DirResponse{Data: res.Data.List, Status: 200, Code: 0, Total: res.Metadata.Total}, nil

}

// Helper function to get the current milliseconds
func getMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Placeholder for setCookie function (implement according to your needs)
func (s *QuarkDriveClient) setCookie() {
	// Set the necessary cookies here
}

// Helper to perform a POST request
func curlPost(url string, data []byte) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	cookie1 := DictService.GetValueByDict("quark", "cookie")
	addGetCookieHeader(req, cookie1)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		fmt.Println("err")
		return nil, err
	}
	fmt.Println(string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Extracts pwd_id, pdir_fid, and passcode from URL
func (s *quarkService) getIdFromUrl(urlStr string) (string, string, string) {
	pwdID, pdirFid, passcode := "", "0", ""

	pattern1 := regexp.MustCompile(`/s/(\w+)(#/list/share.*/(\w+))?`)
	matches := pattern1.FindStringSubmatch(urlStr)
	if len(matches) > 1 {
		pwdID = matches[1]
		if len(matches) > 3 {
			pdirFid = matches[3]
		}
	}

	pattern2 := regexp.MustCompile(`提取码[:：](\S+\d{1,4}\S*)`)
	passcodeMatch := pattern2.FindStringSubmatch(urlStr)
	if len(passcodeMatch) > 1 {
		passcode = passcodeMatch[1]
	}

	return pwdID, pdirFid, passcode
}

// Retrieves stoken using pwd_id and passcode
func (s *quarkService) getStoken(pwdID, passcode string) (bool, string) {
	baseUrl := "https://pan.quark.cn/1/clouddrive/share/sharepage/token"
	query := url.Values{
		"pr":           {"ucpro"},
		"fr":           {"pc"},
		"uc_param_str": {""},
		"__dt":         {strconv.Itoa(rand.Intn(900) + 100)},
		"__t":          {fmt.Sprintf("%d", getMilliseconds())},
	}
	urlWithQuery := fmt.Sprintf("%s?%s", baseUrl, query.Encode())

	postData := map[string]string{
		"passcode": passcode,
		"pwd_id":   pwdID,
	}
	postDataBytes, _ := json.Marshal(postData)

	//s.setCookie()
	res, err := curlPost(urlWithQuery, postDataBytes)
	if err != nil || res == nil {
		return false, "Failed to get stoken"
	}

	if status, ok := res["status"].(float64); ok && status == 200 {
		if code, ok := res["code"].(float64); ok && code == 0 {
			if data, ok := res["data"].(map[string]interface{}); ok {
				return true, data["stoken"].(string)
			}
		}
	}

	return false, res["message"].(string)
}

// curlGet performs a GET request
func curlGet(url string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	addCookieHeader(headers, cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}
	body, err := ioutil.ReadAll(reader)
	fmt.Println("tttd:")
	fmt.Println(string(body))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Retrieves file details using pwd_id, stoken, and pdir_fid
func (s *quarkService) getDetail(pwdID, stoken, pdirFid string) []map[string]interface{} {
	var fileList []map[string]interface{}
	page := 1

	for {
		baseurl := "https://pan.quark.cn/1/clouddrive/share/sharepage/detail"
		query := url.Values{
			"pr":            {"ucpro"},
			"fr":            {"pc"},
			"pwd_id":        {pwdID},
			"stoken":        {stoken},
			"pdir_fid":      {pdirFid},
			"force":         {"0"},
			"_page":         {strconv.Itoa(page)},
			"_size":         {"50"},
			"_fetch_banner": {"0"},
			"_fetch_share":  {"0"},
			"_fetch_total":  {"1"},
			"_sort":         {"file_type:asc,updated_at:desc"},
		}
		urlWithQuery := fmt.Sprintf("%s?%s", baseurl, query.Encode())

		//s.setCookie()
		res, err := curlGet(urlWithQuery)
		if err != nil || res == nil {
			break
		}

		if status, ok := res["status"].(float64); ok && status == 200 {
			if code, ok := res["code"].(float64); ok && code == 0 {
				if data, ok := res["data"].(map[string]interface{}); ok {
					list := data["list"].([]interface{})
					for _, item := range list {
						fileList = append(fileList, item.(map[string]interface{}))
					}
					page++
				}
				if meta, ok := res["metadata"].(map[string]interface{}); ok {
					if total, ok := meta["_total"].(float64); ok && len(fileList) >= int(total) {
						break
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return fileList
}

// Main saveShare function
func (s *quarkService) SaveShare(fid, shareURL string) (bool, interface{}) {
	// Retrieve pwd_id, pdir_fid, and passcode from the URL

	pwdID, pdirFid, passcode := s.getIdFromUrl(shareURL) // Implement getIdFromUrl to parse URL
	// Get stoken
	err, stokenRes := s.getStoken(pwdID, passcode) // Implement getStoken to retrieve stoken
	fmt.Println("toekn res")
	fmt.Println(stokenRes)
	if !err {
		return false, "Failed to get stoken"
	}

	stoken := stokenRes

	// Get share file list
	shareFileList := s.getDetail(pwdID, stoken, pdirFid) // Implement getDetail to get file list
	if len(shareFileList) == 0 {
		return false, "Share directory is empty"
	}

	// Prepare fid and fid_token lists
	var fidList, fidTokenList []string
	for _, saveList := range shareFileList {
		fidList = append(fidList, saveList["fid"].(string))
		fidTokenList = append(fidTokenList, saveList["share_fid_token"].(string))
	}
	fmt.Println("finalPost")
	// Set up the request URL with query parameters
	baseURL := "https://drive.quark.cn/1/clouddrive/share/sharepage/save"
	query := url.Values{
		"pr":           {"ucpro"},
		"fr":           {"pc"},
		"uc_param_str": {""},
		"__dt":         {fmt.Sprintf("%d", rand.Intn(900)+100)},
		"__t":          {fmt.Sprintf("%d", getMilliseconds())},
	}

	urlWithQuery := fmt.Sprintf("%s?%s", baseURL, query.Encode())

	// Prepare the POST data
	postData := map[string]interface{}{
		"fid_list":       fidList,
		"fid_token_list": fidTokenList,
		"to_pdir_fid":    fid,
		"pwd_id":         pwdID,
		"stoken":         stoken,
		"pdir_fid":       "0",
		"scene":          "link",
	}
	postDataBytes, _ := json.Marshal(postData)

	res, _ := curlPost(urlWithQuery, postDataBytes)
	fmt.Println(res)
	if !err {

		return false, "存进网盘失败"
	}

	// Process the response data after ensuring no error occurred
	if status, ok := res["status"].(float64); ok && status == 200 {
		if code, ok := res["code"].(float64); ok && code == 0 {
			if data, ok := res["data"].(map[string]interface{}); ok {
				return true, data
			}
		}
	}
	return false, ""
}
