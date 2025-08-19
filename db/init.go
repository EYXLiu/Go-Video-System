package db

func Init() {
	S3Init()
	PostgresInit()
	RedisInit()
}
