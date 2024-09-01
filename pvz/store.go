package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func (s *store) NewStore(f *os.File, size uint64) (*store, error) {
	file, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	println(file)
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// new position, # of written bytes, error
func (s *store) Append(p []byte) (uint64, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos := s.size
	if err := binary.Write(s.buf, enc, len(p)); err != nil {
		return 0, 0, nil
	}
	nn, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	nn += lenWidth
	writeLen := uint64(nn)
	pos += writeLen
	return pos, writeLen, nil
}

// because i don't know how much is the record size, i don't send an array
func (s *store) Read(pos int64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//flush
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	// find record size
	size := make([]byte, lenWidth)
	if _, err := s.ReadAt(size, pos); err != nil {
		return nil, err
	}
	// get record
	record := make([]byte, enc.Uint64(size))
	if _, err := s.ReadAt(record, pos+lenWidth); err != nil {
		return nil, err
	}
	return record, nil

}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}
