# NameServer Design

## API List

## API

| API                | Method | Parameters                       | Response    | Description |
|--------------------|--------|----------------------------------|-------------|-------------|
| `/ping`            | GET    |                                  |             |             |
| `/v1/auth/login`   | POST   | username, password               |             |             |
| `/v1/auth/refresh` | POST   | jwt                              |             |             |
| `/v1/blob/file`    | GET    | path, mode                       | blob, error | open        |
| `/v1/blob/file`    | POST   | path, blob                       | blob, error | sync        |
| `/v1/blob/io`      | GET    | path, chunkID, chunkOffset, size | error       | read        |
| `/v1/blob/io`      | POST   | path, chunkID, chunkOffset, data | size, error | write       |
| `/v1/blob/path`    | POST   | path                             |             | mkdir       |
| `/v1/blob/path`    | GET    | path                             |             | ls          |
| `/v1/blob/path`    | DELETE | path, recursive                  |             | rm          |
| `/v1/sys/info`     | GET    |                                  |             |             |
| `/v1/sys/session`  | GET    | sessionID                        |             |             |
| `/v1/sys/sessions` | GET    |                                  |             |             |
| `/v1/sys/servers`  | GET    |                                  |             |             |


