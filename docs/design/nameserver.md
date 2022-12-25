# NameServer Design

## API List

## API

| API                     | Parameters         | Response | Description |
|-------------------------|--------------------|----------|-------------|
| `/ping`                 |                    |          |             |
| `/v1/auth/login`        | username, password |          |             |
| `/v1/auth/refresh`      | jwt                |          |             |
| `/v1/blob/open`         | path               |          |             |
| `/v1/blob/close`        | handle             |          |             |
| `/v1/blob/read`         | handle             |          |             |
| `/v1/blob/write`        | handle, data       |          |             |
| `/v1/blob/seek`         | handle, offset     |          |             |
| `/v1/blob/lock`         | path               |          |             |
| `/v1/blob/unlock`       | path               |          |             |
| `/v1/blob/mkdir`        | path               |          |             |
| `/v1/blob/ls`           | path               |          |             |
| `/v1/blob/rm`           | path               |          |             |
| `/v1/blob/rmdir`        | path               |          |             |
| `/v1/dataserver/add`    | ...                |          |             |
| `/v1/dataserver/remove` | ...                |          |             |
| `/v1/sys/info`          |                    |          |             |