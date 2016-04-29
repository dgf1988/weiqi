package db

func Connect(drivername, username, password, hostname string, port int, databasename string) error {
	conn, err := dbGetConnect(drivername, username, password, hostname, port, databasename)
	if err != nil {
		return err
	} else {
		db = conn
		return nil
	}
}

//值设置器
type Set interface {
	Values(values ...interface{}) (int64, error)
}

//单行读取
type Row interface {
	Scan(dest ...interface{}) error
	Struct(object interface{}) error

	Slice() ([]interface{}, error)
	Map() (map[string]interface{}, error)
}

//多行读取
type Rows interface {
	Row
	Close() error
	Err() error
	Next() bool
	Columns() ([]string, error)
}

//表操作接口
type Table interface {

	//输出表结构
	ToSql() string

	//添加记录
	Add(values ...interface{}) (int64, error)

	//删除记录
	Del(args ...interface{}) (int64, error)

	//更新记录
	Update(args ...interface{}) Set

	//获取单条记录
	Get(args ...interface{}) Row

	//查询多条记录
	List(take, skip int) (Rows, error)
	ListDesc(take, skip int) (Rows, error)
	FindAll(args ...interface{}) (Rows, error)
	FindAny(args ...interface{}) (Rows, error)
	Query(query string, args ...interface{}) (Rows, error)

	//统计
	Count(args ...interface{}) (int64, error)
	CountBy(query string, args ...interface{}) (int64, error)
}
