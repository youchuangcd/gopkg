package utils

import (
	"github.com/gogap/errors"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// SliceRemoveZeroValue
// @Description: 切片移除零值
// @param src 必须传切片的指针
func SliceRemoveZeroValue(src interface{}) {
	j := 0
	vrf := reflect.ValueOf(src)
	if vrf.Kind() != reflect.Ptr || vrf.Elem().Kind() != reflect.Slice || vrf.Elem().Len() == 0 {
		return
	}

	num := vrf.Elem().Len()
	vrfElem := vrf.Elem()
	for i := 0; i < num; i++ {
		// 非零值，且可以设置，就把值设置为上一个零值的位置
		if !vrfElem.Index(i).IsZero() && vrfElem.Index(j).CanSet() {
			vrfElem.Index(j).Set(vrfElem.Index(i))
			j++ // 记录有效的位置
		}
	}
	// 只截取非零值的数据
	vrfElem.Set(vrfElem.Slice(0, j))
}

// RemoveRepeatElementSlice
// @Description: 移除切片重复的元素
// @param s
// @return interface{} 返回值需要对应类型的转换
func RemoveRepeatElementSlice(arg interface{}) (res interface{}, err error) {
	s, ok := takeSliceArg(arg)
	if !ok {
		return nil, errors.New("invalid param")
	}
	if len(s) <= 1 {
		return arg, nil
	}
	//if len(s) < 1024 {
	//	res = removeRepeatElementSliceByLoop(s)
	//} else { // 空间换时间
	res = removeRepeatElementSliceByMap(s)
	//}
	return
}

// RemoveRepeatElementSliceString
// @Description: 字符串切片去重
// @param arg
// @return []string
// @return error
func RemoveRepeatElementSliceString(arg []string) ([]string, error) {
	if len(arg) <= 1 {
		return arg, nil
	}
	tmp, err := RemoveRepeatElementSlice(arg)
	if err != nil {
		return nil, err
	}
	tmp2, _ := tmp.([]interface{})
	var newFields = make([]string, 0, len(tmp2))
	for _, v := range tmp2 {
		newFields = append(newFields, v.(string))
	}
	return newFields, nil
}

// RemoveRepeatElementSliceUint
// @Description: uint切片去重
// @param arg
// @return []uint
// @return error
func RemoveRepeatElementSliceUint(arg []uint) ([]uint, error) {
	if len(arg) <= 1 {
		return arg, nil
	}
	tmp, err := RemoveRepeatElementSlice(arg)
	if err != nil {
		return nil, err
	}
	tmp2, _ := tmp.([]interface{})
	var newFields = make([]uint, 0, len(tmp2))
	for _, v := range tmp2 {
		newFields = append(newFields, v.(uint))
	}
	return newFields, nil
}

// RemoveRepeatElementSliceByLoop
// @Description: 通过两重循环过滤重复元素
// @param slc
// @return []int
func removeRepeatElementSliceByLoop(slc []interface{}) []interface{} {
	var result []interface{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// RemoveRepeatElementSliceByMap
// @Description: 通过map主键唯一的特性过滤重复元素
// @param slc
// @return []interface{}
func removeRepeatElementSliceByMap(slc []interface{}) []interface{} {
	var result []interface{}
	tempMap := map[interface{}]struct{}{} // 存放不重复主键
	for _, e := range slc {
		if _, ok := tempMap[e]; !ok {
			tempMap[e] = struct{}{}
			result = append(result, e)
		}
	}
	return result
}

func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}

	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

// DiffArray 求两个切片的差集
func DiffArray(a []int, b []int) []int {
	var diffArray []int
	temp := map[int]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}

func DiffStrArray(a []string, b []string) []string {
	var diffArray []string
	temp := map[string]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}

// SliceTrimSpace
// @Description: 字符串切片去除首尾空格
// @param s
// @return []string
func SliceTrimSpace(s []string) []string {
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	return s
}

// 判断元素是否在切片中存在
func IsSliceExist(s []int, f int) int {
	for _, i := range s {
		if i == f {
			return 1
		}
	}
	return 0
}

// RandShuffle
//
//	@Description: 打乱切片
//	@param slice
func RandShuffle(slice []interface{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
