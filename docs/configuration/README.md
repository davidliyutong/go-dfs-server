# Configuration

## Configurable parameters of servers

| config file         | environment           | command line     | default and example values                                                                                                                   |
|---------------------|-----------------------|------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `config`            | `DFSAPP_CONFIG`       | `--config`       | Path to config, **no default value**                                                                                                         |
| `uuid`              | `DFSAPP_UUID`         |                  | UUID of server, important for data server                                                                                                    |
| `network.port`      | `DFSAPP_PORT`         | `--port,-p`      | Listen port, default to `27903` for name sever and `27904` for dataserver                                                                    |
| `network.interface` | `DFSAPP_INTERFACE`    | `--interface,-i` | \[Name Server Only\] Listen interface, default to `0.0.0.0`                                                                                  |
| `network.endpoint`  | `DFSAPP_ENDPOINT`     | `--endpoint`     | \[Data Server Only\] Endpoint url of this dataserver, nameserver will communicate with dataserver using this url, otherwise automatic detect |
| `volume`            | `DFSAPP_VOLUME`       | `--volume`       | Path to persistence storage, default to `/data`,  **need runtime validation**                                                                |
| `auth.domain`       | `DFSAPP_DOMAIN`       | `--domain`       | \[Name Server Only\] Domain name of DFS cluster, default to `dfs.local`                                                                      |
| `auth.accessKey`    | `DFSAPP_ACCESSKEY`    | `--accessKey`    | \[Name Server Only\] master access key, leave blank to disable authentication                                                                |
| `auth.secretKey`    | `DFSAPP_SECRETKEY`    | `--secretKey`    | \[Name Server Only\] master secret key, leave blank to disable authentication                                                                |
| `debug`             | `DFSAPP_DEBUG`        | `--debug`        | Toggle debug output, will override `log.level`                                                                                               |
| `log.level`         | `DFSAPP_LOG_LEVEL`    |                  | Select log level, default to "info"                                                                                                          |
| `log.path`          | `DFSAPP_LOG_PATH`     |                  | Select log output path                                                                                                                       |
|                     | `DFSAPP_DATA_SERVERS` |                  | Data Server address:port, e.g. `192.168.1.100:27904,192.168.1.100:27905,192.168.1.100:27906,192.168.1.100:27907`                             |


## Configurable parameters of client

| config file | environment      | command line      | default and example values                    |
|-------------|------------------|-------------------|-----------------------------------------------|
| `-`         | `DFSAPP_CONFIG`  | `--config`        | Path to client config, **no default value**   |
| `-`         | `DFSAPP_VERBOSE` | `--verbose`, `-v` | Toggle verbose output                         |

