package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	endWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	indexx := &index{
		file: f,
	}
	fileName, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	indexx.size = uint64(fileName.Size())
	if err := os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}
	if indexx.mmap, err = gommap.Map(indexx.file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED); err != nil {
		return nil, err
	}
	return indexx, nil

}

func (i *index) Read(entries int64) (offset uint32, storePosition uint64, err error) {
	offset = uint32(entries)
	storePosition = uint64(offset) * endWidth
	if i.size < storePosition+endWidth {
		return 0, 0, io.EOF
	}
	offset = enc.Uint32(i.mmap[storePosition : storePosition+offWidth])
	storePosition = enc.Uint64(i.mmap[storePosition+offWidth : storePosition+endWidth])

	return offset, storePosition, nil
}

func (i *index) Write(offset uint32, storePosition uint64) error {
	indexSize := i.size
	if uint64(len(i.mmap)) < indexSize+endWidth {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], offset)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+endWidth], storePosition)
	i.size += endWidth
	return nil
}

func (i *index) Close() error {
	//flush the mmap
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	//flush the file
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	if err := i.file.Close(); err != nil {
		return err
	}
	return nil
}
func (i *index) Name() string {
	return i.file.Name()
}
