package most

import (
	"errors"
	"mostdb/conf"
	"os"
	"sync"

	"github.com/gotoeasy/glang/cmn"
	"github.com/syndtr/goleveldb/leveldb"
)

// 存储结构体
type DataStorage struct {
	dbPath  string      // 数据存放目录
	leveldb *leveldb.DB // leveldb
	opened  bool        // 是否打开状态
	mu      sync.Mutex  // 锁
}

// 存储对象
func newDataStorage(storeName string) (*DataStorage, error) {
	storage := &DataStorage{
		dbPath: conf.GetStorageRoot() + cmn.PathSeparator() + storeName,
	}
	err := os.MkdirAll(storage.dbPath, 0777)
	if err != nil {
		return nil, err
	}

	return storage, storage.open()
}

// 保存
func (s *DataStorage) Put(key []byte, value []byte) error {
	if !s.opened {
		return errors.New("数据库没有打开")
	}
	return s.leveldb.Put(key, value, nil)
}

// 删除
func (s *DataStorage) Del(key []byte) error {
	if !s.opened {
		return errors.New("数据库没有打开")
	}
	return s.leveldb.Delete(key, nil)
}

// 获取
func (s *DataStorage) Get(key []byte) ([]byte, error) {
	if !s.opened {
		return nil, errors.New("数据库没有打开")
	}
	return s.leveldb.Get(key, nil)
}

// 打开数据库
func (s *DataStorage) open() error {
	if s.opened {
		return nil
	}

	s.mu.Lock()         // 锁
	defer s.mu.Unlock() // 解锁
	if s.opened {
		return nil
	}

	db, err := leveldb.OpenFile(s.dbPath, nil)
	if err != nil {
		cmn.Error("打开数据库失败：", s.dbPath)
		return err
	}
	s.leveldb = db

	s.opened = true
	cmn.Info("打开数据库")
	return nil
}

// 关闭数据库
func (s *DataStorage) Close() {
	if !s.opened {
		return
	}

	s.mu.Lock()         // 锁
	defer s.mu.Unlock() // 解锁
	if !s.opened {
		return
	}

	s.opened = false
	s.leveldb.Close()

	cmn.Info("关闭数据库")
}

// 是否关闭状态
func (s *DataStorage) IsOpen() bool {
	return s.opened
}
