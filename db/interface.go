package db

type ISet interface {
	Set(values ...interface{}) (int64, error)
}

type IRow interface {
	Scan(dest ...interface{}) error
	Struct(object interface{}) error

	Slice() ([]interface{}, error)
	Map() (map[string]interface{}, error)
}

type IRows interface {
	IRow
	Close() error
	Err() error
	Next() bool
}

type ITable interface {

	ToSql() string

	//Insert
	Add(values ...interface{}) (int64, error)
	//Delete
	Del(args ...interface{}) (int64, error)
	//Select
	Get(args ...interface{}) IRow
	//update
	Find(args ...interface{}) ISet

	List(take, skip int) (IRows, error)
	Query(query string, args ...interface{}) (IRows, error)

	//other
	// sql = select count(*) from
	Count(query string, args ...interface{}) (int64, error)
}