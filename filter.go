package go_web

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(ctx *Context)

// 检测 MetricFilterBuilder 为 FilterBuilder 类型
var _ FilterBuilder = MetricFilterBuilder

func MetricFilterBuilder(next Filter) Filter {
	return func(ctx *Context) {
		start := time.Now().Nanosecond()
		next(ctx)
		end := time.Now().Nanosecond()
		fmt.Printf("run time：%d\n", end-start)
	}
}
