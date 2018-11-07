package common

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"log"
	"os"
	"reflect"
	"strings"
)

type config struct {
	NsqLookupHost       string
	NsqDHost            string
	ProjectPath         string //项目的根目录
	KongUrl             string
	LogUploadApi        string
	GlobalJwtExp        string
	GlobalJwtISS        string
	GlobalJwtPrivateKey string
}

var Config *config

var Log *logrus.Logger

func init() {

	// 初始化日志系统
	file, err := os.OpenFile("/data/inno/miscroservice/go/gin.log", os.O_CREATE|os.O_WRONLY, 0666)

	if err == nil {
		Log = &logrus.Logger{
			Out:   file,
			Level: logrus.DebugLevel,
			Formatter: &prefixed.TextFormatter{
				DisableColors:   true,
				ForceColors:     true,
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
				ForceFormatting: true,
			},
		}
	} else {
		log.Fatal("Failed to log to file, using default stderr")
	}

	//从.env文件中读取  /data/inno/miscroservice/.env
	err = godotenv.Load("/data/inno/miscroservice/.env")
	if err != nil {
		Log.Error("Error loading .env file, please create .env file in /data/inno/miscroservice/ path.")
		panic("Error loading .env file")
	}

	Config = &config{
		os.Getenv("NSQLOOKUP_URL"),
		os.Getenv("NSQSD_URL"),
		"/data/inno/miscroservice/go",
		os.Getenv("API_GATEWAY"),
		os.Getenv("LOG_UPLOAD_API_URL"),
		os.Getenv("GLOBAL_JWT_EXP_SEC"),
		os.Getenv("GLOBAL_JWT_ISS"),
		os.Getenv("GLOBAL_JWT_PRIVATE_KEY_PKCS8"),
	}

	// Create the job queue.
	JobQueue = make(chan Job, 10)

	// Start the dispatcher.
	dispatcher := NewDispatcher(JobQueue, 2)
	dispatcher.Run()

}

/**
获取文件名
*/
func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

/**
过滤掉空的元素
*/
func Filter(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
