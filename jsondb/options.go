package jsondb

// Option 配置选项函数类型
type Option func(*Database)

// WithAutoSave 设置是否自动持久化
func WithAutoSave(enable bool) Option {
	return func(db *Database) {
		db.autoSave = enable
	}
}

// WithInitialLoad 禁用初始加载（用于全新数据库）
func WithInitialLoad(enable bool) Option {
	return func(db *Database) {
		if !enable {
			db.records = nil
		}
	}
}
