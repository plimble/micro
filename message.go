package micro

type Header map[string]interface{}

func (h Header) Get(key string) interface{} {
	return h[key]
}

func (h Header) GetDefault(key string, defaultValue interface{}) interface{} {
	v, ok := h[key]
	if !ok {
		return defaultValue
	}
	return v
}

func (h Header) Set(key string, val interface{}) {
	h[key] = val
}

type message struct {
	Header Header `msg:"h"`
	Body   []byte `msg:"b"`
}
