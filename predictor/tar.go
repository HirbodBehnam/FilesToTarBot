package predictor

import (
	"archive/tar"
)

// Tar tries to predict the final filesize of a tar format based on its headers
// It only works with headers with tar.FormatGNU format
type Tar struct {
	total int64
}

// blockPadding computes the number of bytes needed to pad offset up to the
// nearest block edge where 0 <= n < blockSize.
// From tar/format.go
func blockPadding(offset int64) (n int64) {
	return -offset & (512 - 1)
}

// AddFile adds one file to header
// The only possible error value is FilenameTooLarge
// header.Format must be tar.FormatGNU
func (t *Tar) AddFile(header tar.Header) {
	if len(header.Name) > 100 {
		t.addFileSize(int64(len(header.Name)) + 1)
	}
	t.addFileSize(header.Size)
}

func (t *Tar) addFileSize(size int64) {
	t.total += 512 + size + blockPadding(size)
}

// Total gets the total filesize which will be written
func (t *Tar) Total() int64 {
	return t.total + 512*2
}
