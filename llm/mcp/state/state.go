package state

import (
	"reflect"
	"sync"
	"time"
)

// StateManager 通用状态管理器
type StateManager struct {
	mu     sync.RWMutex
	states map[string]*StateEntry
}

// StateEntry 状态条目
type StateEntry struct {
	Value     interface{}                `json:"value"`
	Type      string                     `json:"type"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	ExpiresAt *time.Time                 `json:"expires_at,omitempty"`
	OnChange  func(old, new interface{}) `json:"-"`
}

var (
	instance *StateManager
	once     sync.Once
)

// GetInstance 返回单例实例
func GetInstance() *StateManager {
	once.Do(func() {
		instance = &StateManager{
			states: make(map[string]*StateEntry),
		}
	})
	return instance
}

// Set 设置状态值 - 修复版本
func (s *StateManager) Set(key string, value interface{}, options ...StateOption) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// 创建新条目
	entry := &StateEntry{
		Value:     value,
		Type:      getType(value),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 获取旧值（如果存在）
	var oldValue interface{}
	oldEntry, exists := s.states[key]
	if exists {
		oldValue = oldEntry.Value
		entry.CreatedAt = oldEntry.CreatedAt
		entry.OnChange = oldEntry.OnChange
	}

	// 应用选项
	for _, option := range options {
		option(entry)
	}

	// 先保存新状态
	s.states[key] = entry

	// 然后在锁外调用回调（避免死锁）
	if exists && entry.OnChange != nil {
		// 注意：这里在锁外调用回调，避免回调函数中再次调用状态管理器导致死锁
		go func(oldVal, newVal interface{}) {
			entry.OnChange(oldVal, newVal)
		}(oldValue, value)
	}
}

// Get 获取状态值
func (s *StateManager) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.states[key]
	if !exists {
		return nil, false
	}

	// 检查过期
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

// GetWithType 获取状态值并检查类型
func (s *StateManager) GetWithType(key string, expectedType string) (interface{}, bool) {
	value, exists := s.Get(key)
	if !exists {
		return nil, false
	}

	actualType := getType(value)
	if actualType != expectedType {
		return nil, false
	}

	return value, true
}

// Delete 删除状态
func (s *StateManager) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, key)
}

// Exists 检查状态是否存在
func (s *StateManager) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.states[key]
	if !exists {
		return false
	}

	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return false
	}

	return true
}

// GetAll 获取所有状态
func (s *StateManager) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]interface{})
	now := time.Now()

	for key, entry := range s.states {
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
			continue
		}
		result[key] = entry.Value
	}

	return result
}

// Cleanup 清理过期状态
func (s *StateManager) Cleanup() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	now := time.Now()

	for key, entry := range s.states {
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
			delete(s.states, key)
			count++
		}
	}

	return count
}

// Clear 清空所有状态（用于测试）
func (s *StateManager) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states = make(map[string]*StateEntry)
}

// StateOption 状态选项函数
type StateOption func(*StateEntry)

// WithTTL 设置过期时间
func WithTTL(ttl time.Duration) StateOption {
	return func(entry *StateEntry) {
		expiresAt := time.Now().Add(ttl)
		entry.ExpiresAt = &expiresAt
	}
}

// WithOnChange 设置变更回调
func WithOnChange(callback func(old, new interface{})) StateOption {
	return func(entry *StateEntry) {
		entry.OnChange = callback
	}
}

// getType 获取值的类型
func getType(value interface{}) string {
	if value == nil {
		return "nil"
	}
	return reflect.TypeOf(value).String()
}
