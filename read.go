package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func (c *Config) readConfigFile(filename string) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	return c.readConfig(f)
}

func LoadConfigFile(file string, more ...string) (c *Config, err error) {
	filesMap := make([]string, 1, len(more)+1)
	filesMap[0] = file
	if len(filesMap) > 0 {
		filesMap = append(filesMap, more...)
	}
	c = parseConfig(filesMap)
	for _, name := range filesMap {
		if err = c.readConfigFile(name); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Config) readConfig(rd io.Reader) error {
	//创建缓冲区
	buffer := bufio.NewReader(rd)
	group := DEFAULT_GROUP
	counter := 1
	for {
		line, err := buffer.ReadString('\n') //读取一行
		line = strings.TrimSpace(line)
		lineLen := len(line)
		if err != nil {
			if err != io.EOF {
				return err
			}
			if lineLen == 0 {
				break
			}
		}
		if err == io.EOF {
			break
		}
		switch {
		case lineLen == 0:
		case line[0] == '#' || line[0] == ';':
			continue
		case line[0] == '[' && line[lineLen-1] == ']':
			group = strings.TrimSpace(line[1 : lineLen-1])
			counter = 1 //重置自增变量
			continue
		case group == "": //分组为空
			return configReadError{GROUP_FAIL, line}

		default:
			var (
				i                    int
				key, value, valQuote string
			)
			//计算=出现的位置，处理等号左边key
			i = strings.IndexAny(line, "=")
			if i < 0 {
				return configReadError{COULD_NOT_PARSE, line}
			}
			key = strings.TrimSpace(line[0:i])
			if key == "-" {
				key = "#" + fmt.Sprint(counter)
				counter++
			}
			lineRight := strings.TrimSpace(line[i+1:]) //去除等号
			rightLen := len(lineRight)
			if rightLen >= 3 {
				if lineRight[0:1] == `"` && strings.IndexAny(lineRight[1:], `"`) == rightLen-2 {
					valQuote = `"`
				}
			}
			if valQuote != "" {
				value = lineRight[1 : rightLen-1]
			} else {
				value = lineRight
			}
			//添加配置
			c.Set(group, key, value)
		}
	}
	return nil
}

/**
 * 重载配置文件
 */
func (c *Config) Reload() (err error) {
	var newConfig *Config
	if len(c.files) == 1 {
		newConfig, err = LoadConfigFile(c.files[0])
	} else {
		newConfig, err = LoadConfigFile(c.files[0], c.files[1:]...)
	}
	if err == nil {
		*c = *newConfig
	}
	return err
}

func (c *Config) AppendConfigFiles(files ...string) error {
	c.files = append(c.files, files...)
	return c.Reload()
}

type configReadError struct {
	ErrNo   uint
	ErrInfo string
}

/**
 * readConfError实现 Error 接口
 */
func (err configReadError) Error() string {
	switch err.ErrNo {
	case GROUP_FAIL:
		return "空的组名"
	case COULD_NOT_PARSE:
		return fmt.Sprint("无法解析：%s", string(err.ErrInfo))
	}
	return "invalid read error"
}
