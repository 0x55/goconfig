goconfig
========

##About
Go package:The INI configuration file read and write.
[中文文档](README_CN.md)

##Features


##Example

### config.ini
	appname = "WishCMS"
	
	[Demo]
	key1 = "Let"
	
	todo = "恩恩"
	
	中国 = "China"
	
	test = ""China"
	
	[Hi]
	name = "chris"
	
	age = "23"
	
	nu = "-1"
	
	[auto]
	- = "config"
	- = "hello"
	- = "go"
	- = "conf"
	- = "demo"
	
	[New]
	Qint = "123"

=======================
##code:
 
	c, _ := LoadConfigFile("testing/config1.ini", "testing/config2.ini")//加载多个配置
	
	_, _ := c.Get("Demo::key1")//获取Demo组key1值
	
	_, _ := c.Get("appname")//获取默认分组 appname值
	
	_, _ := c.GetGroup("Demo")//获取Demo分组所有配置
	
	c.Set("New", "key1", "999")//设置New组key1值为999
	
	_, _ := c.Int("New::Qint")//获取New组Qint值并转换为int
	
	_ := c.Qint64("New::Qini64", 9)//快速获取值并设置默认值
	
	_, _ := c.Get("auto::#1")//获取auto自增配置组第一个配置项值 config
	
	c.Reload() //重载配置
	
	ok := SaveConfig(c, "save.ini")//保存配置
	
	c, _ := LoadConfigFile("testing/config1.ini")//加载一个配置

##Install
go get github.com/0x55/goconfig

include config package：
	
	import (
		"github.com/0x55/goconfig"
	)


##License
This project is under Apache v2 License.
