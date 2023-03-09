// Package redis 不能在multi中使用，否则脚本无法执行load，导致命令无效
package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"strings"
)

type scriptMapItem struct {
	keyCount int
	script   string
}

const (
	/*
	* 增量计数器，第一次增量计数的时候，给key加上过期时间，解决并发问题
	* @eg: EvalScript(ScriptKeyIncr, "key", expire)
	* @return int
	 */
	ScriptKeyIncr = "incr"
	/*
	* 增量计数器，第一次增量计数的时候，给key加上过期时间，解决并发问题
	* @eg: EvalScript(ScriptKeyIncrBy, "key", increment, expire)
	* @return int
	 */
	ScriptKeyIncrBy = "incrby"
	/*
	* 增量计数器，并在增量值超过最大值时，重置为0;
	* @eg: EvalScript(ScriptKeyIncrReset, "key", max_counter，expire)
	* @return int 当前值
	 */
	ScriptKeyIncrReset = "incr_reset"
	/*
	* 增量计数器，并在超出最大值时，重置为0；或主动重置为0;
	* @eg: EvalScript(ScriptKeyIncrByOutMaxReset, "key", 本次要加的数量, 最大值, 本次是否要重置(0,1), 重置的值)
	* @return int64 是否重置0，1, int64 当前值或重置后的值
	 */
	ScriptKeyIncrByOutMaxReset = "incrby_out_max_reset"
	/*
	* 增量计数器，如果当前值没有大于限定值，才可以加一并返回[1, 累加后的值]，否则返回[0, 当前值]
	* @eg: EvalScript(ScriptKeyIncrMax, "key", max_counter, expire)
	* @return int64 是否达到最大值, int64 累加后的值或当前值
	 */
	ScriptKeyIncrMax = "incr_max"
	/*
	* 存在才将 key 中储存的数字值减一
	* @eg: EvalScript(ScriptKeyDecrExist, "key")
	* @return int64
	 */
	ScriptKeyDecrExist = "decr_exist"
	/*
	* hash结构批量增加计数器，返回每个域增加后的值
	* @eg: EvalScript(ScriptKeyHMIncrBy, "key", "field1", "field1Value", ["field2", "field2Value"]...)
	 */
	ScriptKeyHMIncrBy = "hmincrby"

	/*
	* 与加锁的值相同才解锁
	* @eg: EvalScript(ScriptKeyValueEqualsUnlock, "key", "lock value")
	 */
	ScriptKeyValueEqualsUnlock = "valueEqualsUnlock"
)

var (
	/**
	* 脚本列表
	* lua脚本的返回值是redis原始结构，不可以像php里的redis函数一样使用如（zRevRange等）需要自己处理返回值
	* @var array
	 */
	scriptMap = map[string]scriptMapItem{
		// 增量计数器，第一次增量计数的时候，给key加上过期时间，解决并发问题; eg: EvalScript('incr', 'key', expire)
		ScriptKeyIncr: {
			keyCount: 1, // 表示脚本代码里用了几个key参数; eg: KEYS[1], KEYS[2] 表示用了2个
			script: `
local count = redis.call('incr',KEYS[1]); 
if tonumber(count) == 1 then 
	redis.call('expire',KEYS[1],ARGV[1]); 
end; 
return count;`,
		},
		// 增量计数器，第一次增量计数的时候，给key加上过期时间，解决并发问题; eg: EvalScript("incrby", "key", increment, expire)
		ScriptKeyIncrBy: {
			keyCount: 1,
			script: `
local count = redis.call('incrby',KEYS[1], ARGV[1]); 
if tonumber(count) == 1 then 
	redis.call('expire',KEYS[1],ARGV[2]); 
end; 
return count;
`,
		},
		// 增量计数器，并在增量值超过最大值时，重置为0; eg: EvalScript('incr_reset', 'key', max_counter，expire)
		ScriptKeyIncrReset: {
			keyCount: 1,
			script: `
local count = redis.call('incr',KEYS[1]); 
if count == 1 then 
	redis.call('expire',KEYS[1],ARGV[2]); 
end; 
if count > tonumber(ARGV[1]) then 
	redis.call('set', KEYS[1], 0); 
	return 0; 
end; 
return count;`,
		},
		// 增量计数器，并在超出最大值时，重置为0；或主动重置为0; eg: EvalScript('incr_out_max_reset', 'key', 本次要加的数量, 最大值, 本次是否要重置(0,1), 重置的值)
		ScriptKeyIncrByOutMaxReset: {
			keyCount: 1,
			script: `
local is_reset, count = 0, 0; 
if tonumber(ARGV[3]) == 1 then 
	redis.call('set', KEYS[1], tonumber(ARGV[4])); 
	is_reset = 1; 
	count = tonumber(ARGV[4]); 
else 
	count = redis.call('incrby',KEYS[1], ARGV[1]); 
	if tonumber(count) > tonumber(ARGV[2]) then 
		redis.call('set', KEYS[1], tonumber(ARGV[4])); 
		is_reset = 1; 
		count = tonumber(ARGV[4]); 
	end; 
end; 
return {is_reset, count};`,
		},
		// 增量计数器，如果当前值没有大于限定值，才可以加一并返回[1, 累加后的值]，否则返回[0, 当前值] eg: EvalScript('incr_max'', 'key', [max_counter, expire])
		ScriptKeyIncrMax: {
			keyCount: 1,
			script: `
local count = redis.call('get',KEYS[1]); 
if ( count == false or tonumber(count) < tonumber(ARGV[1]) ) then 
	count = redis.call('incr', KEYS[1]); 
	if count == 1 then 
		redis.call('expire',KEYS[1],ARGV[2]); 
	end; 
	return {1, count}; 
else 
	return {0, tonumber(count)}; 
end;`,
		},
		// 存在才将 key 中储存的数字值减一 eg: EvalScript('decr_exist', 'key')
		ScriptKeyDecrExist: {
			keyCount: 1,
			script: `
local count = redis.call('exists',KEYS[1]); 
if tonumber(count) == 1 then 
	count = redis.call('decr',KEYS[1]); 
end; 
return tonumber(count);`,
		},
		// hash结构批量增加计数器，返回每个域增加后的值
		ScriptKeyHMIncrBy: {
			keyCount: 1,
			script: `
local n = {}; 
for i=1, #(ARGV)/2 do 
	n[i] = redis.call('hincrby',KEYS[1],ARGV[i*2-1],ARGV[i*2]); 
end; 
return n;`,
		},

		// 加锁的值相同才解锁
		ScriptKeyValueEqualsUnlock: {
			keyCount: 1,
			script: `
if redis.call('get', KEYS[1]) == ARGV[1] then 
	return redis.call('del', KEYS[1]) 
else 
	return 0 
end`,
		},
	}
)

// Script
// @Description: 扩展了redigo的脚本执行
type Script struct {
	ctx      context.Context
	rs       *redis.Script
	keyCount int
	src      string
}

// NewScript returns a new script object. If keyCount is greater than or equal
// to zero, then the count is automatically inserted in the EVAL command
// argument list. If keyCount is less than zero, then the application supplies
// the count as the first value in the keysAndArgs argument to the Do, Send and
// SendHash methods.
func NewScript(ctx context.Context, keyCount int, src string) *Script {
	rs := redis.NewScript(keyCount, src)
	return &Script{ctx: ctx, rs: rs, src: src, keyCount: keyCount}
}

// Hash returns the script hash.
func (s *Script) Hash() string {
	return s.rs.Hash()
}

func (s *Script) args(spec string, keysAndArgs []interface{}) []interface{} {
	var args []interface{}
	if s.keyCount < 0 {
		args = make([]interface{}, 1+len(keysAndArgs))
		args[0] = spec
		copy(args[1:], keysAndArgs)
	} else {
		args = make([]interface{}, 2+len(keysAndArgs))
		args[0] = spec
		args[1] = s.keyCount
		copy(args[2:], keysAndArgs)
	}
	return args
}

// Do
// @Description: 执行脚本命令参数
// @receiver s
// @param c
// @param keysAndArgs
// @return interface{}
// @return error
func (s *Script) Do(c redis.Conn, keysAndArgs ...interface{}) (interface{}, error) {
	retryTimes := 1
Retry:
	v, err := c.Do("EVALSHA", s.args(s.rs.Hash(), keysAndArgs)...)
	if e, ok := err.(redis.Error); ok && strings.HasPrefix(string(e), "NOSCRIPT ") {
		retryTimes--
		if retryTimes >= 0 {
			if err = s.Load(c, context.Background()); err != nil {
				return nil, err
			}
			goto Retry
		}
	}
	return v, err
}

func (s *Script) DoContext(c redis.Conn, ctx context.Context, keysAndArgs ...interface{}) (interface{}, error) {
	retryTimes := 1
Retry:
	v, err := redis.DoContext(c, ctx, "EVALSHA", s.args(s.rs.Hash(), keysAndArgs)...)
	if e, ok := err.(redis.Error); ok && strings.HasPrefix(string(e), "NOSCRIPT ") {
		retryTimes--
		if retryTimes >= 0 {
			if err = s.Load(c, ctx); err != nil {
				return nil, err
			}
			goto Retry
		}
	}
	return v, err
}

// EvalScript
// @Description: 执行lua脚本
// @param scriptKey
// @param key
// @param args
// @return res
// @return *Reply
//
// @eg: EvalScript(ctx, 'incr', 'key', expire)
// @eg: EvalScript(ctx, 'incr_reset', 'key', max_counter，expire)
// @eg: EvalScript(ctx, 'incr_out_max_reset', 'key', 本次要加的数量, 最大值, 本次是否要重置(0,1), 重置的值)
// @eg: EvalScript(ctx, 'incr_max”, 'key', max_counter, expire)
// @eg: EvalScript(ctx, 'decr_exist', 'key')
func EvalScript(ctx context.Context, scriptKey string, args ...interface{}) *Reply {
	var (
		info scriptMapItem
		ok   bool
	)
	if info, ok = scriptMap[scriptKey]; !ok {
		mylog.WithError(nil, gopkg.LogRedis, map[string]interface{}{
			"scriptKey": scriptKey,
		}, "请先配置"+scriptKey+"脚本")
		return getReply(nil, gopkg.ErrorRedisInvalidScriptKey)
	}
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	s := NewScript(ctx, info.keyCount, info.script)
	return getReply(s.DoContext(c, ctx, args...))
}

// Load loads the script without evaluating it.
func (s *Script) Load(c redis.Conn, ctx context.Context) error {
	_, err := redis.DoContext(c, ctx, "SCRIPT", "LOAD", s.src)
	return err
}
