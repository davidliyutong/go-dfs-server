package status

import "errors"

var ErrBlobCorrupted = errors.New("blob corrupted")

var ErrChunkIDWrong = errors.New("wrong id")
var ErrChunkCurrentNotPresent = errors.New("current chunk is not present")
var ErrChunkOffsetInvalid = errors.New("invalid chunk offset")
var ErrChunkSizeTooLarge = errors.New("chunk size is too large")
var ErrChunkIDOutOfRange = errors.New("chunkID out of range")
var ErrChunkCannotRemove = errors.New("cannot remove chunk")

var ErrClientNotFound = errors.New("client not found")
var ErrClientNotFoundSome = errors.New("failed to find some clients")

var ErrCreateErrorLocal = errors.New("create file failed, local error")
var ErrCreateErrorRemote = errors.New("create file failed, dataserver error")

var _ = errors.New("cannot get related data server")
var ErrDataServerInsufficient = errors.New("no enough data server available")
var ErrDataServerOfflineSome = errors.New("some dataserver is offline")
var ErrDataServerChecksumMismatch = errors.New("not all servers have the same checksum")
var ErrDataServerUUIDMismatch = errors.New("not all data servers uuid match record")
var ErrDataServerReboot = errors.New("some dataserver has been rebooted")
var ErrDataServerCannotRemoveSome = errors.New("failed to remove from some data servers")
var ErrDataServerCannotCreateSome = errors.New("failed to create from some data servers")

var ErrDirectoryCannotCreate = errors.New("failed to create directory")
var _ = errors.New("failed to remove directory")
var ErrDirectoryCannotRemoveLocal = errors.New("failed to remove directory local")
var ErrDirectoryCannotRemoveRemote = errors.New("failed to remove directory from some data servers")

var ErrFileAlreadyLocked = errors.New("file already locked by this session")
var ErrFileInValid = errors.New("not a valid file")
var ErrFileNotFile = errors.New("not a file")
var ErrFileExists = errors.New("file or directory exists")
var ErrFileNotExist = errors.New("file does not exist")
var ErrFileNotLocked = errors.New("file not locked")

var ErrFileOrDirectoryNotExist = errors.New("file or directory does not exist")
var ErrFileSizeMismatch = errors.New("file size not match")
var ErrFilePresenceNil = errors.New("presence is nil")
var ErrFileOpened = errors.New("file already opened")

var ErrIOReadSizeMismatch = errors.New("read size is not equal to size")
var ErrIOReadLess = errors.New("read less bytes than expected")
var ErrIOWriteLess = errors.New("write less bytes than expected")
var ErrIOReadOnly = errors.New("file is opened in read-only mode")
var ErrIOModeInvalid = errors.New("invalid mode")
var ErrIOOffsetInvalid = errors.New("invalid offset")

var ErrMetaDataCannotLoad = errors.New("cannot load metadata")
var ErrMetaDataCannotDump = errors.New("cannot dump metadata")
var ErrMetaVersionConflict = errors.New("version conflict")
var ErrMetaNotFound = errors.New("no meta data found for this chunk id")
var ErrMetaTruncateFailedRemote = errors.New("remote truncate failed")

var ErrNameServerReboot = errors.New("nameserver has been rebooted")

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionDeleting = errors.New("session is deleting")

var ErrDirectoryCannotDelete = errors.New("cannot delete this directory")
var _ = errors.New("cannot create this directory")

var ErrClientNotPong = errors.New("response not pong")
var ErrClientLoginFailed = errors.New("login failed, access denied, try logout then login")
