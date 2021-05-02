package lego

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mjiulee/lego/utils"
	"sync"
	"time"
)

const ()

var _pool *redis.Pool
var _once sync.Once

type RedisHelper struct {
}

/* 获取Redis连接
* params:
  ---
*/
func GetRedisConn() (redis.Conn, error) {
	if _pool == nil {
		_once.Do(func() {
			host := GetIniByKey("REDIS", "REDIS_HOST")
			port := GetIniByKey("REDIS", "REDIS_PORT")
			password := GetIniByKey("REDIS", "REDIS_PSWD")
			_pool = &redis.Pool{
				MaxIdle:     10,
				IdleTimeout: 240 * time.Second,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
					if err != nil {
						return nil, err
					}
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
					return c, err
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
		})
	}
	return _pool.Get(), nil
}

/*
 * 再多封装一次，查key是否存在的用法.
 */
func RedisKeyExists(key string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("EXISTS", key))
}

/*
 * 再多封装一次，get、set的用法.
 */
func RedisGetKey(key string) (string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return "", err
	}
	defer redisconn.Close()
	return redis.String(redisconn.Do("GET", key))
}

/*
 * 再多封装一次，get、set的用法.
 */
func RedisSetKey(key, value, expiretime string, ifexpire bool) (string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return "", err
	}
	defer redisconn.Close()

	if ifexpire {
		expseconds, _ := utils.StringToInt(expiretime)
		return redis.String(redisconn.Do("SET", key, value, "EX", expseconds))
	} else {
		return redis.String(redisconn.Do("SET", key, value))
	}
}

/*
 * 一次性设置多个key、val.
 */
func RedisSetMultiKey(key, value []string) (string, error) {
	if len(key) != len(value) {
		return "", errors.New("key和val长度不一致")
	}

	redisconn, err := GetRedisConn()
	if err != nil {
		return "", err
	}
	defer redisconn.Close()

	params := make([]string, len(key)*2)
	for i := 0; i < len(key); i++ {
		params = append(params, key[i], value[i])
	}
	return redis.String(redisconn.Do("MSET", redis.Args{}.AddFlat(params)...))
}

/*
 * 再多封装一次，DEL的用法.
 */
func RedisDeleteKey(key string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("DEL", key))
}

/*
 * SET集合设置成员
 */
func RedisSetAdd(setkey string, members ...string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("SADD", redis.Args{}.Add(setkey).AddFlat(members)...))
}

/*
 * SET集合删除成员
 */
func RedisSetRemove(setkey string, members ...string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("SREM", redis.Args{}.Add(setkey).AddFlat(members)...))
}

/*
 * SET集合元素个数
 */
func RedisSetMemberCount(setkey string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("SCARD", setkey))
}

/*
 * SET集合元素读取（所有元素）
 */
func RedisSetScan(setkey string) ([]string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return []string{}, err
	}
	defer redisconn.Close()
	var (
		cursor int64
		items  []string
	)

	results := make([]string, 0)
	for {
		//values, err := redis.Values(redisconn.Do("SSCAN", setkey, cursor))
		// 一次读5000条
		values, err := redis.Values(redisconn.Do("SSCAN", setkey, cursor, "COUNT", 5000))
		if err != nil {
			return []string{}, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			return []string{}, err
		}

		results = append(results, items...)

		if cursor == 0 {
			break
		}
	}
	return results, nil
}

/*
 * SET集合-- 判断元素是否在集合中
 */
func RedisSetExist(setkey string, member string) (bool, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return false, err
	}
	defer redisconn.Close()
	return redis.Bool(redisconn.Do("SISMEMBER", setkey, member))
}

/*
 * SET集合-- 2个集合的交集
 */
func RedisSetSinter(setkey1, setkey2 string) ([]string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return []string{}, err
	}
	defer redisconn.Close()
	return redis.Strings(redisconn.Do("SINTER", setkey1, setkey2))
}

/*
 * SET集合-- 2个集合的差集
 */
func RedisSetSDIFF(setkey1, setkey2 string) ([]string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return []string{}, err
	}
	defer redisconn.Close()
	return redis.Strings(redisconn.Do("SDIFF", setkey1, setkey2))
}

/*
 * 发布及订阅者模式
 */
func RedisPublish(chenel, message string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("PUBLISH", chenel, message))
}

/*
 * 分布式锁
 */
func RedisSETNX(key, val string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("SETNX", key, val))
}

// GEO 相关
func RedisGeoAdd(key string, lat, lng float64, id string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("GEOADD", key, lng, lat, id))
}

func RedisGeoPos(key, id string) (lat, lng float64, err error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, 0, err
	}
	defer redisconn.Close()
	rt, err := redis.Positions(redisconn.Do("GEOPOS", key, id))
	//TODO: 这地方有bug，
	fmt.Println(err )
	fmt.Println(rt )
	if err == nil && rt != nil && len(rt) > 0 && rt[0] != nil {
		return rt[0][1], rt[0][0], nil
	} else {
		return 0, 0, errors.New("RedisGeoPos error ")
	}
}

func RedisGeoDist(key, id1, id2 string) (string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return "", err
	}
	defer redisconn.Close()
	return redis.String(redisconn.Do("GEODIST", key, id1,id2))
}

// 半径内元素个数
func RedisGeoRadius(key string, lat, lng float64, radius int) ([]string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return []string{}, err
	}
	defer redisconn.Close()
	return redis.Strings(redisconn.Do("GEORADIUS", key, lng, lat, radius, "km", "ASC", "COUNT", 20))
}

/*
 * SET集合删除成员
 */
func RedisGeoRemove(setkey string, members ...string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("ZREM", redis.Args{}.Add(setkey).AddFlat(members)...))
}

// 列表相关
func RedisListLPush(key string, val interface{}) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("LPUSH", key, val))
}

func RedisListRPush(key string, val interface{}) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("RPUSH", key, val))
}

func RedisListLen(key string) (int, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}
	defer redisconn.Close()
	return redis.Int(redisconn.Do("LLEN", key))
}

func RedisListRange(key string, begin, end int) ([]string, error) {
	redisconn, err := GetRedisConn()
	if err != nil {
		return []string{}, err
	}
	defer redisconn.Close()
	return redis.Strings(redisconn.Do("LRANGE", key, begin, end))
}
