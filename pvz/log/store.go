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

func newStore(file *os.File) (*store, error) {
	fileName, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fileName.Size())
	return &store{
		File: file,
		size: size,
		buf:  bufio.NewWriter(file),
	}, nil
}

// , # of written bytes, start position ,error
func (s *store) Append(record []byte) (uint64, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos := s.size
	//write the size into the buffer
	if err := binary.Write(s.buf, enc, uint64(len(record))); err != nil {
		return 0, 0, nil
	}
	//write the record
	nn, err := s.buf.Write(record)
	if err != nil {
		return 0, 0, err
	}
	nn += lenWidth
	writeLen := uint64(nn)
	s.size += writeLen
	return writeLen, pos, nil
}

// because i don't know how much is the record size, i don't send an array
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//flush
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	// find record size
	size := make([]byte, lenWidth)
	if _, err := s.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	// get record
	record := make([]byte, enc.Uint64(size))
	if _, err := s.ReadAt(record, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return record, nil

}

// it reads until the slice is filled
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// flush and close
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
