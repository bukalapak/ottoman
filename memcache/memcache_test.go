package memcache_test

import (
	"bytes"
	"compress/zlib"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	gomemcache "github.com/bradfitz/gomemcache/memcache"
	"github.com/bukalapak/ottoman/memcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

type netAddr struct {
	network string
	str     string
}

func (a *netAddr) Network() string {
	return a.network
}

func (a *netAddr) String() string {
	return a.str
}

type MockMemcacheClient struct {
	mock.Mock
}

func (mc *MockMemcacheClient) Set(i *gomemcache.Item) error {
	args := mc.Called(i)
	err, _ := args.Get(0).(error)

	return err
}

func (mc *MockMemcacheClient) Get(key string) (*gomemcache.Item, error) {
	args := mc.Called(key)

	item, _ := args.Get(0).(*gomemcache.Item)
	err, _ := args.Get(1).(error)

	return item, err
}

func (mc *MockMemcacheClient) GetMulti(keys []string) (map[string]*gomemcache.Item, error) {
	args := mc.Called(keys)
	items, _ := args.Get(0).(map[string]*gomemcache.Item)
	err, _ := args.Get(1).(error)

	return items, err
}

func (mc *MockMemcacheClient) Delete(key string) error {
	args := mc.Called(key)
	err, _ := args.Get(0).(error)

	return err
}

func TestMemcache(t *testing.T) {
	addr := memcachedAddr()
	client := gomemcache.New(addr)

	t.Run("Name", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{})
		assert.Equal(t, "Memcached", c.Name())
	})

	t.Run("MaxIdleConns", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{})
		assert.Equal(t, 2, c.MaxIdleConns())

		c2 := memcache.New([]string{addr}, memcache.Option{
			MaxIdleConns: 20,
		})
		assert.Equal(t, 20, c2.MaxIdleConns())
	})

	t.Run("Write", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{
			Compress: false,
		})

		err := c.Write("foo", []byte("bar"), 10*time.Second)
		assert.Nil(t, err)

		item, err := client.Get("foo")
		assert.Nil(t, err)
		assert.Equal(t, []byte("bar"), item.Value)

		cleanFixtures(client)
	})

	t.Run("Write-Compress", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{
			Compress: true,
		})

		err := c.Write("foo", []byte("bar"), 10*time.Second)
		assert.Nil(t, err)

		item, err := client.Get("foo")
		assert.Nil(t, err)
		assert.Equal(t, []byte("bar"), item.Value)

		cleanFixtures(client)
	})

	t.Run("Write-Compress-Zlib", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{
			Compress: true,
		})

		err := c.Write("foo", []byte(strings.Repeat("bar", 1000)), 10*time.Second)
		assert.Nil(t, err)

		item, err := client.Get("foo")
		assert.Nil(t, err)

		var b bytes.Buffer

		w := zlib.NewWriter(&b)
		w.Write([]byte(strings.Repeat("bar", 1000)))
		w.Close()

		assert.Equal(t, b.Bytes(), item.Value)

		cleanFixtures(client)
	})

	t.Run("Read", func(t *testing.T) {
		loadUncompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{
			Compress: false,
		})

		b, err := c.Read("foo")
		assert.Nil(t, err)
		assert.Equal(t, `{"foo":"bar"}`, string(b))

		cleanFixtures(client)
	})

	t.Run("Read-Miss", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{})

		b, err := c.Read("boo")
		assert.Error(t, gomemcache.ErrCacheMiss, err.Error())
		assert.Nil(t, b)
	})

	t.Run("Read-Zlib", func(t *testing.T) {
		loadCompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{
			Compress: true,
		})

		b, err := c.Read("foo")
		assert.Nil(t, err)
		assert.Equal(t, `{"foo":"bar"}`, string(b))

		cleanFixtures(client)
	})

	t.Run("Read-Zlib-Uncompressed-Cache", func(t *testing.T) {
		loadUncompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{
			Compress: true,
		})

		b, err := c.Read("foo")
		assert.Nil(t, err)
		assert.Equal(t, `{"foo":"bar"}`, string(b))

		cleanFixtures(client)
	})

	t.Run("ReadMulti", func(t *testing.T) {
		loadCompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{
			Compress: true,
		})

		keys := []string{
			"foo",
			"boo",
			"fox",
		}

		z, err := c.ReadMulti(keys)
		assert.Nil(t, err)
		assert.Len(t, z, 2)

		for _, key := range keys {
			switch key {
			case "boo":
				assert.Nil(t, z[key])
			case "foo":
				assert.Equal(t, []byte(`{"foo":"bar"}`), z[key])
			case "fox":
				assert.Equal(t, []byte(`{"fox":"baz"}`), z[key])
			}
		}

		cleanFixtures(client)
	})

	t.Run("ReadMulti-Failure", func(t *testing.T) {
		loadCompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{
			Timeout: 1 * time.Microsecond,
		})

		m, err := c.ReadMulti([]string{"foo", "boo", "fox"})
		assert.NotNil(t, err)
		assert.Len(t, m, 0)

		cleanFixtures(client)
	})

	t.Run("Delete", func(t *testing.T) {
		loadUncompressedFixtures(client)

		c := memcache.New([]string{addr}, memcache.Option{})

		err := c.Delete("foo")
		assert.Nil(t, err)

		cleanFixtures(client)
	})

	t.Run("Delete-Miss", func(t *testing.T) {
		c := memcache.New([]string{addr}, memcache.Option{})

		err := c.Delete("boo")
		assert.Error(t, gomemcache.ErrCacheMiss, err.Error())
	})

	t.Run("Retry-On-Timeout", func(t *testing.T) {
		mc := &MockMemcacheClient{}

		c := memcache.NewWithClient(mc, memcache.Option{
			Timeout:    100 * time.Millisecond,
			Compress:   false,
			MaxAttempt: 3,
		})

		errTimeout := &net.OpError{
			Op:     "read",
			Net:    "tcp",
			Source: &netAddr{network: "tcp", str: "10.90.142.64:43080"},
			Addr:   &netAddr{network: "tcp", str: "10.49.42.100:11212"},
			Err:    &timeoutError{},
		}

		mc.On("Set", &gomemcache.Item{
			Key:        "foo",
			Value:      []byte("bar"),
			Expiration: 10,
		}).Return(errTimeout)
		mc.On("Get", "foo").Return(nil, errTimeout)
		mc.On("GetMulti", []string{"foo", "bar"}).Return(nil, errTimeout)
		mc.On("Delete", "foo").Return(errTimeout)

		err := c.Write("foo", []byte("bar"), 10*time.Second)
		mc.AssertNumberOfCalls(t, "Set", 3)
		assert.Equal(t, errTimeout, err)

		_, err = c.Read("foo")
		mc.AssertNumberOfCalls(t, "Get", 3)
		assert.Equal(t, errTimeout, err)

		_, err = c.ReadMulti([]string{"foo", "bar"})
		mc.AssertNumberOfCalls(t, "GetMulti", 3)
		assert.Equal(t, errTimeout, err)

		err = c.Delete("foo")
		mc.AssertNumberOfCalls(t, "Delete", 3)
		assert.Equal(t, errTimeout, err)
	})
}

func loadCompressedFixtures(client *gomemcache.Client) {
	loadFixtures(client, true)
}

func loadUncompressedFixtures(client *gomemcache.Client) {
	loadFixtures(client, false)
}

func loadFixtures(client *gomemcache.Client, compress bool) {
	data := map[string]string{
		"foo": `{"foo":"bar"}`,
		"fox": `{"fox":"baz"}`,
		"baz": `x`,
	}

	for k, v := range data {
		var z []byte

		if compress {
			var b bytes.Buffer

			w := zlib.NewWriter(&b)
			io.WriteString(w, v)
			w.Close()

			z = b.Bytes()
		} else {
			z = []byte(v)
		}

		x := &gomemcache.Item{
			Key:        k,
			Value:      z,
			Expiration: int32(time.Minute.Seconds()),
		}

		err := client.Set(x)
		if err != nil {
			panic(err)
		}
	}
}

func cleanFixtures(client *gomemcache.Client) {
	client.Delete("foo")
	client.Delete("baz")
}

func memcachedAddr() string {
	if addr := os.Getenv("MEMCACHE_ADDR"); addr != "" {
		return addr
	}

	return "127.0.0.1:11211"
}
