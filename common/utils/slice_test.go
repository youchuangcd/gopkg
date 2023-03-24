package utils

import (
	"testing"
)

func getSlice() []string {
	var a []string
	num := RandInt(100, 200)
	for i := 0; i < num; i++ {
		if i%2 == 0 {
			a = append(a, "")
		} else {
			a = append(a, "1")
		}
	}
	return a
}

func BenchmarkStringSliceRemoveEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := getSlice()
		SliceRemoveZeroValue(&slice)
	}
	//old := getSlice()
	////old := []string{"1","","3","5","","6"}
	////fmt.Println(old)
	//fmt.Println(len(old))
	//SliceRemoveZeroValue(&old)
	////fmt.Println(old)
	//fmt.Println(len(old))
}

func TestRandShuffle(t *testing.T) {
	var arr = []uint{1, 4, 5, 6, 7}
	RandShuffle(arr)
	t.Log(arr)
	var arr2 = []int{1, 4, 5, 6, 7}
	RandShuffle(arr2)
	t.Log(arr2)
}
