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

type FsStream interface {
	Stream
	Owner() bool
	SetOwner(bool)
}

type aFsStream struct {
	*aStream
}

func castFsStream(cfs *C.GMimeStreamFs) *aFsStream {
	fs := (*C.GMimeStream)(unsafe.Pointer(cfs))
	s := CastStream(fs)
	return &aFsStream{s}
}

func NewFsStream(fd int) FsStream {
	s := C.g_mime_stream_fs_new(C.int(fd))
	fStream := (*C.GMimeStreamFs)(unsafe.Pointer(s))
	defer unref(C.gpointer(fStream))
	return castFsStream(fStream)
}

func NewFsStreamWithBounds(fd int, start int64, end int64) FsStream {
	sBound := C.g_mime_stream_fs_new_with_bounds(C.int(fd), (C.gint64)(start), (C.gint64)(end))
	fStream := (*C.GMimeStreamFs)(unsafe.Pointer(sBound))
	defer unref(C.gpointer(fStream))
	return castFsStream(fStream)
}


func (f *aFsStream) rawFsStream() *C.GMimeStreamFs {
	return (*C.GMimeStreamFs)(f.pointer())
}

func (f *aFsStream) Owner() bool {
	result := C.g_mime_stream_fs_get_owner(f.rawFsStream())
	return gobool(result)
}

func (f *aFsStream) SetOwner(owner bool) {
	C.g_mime_stream_fs_set_owner(f.rawFsStream(), gbool(owner))
}
