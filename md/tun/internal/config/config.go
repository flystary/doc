package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Path  string
	viper *viper.Viper
}

func NewConfigPath(path string) (C Config, err error) {
	C.Path = path
	C.viper = viper.New()
	C.viper.SetConfigFile(path)
	C.viper.SetConfigType("yml")
	err = C.viper.ReadInConfig()
	return
}

func (c Config) Get(key string, defalut any) (ret any) {
	if c.viper.IsSet(key) {
		return c.viper.Get(key)
	}
	return defalut
}

func (c Config) GetString(key string, defalut string) (ret string) {
	if c.viper.IsSet(key) {
		return c.viper.GetString(key)
	}
	return defalut
}

func (c Config) GetBool(key string, defalut bool) (ret bool) {
	if c.viper.IsSet(key) {
		return c.viper.GetBool(key)
	}
	return defalut
}

func (c Config) GetInt(key string, defalut int) (ret int) {
	if c.viper.IsSet(key) {
		return c.viper.GetInt(key)
	}
	return defalut
}

func (c Config) GetInt64(key string, defalut int64) (ret int64) {
	if c.viper.IsSet(key) {
		return c.viper.GetInt64(key)
	}
	return defalut
}

func (c Config) get(k string, v interface{}) interface{} {
	parts := strings.Split(k, ".")
	for _, p := range parts {
		m, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil
		}

		v, ok = m[p]
		if !ok {
			return nil
		}
	}

	return v
}
