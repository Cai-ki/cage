# 包功能说明

jsondb 是一个轻量级的 JSON 文档数据库，提供基于内存的快速数据存储和查询功能，支持自动持久化到本地文件系统。该包采用 Go 语言编写，设计目标是为小型应用提供简单易用的数据存储解决方案，无需依赖外部数据库服务。核心特性包括记录按时间戳自动排序、条件查询、范围查询、数据更新和删除等操作，所有操作都是并发安全的。典型使用场景包括日志记录、小型应用配置存储、临时数据缓存等需要简单持久化功能的场合。

## 结构体与接口

```go
type Record struct {
	Timestamp time.Time
	RawData   []byte
}
```

Record 结构体表示数据库中的单条记录。Timestamp 字段是记录添加时的精确时间戳，由系统自动生成，用于内部排序和查询；RawData 字段存储用户原始 JSON 数据，保持数据的原始格式。

```go
type Database struct {
	// 包含未导出字段：mu, records, filePath, autoSave
}
```

Database 是核心数据库结构，封装了所有数据操作功能。通过读写锁保证并发安全，内部记录按时间戳升序排列，支持自动持久化到指定文件路径。

```go
type Option func(*Database)
```

Option 是配置选项函数类型，用于在创建 Database 实例时提供可选的配置参数，支持自定义数据库行为。

## 函数

```go
func NewDatabase(filePath string, opts ...Option) (*Database, error)
```

NewDatabase 创建新的数据库实例。filePath 参数指定数据文件的存储路径，opts 是可选配置项。函数会自动创建文件所在目录，如果文件已存在则会加载现有数据，返回初始化后的数据库实例或错误信息。

```go
func WithAutoSave(enable bool) Option
```

WithAutoSave 设置是否启用自动持久化功能。当 enable 为 true 时，每次数据修改后都会立即保存到文件；为 false 时需要手动调用 Save 方法。默认情况下自动保存功能是开启的。

```go
func WithInitialLoad(enable bool) Option
```

WithInitialLoad 设置是否在创建数据库时从文件加载现有数据。当 enable 为 false 时会禁用初始加载，用于创建全新的数据库实例。默认情况下会尝试加载现有数据。

```go
func (db *Database) Add(record interface{}) error
```

Add 方法向数据库添加新记录。record 参数可以是任意可序列化为 JSON 的结构体，方法会自动为记录添加当前时间戳，并按时间顺序插入到内存中的记录列表。如果启用了自动保存，会立即持久化到文件。

```go
func (db *Database) GetLatest(n int, result interface{}) error
```

GetLatest 方法获取最近 N 条记录。n 参数指定要获取的记录数量，result 必须是指向切片的指针，用于接收反序列化后的结果。如果请求的数量超过总记录数，会返回全部可用记录。

```go
func (db *Database) GetByTimeRange(start, end time.Time, result interface{}) error
```

GetByTimeRange 方法获取指定时间范围内的记录。start 和 end 参数定义时间范围（包含 start，不包含 end），result 必须是指向切片的指针。方法使用二分查找高效定位时间范围内的记录。

```go
func (db *Database) GetByCondition(condition func(*Record) bool, result interface{}) error
```

GetByCondition 方法获取满足自定义条件的记录。condition 是用户定义的过滤函数，对每条记录进行评估，返回 true 表示包含该记录。result 必须是指向切片的指针，用于接收筛选结果。

```go
func (db *Database) DeleteByCondition(condition func(*Record) bool) error
```

DeleteByCondition 方法删除满足自定义条件的记录。condition 是用户定义的过滤函数，返回 true 的记录将被删除。如果启用了自动保存，删除操作后会立即持久化到文件。

```go
func (db *Database) UpdateByCondition(condition func(*Record) bool, updateFunc func(interface{}) interface{}) error
```

UpdateByCondition 方法更新满足条件的记录。condition 函数识别需要更新的记录，updateFunc 函数接收反序列化后的数据并返回更新后的数据。记录的时间戳保持不变，只有数据内容被更新。

```go
func (db *Database) Count(condition func(*Record) bool) int
```

Count 方法返回满足条件的记录数量。condition 是过滤函数，方法遍历所有记录并统计满足条件的记录数，返回整型计数结果。

```go
func (db *Database) Exists(condition func(*Record) bool) bool
```

Exists 方法检查是否存在满足条件的记录。condition 是过滤函数，方法在找到第一条满足条件的记录时立即返回 true，如果遍历所有记录都没有找到则返回 false。

```go
func (db *Database) First(condition func(*Record) bool, result interface{}) error
```

First 方法获取第一个满足条件的记录。condition 是过滤函数，result 必须是指向单个结构体的指针（不是切片）。方法找到第一条满足条件的记录后立即返回，不会继续查找后续记录。

```go
func (db *Database) DeleteBefore(t time.Time) error
```

DeleteBefore 方法删除指定时间之前的所有记录。t 参数是时间阈值，所有时间戳早于该时间的记录都会被删除。方法使用二分查找高效定位删除位置，保持剩余记录的时间顺序。

```go
func (db *Database) DeleteAll() error
```

DeleteAll 方法清空数据库中的所有记录。如果启用了自动保存，会立即将空状态持久化到文件，相当于重置整个数据库。

```go
func (db *Database) Save() error
```

Save 方法手动将当前内存中的数据持久化到文件。无论是否启用自动保存，都可以调用此方法强制保存数据，适用于批量操作后的一次性保存场景。

```go
func (db *Database) Close() error
```

Close 方法保存并关闭数据库。这是数据库的清理方法，确保所有数据都被正确持久化，通常在程序退出前调用。

## 变量与常量

```go
var ErrInvalidResultType = &jsonError{"result must be a pointer to slice"}
```

ErrInvalidResultType 是预定义错误变量，当查询方法的 result 参数不是指向切片的指针时返回此错误。用于提示用户正确使用查询接口。

```go
var ErrMissingTimeField = &jsonError{"time field missing in JSON"}
```

ErrMissingTimeField 是预定义错误变量，在需要时间字段但数据中缺少该字段时返回。虽然当前代码中未直接使用，但作为包的错误类型定义保留。

```go
var ErrInvalidTimeFormat = &jsonError{"invalid time format"}
```

ErrInvalidTimeFormat 是预定义错误变量，在时间字段格式不正确时返回。虽然当前代码中未直接使用，但作为包的错误类型定义保留。