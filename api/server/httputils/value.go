package httputils

const (
	globalPrefix = "in-single-request-"

	// CLI_USER_AGENT is used as key type for user-agent string in net/context struct
	CLI_USER_AGENT  = "upstream-user-agent"
	CLI_API_VERSION = "upstream-api-version"
)

func Exist(vars map[string]string, key string) bool {
	if vars == nil {
		return false
	}
	_, exist := vars[globalPrefix+key]
	return exist
}

func Get(vars map[string]string, key string) string {
	if vars == nil {
		return ""
	}

	return vars[globalPrefix+key]
}

func Put(vars map[string]string, key, val string) {
	if vars != nil {
		vars[globalPrefix+key] = val
	}
}

func Del(vars map[string]string, key string) {
	if vars != nil {
		delete(vars, globalPrefix+key)
	}
}
