package console

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"io"
	"os"
	"os/exec"
	. "shell-exec/lib"
	"strings"
	"sync"
)

type NSQHandler struct {
}

/**
启动nsq consumer
*/
func ConsumerLogUpload() {
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go func() {

		defer waiter.Done()

		config := nsq.NewConfig()
		// 设置一个consumer最大接受的消息数量
		config.MaxInFlight = 1

		// 启动多个连接
		for i := 0; i < 5; i++ {
			consumer, err := nsq.NewConsumer("log-upload", "log-upload", config)

			if nil != err {
				Log.Error("init consumer err:", err)
				return
			}

			consumer.AddHandler(&NSQHandler{})

			err = consumer.ConnectToNSQLookupd(Config.NsqLookupHost)
			if nil != err {
				Log.Error("error connect to nsq lookup:", err)
				return
			}
		}
		select {}
	}()

	waiter.Wait()

}

/**
接受到nsq的request id 然后，执行shell脚本，请求es的api接口
*/
func (handler *NSQHandler) HandleMessage(msg *nsq.Message) error {

	// 当前消息的id，防止并发重复
	msgId := string(msg.ID[:len(msg.ID)])
	// 当前请求的request id
	var message map[string]string
	json.Unmarshal([]byte(string(msg.Body)), &message)
	requestId := message["request_id"]

	Log.Info("receive msg from log upload nsq:", message)

	// 执行shell，根据request id过滤出日志
	cmd := exec.Command(Config.ProjectPath+"/filter.sh", requestId, msgId)

	// 获取shell脚本的输出内容
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行shell脚本
	err := cmd.Run()

	if err != nil {
		Log.Error("failed exec filter shell: %s", err)
		return nil
	}

	// 错误和正常输出分开获取
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if len(errStr) != 0 {
		Log.Error("exec filter shell has an error: %s", string(errStr))
		return nil
	}
	// 过滤掉换行和空格，防止在读取时候报错
	files := strings.Split(string(outStr), "\n")

	// 声明一个切片map
	result := make(map[string][]string)

	Log.Info("find log by request id files:", Filter(files))

	// 逐行读取，组合成map，请求接口
	for _, file := range Filter(files) {
		content := readFileLine(file)

		result[Basename(file)] = content
	}

	Log.Info("combine log msg success")

	url := strings.TrimRight(Config.KongUrl, "/") + "/" + Config.LogUploadApi

	for key, messages := range result {

		for _, msg := range messages {
			// 组合出所有的待上传的日志，开始上传，用异步来做
			// 这边由于map是引用类型，所以每次重新初始化一个
			params := make(map[string]interface{})

			params["message"] = msg
			params["type"] = key
			Log.Info("ready to request es api, params:", key, string([]rune(msg)[:10]))
			// Create Job and push the work onto the jobQueue.
			job := &task{url, "POST", params}

			// 入列消费
			JobQueue <- job
		}

	}

	// fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05.158Z"))

	return nil
}

/**
逐行读取文件内容
*/
func readFileLine(file string) []string {

	file = "./" + file
	// 逐行读取文件
	fi, err := os.Open(file)
	if err != nil {
		Log.Error("error for read log line:", err)
	}
	// 读取完成后就删除
	defer fi.Close()

	br := bufio.NewReader(fi)

	content := make([]string, 1)

	// 声明一个slice，用于存储逐行读取的内容
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(string(a)) != 0 {
			content = append(content, string(a))
		}

	}

	// 读取完了就删了
	os.Remove(file)

	return content
}
