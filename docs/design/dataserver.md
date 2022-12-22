# Dataserver Design

## API

| API                        | Parameters | Response | Description |
|----------------------------|------------|----------|-------------|
| `/ping`                    | -          |          |             |
| `/v1/blob/createChunk`     | path, id   |          |             |
| `/v1/blob/createFile`      | path       |          |             |
| `/v1/blob/createDirectory` | path       |          |             |
| `/v1/blob/deleteChunk`     | path, id   |          |             |
| `/v1/blob/deleteFile`      | path       |          |             |
| `/v1/blob/deleteDirectory` | path       |          |             |
| `/v1/blob/lockFile`        | path       |          |             |
| `/v1/blob/readChunk`       | path, id   |          |             |
| `/v1/blob/readMeta`        | path, id   |          |             |
| `/v1/blob/unlockFile`      | path       |          |             |
| `/v1/blob/writeChunk`      | path, id   |          |             |
| `/v1/sys/info`             |            |          |             |
