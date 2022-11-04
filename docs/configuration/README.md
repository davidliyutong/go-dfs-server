# Configuration

## server basics

| config file         | environment        | command line     | default and example values                                                                                                              |
|---------------------|--------------------|------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `config`            | `DFSAPP_CONFIG`    | `--config`       | Path to config, **no default value**                                                                                                    |
| `network.port`      | `DFSAPP_PORT`      | `--port,-p`      | Listen port, default to `27903` for name sever and `27904` for data server                                                              |
| `network.interface` | `DFSAPP_INTERFACE` | `--interface,-i` | \[Name Server\] Listen interface, default to `0.0.0.0`                                                                                  |
| `network.remote`    | `DFSAPP_REMOTE`    | `--remote`       | \[Data Server\] Url to reach nameserver, e.g. dfs://192.168.1.1:27903,  **need runtime validation**                                     |
| `network.endpoint`  | `DFSAPP_ENDPOINT`  | `--endpoint`     | \[Data Server\] Endpoint url of this dataserver, nameserver will communicate with dataserver using this url, otherwise automatic detect |
| `volume`            | `DFSAPP_VOLUME`    | `--volume`       | Path to persistence storage, default to `/data`,  **need runtime validation**                                                           |
| `auth.domain`       | `DFSAPP_DOMAIN`    | `--domain`       | Domain name of DFS cluster, default to `dfs.local`                                                                                      |
| `auth.accessKey`    | `DFSAPP_ACCESSKEY` | `--accessKey`    | master access key, leave blank to disable authentication                                                                                |
| `auth.secretKey`    | `DFSAPP_SECRETKEY` | `--secretKey`    | master secret key, leave blank to disable authentication                                                                                |
| `debug`             | `DFSAPP_DEBUG`     | `--debug`        | Toggle debug output, will override `log.level`                                                                                          |
| `log.level`         | `DFSAPP_LOG_LEVEL` |                  | Select log level, default to "info"                                                                                                     |
| `log.path`          | `DFSAPP_LOG_PATH`  |                  | Select log output path                                                                                                                  |



##  client commands
                                                           
client use `local:` and `remote:` prefix to decide if path is remote path (on the DFS cluster) or local path. A path is assumed to be remote path if no prefix is provided.

| Name    | Argument(SDK)           | Argument (CLI)  | Options                                                | Functional                                             |
|---------|-------------------------|-----------------|--------------------------------------------------------|--------------------------------------------------------|
| `ls`    | path                    | \[path\]        |                                                        |                                                        |
| `cat`   | path                    | \[path\]        |                                                        |                                                        |
| `mkdir` | path                    | \[path\]        |                                                        |                                                        |
| `rmdir` | path                    | \[path\]        |                                                        |                                                        |
| `touch` | path                    | \[path\]        |                                                        |                                                        |
| `rm`    | path recursive          | \[path\]        | `-r`                                                   |                                                        |
| `cp`    | src dst                 | \[src\] \[dst\] |                                                        |                                                        |
| `mv`    | src dst                 | \[src\] \[dst\] |                                                        |                                                        |
| `read`  | path offset size output |                 | `--offset=0x0` `--size=1` `--output=/path/to/output`   |                                                        |
| `write` | path offset size output |                 | `--offset=0x0`` --size=1` `--output=/path/to/output`   |                                                        |
| `open`  | -                       | \[path\]        |                                                        | open a file with an interactive prompt **interactive** |
                                   
### Interactive Prompt

| Name      | Argument (SDK)  | Argument (CLI)                                   | Functional        |
|-----------|-----------------|--------------------------------------------------|-------------------|
| `fseek`   | offset (handle) | \[offset\]                                       | move file pointer |
| `fread`   | size (handle)   | \[size\]                                         |                   |
| `fwrite`  | data (handle)   | `--data=/path/to/file` or `--data=/path/to/file` |                   | 
| `fputc`   | c (handle)      | c                                                |                   |   
| `fputs`   | str (handle)    | str                                              |                   |   
| `fgetc`   | (handle)        |                                                  |                   |   
| `fgets`   | size (handle)   |                                                  |                   |     
| `fflush`  | (handle)        |                                                  |                   | 
| `fclose`  | (handle)        |                                                  |                   |  










