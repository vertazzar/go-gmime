package gmime

/*
#cgo pkg-config: gmime-3.0
#include <stdlib.h>
#include <gmime/gmime.h>
static gboolean object_is_part(GTypeInstance *obj) {
    return GMIME_IS_PART(obj);
}
static gboolean object_is_multipart(GTypeInstance *obj) {
    return GMIME_IS_MULTIPART(obj);
}
static gboolean object_is_message(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE(obj);
}
static gboolean object_is_message_part(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE_PART(obj);
}
static gboolean object_is_partial(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE_PARTIAL(obj);
}
static GMimePart * gmime_part(GMimeObject *obj) {
	return GMIME_PART(obj);
}
*/
import "C"
import (
	"unsafe"
)

type Object interface {
	Janitor
	ContentType() ContentType
	ContentID() (string, bool)
	Header(string) (string, bool)
	ToString() string
	ContentDisposition() ContentDisposition
	Headers() string
	WalkHeaders(cb func(string, string) error) error
}

type anObject struct {
	*PointerMixin
}

type rawObject interface {
	Object
	rawObject() *C.GMimeObject
}

func CastObject(o *C.GMimeObject) *anObject {
	return &anObject{CastPointer(C.gpointer(o))}
}

func NewObject(contentType ContentType) Object {
	rawContentType := contentType.(rawContentType)
	object := C.g_mime_object_new(nil, rawContentType.rawContentType())
	defer unref(C.gpointer(object))
	o := objectAsSubclass(object)
	return o
}

func NewObjectWithType(ctype string, csubtype string) Object {
	var _ctype *C.char = C.CString(ctype)
	var _csubtype *C.char = C.CString(csubtype)
	defer C.free(unsafe.Pointer(_ctype))
	defer C.free(unsafe.Pointer(_csubtype))

	object := C.g_mime_object_new_type(nil, _ctype, _csubtype)
	defer unref(C.gpointer(object))
	o := objectAsSubclass(object)

	return o
}

func (o *anObject) ContentID() (string, bool) {
	cid := C.g_mime_object_get_content_id(o.rawObject())
	return maybeGoString(cid)
}


func (o *anObject) ContentType() ContentType {
	if ct := C.g_mime_object_get_content_type(o.rawObject()); ct != nil {
		return CastContentType(ct)
	}
	return nil
}

func (o *anObject) Header(name string) (string, bool) {
	var _name *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(_name))

	return maybeGoString(C.g_mime_object_get_header(o.rawObject(), _name))
}

func (o *anObject) ToString() string {
	str := C.g_mime_object_to_string(o.rawObject(), nil)
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str)
}

func (o *anObject) ContentDisposition() ContentDisposition {
	cd := C.g_mime_object_get_content_disposition(o.rawObject())
	if cd != nil {
		return CastContentDisposition(cd)
	}
	return nil
}

func (o *anObject) Headers() string {
	headers := C.g_mime_object_get_headers(o.rawObject(), nil)
	defer C.free(unsafe.Pointer(headers))

	return C.GoString(headers)
}


func (o *anObject) rawObject() *C.GMimeObject {
	return (*C.GMimeObject)(o.pointer())
}

func objectAsSubclass(o *C.GMimeObject) Object {
	partType := (*C.GTypeInstance)(unsafe.Pointer(o))

	if gobool(C.object_is_message_part(partType)) {
		return CastMessagePart((*C.GMimeMessagePart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_partial(partType)) {
		return CastMessagePartial((*C.GMimeMessagePartial)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_multipart(partType)) {
		return CastMultipart((*C.GMimeMultipart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_part(partType)) {
		return CastPart((*C.GMimePart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_message(partType)) {
		return CastMessage((*C.GMimeMessage)(unsafe.Pointer(o)))
	} else {
		return CastObject(o)
	}
}

func (o *anObject) WalkHeaders(cb func(string, string) error) error {
	ghl := C.g_mime_object_get_header_list(o.rawObject())
	defer C.free(unsafe.Pointer(ghl))
	i := 0
	for {
		header := C.g_mime_header_list_get_header_at(ghl, C.int(i))
		if header == nil {
			return nil
		}
		name := C.GoString(C.g_mime_header_get_name(header))
		value := C.GoString(C.g_mime_header_get_value(header))

		C.free(unsafe.Pointer(header))

		err := cb(name, value)
		if err != nil {
			return err
		}
		i ++
	}
}

// Very minimal interface, to inspection only
type HeaderIterator interface {
	Janitor
	Name() string
	Value() string
	Next() bool
}

type aHeader struct {
}
