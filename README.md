# gogo
based on the go lang branch dev.go2go
## 修改 ast 支持 gogo 关键字
## 修改 checker(types/stmt.go)
## Example
```
package main

import (
  "fmt"
  "time"
)

func main() {
	gogo fmt.Println("vim-go","123")
	gogo fmt.Println("vim-go2","321")
	time.Sleep(time.Second)
	return
}
```
