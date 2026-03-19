package inits

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gospacex-tz/srv/basic/config"
)

func RedisInit() {

	data := config.GlobalConfig.Redis
	Addr := fmt.Sprintf("%s:%d", data.Host, data.Port)
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: data.Password, // no password set
		DB:       data.Database, // use default DB
	})

	err := config.Rdb.Ping(config.Ctx).Err()
	if err != nil {
		panic("redis连接失败: " + err.Error())
	}
	fmt.Println("redis连接成功")
}
