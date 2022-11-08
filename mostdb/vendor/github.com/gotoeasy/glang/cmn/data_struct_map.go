package cmn

// 值都为string类型的map
type MapString map[string]string

// 值都为any类型的map
type Map map[string]any

// 创建MapString对象
func NewMapString() MapString {
	return make(MapString)
}

// 创建MapString对象
func NewMap() Map {
	return make(Map)
}

// 设定
func (m MapString) Put(key string, value string) MapString {
	m[key] = value
	return m
}

// 获取
func (m MapString) Get(key string) string {
	return m[key]
}

// 设定
func (m Map) Put(key string, value any) Map {
	m[key] = value
	return m
}

// 获取
func (m Map) Get(key string) any {
	return m[key]
}
