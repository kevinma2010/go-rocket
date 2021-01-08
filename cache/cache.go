package cache

import "time"

type Cache interface {
	Get(key string) (string, error)
	GetObj(key string, model interface{}) (bool, error)
	Set(key, value string, expiration time.Duration) error
	SetObj(key string, obj interface{}, expiration time.Duration) error
	Del(key string) error
}
