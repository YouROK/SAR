package config

var (
	cfg map[string]string
)

func Set(name, val string) {
	if cfg == nil {
		cfg = make(map[string]string)
	}
	cfg[name] = val
}

func Get(name string) string {
	if cfg == nil {
		return ""
	}
	if ret, ok := cfg[name]; ok {
		return ret
	}
	return ""
}
