package db

type IRow interface {
	Scan(dest ...interface{}) error
}

type IRows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

type ITable interface {

	Add(args ...interface{}) (int64, error)
	Del(key interface{}) (int64, error)

	Set(key interface{}, args ...interface{}) error
	Get(key interface{}, dest ...interface{}) error

	Query(query string, args ...interface{}) (IRows, error)
	Update(datas map[string]interface{}, query string, args ...interface{}) (int64, error)

}
