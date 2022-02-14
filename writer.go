package pngchunks

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) (*Writer, error) {
	if _, err := io.WriteString(w, header); err != nil {
		return nil, err
	}
	return &Writer{w}, nil
}

func (w *Writer) WriteChunk(length int32, typ string, r io.Reader) error {
	if err := binary.Write(w.w, binary.BigEndian, length); err != nil {
		return err
	}

	if _, err := w.w.Write([]byte(typ)); err != nil {
		return err
	}

	checksummer := crc32.NewIEEE()
	checksummer.Write([]byte(typ))

	if _, err := io.CopyN(io.MultiWriter(w.w, checksummer), r, int64(length)); err != nil {
		return err
	}

	if err := binary.Write(w.w, binary.BigEndian, checksummer.Sum32()); err != nil {
		return err
	}

	return nil
}
