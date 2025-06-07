package core

import (
	"errors"

	"github.com/InazumaV/V2bX/conf"
)

var (
	cores = map[string]func(c *conf.CoreConfig) (Core, error){}
)

func NewCore(c []conf.CoreConfig) (Core, error) {
	if len(c) < 0 {
		return nil, errors.New("no have vail core")
	}
	// multi core
	if len(c) > 1 {
		return NewSelector(c)
	}
	// one core
	typeKey := normalizeCoreType(c[0].Type)
	if f, ok := cores[typeKey]; ok {
		return f(&c[0])
	} else {
		return nil, errors.New("unknown core type: " + c[0].Type)
	}
}

func RegisterCore(t string, f func(c *conf.CoreConfig) (Core, error)) {
	cores[t] = f
}

func RegisteredCore() []string {
	cs := make([]string, 0, len(cores))
	for k := range cores {
		cs = append(cs, k)
	}
	return cs
}

// normalizeCoreType 支持常见别名
func normalizeCoreType(t string) string {
	switch t {
	case "sing", "sing-box", "singbox":
		return "sing"
	case "xray", "v2ray":
		return "xray"
	case "hysteria2", "hy2":
		return "hysteria2"
	default:
		return t
	}
}
