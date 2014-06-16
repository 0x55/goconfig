package config

import (
	// "fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	Convey("加载一个配置信息测试", t,
		func() {
			c, err := LoadConfigFile("testing/config1.ini")
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)

			Convey("测试获取配置", func() {
				val, err := c.Get("Demo::key1")
				So(err, ShouldBeNil)
				So(val, ShouldEqual, "Let")
			})

			Convey("测试获取不存在的配置", func() {
				_, err := c.Get("Demo::test404")
				So(err, ShouldNotBeNil)
			})

			Convey("删除配置项", func() {
				ok := c.Del("Demo::key1")
				So(ok, ShouldBeTrue)
				ok = c.Del("Demo::key404")
				So(ok, ShouldBeFalse)
				ok = c.DelGroup("demo")
				So(ok, ShouldBeFalse)
				ok = c.DelGroup("Demo")
				So(ok, ShouldBeTrue)
			})

		})
	Convey("加载不存在的配置信息", t,
		func() {
			_, err := LoadConfigFile("testing/conf404.ini")
			So(err, ShouldNotBeNil)
		})

	Convey("加载多个配置文件", t,
		func() {
			c, err := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
			Convey("获取两个配置被替换的参数", func() {
				v, err := c.Get("Demo::中国")
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "China")
			})
			Convey("获取新配置 for config1.ini", func() {
				v := c.Qint("New::Qint", 99)
				So(v, ShouldEqual, 123)
				v64 := c.Qint64("New::Qini64", 9)
				So(v64, ShouldEqual, 9)
			})
		})
}

func TestReload(t *testing.T) {
	Convey("测试Reload", t, func() {
		c, err := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		//删除一条
		ok := c.Del("Demo::key1")
		So(ok, ShouldBeTrue)
		//获取删除的配置
		val, err := c.Get("Demo::key1")
		So(err, ShouldNotBeNil) //错误
		So(val, ShouldBeEmpty)  //空

		c.Reload() //重载

		//获取重载后的配置
		v, err := c.Get("Demo::key1")
		So(err, ShouldBeNil) //错误
		So(v, ShouldEqual, "Let")
	})
}

func TestSaveConfig(t *testing.T) {
	Convey("测试保存配置文件", t, func() {
		c, err := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		ok := SaveConfig(c, "save.ini")
		So(ok, ShouldBeNil)
	})
}

func Benchmark_Reload(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
	for i := 0; i < b.N; i++ {
		c.Reload()
	}
}

func Benchmark_LoadConfigFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = LoadConfigFile("testing/config1.ini", "testing/config2.ini")
	}
}

func Benchmark_Set(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
	for i := 0; i < b.N; i++ {
		c.Set("HSet", "key1", "set")
	}
}
func Benchmark_Get(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini", "testing/config2.ini")
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("Demo::key1")
	}
}

func Benchmark_GetGroup(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini")
	for i := 0; i < b.N; i++ {
		_, _ = c.GetGroup("Demo")
	}
}

func Benchmark_Del(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini")
	for i := 0; i < b.N; i++ {
		_ = c.DelGroup("Demo::key1")
	}
}

func Benchmark_DelGroup(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini")
	for i := 0; i < b.N; i++ {
		_ = c.DelGroup("Demo")
	}
}

func Benchmark_Int(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini")
	for i := 0; i < b.N; i++ {
		_, _ = c.Int("New::Qint")
	}
}

func Benchmark_QInt(b *testing.B) {
	c, _ := LoadConfigFile("testing/config1.ini")
	for i := 0; i < b.N; i++ {
		_ = c.Qint("New::Qint", 0)
	}
}
