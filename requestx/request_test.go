package requestx

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestTableDriven(t *testing.T) {
	bodyStr := `{"title": "foo", "body": "bar", "userId": 1}`
	path1 := "/test1"
	path2 := "/test2"

	// 创建 Mock 服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path1 {
			w.WriteHeader(http.StatusOK)
			body := []byte(bodyStr)
			l, err := w.Write(body)
			assert.NoError(t, err)
			assert.Equal(t, len(body), l)
		} else if r.URL.Path == path2 {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Test 1",
			url:          mockServer.URL + path1,
			expectedCode: http.StatusOK,
			expectedBody: bodyStr,
		},
		{
			name:         "Test 2",
			url:          mockServer.URL + path2,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().Get(tt.url)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, resp.String())
			}
		})
	}
}

func TestPostRequestTableDriven(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		if string(body) == `{"title": "foo", "body": "bar", "userId": 1}` {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "Valid Request",
			body:         `{"title": "foo", "body": "bar", "userId": 1}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid Request",
			body:         `{"title": "invalid"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				Post(mockServer.URL)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())
		})
	}
}

func TestRequest(t *testing.T) {
	type Channel struct {
		ChannelList []struct {
			Channel     string `json:"channel"`
			VersionCode int64  `json:"version_code"`
		} `json:"channel_list"`
	}

	type ChannelList struct {
		Channel           string `json:"channel"`
		OtaUUID           string `json:"ota_uuid"`
		ReleaseStamp      int64  `json:"release_stamp"`
		TargetVersionCode int64  `json:"target_version_code"`
		VersionCode       int64  `json:"version_code"`
		VersionName       string `json:"version_name"`
		VersionType       int64  `json:"version_type"`
	}

	type Datachan struct {
		ChannelList []ChannelList `json:"channel_list"`
	}

	type ResponseChan struct {
		Code    int64    `json:"code"`
		Data    Datachan `json:"data"`
		Message string   `json:"message"`
	}

	var r ResponseChan

	req := NewRequestx()
	assert.NotNil(t, req)
	auth := "MTA6YjNhZGMxN2JlY2EwMDhjMGZiY2MyNTg5MGZmNjY2NWQ="

	req.SetBaseURL("https://v2-fw-test.xag.cn")
	res, err := req.R().
		SetAuthToken(auth).
		SetAuthScheme("Basic").
		SetQueryParams(map[string]string{
			"guid":         "0CEB4B3EB55664723BB4ED423153FF79",
			"access_token": "b23739d9cff3ed106a09ced14b829e4c",
			"area":         "zh-CN",
		}).
		SetBody(Channel{
			ChannelList: []struct {
				Channel     string `json:"channel"`
				VersionCode int64  `json:"version_code"`
			}{
				{
					Channel:     "UAV35",
					VersionCode: 0,
				},
				{
					Channel:     "UAV34",
					VersionCode: 0,
				},
			},
		}).
		SetResult(&r).
		Post("/firmware_system_api/v2.2/appsync/ota/list")
	assert.NoError(t, err)

	fmt.Println(res.StatusCode())

	fmt.Println(res.String())
	fmt.Println(r)
}

// User 结构体用于存储响应数据
type User struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type BaseResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func TestGetUser(t *testing.T) {
	// 创建 Mock 服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法和路径
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/posts/1", r.URL.Path)

		// 设置content type
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 返回模拟的 JSON 响应
		assert.NoError(t, jsoniter.NewEncoder(w).Encode(&BaseResponse[User]{
			Data:    User{ID: 1, Title: "mock title", Body: "mock body"},
			Code:    200,
			Message: "ok",
		}))
	}))
	defer mockServer.Close()

	client := resty.New()

	// 创建 User 变量用于存储响应
	var user BaseResponse[User]

	// 发起 GET 请求并将结果解析到 user 结构体
	resp, err := client.R().
		SetResult(&user). // 使用 SetResult 将响应解析到 user 结构体
		Get(mockServer.URL + "/posts/1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestPostUser(t *testing.T) {
	user := User{ID: 1, Title: "req title", Body: "req body"}
	// 创建 Mock 服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法和路径
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/posts", r.URL.Path)

		// 检查post的body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var getUser User
		err = jsoniter.Unmarshal(body, &getUser)
		assert.NoError(t, err)
		assert.Equal(t, getUser, user)

		// 设置content type
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// 返回模拟的 JSON 响应
		assert.NoError(t, jsoniter.NewEncoder(w).Encode(&BaseResponse[string]{
			Data:    "success",
			Code:    200,
			Message: "ok",
		}))
	}))
	defer mockServer.Close()

	bYte, err := jsoniter.Marshal(user)
	assert.NoError(t, err)

	// 创建 User 变量用于存储响应
	var respUser BaseResponse[string]

	// 发起 POST 请求并将结果解析到 user 结构体
	resp, err := resty.New().R().
		SetBody(bYte).
		SetResult(&respUser). // 使用 SetResult 将响应解析到 user 结构体
		Post(mockServer.URL + "/posts")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}
