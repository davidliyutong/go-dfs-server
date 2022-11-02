# Configuration

## common basics

- `DFSAPP_CONFIG`

## nameserver basics

| config file | environment         | command line     | default and example values                                           |
|-------------|---------------------|------------------|----------------------------------------------------------------------|
| `config`    | `DFSAPP_CONFIG`     | `--config`       | Path to config, **no default value**                                 |
| `port`      | `DFSAPP_PORT`       | `--port,-p`      | Listen port, default to `27903`                                      |
| `interface` | `DFSAPP_INTERFACE`  | `--interface,-i` | Listen interface, default to `0.0.0.0`                               |
| `volume`    | `DFSAPP_VOLUME`     | `--volume`       | Path to persistence storage, default to `/data`. The path must exist |
| `accessKey` | `DFSAPP_ACCESSKEY`  | `--accessKey`    | master access key                                                    |
| `secretKey` | `DFSAPP_SECRETKEY`  | `--secretKey`    | master secret key                                                    |


## dataserver basics

| config file     | environment             | command line      | default and example values                                                                                         |
|-----------------|-------------------------|-------------------|--------------------------------------------------------------------------------------------------------------------|
| `config`        | `DFSAPP_CONFIG`         | `--config`        | Path to config, **no default value**                                                                               |
| `port`          | `DFSAPP_PORT`           | `--port`          | Listen port, default to `27903`                                                                                    |
| `nameServerURL` | `DFSAPP_NAMESERVER_URL` | `--nameserverUrl` | Url to reach name server, for data server only, must specify value                                                 |
| `overwriteURL`  | `DFSAPP_OVERWRITE_URL`  | `--overwriteUrl`  | Reported url to nameserver, nameserver will communicate with dataserver using this url, otherwise automatic detect |
| `volume`        | `DFSAPP_VOLUME`         | `--volume`        | Path to persistence storage, default to `/data`. The path must exist                                               |
| `accessKey`     | `DFSAPP_ACCESSKEY`      | `--accessKey`     | master access key                                                                                                  |
| `secretKey`     | `DFSAPP_SECRETKEY`      | `--secretKey`     | master secret key                                                                                                  |







