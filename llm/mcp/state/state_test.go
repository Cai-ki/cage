package state_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Cai-ki/cage/llm/mcp/state"
)

func TestStateManagerSingleton(t *testing.T) {
	manager1 := state.GetInstance()
	manager2 := state.GetInstance()

	if manager1 != manager2 {
		t.Error("StateManager should be singleton")
	}
}

func TestStateManagerSetGet(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	key := "test_key"
	value := "test_value"

	manager.Set(key, value)
	retrieved, exists := manager.Get(key)

	if !exists {
		t.Error("Key should exist after setting")
	}

	if retrieved != value {
		t.Errorf("Expected %v, got %v", value, retrieved)
	}
}

func TestStateManagerDelete(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	key := "to_delete"
	value := "value"

	manager.Set(key, value)
	manager.Delete(key)

	_, exists := manager.Get(key)
	if exists {
		t.Error("Key should not exist after deletion")
	}
}

func TestStateManagerTTL(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	key := "temp_key"
	value := "temp_value"

	// 设置100毫秒过期
	manager.Set(key, value, state.WithTTL(100*time.Millisecond))

	// 立即获取应该存在
	_, exists := manager.Get(key)
	if !exists {
		t.Error("Key should exist immediately after setting")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 再次获取应该不存在
	_, exists = manager.Get(key)
	if exists {
		t.Error("Key should not exist after TTL expiration")
	}
}

func TestStateManagerConcurrentAccess(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	var wg sync.WaitGroup
	iterations := 1000

	// 并发写入
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			defer wg.Done()
			key := "key_" + string(rune(index))
			value := "value_" + string(rune(index))
			manager.Set(key, value)
		}(i)
	}

	// 并发读取
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			defer wg.Done()
			key := "key_" + string(rune(index))
			manager.Get(key)
		}(i)
	}

	wg.Wait()

	// 验证最终状态
	allStates := manager.GetAll()
	if len(allStates) != iterations {
		t.Errorf("Expected %d states, got %d", iterations, len(allStates))
	}
}

func TestStateManagerGetAll(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	// 设置多个状态
	states := map[string]interface{}{
		"string_key": "string_value",
		"int_key":    42,
		"bool_key":   true,
		"float_key":  3.14,
	}

	for key, value := range states {
		manager.Set(key, value)
	}

	// 获取所有状态
	allStates := manager.GetAll()

	if len(allStates) != len(states) {
		t.Errorf("Expected %d states, got %d", len(states), len(allStates))
	}

	for key, expectedValue := range states {
		actualValue, exists := allStates[key]
		if !exists {
			t.Errorf("Key %s should exist in GetAll result", key)
		}
		if actualValue != expectedValue {
			t.Errorf("For key %s, expected %v, got %v", key, expectedValue, actualValue)
		}
	}
}

func TestStateManagerCleanup(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	// 设置一些永久状态
	manager.Set("permanent1", "value1")
	manager.Set("permanent2", "value2")

	// 设置一些临时状态
	manager.Set("temp1", "value1", state.WithTTL(1*time.Millisecond))
	manager.Set("temp2", "value2", state.WithTTL(1*time.Millisecond))

	// 等待临时状态过期
	time.Sleep(10 * time.Millisecond)

	// 执行清理
	cleaned := manager.Cleanup()

	if cleaned != 2 {
		t.Errorf("Expected 2 cleaned items, got %d", cleaned)
	}

	// 验证永久状态仍然存在
	_, exists := manager.Get("permanent1")
	if !exists {
		t.Error("Permanent state should not be cleaned")
	}

	// 验证临时状态已被清理
	_, exists = manager.Get("temp1")
	if exists {
		t.Error("Temporary state should be cleaned")
	}
}

func TestStateManagerOnChangeCallback(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	// 使用通道来同步回调调用
	callbackCalled := make(chan bool, 1)
	var callbackOldValue, callbackNewValue interface{}

	// 设置变更回调
	onChange := func(old, new interface{}) {
		callbackOldValue = old
		callbackNewValue = new
		callbackCalled <- true
	}

	// 第一次设置（没有旧值）
	manager.Set("callback_key", "initial_value")

	// 第二次设置，应该触发回调 - 这次带上回调函数
	manager.Set("callback_key", "updated_value", state.WithOnChange(onChange))

	// 等待回调被调用或超时
	select {
	case <-callbackCalled:
		// 回调被调用，检查参数
		if callbackOldValue != "initial_value" {
			t.Errorf("Expected old value 'initial_value', got %v", callbackOldValue)
		}
		if callbackNewValue != "updated_value" {
			t.Errorf("Expected new value 'updated_value', got %v", callbackNewValue)
		}
	case <-time.After(1 * time.Second):
		t.Error("OnChange callback should be called within 1 second")
	}
}

// 修复后的测试 - 使用正确的回调设置顺序
func TestStateManagerOnChangeCallbackCorrect(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	// 使用通道来同步回调调用
	callbackCalled := make(chan bool, 1)
	var callbackOldValue, callbackNewValue interface{}

	// 设置变更回调
	onChange := func(old, new interface{}) {
		callbackOldValue = old
		callbackNewValue = new
		callbackCalled <- true
	}

	// 第一次设置时就设置回调
	manager.Set("callback_key", "initial_value", state.WithOnChange(onChange))

	// 第二次设置，应该触发回调
	manager.Set("callback_key", "updated_value")

	// 等待回调被调用或超时
	select {
	case <-callbackCalled:
		// 回调被调用，检查参数
		if callbackOldValue != "initial_value" {
			t.Errorf("Expected old value 'initial_value', got %v", callbackOldValue)
		}
		if callbackNewValue != "updated_value" {
			t.Errorf("Expected new value 'updated_value', got %v", callbackNewValue)
		}
	case <-time.After(1 * time.Second):
		t.Error("OnChange callback should be called within 1 second")
	}
}

func TestStateManagerGetWithType(t *testing.T) {
	manager := state.GetInstance()
	manager.Clear()

	manager.Set("string_key", "hello")
	manager.Set("int_key", 42)
	manager.Set("float_key", 3.14)

	// 测试类型匹配
	value, exists := manager.GetWithType("string_key", "string")
	if !exists || value != "hello" {
		t.Error("Should retrieve string value with correct type")
	}

	// 测试类型不匹配
	_, exists = manager.GetWithType("string_key", "int")
	if exists {
		t.Error("Should not retrieve value when type doesn't match")
	}

	// 测试不存在的键
	_, exists = manager.GetWithType("nonexistent", "string")
	if exists {
		t.Error("Should not retrieve nonexistent key")
	}
}

func BenchmarkStateManagerSet(b *testing.B) {
	manager := state.GetInstance()
	manager.Clear()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key_" + string(rune(i))
		manager.Set(key, "value")
	}
}

func BenchmarkStateManagerGet(b *testing.B) {
	manager := state.GetInstance()
	manager.Clear()

	// 先设置一些数据
	for i := 0; i < 1000; i++ {
		key := "key_" + string(rune(i))
		manager.Set(key, "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key_" + string(rune(i%1000))
		manager.Get(key)
	}
}
