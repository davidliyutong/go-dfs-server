# Client


##  Regulations

client use `local:` and `remote:` prefix to decide if path is remote path (on the DFS cluster) or local path. A path is assumed to be remote path if no prefix is provided.

##  Commands

Client support multiple genres of actions:

- Authentication
- Maneuver

### Authentication

| Name     | Argument(SDK)           | Argument (CLI)  | Options                                              | Functional                                             |
|----------|-------------------------|-----------------|------------------------------------------------------|--------------------------------------------------------|
| `login`  | -                       | -               | `--accessKey=12345678` `--secretKey=xxxxxxxx`        |                                                        |
| `logout` | -                       | -               |                                                      |                                                        |

### Maneuver

| Name     | Argument(SDK)           | Argument (CLI)  | Options                                              | Functional                                             |
|----------|-------------------------|-----------------|------------------------------------------------------|--------------------------------------------------------|
| `info`   | -                       | -               |                                                      |                                                        |
| `ls`     | path                    | \[path\]        |                                                      |                                                        |
| `cat`    | path                    | \[path\]        |                                                      |                                                        |
| `mkdir`  | path                    | \[path\]        |                                                      |                                                        |
| `rmdir`  | path                    | \[path\]        |                                                      |                                                        |
| `touch`  | path                    | \[path\]        |                                                      |                                                        |
| `rm`     | path recursive          | \[path\]        | `-r`                                                 |                                                        |
| `cp`     | src dst                 | \[src\] \[dst\] |                                                      |                                                        |
| `mv`     | src dst                 | \[src\] \[dst\] |                                                      |                                                        |
| `read`   | path offset size output |                 | `--offset=0x0` `--size=1` `--output=/path/to/output` |                                                        |
| `write`  | path offset size output |                 | `--offset=0x0`` --size=1` `--output=/path/to/output` |                                                        |
| `open`   | -                       | \[path\]        |                                                      | open a file with an interactive prompt **interactive** |

## Interactive Prompt

The client can open an interactive prompt for demo propose

| Name      | Argument (SDK)  | Argument (CLI)                                   | Functional        |
|-----------|-----------------|--------------------------------------------------|-------------------|
| `fopen`   | path            |                                                  | open file         |
| `fseek`   | offset (handle) | \[offset\]                                       | move file pointer |
| `fread`   | size (handle)   | \[size\]                                         |                   |
| `fwrite`  | data (handle)   | `--data=/path/to/file` or `--data=/path/to/file` |                   | 
| `fputc`   | c (handle)      | c                                                |                   |   
| `fputs`   | str (handle)    | str                                              |                   |   
| `fgetc`   | (handle)        |                                                  |                   |   
| `fgets`   | size (handle)   |                                                  |                   |     
| `fflush`  | (handle)        |                                                  |                   | 
| `fclose`  | (handle)        |                                                  |                   |  










