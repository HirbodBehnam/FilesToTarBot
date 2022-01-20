package database

import "errors"

type Interface interface {
	// AddFile must add one file the files of user
	// If the user doesn't have any files it must create it
	AddFile(userID int64, file File) error
	// Reset must delete all files of a user
	Reset(userID int64)
	// GetFiles must return a list of files for user which must be tarred
	GetFiles(userID int64) (files []File, tarSize int64)
}

// File holds the info needed for a file
// See tg.InputDocumentFileLocation for more info about fields
type File struct {
	FileReference []byte
	// Simply the file name
	Name       string
	ID         int64
	AccessHash int64
	Size       int64
}

// MainDatabase is the database which we should use
var MainDatabase = NewMemoryCache()

// MaxFileSize which can be uploaded to Telegram
const MaxFileSize = 1.5 * 1000 * 1000 * 1000

// TooBigFileError indicates that the result file is too big to be uploaded to telegram
var TooBigFileError = errors.New("file is too big to upload")

// FileAlreadyExistsError indicates that the filename already exists in tar archive
var FileAlreadyExistsError = errors.New("this filename already exists")
