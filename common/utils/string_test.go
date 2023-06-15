package utils

import (
	"context"
	"testing"
)

func TestCutStr(t *testing.T) {
	s := RandSeq(80)
	s += "中文额维护费我返回二维复婚非法额会务if额外复合瓦"
	s += "真的何物的核武的让我回答饿汉武帝回调2hd782的话du8ew问问"
	t.Log("新的内容: ", CutStr(s, 20, "..."))
}

func BenchmarkUUID(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < 1000000; i++ {
		GenTraceId(ctx)
	}
}

func BenchmarkCutStr(b *testing.B) {
	s := RandSeq(1000)
	s += "中文额维护费我返回二维复婚非法额会务if额外复合瓦"
	s += RandSeq(1000)
	for i := 0; i < b.N; i++ {
		CutStr(s, 2000, "...")
	}
}
