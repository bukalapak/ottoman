package redis_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/redis"
	envx "github.com/bukalapak/ottoman/x/env"
	redisc "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

type Connector interface {
	Get(key string) *redisc.StringCmd
	MGet(keys ...string) *redisc.SliceCmd
	Incr(key string) *redisc.IntCmd
	Expire(key string, expiration time.Duration) *redisc.BoolCmd
	Del(keys ...string) *redisc.IntCmd
	Set(key string, value interface{}, expiration time.Duration) *redisc.StatusCmd
	TTL(key string) *redisc.DurationCmd
}

func TestRedis(t *testing.T) {
	t.Run("Standalone", func(t *testing.T) {
		client := NewRedisConnector()
		c := NewRedis()

		t.Run("Name", func(t *testing.T) { assert.Equal(t, "Redis", c.Name()) })
		t.Run("Write", func(t *testing.T) { testWrite(t, client, c) })
		t.Run("Read", func(t *testing.T) { testRead(t, client, c) })
		t.Run("Read-Unknown-Cache", func(t *testing.T) { testReadUnknown(t, c) })
		t.Run("ReadMulti", func(t *testing.T) { testReadMulti(t, client, c) })
		t.Run("Incr", func(t *testing.T) { testIncr(t, c) })
		t.Run("Expire", func(t *testing.T) { testExpire(t, client, c) })
		t.Run("Delete", func(t *testing.T) { testDelete(t, client, c) })
		t.Run("Delete-Unknown", func(t *testing.T) { testDeleteUnknown(t, c) })
	})

	t.Run("RedisCluster", func(t *testing.T) {
		client := NewRedisClusterConnector()
		c := NewRedisCluster()

		t.Run("Name", func(t *testing.T) { assert.Equal(t, "Redis Cluster", c.Name()) })
		t.Run("Write", func(t *testing.T) { testWrite(t, client, c) })
		t.Run("Read", func(t *testing.T) { testRead(t, client, c) })
		t.Run("Read-Unknown-Cache", func(t *testing.T) { testReadUnknown(t, c) })
		t.Run("ReadMulti", func(t *testing.T) { testReadMulti(t, client, c) })
		t.Run("Incr", func(t *testing.T) { testIncr(t, c) })
		t.Run("Expire", func(t *testing.T) { testExpire(t, client, c) })
		t.Run("Delete", func(t *testing.T) { testDelete(t, client, c) })
		t.Run("Delete-Unknown", func(t *testing.T) { testDeleteUnknown(t, c) })
		t.Run("ReadMulti-CROSSSLOT", func(t *testing.T) {
			loadFixtures(client)

			keys := []string{
				"foo",
				"{x}.fox",
			}

			m, err := c.ReadMulti(keys)
			assert.Contains(t, err.Error(), "CROSSSLOT")
			assert.Nil(t, m)

			cleanFixtures(client)
		})
	})

	t.Run("RedisSentinel", func(t *testing.T) {
		client := NewRedisSentinelConnector()
		c := NewRedisSentinel()

		t.Run("Name", func(t *testing.T) { assert.Equal(t, "Redis Sentinel", c.Name()) })
		t.Run("Write", func(t *testing.T) { testWrite(t, client, c) })
		t.Run("Read", func(t *testing.T) { testRead(t, client, c) })
		t.Run("Read-Unknown-Cache", func(t *testing.T) { testReadUnknown(t, c) })
		t.Run("ReadMulti", func(t *testing.T) { testReadMulti(t, client, c) })
		t.Run("Incr", func(t *testing.T) { testIncr(t, c) })
		t.Run("Expire", func(t *testing.T) { testExpire(t, client, c) })
		t.Run("Delete", func(t *testing.T) { testDelete(t, client, c) })
		t.Run("Delete-Unknown", func(t *testing.T) { testDeleteUnknown(t, c) })
	})

	os.Clearenv()
}

func testWrite(t *testing.T, client Connector, c *redis.Redis) {
	err := c.Write("foo", []byte("bar"), 10*time.Second)
	assert.Nil(t, err)

	cleanFixtures(client)
}

func testRead(t *testing.T, client Connector, c *redis.Redis) {
	loadFixtures(client)

	b, err := c.Read("foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), b)

	cleanFixtures(client)
}

func testReadUnknown(t *testing.T, c *redis.Redis) {
	b, err := c.Read("boo")
	assert.NotNil(t, err)
	assert.Nil(t, b)
}

func testReadMulti(t *testing.T, client Connector, c *redis.Redis) {
	loadFixtures(client)

	keys := []string{
		"{x}.foo",
		"{x}.boo",
		"{x}.fox",
	}

	z, err := c.ReadMulti(keys)
	assert.Nil(t, err)
	assert.Len(t, z, 2)

	for _, key := range keys {
		switch key {
		case "{x}.boo":
			assert.Nil(t, z[key])
		case "{x}.foo":
			assert.Equal(t, []byte(`{"foo":"bar"}`), z[key])
		case "{x}.fox":
			assert.Equal(t, []byte(`{"fox":"baz"}`), z[key])
		}
	}

	cleanFixtures(client)
}

func testIncr(t *testing.T, c *redis.Redis) {
	n, err := c.Incr("foo")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), n)
}

func testExpire(t *testing.T, client Connector, c *redis.Redis) {
	err := client.Set("foo", "bar", time.Minute).Err()
	assert.Nil(t, err)

	b, err := c.Expire("foo", time.Hour)
	assert.Nil(t, err)
	assert.True(t, b)

	cmd := client.TTL("foo")
	assert.Nil(t, cmd.Err())
	assert.Equal(t, time.Hour, cmd.Val())

	cleanFixtures(client)
}

func testDelete(t *testing.T, client Connector, c *redis.Redis) {
	loadFixtures(client)

	err := c.Delete("foo")
	assert.Nil(t, err)

	cleanFixtures(client)
}

func testDeleteUnknown(t *testing.T, c *redis.Redis) {
	err := c.Delete("boo")
	assert.NotNil(t, err)
}

func NewRedisConnector() Connector {
	return redisc.NewClient(&redisc.Options{
		Addr: redisAddr(),
		DB:   redisDb(),
	})
}

func NewRedis() *redis.Redis {
	opts := &redis.Option{
		Addrs: []string{redisAddr()},
		DB:    redisDb(),
	}

	return redis.New(opts)
}

func redisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}

	return "127.0.0.1:6379"
}

func redisDb() int {
	return envx.Int("REDIS_DB")
}

func redisClusterAddrs() []string {
	if addr := os.Getenv("REDIS_CLUSTER_ADDR"); addr != "" {
		return strings.Split(addr, ",")
	}

	return []string{
		"127.0.0.1:7000",
		"127.0.0.1:7001",
		"127.0.0.1:7002",
		"127.0.0.1:7003",
		"127.0.0.1:7004",
		"127.0.0.1:7005",
	}
}

func NewRedisClusterConnector() Connector {
	return redisc.NewClusterClient(&redisc.ClusterOptions{Addrs: redisClusterAddrs()})
}

func NewRedisCluster() *redis.Redis {
	return redis.New(&redis.Option{Addrs: redisClusterAddrs(), ReadOnly: true})
}

func redisSentinelMaster() string {
	if name := os.Getenv("REDIS_SENTINEL_MASTER"); name != "" {
		return name
	}

	return "redis-master"
}

func redisSentinelAddrs() []string {
	if addr := os.Getenv("REDIS_SENTINEL_ADDR"); addr != "" {
		return strings.Split(addr, ",")
	}

	return []string{
		"127.0.0.1:26379",
	}
}

func NewRedisSentinelConnector() Connector {
	return redisc.NewFailoverClient(&redisc.FailoverOptions{MasterName: redisSentinelMaster(), SentinelAddrs: redisSentinelAddrs()})
}

func NewRedisSentinel() *redis.Redis {
	return redis.New(&redis.Option{MasterName: redisSentinelMaster(), Addrs: redisSentinelAddrs()})
}

func loadFixtures(client Connector) {
	data := map[string]string{
		"foo":     `{"foo":"bar"}`,
		"fox":     `{"fox":"baz"}`,
		"{x}.foo": `{"foo":"bar"}`,
		"{x}.fox": `{"fox":"baz"}`,
		"baz":     `x`,
	}

	for k, v := range data {
		err := client.Set(k, []byte(v), time.Minute).Err()
		if err != nil {
			panic(err)
		}
	}
}

func cleanFixtures(client Connector) {
	client.Del("foo")
}
