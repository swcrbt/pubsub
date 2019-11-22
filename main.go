package main

import (
	"gitlab.orayer.com/golang/issue/cmd"
)

var (
	// AppVersion 应用版本
	AppVersion string
	// BuildDate 构建日期
	BuildDate string
	// GitCommit 最后提交的git commit
	GitCommit string
)
 
func main(){
	cmd.SetVersion(AppVersion, BuildDate, GitCommit)
	cmd.Execute()
}