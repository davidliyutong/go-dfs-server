# NameServer Design

## API List

## API

| API                | Method | Parameters         | Response | Description |
|--------------------|--------|--------------------|----------|-------------|
| `/ping`            | GET    |                    |          |             |
| `/v1/auth/login`   | POST   | username, password |          |             |
| `/v1/auth/refresh` | POST   | jwt                |          |             |
| `/v1/blob/open`    | POST   | path               |          |             |
| `/v1/blob/close`   | POST   | handle             |          |             |
| `/v1/blob/flush`   | POST   | handle             |          |             |
| `/v1/blob/read`    | GET    | handle, size       |          |             |
| `/v1/blob/write`   | POST   | handle, data       |          |             |
| `/v1/blob/seek`    | POST   | handle, offset     |          |             |
| `/v1/blob/lock`    | POST   | path               |          |             |
| `/v1/blob/unlock`  | POST   | path               |          |             |
| `/v1/blob/mkdir`   | POST   | path               |          |             |
| `/v1/blob/ls`      | GET    | path               |          |             |
| `/v1/blob/rm`      | POST   | path               |          |             |
| `/v1/blob/rmdir`   | POST   | path               |          |             |
| `/v1/sys/info`     | GET    |                    |          |             |