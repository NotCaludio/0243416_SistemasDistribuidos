package log

import (
	"fmt"
	"os"
	"path"

	api "github.com/notcaludio/pvz/api/v1"
	"google.golang.org/protobuf/proto"
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint64
	config                 Config
}

func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}
	var err error
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

// idk
func (s *segment) Append(record *api.Record) (uint64, error) {
	offset := s.nextOffset
	record.Offset = offset
	marshaledValue, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}
	_, storePosition, err := s.store.Append(marshaledValue)
	if err != nil {
		return 0, nil
	}
	if err := s.index.Write(
		// Index offsets are relative to the base offset on the store file
		uint32(s.nextOffset-uint64(s.baseOffset)),
		storePosition,
	); err != nil {
		return 0, err
	}
	s.nextOffset++
	return offset, nil
}

// idk
func (s *segment) Read(offset uint64) (record *api.Record, funcErr error) {

	_, pos, err := s.index.Read(int64(offset - s.baseOffset))
	if err != nil {
		return nil, err
	}
	value, err := s.store.Read(pos)
	if err != nil {
		return nil, err
	}
	funcErr = proto.Unmarshal(value, record)
	return record, funcErr

}

// idk
func (s *segment) IsMaxed() bool {
	return s.index.size >= s.config.Segment.MaxIndexBytes || s.store.size >= s.config.Segment.MaxStoreBytes
}

func (s *segment) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		return nil
	}
	if err := os.Remove(s.index.Name()); err != nil {
		return nil
	}
	return nil
}

func (s *segment) Close() error {
	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	return nil

}
