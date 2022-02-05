package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const chunkSize = 8092

var (
	riff    = []byte{'\x52', '\x49', '\x46', '\x46'}
	waveFmt = []byte{'\x57', '\x41', '\x56', '\x45', '\x66', '\x6d', '\x74'}
	riffLen = 4
)

func findFirst(r io.ReaderAt, toFind []byte, offset int64) (int64, error) {
	chunk := make([]byte, chunkSize+len(toFind))
	for {
		n, err := r.ReadAt(chunk, offset)
		if err != nil && err != io.EOF {
			return -1, err
		}
		index := bytes.Index(chunk[:n], toFind)
		if index != -1 {
			return offset + int64(index), nil
		}
		if err == io.EOF {
			break
		}
		offset += int64(chunkSize)
	}
	return -1, nil
}

func getSize(f *os.File, offset int64) (int64, error) {
	sizeBuf := make([]byte, 4)
	waveBuf := make([]byte, 7)
	_, err := f.Seek(offset+4, io.SeekStart)
	if err != nil {
		return -1, err
	}
	_, err = f.Read(sizeBuf)
	if err != nil {
		return -1, err
	}
	_, err = f.Read(waveBuf)
	if err != nil {
		return -1, err
	}
	if !bytes.Equal(waveBuf, waveFmt) {
		return -1, nil
	}
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return -1, err
	}
	size := binary.LittleEndian.Uint32(sizeBuf)
	actualSize := int64(size) + 8
	return actualSize, nil
}

func writeRiff(outPath string, f *os.File, size int64) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.CopyN(outFile, f, size)
	return err
}

func main() {
	path := os.Args[1]
	var (
		i      int   = 1
		offset int64 = 0
	)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for {
		offset, err = findFirst(f, riff, offset)
		if err != nil {
			panic(err)
		}
		if offset == -1 {
			break
		}
		startOffset := offset
		fnameSuffix := "_" + strconv.Itoa(i) + ".wem"
		size, err := getSize(f, offset)
		if err != nil {
			panic(err)
		}
		if size == -1 {
			offset += int64(riffLen)
			continue
		}
		offset += size
		fmt.Println(filepath.Base(path) + fnameSuffix)
		fmt.Printf("%d bytes\nStart: 0x%X\nEnd: 0x%X\n\n", size, startOffset, offset)
		err = writeRiff(path+fnameSuffix, f, size)
		if err != nil {
			panic(err)
		}
		i++
	}
}
