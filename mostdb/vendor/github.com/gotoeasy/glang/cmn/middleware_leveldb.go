package cmn

import (
	"os"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDB结构体
type LevelDB struct {
	dbPath        string      // 数据存放目录
	leveldb       *leveldb.DB // leveldb
	opened        bool        // 是否打开状态
	lastTime      int64       // 最后一次访问时间
	maxIdleSecond int64       // 最大空闲时间，达空闲时间后自动关闭，小于等于0时不自动关闭
	mu            sync.Mutex  // 锁
}

// LevelDB选项
type OptionLevelDB struct {
	// 最大空闲时间（分钟，0~60），达空闲时间后自动关闭，小于等于0时不自动关闭，默认5分钟
	MaxIdleMinute int
}

var _leveldbClient *LevelDB     // 客户端实例
var _leveldbClientMu sync.Mutex // 锁

// 创建LevelDB对象，参数dbPath为数据库名目录
// 实际每次调用都返回同一对象，opt仅首次调用时有效，opt为nil时使用默认值
func NewLevelDB(dbPath string, opt *OptionLevelDB) *LevelDB {

	// 确保仅一个客户端实例
	if _leveldbClient != nil {
		return _leveldbClient
	}
	_leveldbClientMu.Lock()         // 锁
	defer _leveldbClientMu.Unlock() // 解锁
	if _leveldbClient != nil {
		return _leveldbClient
	}

	// 开始创建客户端实例
	maxIdleTime := int64(5 * 60) // 默认5分钟没操作则自动关闭
	if opt != nil {
		// 最大空闲时间，不超60分钟
		if opt.MaxIdleMinute > 60 {
			maxIdleTime = int64(60 * 60)
		} else {
			maxIdleTime = int64(opt.MaxIdleMinute * 60)
		}
	}
	db := &LevelDB{
		dbPath:        dbPath,
		maxIdleSecond: maxIdleTime,
	}

	err := os.MkdirAll(db.dbPath, 0777)
	if err != nil {
		Error("创建LevelDB对象时发生错误", dbPath, err)
	}

	_leveldbClient = db
	return db
}

// 保存
func (s *LevelDB) Put(key []byte, value []byte) error {
	s.lastTime = time.Now().Unix()
	if !s.opened {
		s.Open() // 自动打开数据库
	}
	return s.leveldb.Put(key, value, nil)
}

// 删除
func (s *LevelDB) Del(key []byte) error {
	s.lastTime = time.Now().Unix()
	if !s.opened {
		s.Open() // 自动打开数据库
	}
	return s.leveldb.Delete(key, nil)
}

// 获取
func (s *LevelDB) Get(key []byte) ([]byte, error) {
	s.lastTime = time.Now().Unix()
	if !s.opened {
		s.Open() // 自动打开数据库
	}
	return s.leveldb.Get(key, nil)
}

// 快照
func (s *LevelDB) GetSnapshot() (*leveldb.Snapshot, error) {
	s.lastTime = time.Now().Unix()
	if !s.opened {
		s.Open() // 自动打开数据库
	}
	return s.leveldb.GetSnapshot()
}

func (s *LevelDB) autoCloseWhenMaxIdle() {
	if s.maxIdleSecond > 0 {
		ticker := time.NewTicker(time.Minute) // 每分钟判断一次
		for {
			<-ticker.C
			if time.Now().Unix()-s.lastTime > s.maxIdleSecond {
				s.Close()
				ticker.Stop()
				break
			}
		}
	}
}

// 打开数据库
func (s *LevelDB) Open() error {
	s.lastTime = time.Now().Unix()
	if s.opened {
		return nil
	}

	s.mu.Lock()         // 锁
	defer s.mu.Unlock() // 解锁
	if s.opened {
		return nil
	}

	option := new(opt.Options)                    // leveldb选项
	option.Filter = filter.NewBloomFilter(10)     // 使用布隆过滤器
	db, err := leveldb.OpenFile(s.dbPath, option) // 打开数据库
	if err != nil {
		Error("打开数据库失败：", s.dbPath)
		return err
	}
	s.leveldb = db

	s.opened = true
	go s.autoCloseWhenMaxIdle() // 开启自动关闭

	Info("打开数据库")
	return nil
}

// 关闭数据库
func (s *LevelDB) Close() {
	if !s.opened {
		return
	}

	s.mu.Lock()         // 锁
	defer s.mu.Unlock() // 解锁
	if !s.opened {
		return
	}

	s.opened = false
	err := s.leveldb.Close()
	if err != nil {
		Error("关闭数据库失败", err)
	} else {
		Info("关闭数据库")
	}

}
