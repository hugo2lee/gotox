package requestx

import (
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
		_, err := r.Body.Read(body)
		assert.NoError(t, err)

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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/posts/1", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		assert.NoError(t, jsoniter.NewEncoder(w).Encode(&BaseResponse[User]{
			Data:    User{ID: 1, Title: "mock title", Body: "mock body"},
			Code:    200,
			Message: "ok",
		}))
	}))
	defer mockServer.Close()

	client := resty.New()
	var user BaseResponse[User]
	resp, err := client.R().
		SetResult(&user).
		Get(mockServer.URL + "/posts/1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestPostUser(t *testing.T) {
	user := User{ID: 1, Title: "req title", Body: "req body"}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/posts", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var getUser User
		err = jsoniter.Unmarshal(body, &getUser)
		assert.NoError(t, err)
		assert.Equal(t, getUser, user)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		assert.NoError(t, jsoniter.NewEncoder(w).Encode(&BaseResponse[string]{
			Data:    "success",
			Code:    200,
			Message: "ok",
		}))
	}))
	defer mockServer.Close()

	bYte, err := jsoniter.Marshal(user)
	assert.NoError(t, err)

	var respUser BaseResponse[string]
	resp, err := resty.New().R().
		SetBody(bYte).
		SetResult(&respUser).
		Post(mockServer.URL + "/posts")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}
