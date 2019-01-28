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

type Message interface {
	Object
	SetSubject(string)
	Subject() (string, bool)
	SetMessageId(string)
	MessageId() (string, bool)
	To() *InternetAddressList
	Cc() *InternetAddressList
	Bcc() *InternetAddressList
	AllRecipients() *InternetAddressList
	SetMimePart(Object)
	MimePart() Object
	Body() Object
	Date() string
}

type aMessage struct {
	*anObject
}

type rawMessage interface {
	Message
	rawMessage() *C.GMimeMessage
}

func CastMessage(m *C.GMimeMessage) *aMessage {
	return &aMessage{CastObject((*C.GMimeObject)(unsafe.Pointer(m)))}
}

func NewMessage() Message {
	m := C.g_mime_message_new(0)
	defer unref(C.gpointer(m))
	return CastMessage(m)
}



func (m *aMessage) SetSubject(subject string) {
	var cSubject *C.char = C.CString(subject)
	C.g_mime_message_set_subject(m.rawMessage(), cSubject, nil)
	C.free(unsafe.Pointer(cSubject))
}

func (m *aMessage) Subject() (string, bool) {
	subject := C.g_mime_message_get_subject(m.rawMessage())
	return maybeGoString(subject)
}

func (m *aMessage) SetMessageId(messageId string) {
	var cMessageId *C.char = C.CString(messageId)
	C.g_mime_message_set_message_id(m.rawMessage(), cMessageId)
	C.free(unsafe.Pointer(cMessageId))
}

func (m *aMessage) MessageId() (string, bool) {
	messageId := C.g_mime_message_get_message_id(m.rawMessage())
	return maybeGoString(messageId)
}




func (m *aMessage) To() *InternetAddressList {
	cList := C.g_mime_message_get_addresses(m.rawMessage(), C.GMIME_ADDRESS_TYPE_TO)
	return CastInternetAddressList(cList)
}


func (m *aMessage) Cc() *InternetAddressList {
	cList := C.g_mime_message_get_addresses(m.rawMessage(), C.GMIME_ADDRESS_TYPE_CC)
	return CastInternetAddressList(cList)
}


func (m *aMessage) Bcc() *InternetAddressList {
	cList := C.g_mime_message_get_addresses(m.rawMessage(), C.GMIME_ADDRESS_TYPE_BCC)
	return CastInternetAddressList(cList)
}

func (m *aMessage) AllRecipients() *InternetAddressList {
	// This is major exception: we have newly allocated list here
	cList := C.g_mime_message_get_all_recipients(m.rawMessage())
	defer unref(C.gpointer(cList))
	return CastInternetAddressList(cList)
}

func (m *aMessage) SetMimePart(mimePart Object) {
	part := mimePart.(rawObject)
	switch mimePart.(type) {
	case Part:
		C.g_mime_message_set_mime_part(m.rawMessage(), part.rawObject())
	case Multipart:
		C.g_mime_message_set_mime_part(m.rawMessage(), part.rawObject())
	}
}

func (m *aMessage) MimePart() Object {
	object := C.g_mime_message_get_mime_part(m.rawMessage())
	return objectAsSubclass(object)
}

func (m *aMessage) Body() Object {
	object := C.g_mime_message_get_body(m.rawMessage())
	return objectAsSubclass(object)
}


func (m *aMessage) rawMessage() *C.GMimeMessage {
	return (*C.GMimeMessage)(m.pointer())
}
