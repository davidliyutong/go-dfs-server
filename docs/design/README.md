# System Design

## Nameserver APIs

### File Management

- `OpenFile(path)`
- `CloseFile(handle)`
- `Read(handle)`
- `Write(handle)`
- `SeekFile(handle, offset)`
- `LockFile(handle, offset)`

### Directory Management

- `MakeDirectory(directory)`
- `ListDirectory(directory)`
- `DirectoryExists(directory)`
- `RemoveDirectory(directory)`

### Data Server Management

- `AlivenessProbe(server)`
- `DeleteServer(server)`

## File Metadata

```json
{
  "uuid": "ddsfa-safsa",
  "filename": "xxx",
  "ctime": "xxx",
  "length": 1000,
  "n_replicas": 3,
  "chunks": [
    
  ],
  "distrubution": [
    
  ]
}

```