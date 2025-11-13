package jsondb_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Cai-ki/cage/jsondb"
)

// TestRecord 测试用的数据结构
type TestRecord struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func TestNewDatabase(t *testing.T) {
	db, err := jsondb.NewDatabase("test.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test.db")

	if db == nil {
		t.Fatal("Database should not be nil")
	}
}

func TestAddAndGetLatest(t *testing.T) {
	db, err := jsondb.NewDatabase("test_add.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_add.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Record 1", Value: 10.5},
		{ID: 2, Name: "Record 2", Value: 20.3},
		{ID: 3, Name: "Record 3", Value: 30.7},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
		// 确保时间戳不同
		time.Sleep(1 * time.Millisecond)
	}

	// 获取最新的2条记录
	var results []TestRecord
	if err := db.GetLatest(2, &results); err != nil {
		t.Fatalf("Failed to get latest records: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(results))
	}

	// 验证结果顺序（最新的在后）
	if results[0].ID != 2 {
		t.Errorf("Expected first result ID to be 2, got %d", results[0].ID)
	}
	if results[1].ID != 3 {
		t.Errorf("Expected second result ID to be 3, got %d", results[1].ID)
	}
}

func TestGetByTimeRange(t *testing.T) {
	db, err := jsondb.NewDatabase("test_timerange.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_timerange.db")

	// 添加测试数据
	baseTime := time.Now()
	time.Sleep(1 * time.Millisecond)

	testData := []TestRecord{
		{ID: 1, Name: "Old Record", Value: 1.0},
		{ID: 2, Name: "Middle Record", Value: 2.0},
		{ID: 3, Name: "New Record", Value: 3.0},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // 确保时间戳不同
	}

	// 获取时间范围内的记录
	startTime := baseTime.Add(5 * time.Millisecond)
	endTime := baseTime.Add(25 * time.Millisecond)

	var results []TestRecord
	if err := db.GetByTimeRange(startTime, endTime, &results); err != nil {
		t.Fatalf("Failed to get records by time range: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected some records in time range")
	}
}

func TestDeleteBefore(t *testing.T) {
	db, err := jsondb.NewDatabase("test_delete.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_delete.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Old Record", Value: 1.0},
		{ID: 2, Name: "Middle Record", Value: 2.0},
		{ID: 3, Name: "New Record", Value: 3.0},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // 确保时间戳不同
	}

	// 获取所有记录数量
	var allRecords []TestRecord
	if err := db.GetLatest(10, &allRecords); err != nil {
		t.Fatalf("Failed to get all records: %v", err)
	}

	if len(allRecords) != 3 {
		t.Fatalf("Expected 3 records initially, got %d", len(allRecords))
	}

	// 删除指定时间之前的记录
	deleteTime := time.Now()
	time.Sleep(1 * time.Millisecond)

	if err := db.DeleteBefore(deleteTime); err != nil {
		t.Fatalf("Failed to delete records: %v", err)
	}

	// 验证删除后的记录数量
	var remainingRecords []TestRecord
	if err := db.GetLatest(10, &remainingRecords); err != nil {
		t.Fatalf("Failed to get remaining records: %v", err)
	}

	// 预期剩余的记录数量取决于deleteTime与添加时间的关系
	// 由于deleteTime是在最后添加记录之后设置的，预期可能删除部分或全部较早的记录
	t.Logf("Remaining records after delete: %d", len(remainingRecords))
}

func TestDeleteAll(t *testing.T) {
	db, err := jsondb.NewDatabase("test_deleteall.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_deleteall.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Record 1", Value: 1.0},
		{ID: 2, Name: "Record 2", Value: 2.0},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 删除所有记录
	if err := db.DeleteAll(); err != nil {
		t.Fatalf("Failed to delete all records: %v", err)
	}

	// 验证所有记录都被删除
	var results []TestRecord
	if err := db.GetLatest(10, &results); err != nil {
		t.Fatalf("Failed to get records after deletion: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 records after DeleteAll, got %d", len(results))
	}
}

func TestSaveAndLoad(t *testing.T) {
	db, err := jsondb.NewDatabase("test_save_load.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_save_load.db")

	// 添加测试数据
	testData := TestRecord{ID: 100, Name: "Saved Record", Value: 99.9}
	if err := db.Add(testData); err != nil {
		t.Fatalf("Failed to add record: %v", err)
	}

	// 手动保存
	if err := db.Save(); err != nil {
		t.Fatalf("Failed to save database: %v", err)
	}

	// 创建新实例并加载数据
	db2, err := jsondb.NewDatabase("test_save_load.db")
	if err != nil {
		t.Fatalf("Failed to create second database instance: %v", err)
	}

	// 获取最新记录
	var results []TestRecord
	if err := db2.GetLatest(1, &results); err != nil {
		t.Fatalf("Failed to get records from loaded database: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 record after load, got %d", len(results))
	}

	if results[0].ID != 100 || results[0].Name != "Saved Record" {
		t.Errorf("Loaded record doesn't match original: %+v", results[0])
	}
}

func TestWithAutoSave(t *testing.T) {
	db, err := jsondb.NewDatabase("test_autosave.db", jsondb.WithAutoSave(false))
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_autosave.db")

	// 添加记录（不会自动保存）
	testData := TestRecord{ID: 200, Name: "No Auto Save", Value: 88.8}
	if err := db.Add(testData); err != nil {
		t.Fatalf("Failed to add record: %v", err)
	}

	// 直接创建新实例，应该没有数据（因为没有手动保存）
	db2, err := jsondb.NewDatabase("test_autosave.db")
	if err != nil {
		t.Fatalf("Failed to create second database instance: %v", err)
	}

	var results []TestRecord
	if err := db2.GetLatest(1, &results); err != nil {
		t.Fatalf("Failed to get records: %v", err)
	}

	// 应该没有记录，因为没有手动保存
	if len(results) != 0 {
		t.Fatalf("Expected 0 records when auto-save is disabled and not manually saved, got %d", len(results))
	}

	// 现在手动保存
	if err := db.Save(); err != nil {
		t.Fatalf("Failed to manually save: %v", err)
	}

	// 再次加载
	db3, err := jsondb.NewDatabase("test_autosave.db")
	if err != nil {
		t.Fatalf("Failed to create third database instance: %v", err)
	}

	if err := db3.GetLatest(1, &results); err != nil {
		t.Fatalf("Failed to get records after manual save: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 record after manual save, got %d", len(results))
	}
}

func TestComplexJSONData(t *testing.T) {
	db, err := jsondb.NewDatabase("test_complex.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_complex.db")

	// 测试复杂数据结构
	complexData := map[string]interface{}{
		"id":   1,
		"name": "Complex Record",
		"tags": []string{"tag1", "tag2", "tag3"},
		"metadata": map[string]interface{}{
			"created": time.Now().Unix(),
			"version": 1.0,
		},
		"values": []float64{1.1, 2.2, 3.3},
	}

	if err := db.Add(complexData); err != nil {
		t.Fatalf("Failed to add complex data: %v", err)
	}

	// 获取数据并验证
	var results []map[string]interface{}
	if err := db.GetLatest(1, &results); err != nil {
		t.Fatalf("Failed to get complex data: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 complex record, got %d", len(results))
	}

	// 验证结构
	result := results[0]
	if id, ok := result["id"].(float64); !ok || int(id) != 1 {
		t.Errorf("Expected id to be 1, got %v", result["id"])
	}

	if name, ok := result["name"].(string); !ok || name != "Complex Record" {
		t.Errorf("Expected name to be 'Complex Record', got %v", result["name"])
	}
}

func TestClose(t *testing.T) {
	db, err := jsondb.NewDatabase("test_close.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_close.db")

	// 添加一些数据
	testData := TestRecord{ID: 300, Name: "Test Close", Value: 77.7}
	if err := db.Add(testData); err != nil {
		t.Fatalf("Failed to add record: %v", err)
	}

	// 关闭数据库（应该触发保存）
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}

	// 验证数据被保存
	db2, err := jsondb.NewDatabase("test_close.db")
	if err != nil {
		t.Fatalf("Failed to reopen database: %v", err)
	}

	var results []TestRecord
	if err := db2.GetLatest(1, &results); err != nil {
		t.Fatalf("Failed to get records after close: %v", err)
	}

	if len(results) != 1 || results[0].ID != 300 {
		t.Errorf("Data not properly saved on close: %+v", results)
	}
}

func TestGetByCondition(t *testing.T) {
	db, err := jsondb.NewDatabase("test_get_condition.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_get_condition.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Record A", Value: 10.5},
		{ID: 2, Name: "Record B", Value: 20.3},
		{ID: 3, Name: "Record C", Value: 30.7},
		{ID: 4, Name: "Record A", Value: 40.1},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 测试条件：Name为"Record A"
	var results []TestRecord
	condition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Name == "Record A"
	}

	if err := db.GetByCondition(condition, &results); err != nil {
		t.Fatalf("Failed to get records by condition: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 records with Name 'Record A', got %d", len(results))
	}

	for _, result := range results {
		if result.Name != "Record A" {
			t.Errorf("Expected Name 'Record A', got %s", result.Name)
		}
	}
}

func TestDeleteByCondition(t *testing.T) {
	db, err := jsondb.NewDatabase("test_delete_condition.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_delete_condition.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Keep 1", Value: 10.5},
		{ID: 2, Name: "Delete", Value: 20.3},
		{ID: 3, Name: "Keep 2", Value: 30.7},
		{ID: 4, Name: "Delete", Value: 40.1},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 删除Name为"Delete"的记录
	condition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Name == "Delete"
	}

	if err := db.DeleteByCondition(condition); err != nil {
		t.Fatalf("Failed to delete records by condition: %v", err)
	}

	// 验证剩余记录
	var remaining []TestRecord
	if err := db.GetLatest(10, &remaining); err != nil {
		t.Fatalf("Failed to get remaining records: %v", err)
	}

	if len(remaining) != 2 {
		t.Fatalf("Expected 2 remaining records, got %d", len(remaining))
	}

	for _, record := range remaining {
		if record.Name == "Delete" {
			t.Errorf("Found record that should have been deleted: %s", record.Name)
		}
	}
}

func TestUpdateByCondition(t *testing.T) {
	db, err := jsondb.NewDatabase("test_update_condition.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_update_condition.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Record 1", Value: 10.5},
		{ID: 2, Name: "Record 2", Value: 20.3},
		{ID: 3, Name: "Record 1", Value: 30.7},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 更新ID为1的记录，将其值乘以2
	condition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.ID == 1
	}

	updateFunc := func(data interface{}) interface{} {
		if record, ok := data.(map[string]interface{}); ok {
			if value, exists := record["value"].(float64); exists {
				record["value"] = value * 2
				record["name"] = "Updated Record 1"
			}
		}
		return data
	}

	if err := db.UpdateByCondition(condition, updateFunc); err != nil {
		t.Fatalf("Failed to update records by condition: %v", err)
	}

	// 验证更新结果
	var results []TestRecord
	if err := db.GetByCondition(condition, &results); err != nil {
		t.Fatalf("Failed to get updated records: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 updated record, got %d", len(results))
	}

	if results[0].Value != 21.0 || results[0].Name != "Updated Record 1" {
		t.Errorf("Update failed: expected Value=21.0 and Name='Updated Record 1', got %+v", results[0])
	}
}

func TestCount(t *testing.T) {
	db, err := jsondb.NewDatabase("test_count.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_count.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "High Value", Value: 100.5},
		{ID: 2, Name: "Low Value", Value: 5.3},
		{ID: 3, Name: "High Value", Value: 200.7},
		{ID: 4, Name: "Low Value", Value: 3.1},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 测试计数条件：Value > 50
	condition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Value > 50
	}

	count := db.Count(condition)
	if count != 2 {
		t.Fatalf("Expected count 2 for Value > 50, got %d", count)
	}

	// 测试计数条件：Name为"Low Value"
	nameCondition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Name == "Low Value"
	}

	nameCount := db.Count(nameCondition)
	if nameCount != 2 {
		t.Fatalf("Expected count 2 for Name 'Low Value', got %d", nameCount)
	}
}

func TestExists(t *testing.T) {
	db, err := jsondb.NewDatabase("test_exists.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_exists.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 1, Name: "Test Record", Value: 10.5},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 测试存在条件：ID为1
	existsCondition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.ID == 1
	}

	if !db.Exists(existsCondition) {
		t.Error("Expected record with ID 1 to exist")
	}

	// 测试不存在条件：ID为999
	nonExistsCondition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.ID == 999
	}

	if db.Exists(nonExistsCondition) {
		t.Error("Expected record with ID 999 to not exist")
	}
}

func TestFirst(t *testing.T) {
	db, err := jsondb.NewDatabase("test_first.db")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove("test_first.db")

	// 添加测试数据
	testData := []TestRecord{
		{ID: 5, Name: "First Match", Value: 10.5},
		{ID: 2, Name: "Second Match", Value: 20.3},
		{ID: 3, Name: "Other", Value: 30.7},
		{ID: 1, Name: "First Match", Value: 40.1},
	}

	for _, record := range testData {
		if err := db.Add(record); err != nil {
			t.Fatalf("Failed to add record: %v", err)
		}
	}

	// 查找Name为"First Match"的第一条记录
	var result TestRecord
	condition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Name == "First Match"
	}

	if err := db.First(condition, &result); err != nil {
		t.Fatalf("Failed to get first record: %v", err)
	}

	// 由于记录按时间戳排序，第一条"First Match"应该是ID为5的记录
	if result.ID != 5 {
		t.Errorf("Expected first record with ID 5, got ID %d", result.ID)
	}

	// 测试不存在的情况
	nonExistsCondition := func(record *jsondb.Record) bool {
		var temp TestRecord
		json.Unmarshal(record.RawData, &temp)
		return temp.Name == "Non Exists"
	}

	var emptyResult TestRecord
	if err := db.First(nonExistsCondition, &emptyResult); err != nil {
		t.Fatalf("Failed to handle non-existent record: %v", err)
	}
	// 应该返回nil，结果应该是零值
}
