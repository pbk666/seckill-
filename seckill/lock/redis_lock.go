package lock

import (
	"fmt"
	"time"

	"decleration/dao"
	"github.com/go-redsync/redsync/v4"
	goredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

var redsyncClient *redsync.Redsync

func InitLock() {
	pool := goredis.NewPool(dao.GetRedisClient())
	redsyncClient = redsync.New(pool)
}

// 封装后的获取锁方法
func AcquireLock(key string, expiry time.Duration) (*redsync.Mutex, error) {
	mutex := redsyncClient.NewMutex(
		key,
		redsync.WithExpiry(expiry),
		redsync.WithTries(3),
	)
	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("获取锁失败: %v", err)
	}
	return mutex, nil
}

// 封装后的释放锁方法
func ReleaseLock(mutex *redsync.Mutex) {
	if ok, err := mutex.Unlock(); !ok || err != nil {
		fmt.Printf("释放锁失败: %v\n", err)
	}
}
