## SimpleLanzou

### Go 引入包使用

`go get -u github.com/psyngw/simpleLanzou`

```go
package main

import (
      "github.com/psyngw/simpleLanzou/lanzou"
      "fmt"
  )

func main() {
	direct_url, err := lanzou.Lanzou("网址", "密码(无密码为空)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(direct_url)
}
```
### 下载使用

[下载地址](https://github.com/psyngw/simpleLanzou/releases)

