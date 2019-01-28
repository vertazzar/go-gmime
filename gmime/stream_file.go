package gmime

/*
#cgo pkg-config: gmime-3.0
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

import (
	"unsafe"
)

type FileStream interface {
	Stream
	Owner() bool
	SetOwner(bool)
}

type aFileStream struct {
	*aStream
}

func CastFileStream(cfs *C.GMimeStreamFile) *aFileStream {
	fs := (*C.GMimeStream)(unsafe.Pointer(cfs))
	s := CastStream(fs)
	return &aFileStream{s}
}

func NewFileStream(f *C.FILE) FileStream {
	return NewFileStreamWithMode(f, "a")
}

func NewFileStreamWithMode(f *C.FILE, mode string) FileStream {
	cMode := C.CString(mode)
	defer C.free(unsafe.Pointer(cMode))
	s := C.g_mime_stream_file_new(f)
	fileStream := (*C.GMimeStreamFile)(unsafe.Pointer(s))
	defer unref(C.gpointer(fileStream))
	return CastFileStream(fileStream)
}

func NewFileStreamWithBounds(f *C.FILE, start int64, end int64) FileStream {
	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))
	sBound := C.g_mime_stream_file_new_with_bounds(f, (C.gint64)(start), (C.gint64)(end))
	fileStream := (*C.GMimeStreamFile)(unsafe.Pointer(sBound))
	defer unref(C.gpointer(fileStream))
	return CastFileStream(fileStream)
}

func (f *aFileStream) rawFileStream() *C.GMimeStreamFile {
	return (*C.GMimeStreamFile)(f.pointer())
}

func (f *aFileStream) Owner() bool {
	result := C.g_mime_stream_file_get_owner(f.rawFileStream())
	return gobool(result)
}

func (f *aFileStream) SetOwner(owner bool) {
	C.g_mime_stream_file_set_owner(f.rawFileStream(), gbool(owner))
}
