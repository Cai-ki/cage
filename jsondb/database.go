package jsondb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Record 表示数据库中的单条记录
type Record struct {
	Timestamp time.Time // 内部时间戳（记录添加时的精确时间，非用户数据中的时间）
	RawData   []byte    // 用户原始 JSON 数据
}

// Database 是核心数据库结构
type Database struct {
	mu       sync.RWMutex // 读写锁，保证并发安全
	records  []*Record    // 内存中的记录列表，按时间戳升序排列
	filePath string       // 数据文件的存储路径
	autoSave bool         // 是否自动持久化（每次修改后立即保存）
}

// NewDatabase 创建新数据库实例
func NewDatabase(filePath string, opts ...Option) (*Database, error) {
	db := &Database{
		filePath: filepath.Clean(filePath), // 清理文件路径
		autoSave: true,                     // 默认开启自动保存
		records:  make([]*Record, 0),       // 初始化空记录列表
	}

	// 应用用户传入的配置选项
	for _, opt := range opts {
		opt(db)
	}

	// 确保文件所在目录存在
	dir := filepath.Dir(db.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 从文件加载已存在的数据
	if err := db.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return db, nil
}

// Add 添加新记录（自动持久化）
func (db *Database) Add(record interface{}) error {
	// 将用户传入的结构体序列化为 JSON 字节数组
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// 创建新的记录，使用当前时间作为时间戳
	newRec := &Record{
		Timestamp: time.Now(), // 自动获取当前时间，不依赖用户数据中的时间字段
		RawData:   data,       // 保存用户数据的原始 JSON
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// 二分查找插入位置，保持按时间戳升序排列
	// 找到第一个时间戳大于 newRec.Timestamp 的位置
	idx := sort.Search(len(db.records), func(i int) bool {
		return db.records[i].Timestamp.After(newRec.Timestamp)
	})

	// 在指定位置插入新记录
	db.records = append(db.records, nil)       // 扩容切片
	copy(db.records[idx+1:], db.records[idx:]) // 向后移动元素
	db.records[idx] = newRec                   // 插入新记录

	// 如果启用自动保存，则立即持久化到文件
	if db.autoSave {
		return db.saveUnsafe()
	}
	return nil
}

// GetLatest 获取最近 N 条记录
func (db *Database) GetLatest(n int, result interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 检查是否有足够记录
	count := len(db.records)
	if n > count {
		n = count // 如果请求的数量超过总记录数，则返回全部记录
	}
	if n <= 0 {
		return nil // 如果请求的数量小于等于0，返回空结果
	}

	// 由于记录按时间戳升序排列，最后 n 条即为最近的 n 条
	// 例如：records = [t1, t2, t3, t4, t5]，取最新的2条，得到 [t4, t5]
	latestRecords := db.records[count-n:]

	// 将原始 JSON 数据反序列化到用户提供的结果变量中
	return unmarshalRecords(latestRecords, result)
}

// GetByTimeRange 获取指定时间范围内的记录
func (db *Database) GetByTimeRange(start, end time.Time, result interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 二分查找起始时间点在记录列表中的索引
	// 找到第一个时间戳 >= start 的位置
	startIdx := sort.Search(len(db.records), func(i int) bool {
		return !db.records[i].Timestamp.Before(start) // 等价于 db.records[i].Timestamp >= start
	})

	// 二分查找结束时间点之后的索引
	// 找到第一个时间戳 > end 的位置
	endIdx := sort.Search(len(db.records), func(i int) bool {
		return db.records[i].Timestamp.After(end) // 等价于 db.records[i].Timestamp > end
	})

	// 如果起始索引大于等于结束索引，说明没有记录在该时间范围内
	if startIdx >= endIdx {
		return nil
	}

	// 提取时间范围内的记录并反序列化
	return unmarshalRecords(db.records[startIdx:endIdx], result)
}

// GetByCondition 获取满足条件的记录
func (db *Database) GetByCondition(condition func(*Record) bool, result interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 筛选满足条件的记录
	var filteredRecords []*Record
	for _, rec := range db.records {
		if condition(rec) {
			filteredRecords = append(filteredRecords, rec)
		}
	}

	// 将筛选后的记录反序列化到结果变量
	return unmarshalRecords(filteredRecords, result)
}

// DeleteByCondition 删除满足条件的记录
func (db *Database) DeleteByCondition(condition func(*Record) bool) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 筛选不满足条件的记录
	remainingRecords := make([]*Record, 0)
	for _, rec := range db.records {
		if !condition(rec) {
			remainingRecords = append(remainingRecords, rec)
		}
	}

	db.records = remainingRecords

	// 如果启用自动保存，持久化修改
	if db.autoSave {
		return db.saveUnsafe()
	}
	return nil
}

// UpdateByCondition 更新满足条件的记录
func (db *Database) UpdateByCondition(condition func(*Record) bool, updateFunc func(interface{}) interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	updated := false
	for i, rec := range db.records {
		if condition(rec) {
			// 反序列化当前记录
			var tempData interface{}
			if err := json.Unmarshal(rec.RawData, &tempData); err != nil {
				return err
			}

			// 应用更新函数
			updatedData := updateFunc(tempData)

			// 重新序列化
			newRawData, err := json.Marshal(updatedData)
			if err != nil {
				return err
			}

			// 更新记录，保持时间戳不变
			db.records[i] = &Record{
				Timestamp: rec.Timestamp,
				RawData:   newRawData,
			}
			updated = true
		}
	}

	// 如果有更新，保存到文件
	if updated && db.autoSave {
		return db.saveUnsafe()
	}
	return nil
}

// Count 返回满足条件的记录数量
func (db *Database) Count(condition func(*Record) bool) int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	count := 0
	for _, rec := range db.records {
		if condition(rec) {
			count++
		}
	}
	return count
}

// Exists 检查是否存在满足条件的记录
func (db *Database) Exists(condition func(*Record) bool) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, rec := range db.records {
		if condition(rec) {
			return true
		}
	}
	return false
}

// First 获取第一个满足条件的记录
func (db *Database) First(condition func(*Record) bool, result interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, rec := range db.records {
		if condition(rec) {
			// 直接反序列化到结果变量，不使用unmarshalRecords
			return json.Unmarshal(rec.RawData, result)
		}
	}
	return nil // 没有找到符合条件的记录
}

// DeleteBefore 删除指定时间之前的所有记录
func (db *Database) DeleteBefore(t time.Time) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 二分查找第一个不小于指定时间的记录索引
	// 例如：records = [t1, t2, t3, t4, t5]，t = t3，则找到 t3 的位置
	idx := sort.Search(len(db.records), func(i int) bool {
		return !db.records[i].Timestamp.Before(t) // 等价于 db.records[i].Timestamp >= t
	})

	// 保留从 idx 开始的所有记录，删除 idx 之前的记录
	db.records = db.records[idx:]

	// 如果启用自动保存，持久化修改
	if db.autoSave {
		return db.saveUnsafe()
	}
	return nil
}

// DeleteAll 清空所有记录
func (db *Database) DeleteAll() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 清空记录列表
	db.records = nil

	// 如果启用自动保存，持久化修改（即保存一个空数组）
	if db.autoSave {
		return db.saveUnsafe()
	}
	return nil
}

// Save 手动持久化到文件
func (db *Database) Save() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.saveUnsafe()
}

// Close 保存并关闭数据库
func (db *Database) Close() error {
	return db.Save()
}

// load 从文件加载数据（内部使用）
func (db *Database) load() error {
	data, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	// 如果文件为空，初始化为空数组
	if len(data) == 0 {
		db.mu.Lock()
		db.records = make([]*Record, 0)
		db.mu.Unlock()
		return nil
	}

	// 解析文件内容为内部格式（包含时间戳）
	type internalRecord struct {
		Timestamp time.Time       `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	}

	var internalRecords []internalRecord
	if err := json.Unmarshal(data, &internalRecords); err != nil {
		return err
	}

	// 重建内存中的记录列表，保持时间戳信息
	records := make([]*Record, len(internalRecords))
	for i, internalRec := range internalRecords {
		records[i] = &Record{
			Timestamp: internalRec.Timestamp,
			RawData:   []byte(internalRec.Data),
		}
	}

	// 按时间戳排序（确保数据一致性）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.Before(records[j].Timestamp)
	})

	db.mu.Lock()
	db.records = records
	db.mu.Unlock()
	return nil
}

// saveUnsafe 不加锁的保存方法（内部使用）
func (db *Database) saveUnsafe() error {
	// 创建包含时间戳的内部格式
	type internalRecord struct {
		Timestamp time.Time       `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	}

	internalRecords := make([]internalRecord, len(db.records))
	for i, rec := range db.records {
		internalRecords[i] = internalRecord{
			Timestamp: rec.Timestamp,
			Data:      json.RawMessage(rec.RawData),
		}
	}

	buf, err := json.Marshal(internalRecords)
	if err != nil {
		return err
	}

	return os.WriteFile(db.filePath, buf, 0644)
}
