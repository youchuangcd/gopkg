package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type Reply struct {
	reply interface{}
	error error
}

var ErrNil = redis.ErrNil

func getReply(rp interface{}, err error) *Reply {
	return &Reply{reply: rp, error: err}
}

func (r *Reply) Result() (interface{}, error) {
	return r.reply, r.error
}

func (r *Reply) Error() error {
	return r.error
}

func (r *Reply) Int64() (int64, error) {
	return redis.Int64(r.reply, r.error)
}

func (r *Reply) Int() (int, error) {
	return redis.Int(r.reply, r.error)
}

func (r *Reply) Uint64() (uint64, error) {
	return redis.Uint64(r.reply, r.error)
}

func (r *Reply) Uint() (uint, error) {
	res, err := redis.Uint64(r.reply, r.error)
	return uint(res), err
}

func (r *Reply) Float64() (float64, error) {
	return redis.Float64(r.reply, r.error)
}

func (r *Reply) String() (string, error) {
	return redis.String(r.reply, r.error)
}

func (r *Reply) Bytes() ([]byte, error) {
	return redis.Bytes(r.reply, r.error)
}

func (r *Reply) Bool() (bool, error) {
	return redis.Bool(r.reply, r.error)
}

func (r *Reply) Float64s() ([]float64, error) {
	return redis.Float64s(r.reply, r.error)
}

func (r *Reply) Strings() ([]string, error) {
	return redis.Strings(r.reply, r.error)
}

func (r *Reply) ByteSlices() ([][]byte, error) {
	return redis.ByteSlices(r.reply, r.error)
}

func (r *Reply) Int64s() ([]int64, error) {
	return redis.Int64s(r.reply, r.error)
}

func (r *Reply) Ints() ([]int, error) {
	return redis.Ints(r.reply, r.error)
}

func (r *Reply) StringMap() (map[string]string, error) {
	return redis.StringMap(r.reply, r.error)
}

func (r *Reply) IntMap() (map[string]int, error) {
	return redis.IntMap(r.reply, r.error)
}

func (r *Reply) Int64Map() (map[string]int64, error) {
	return redis.Int64Map(r.reply, r.error)
}

func (r *Reply) Positions() ([]*[2]float64, error) {
	return redis.Positions(r.reply, r.error)
}

func (r *Reply) Values() ([]interface{}, error) {
	return redis.Values(r.reply, r.error)
}

// obj为一个指针对象
func (r *Reply) ScanStruct(obj interface{}) error {
	v, err := r.Values()
	if err != nil {
		return err
	}
	return redis.ScanStruct(v, obj)
}

// UintKeyMap
// @Description:
// @receiver r
// @return map[uint]uint
// @return error
func (r *Reply) UintKeyMap() (map[uint]uint, error) {
	values, err := redis.Values(r.reply, r.error)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: Int64Map expects even number of values result")
	}
	m := make(map[uint]uint, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: Int64Map key not a bulk string value")
		}
		value, err := redis.Uint64(values[i+1], nil)
		if err != nil {
			return nil, err
		}
		ck, err := strconv.ParseUint(string(key), 10, 64)
		if err != nil {
			return nil, err
		}
		m[uint(ck)] = uint(value)
	}
	return m, nil
}
