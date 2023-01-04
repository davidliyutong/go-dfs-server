# Client


##  Regulations

client use `local:` and `remote:` prefix to decide if path is remote path (on the DFS cluster) or local path. A path is assumed to be remote path if no prefix is provided.

##  Commands

Client support multiple genres of actions:

- Authentication
- Maneuver

### Authentication

| Name     | Argument (CLI)  | Options                                              | Functional                                             |
|----------|-----------------|------------------------------------------------------|--------------------------------------------------------|
| `login`  | -               | `--accessKey=12345678` `--secretKey=xxxxxxxx`        |                                                        |
| `logout` | -               |                                                      |                                                        |

### Maneuver

| Name    | Argument (CLI)  | Options         | Functional                                             |
|---------|-----------------|-----------------|--------------------------------------------------------|
| `info`  | -               |                 |                                                        |
| `cat`   | \[path\]        |                 |                                                        |
| `pipe`  | \[path\]        | `--mode=in/out` | open a file and write file with data from stdin        |
| `ls`    | \[path\]        |                 | list a remote directory of file                        |
| `mkdir` | \[path\]        |                 | create a remote directory                              |
| `rm`    | \[path\]        | `-r`            | delete a remote file                                   |
| `touch` | \[path\]        |                 | create a remote file                                   |
| `get`   | \[src\] \[dst\] | `-r`            | download file from remote                              |
| `put`   | \[src\] \[dst\] | `-r`            | upload file to remote                                  |
| `mv`    | \[src\] \[dst\] |                 | move remote file                                       |
| `shell` |                 | `--mode=0x0`    | open a file with an interactive prompt **interactive** |

## Interactive Prompt

The client can open an interactive prompt for demo propose

| Name     | Argument (SDK)  | Argument (CLI) | Options                                          | Functional        |
|----------|-----------------|----------------|--------------------------------------------------|-------------------|
| `fopen`  | path            | path           |                                                  | open file         |
| `fseek`  | offset (handle) | \[offset\]     |                                                  | move file pointer |
| `fread`  | size (handle)   | \[size\]       |                                                  |                   |
| `fwrite` | data (handle)   |                | `--data=/path/to/file` or `--data=/path/to/file` |                   | 
| `fputc`  | c (handle)      | c              |                                                  |                   |   
| `fputs`  | str (handle)    | str            |                                                  |                   |   
| `fgetc`  | (handle)        |                |                                                  |                   |   
| `fgets`  | size (handle)   |                |                                                  |                   |     
| `fflush` | (handle)        |                |                                                  |                   | 
| `fclose` | (handle)        |                |                                                  |                   |  
| `cd`     |                 |                |                                                  |                   |  
| `ls`     |                 |                |                                                  |                   |  
| `mkdir`  |                 |                |                                                  |                   |  
| `rm`     |                 |                |                                                  |                   |  











