# Changelog

## [1.10.0] - 2020-05-28

### Added

- add option to configure `maxIdleConns` on Memcached client

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
