# DataServer Design

## Basics

DataServers are managed by the nameserver.

DataServers does not use authentication.

## API

| API                        | Method | Parameters             | Response                 | Description                       |
|----------------------------|--------|------------------------|--------------------------|-----------------------------------|
| `/ping`                    | GET    | -                      |                          |                                   |
| `/v1/blob/createChunk`     | POST   | path, id               | code, msg                |                                   |
| `/v1/blob/createFile`      | POST   | path                   | code, msg                |                                   |
| `/v1/blob/createDirectory` | POST   | path                   | code, msg                |                                   |
| `/v1/blob/deleteChunk`     | POST   | path, id               | code, msg                |                                   |
| `/v1/blob/deleteFile`      | POST   | path                   | code, msg                |                                   |
| `/v1/blob/deleteDirectory` | POST   | path                   | code, msg                |                                   |
| `/v1/blob/lockFile`        | POST   | path, session          | code, msg                |                                   |
| `/v1/blob/readChunk`       | GET    | path, id, offset, size | binary                   |                                   |
| `/v1/blob/readFileLock`    | GET    | path                   | code, msg, id(session)[] |                                   |
| `/v1/blob/readFileMeta`    | GET    | path                   | code, msg, md5[]         |                                   |
| `/v1/blob/readChunkMeta`   | GET    | path, id               | code, msg, md5           |                                   |
| `/v1/blob/unlockFile`      | POST   | path                   | code, msg                |                                   |
| `/v1/blob/writeChunk`      | PUT    | form: path, id, file   | code, msg                |                                   |
| `/v1/sys/config`           | GET    | -                      | code, msg, config        |                                   |
| `/v1/sys/info`             | GET    | -                      | code, msg, info          |                                   |
| `/v1/sys/uuid`             | GET    | -                      | code, msg, uuid          |                                   |
| `/v1/sys/register`         | POST   | uuid                   | code, msg, uuid          | Nameserver Register to Dataserver |

## FileLocks


