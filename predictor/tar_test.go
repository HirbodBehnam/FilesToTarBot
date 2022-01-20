package predictor

import (
	"archive/tar"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var allChars = []rune(".abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ضصثقفغعهخحمنتالبیسشظطزرذدئ")
var unicodeOnly = []rune("ضصثقفغعهخحمنتالبیسشظطزرذدئ")

type writerCounter struct {
	total int64
}

func (w *writerCounter) Write(b []byte) (int, error) {
	w.total += int64(len(b))
	return len(b), nil
}

func TestUnicode(t *testing.T) {
	const tests = 100
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < tests; i++ {
		counter := new(writerCounter)
		realTar := tar.NewWriter(counter)
		fakeTar := new(Tar)
		name := randomString(rng, unicodeOnly)
		size := rng.Int63n(1024 * 1024)
		writeTar(t, realTar, fakeTar, name, size)
		_ = realTar.Close()
		if fakeTar.Total() != counter.total {
			t.Errorf("mismach total %d vs %d\nname = %s, size = %d\n", fakeTar.Total(), counter.total, name, size)
		}
	}
}

func TestFuzz(t *testing.T) {
	const tests = 1000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for c := 0; c < tests; c++ {
		totalFiles := rng.Intn(100)
		counter := new(writerCounter)
		realTar := tar.NewWriter(counter)
		fakeTar := new(Tar)
		for i := 0; i < totalFiles; i++ {
			writeTar(t, realTar, fakeTar, randomString(rng, allChars), rng.Int63n(1024*1024))
		}
		if err := realTar.Close(); err != nil {
			t.Fatalf("cannot close tar: %s\n", err)
		}
		if fakeTar.Total() != counter.total {
			t.Errorf("mismach total %d vs %d\n", fakeTar.Total(), counter.total)
		}
	}
}

func writeTar(t *testing.T, realTar *tar.Writer, fakeTar *Tar, name string, size int64) {
	header := &tar.Header{
		Name:   name,
		Size:   size,
		Mode:   0600,
		Format: tar.FormatGNU,
	}
	if err := realTar.WriteHeader(header); err != nil {
		t.Fatalf("cannot write header %v: %s\n", *header, err)
	}
	if _, err := realTar.Write(make([]byte, size)); err != nil {
		t.Fatalf("cannot write data %v: %s\n", *header, err)
	}
	fakeTar.AddFile(*header)
}

func randomString(rng *rand.Rand, source []rune) string {
	size := 1
	if rng.Int()%2 == 0 {
		size += rng.Intn(10)
	} else {
		size += rng.Intn(1000)
	}
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		sb.WriteRune(source[rng.Intn(len(source))])
	}
	return sb.String()
}
