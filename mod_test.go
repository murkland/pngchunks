package pngchunks_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/yumland/pngchunks"
)

func TestReadPNG(t *testing.T) {
	f, err := os.Open("testdata/test.png")
	if err != nil {
		t.Fatalf("Open(): %s", err)
	}
	defer f.Close()

	pngr, err := pngchunks.NewReader(f)
	if err != nil {
		t.Errorf("NewReader(): %s", err)
	}

	for {
		chunk, err := pngr.NextChunk()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Errorf("ReadChunk(): %s", err)
		}

		if chunk.Type() != "tEXt" {
			if _, err := io.Copy(ioutil.Discard, chunk); err != nil {
				t.Errorf("io.Copy(): %s", err)
			}
		} else {
			buf, err := ioutil.ReadAll(chunk)
			if err != nil {
				t.Errorf("ioutil.ReadAll(): %s", err)
			}
			if expected := "comment\x00hello there!"; string(buf) != expected {
				t.Errorf("expected tEXt = %s, got %s", expected, string(buf))
			}
		}

		if err := chunk.Close(); err != nil {
			t.Errorf("Close(): %s", err)
		}
	}
}

func TestWritePNG(t *testing.T) {
	var w bytes.Buffer

	f, err := os.Open("testdata/test.png")
	if err != nil {
		t.Fatalf("Open(): %s", err)
	}
	defer f.Close()

	pngr, err := pngchunks.NewReader(f)
	if err != nil {
		t.Errorf("NewReader(): %s", err)
	}

	pngw, err := pngchunks.NewWriter(&w)
	if err != nil {
		t.Errorf("NewWriter(): %s", err)
	}

	for {
		chunk, err := pngr.NextChunk()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Errorf("ReadChunk(): %s", err)
		}

		if chunk.Type() != "tEXt" {
			if err := pngw.WriteChunk(chunk.Length(), chunk.Type(), chunk); err != nil {
				t.Errorf("WriteChunk(): %s", err)
			}
		} else {
			if _, err := io.Copy(ioutil.Discard, chunk); err != nil {
				t.Errorf("io.Copy(): %s", err)
			}

			newComment := []byte("comment\x00hi everyone!")
			if err := pngw.WriteChunk(int32(len(newComment)), chunk.Type(), bytes.NewBuffer(newComment)); err != nil {
				t.Errorf("WriteChunk(): %s", err)
			}
		}

		if err := chunk.Close(); err != nil {
			t.Errorf("Close(): %s", err)
		}
	}

	r := bytes.NewBuffer(w.Bytes())

	pngr, err = pngchunks.NewReader(r)
	if err != nil {
		t.Errorf("NewReader(): %s", err)
	}

	for {
		chunk, err := pngr.NextChunk()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Errorf("ReadChunk(): %s", err)
		}

		if chunk.Type() != "tEXt" {
			if _, err := io.Copy(ioutil.Discard, chunk); err != nil {
				t.Errorf("io.Copy(): %s", err)
			}
		} else {
			buf, err := ioutil.ReadAll(chunk)
			if err != nil {
				t.Errorf("ioutil.ReadAll(): %s", err)
			}
			if expected := "comment\x00hi everyone!"; string(buf) != expected {
				t.Errorf("expected tEXt = %s, got %s", expected, string(buf))
			}
		}

		if err := chunk.Close(); err != nil {
			t.Errorf("Close(): %s", err)
		}
	}
}
