package db

/*

	1、数据类型 - Scan
		驱动里返回的源数据类型均是[]byte
		数据经过sql标准库时，跟据目标数据类型和源数据类型进行断言赋值或反射赋值或字符串和[]byte的拷贝。
			如果目标类型是interface{}, sql库直接拷贝驱动库出来的数据，所以出现了大量的[]byte串
			如果目标类型不是基本类型，则sql库断言你为scanner，并调用你的Scan(v interface{}) error 方法。
			其它类型报错。

	2、把数据做为参数传给查询函数
		如果不是基本数据类型，则断言你为valuer，并调用你的Value() driver.value, error 方法。
			如果断言失败，报错。

*/
