package db

func Close() {
	Db.Close()
	Rdb.Close()
}
