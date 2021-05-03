# Changelog

## [1.16.0] - 2021-05-03

### Added

- Add new methods `String`, `Int64`, and `Float64` for type `Number`

## [1.15.0] - 2021-04-30

### Added

- Add new type `Number` as alias from `json.Number`
- Add custom marshaler and unmarshaler for type `Number`

## [1.14.0] - 2021-02-26

### Updated

- Use go 1.16

## [1.13.0] - 2021-02-11

### Added

- Add `MemcacheClient` interface
- Add configurable `MaxAttempt` on memcached
- Add retry when timeout on memcached

## [1.12.0] - 2020-09-10

### Updated

- Update json.Boolean field member as public

## [1.11.0] - 2020-08-24

### Updated

- Bump go-redis client to v7 to support redis6 cluster

## [1.10.0] - 2020-05-28

### Added

- Add option to configure `maxIdleConns` on Memcached client

## [1.9.0] - 2019-07-26

### Added

- `Tracker` interface
- Datadog tracker implementation

## [1.8.0] - 2019-06-24

### Added

- Support on `json.Timestamp` for time format "2006-01-02 15:04:05 -0700"

## [1.7.0] - 2019-05-31

### Updated

- Signature of `jose.Decode` and `jose.Decrypt`, both using *rsa.PublicKey and *rsa.PrivateKey directly.

## [1.6.0] - 2019-05-28

### Added

- Helper `jose.Decode` and `jose.Decrypt`
- New interface `jose.Encryption` and `jose.Signature` as the result of `jose.Stamper` expansion.

## [1.5.0] - 2019-05-15

### Added

- Add ErrorHandler to Proxy

## [1.4.0] - 2019-04-24

### Added

- Redis: sentinel support

## [1.3.0] - 2019-04-15

### Changed

- Updated RecoveryLogger signature

## [1.2.1] - 2018-12-10

### Added

- `MaxRetries` and `IdleTimeout` options for redis package.

## [1.2.0] - 2018-11-12

### Added

- Delete method for cache provider, memcache and redis.

## [1.1.0] - 2018-11-06

### Added

- Memcache: automatically compress large cache value.

## [1.0.0] - 2018-10-09

Initial release
