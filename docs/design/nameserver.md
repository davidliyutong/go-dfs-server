# NameServer Design

## API List

## API

| API                | Method | Parameters         | Response | Description |
|--------------------|--------|--------------------|----------|-------------|
| `/ping`            | GET    |                    |          |             |
| `/v1/auth/login`   | POST   | username, password |          |             |
| `/v1/auth/refresh` | POST   | jwt                |          |             |
| `/v1/blob/session` | GET    | path, mode         | session  | open        |
| `/v1/blob/session` | DELETE | sessionID          | -        | close       |
| `/v1/blob/session` | POST   | sessionID          |          | flush       |
| `/v1/blob/io`      | GET    | sessionID, size    |          | read        |
| `/v1/blob/io`      | POST   | sessionID, data    | size     | write       |
| `/v1/blob/io`      | DELETE | sessionID, size    |          | truncate    |
| `/v1/blob/seek`    | POST   | sessionID, offset  |          |             |
| `/v1/blob/lock`    | GET    | path               |          | getLock     |
| `/v1/blob/lock`    | POST   | sessionID          |          | lock        |
| `/v1/blob/lock`    | DELETE | sessionID          |          | unlock      |
| `/v1/blob/path`    | POST   | path               |          | mkdir       |
| `/v1/blob/path`    | GET    | path               |          | ls          |
| `/v1/blob/path`    | DELETE | path, recursive    |          | rm          |
| `/v1/blob/meta`    | GET    | path               |          | getFileMeta |
| `/v1/sys/info`     | GET    |                    |          |             |
| `/v1/sys/session`  | GET    | sessionID          |          |             |
| `/v1/sys/sessions` | GET    |                    |          |             |
| `/v1/sys/servers`  | GET    |                    |          |             |


