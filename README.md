# chix

一个基于 chi 的轻量 Web 框架封装，提供 Echo 风格 API（Use/Group）与默认恢复中间件。

快速开始：

```go
package main

import (
	"github.com/lemoba/chix"
)

func main() {
	r := chix.New()

	r.Use(/* middlewares */)

	r.Group("/api", func(g *chix.Chix) {
		g.Get("/health", func(c chix.Context) { c.Success("ok") })
	})

	r.Run(":3000")
}
```
