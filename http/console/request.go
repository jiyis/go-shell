package console

import (
	"bytes"
	"github.com/gin-gonic/gin/json"
	"io/ioutil"
	"net/http"
	. "shell-exec/lib"
	"time"
)

// 待请求的任务
type task struct {
	Url    string
	Method string
	Params map[string]interface{}
}

var client *http.Client

//var Task *task

func init() {
	client = &http.Client{
		Timeout: time.Second * 30,
	}
}

// 实现handle接口，用于投递任务
func (t *task) Handle(i interface{}) error {
	request(t.Url, t.Method, t.Params)
	return nil
}

/**
请求es api，写入搜索
*/
func request(url string, method string, params map[string]interface{}) error {

	bytesParams, err := json.Marshal(params)
	if err != nil {
		Log.Error("unable to encode json:", err.Error())
		return nil
	}
	body := bytes.NewReader(bytesParams)

	request, err := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-App-Key", SignGlobalJwt())

	if err != nil {
		Log.Error(err.Error())
		return nil
	}

	response, err := client.Do(request)

	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}

	content := string(respBytes)

	if check, _ := InArray(response.StatusCode, []int{200, 201}); !check {
		Log.Error("error to request es api", content, response.StatusCode)
	} else {
		Log.Info("success request", content)
	}

	return nil
}
