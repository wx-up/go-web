package go_web

import (
	"fmt"
	"time"
)

type FilterBuilder func(next HandlerFunc) HandlerFunc

// 检测 MetricFilterBuilder 为 FilterBuilder 类型
var _ FilterBuilder = MetricFilterBuilder

func MetricFilterBuilder(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		start := time.Now().Nanosecond()
		next(ctx)
		end := time.Now().Nanosecond()
		fmt.Printf("run time：%d\n", end-start)
	}
}
