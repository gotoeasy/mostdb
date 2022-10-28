package most

import "github.com/gotoeasy/glang/cmn"

var mapStorageHandle map[string](*DataStorage)

func init() {
	mapStorageHandle = make(map[string](*DataStorage))
}

func InitDb() {
	NewDataStorageHandle("")
}

func NewDataStorageHandle(storeName string) *DataStorage {

	if storeName == "" {
		storeName = "store"
	}

	// TODO 考虑不受控制的Map扩张
	cacheStore := mapStorageHandle[storeName] // 缓存中的存储对象
	if cacheStore != nil {
		if cacheStore.IsOpen() {
			return cacheStore
		}
	}

	store, err := newDataStorage(storeName)
	if err != nil {
		cmn.Error(err)
	}
	mapStorageHandle[storeName] = store
	return store
}
