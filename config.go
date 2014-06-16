package config

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	DEFAULT_GROUP = "default"
	//group未找到
	GROUP_NOT_FOUND = iota
	//选项未找到
	KEY_NOT_FOUND
	//组错误
	GROUP_FAIL
	//解析失败
	COULD_NOT_PARSE
)

var (
	lineBreak = "\n"
)

func init() {
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}
}

type Config struct {
	files []string
	data  map[string]map[string]string //group:key=>val
	lock  sync.RWMutex
	Mode  bool
}

/**
 * 初始化配置类型
 */
func parseConfig(files []string) *Config {
	c := new(Config)
	c.files = files
	c.data = make(map[string]map[string]string)
	c.Mode = true
	return c
}

/**
 * 预处理key
 */
func parseKey(k string) (group, key string) {
	i := strings.IndexAny(k, "::")
	group = DEFAULT_GROUP
	key = k
	if i > 0 {
		group = k[0:i]
		key = k[i+2:]
	}
	return group, key
}

/**
 * 设置配置
 */
func (c *Config) Set(group, key, val string) {
	if c.Mode {
		defer c.lock.Unlock()
		c.lock.Lock()
	}
	// g, k := parseKey(key)
	if _, ok := c.data[group]; !ok {
		c.data[group] = make(map[string]string)
	}
	c.data[group][key] = val //设置值
}

/**
 * 获取配置
 */
func (c *Config) Get(key string) (string, error) {
	if c.Mode {
		defer c.lock.Unlock()
		c.lock.Lock()
	}
	g, k := parseKey(key)
	if _, ok := c.data[g]; !ok {
		return "", configError{GROUP_NOT_FOUND, g}
	}
	val, ok := c.data[g][k]
	if !ok || len(val) == 0 {
		return "", configError{KEY_NOT_FOUND, key}
	}
	return val, nil
}

/**
 * 获取一组配置
 */
func (c *Config) GetGroup(group string) (map[string]string, error) {
	if _, ok := c.data[group]; !ok {
		return nil, configError{GROUP_NOT_FOUND, group}
	}
	valsMsp := c.data[group]
	return valsMsp, nil
}

/**
 * 删除配置项
 */
func (c *Config) Del(key string) bool {
	g, k := parseKey(key)
	if _, ok := c.data[g]; !ok {
		return false
	}
	if _, ok := c.data[g][k]; ok {
		delete(c.data[g], k)
		return true
	}
	return false
}

/**
 * 删除一组配置项
 */
func (c *Config) DelGroup(group string) bool {
	if _, ok := c.data[group]; !ok {
		return false
	}
	delete(c.data, group)
	return true
}

///////////////////////获取指定类型的配置///////////////////////
func (c *Config) Bool(key string) (bool, error) {
	val, err := c.Get(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(val)
}

func (c *Config) Float64(key string) (float64, error) {
	val, err := c.Get(key)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(val, 64)
}
func (c *Config) Int(key string) (int, error) {
	val, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (c *Config) Int64(key string) (int64, error) {
	val, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

////////////////////快速获取//////////////////////

func (c *Config) Qbool(key string, def bool) bool {
	val, err := c.Bool(key)
	if err != nil {
		return def
	}
	return val
}
func (c *Config) Qfloat64(key string, def float64) float64 {
	val, err := c.Float64(key)
	if err != nil {
		return def
	}
	return val
}
func (c *Config) Qint(key string, def int) int {
	val, err := c.Int(key)
	if err != nil {
		return def
	}
	return val
}
func (c *Config) Qint64(key string, def int64) int64 {
	val, err := c.Int64(key)
	if err != nil {
		return def
	}
	return val
}

/////////////////////////////////////////////////
type configError struct {
	Reason  uint
	Message string
}

func (err configError) Error() string {
	switch err.Reason {
	case GROUP_FAIL:
		// return "组名错误"
		return fmt.Sprint("组名错误：'%s'", string(err.Message))
	case GROUP_NOT_FOUND:
		// return "组未找到"
		return fmt.Sprint("组未找到：'%s'", string(err.Message))
	case COULD_NOT_PARSE:
		return fmt.Sprint("无法解析：'%s'", string(err.Message))
	}
	return "[Config] invalid read error"
}
