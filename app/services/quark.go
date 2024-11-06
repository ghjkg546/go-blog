package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/global"

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

	} else {
		fmt.Println("Share failed:", message)
	}

	return message
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
