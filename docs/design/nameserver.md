# NameServer Design

## API List

## API

| API                 | Method | Parameters         | Response | Description |
|---------------------|--------|--------------------|----------|-------------|
| `/ping`             | GET    |                    |          |             |
| `/v1/auth/login`    | POST   | username, password |          |             |
| `/v1/auth/refresh`  | POST   | jwt                |          |             |
| `/v1/blob/open`     | POST   | path, mode         | session  |             |
| `/v1/blob/close`    | POST   | sessionID          | -        |             |
| `/v1/blob/flush`    | POST   | sessionID          |          |             |
| `/v1/blob/read`     | GET    | sessionID, size    |          |             |
| `/v1/blob/write`    | POST   | sessionID, data    | size     |             |
| `/v1/blob/truncate` | POST   | sessionID, size    |          |             |
| `/v1/blob/seek`     | POST   | sessionID, offset  |          |             |
| `/v1/blob/lock`     | POST   | sessionID          |          |             |
| `/v1/blob/unlock`   | POST   | sessionID          |          |             |
| `/v1/blob/mkdir`    | POST   | path               |          |             |
| `/v1/blob/ls`       | GET    | path               |          |             |
| `/v1/blob/rm`       | POST   | path               |          |             |
| `/v1/blob/rmdir`    | POST   | path               |          |             |
| `/v1/sys/info`      | GET    |                    |          |             |
| `/v1/sys/session`   | GET    | sessionID          |          |             |
| `/v1/sys/sessions`  | GET    |                    |          |             |

