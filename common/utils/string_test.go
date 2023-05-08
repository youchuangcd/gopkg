package utils

import (
	"context"
	"fmt"
	"testing"
)

func TestCutStr(t *testing.T) {
	s := RandSeq(80)
	s += "中文额维护费我返回二维复婚非法额会务if额外复合瓦"
	s += "真的何物的核武的让我回答饿汉武帝回调2hd782的话du8ew问问"
	t.Log("新的内容: ", CutStr(s, 20, "..."))
}

func TestUUID(t *testing.T) {
	s := GenTraceId(context.Background())
	fmt.Println(s)
	t.Log("新的内容: ", s)
}

func BenchmarkCutStr(b *testing.B) {
	s := RandSeq(1000)
	s += "中文额维护费我返回二维复婚非法额会务if额外复合瓦"
	s += RandSeq(1000)
	for i := 0; i < b.N; i++ {
		CutStr(s, 2000, "...")
	}
}
