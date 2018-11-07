package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"os/exec"
	. "shell-exec/lib"
)

/**
打包日志下载
*/
func Packing(c *gin.Context) {
	// 获取待打包日志的天数，默认为1
	day := c.DefaultQuery("day", "1")

	cmd := exec.Command(Config.ProjectPath+" /log.sh", day)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		c.String(422, "failed with"+err.Error())
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if len(errStr) != 0 {
		c.String(422, string(errStr))
		return
	}

	c.String(200, string(outStr))
}
