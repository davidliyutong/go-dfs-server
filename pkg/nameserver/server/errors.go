package server

import "errors"

var ErrBlobCorrupted = errors.New("blob corrupted")

var ErrChunkIDWrong = errors.New("wrong id")
var ErrChunkCurrentNotPresent = errors.New("current chunk is not present")
var ErrChunkOffsetInvalid = errors.New("invalid chunk offset")
var ErrChunkSizeTooLarge = errors.New("chunk size is too large")
var ErrChunkIDOutOfRange = errors.New("chunkID out of range")

var ErrClientNotFound = errors.New("client not found")
var ErrClientNotFoundSome = errors.New("failed to find some clients")

var ErrCreateErrorLocal = errors.New("create file failed, local error")
var ErrCreateErrorRemote = errors.New("create file failed, dataserver error")

var ErrDataserverCannotGet = errors.New("cannot get related data server")
var ErrDataServerInsufficient = errors.New("no enough data server available")
var ErrDataServerOffline = errors.New("some dataserver is offline")
var ErrDataServerChecksumMismatch = errors.New("not all servers have the same checksum")
var ErrDataServerUUIDMismatch = errors.New("not all data servers uuid match record")
var ErrDataServerReboot = errors.New("some dataserver has been rebooted")

var ErrDirectoryCannotRemoveLocal = errors.New("failed to remove directory")
var ErrDirectoryCannotRemoveRemote = errors.New("failed to remove directory from some data servers")

var ErrFileInValid = errors.New("not a valid file")
var ErrFileNotFile = errors.New("not a file")
var ErrFileNotExist = errors.New("file does not exist")
var ErrFileOrDirectoryNotExist = errors.New("file or directory does not exist")
var ErrFileSizeMismatch = errors.New("file size not match")
var ErrFilePresenceNil = errors.New("presence is nil")
var ErrFileOpened = errors.New("file already opened")

var IOReadSizeMismatch = errors.New("read size is not equal to size")

var ErrMetaDataCannotRead = errors.New("cannot open metadata")
var ErrModeInvalid = errors.New("invalid mode")

var ErrNameServerReboot = errors.New("nameserver has been rebooted")

var ErrSessionNotFound = errors.New("session not found")
var _ErrSessionDeleting = errors.New("session is deleting")

var ErrDirectoryCannotDelete = errors.New("cannot delete this directory")
var ErrDirectoryCannotCreate = errors.New("cannot create this directory")
