package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/psyngw/simpleLanzou/lanzou"
)

func main() {
	var url string
	var pwd string
	var download string
	fmt.Println(runtime.GOOS)
	fmt.Println("请输入单文件网址：")
	fmt.Scanln(&url)
	direct_url, err := lanzou.Lanzou(url, pwd, "pub")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(`完成解析，直链地址为：
` + direct_url + `
是否需要直接打开下载？[Y/N]`)
		fmt.Scanln(&download)
		if download == "y" || download == "Y" {
			openUrl(direct_url)
		}
	}
	fmt.Println("按任意键退出...")
	fmt.Scanln()
}

func openUrl(url string) {
	var cmd *exec.Cmd
	os_version := runtime.GOOS
	if os_version == "linux" {
		cmd = exec.Command("xdg-open", url)
	} else if os_version == "windows" {
		cmd = exec.Command("cmd", "/c", "start", strings.Replace(url, "&", "^&", -1))
	} else {
		cmd = exec.Command("open", url)
	}
	cmd.Run()
	return
}
