package migrator

// Migrator — общий контракт применения миграций к базе данных
type Migrator interface {
	Up() error
}
