package fap

/*
#cgo LDFLAGS: -lfap
#include <stdlib.h>
#include <string.h>
#include <fap.h>

char* str_at(char **lst, int idx) {
    return lst[idx];
}
char* new_c_str(uint size) {
    return (char*) malloc( size*sizeof(char) );
}
fap_packet_type_t* packet_type(fap_packet_t *p) {
    return p->type;
}
*/
import "C"

import (
	"errors"
	"log"
	"time"
	"unsafe"
)

const (
	UNKNOWN          = 0
	LOCATION         = C.fapLOCATION
	OBJECT           = C.fapOBJECT
	ITEM             = C.fapITEM
	MICE             = C.fapMICE
	NMEA             = C.fapNMEA
	WX               = C.fapWX
	MESSAGE          = C.fapMESSAGE
	CAPABILITIES     = C.fapCAPABILITIES
	STATUS           = C.fapSTATUS
	TELEMETRY        = C.fapTELEMETRY
	TELEMETRYMESSAGE = C.fapTELEMETRY_MESSAGE
	DXSPOT           = C.fapDX_SPOT
	EXPERIMENTAL     = C.fapEXPERIMENTAL
)

const (
	POSUNKNOWN      = 0
	POSCOMPRESSED   = C.fapPOS_COMPRESSED
	POSUNCOMPRESSED = C.fapPOS_UNCOMPRESSED
	POSMICE         = C.fapPOS_MICE
	POSNMEA         = C.fapPOS_NMEA
)

func init() {
	C.fap_init()
}

func Cleanup() {
	C.fap_cleanup()
}

func ParseAprs(input string, isAX25 bool) (*FapPacket, error) {
	c_input := C.CString(input)
	defer C.free(unsafe.Pointer(c_input))

	c_len := C.uint(C.strlen(c_input))

	var c_isAX25 C.short
	if isAX25 {
		c_isAX25 = 1
	}

	c_fapPacket := C.fap_parseaprs(c_input, c_len, c_isAX25)
	defer C.fap_free(c_fapPacket)

	if c_fapPacket == nil {
		log.Fatal("fap_parseaprs returned nil. Is libfap initialized?")
	}

	fapPacket, err := c_fapPacket.goFapPacket()

	return fapPacket, err
}

func Distance(lon0, lat0, lon1, lat1 float64) float64 {
	c_dist := C.fap_distance(
		C.double(lon0), C.double(lat0),
		C.double(lon1), C.double(lat1),
	)

	return float64(c_dist)
}

func Direction(lon0, lat0, lon1, lat1 float64) float64 {
	c_dir := C.fap_direction(
		C.double(lon0), C.double(lat0),
		C.double(lon1), C.double(lat1),
	)

	return float64(c_dir)
}

func MicEMbitsToMessage(mbits string) string {
	if mbits == "" {
		log.Fatal("MicEMbitsToMessage() called with empty string")
	}

	buffer := C.new_c_str(60)
	defer C.free(unsafe.Pointer(buffer))

	C.fap_mice_mbits_to_message(C.CString(mbits), buffer)

	return C.GoString(buffer)
}

func (c *_Ctype_fap_packet_t) goFapPacket() (*FapPacket, error) {
	err := c.error()

	packet := FapPacket{
		// error_code (removed)
		// type -> PacketType (set below)

		OrigPacket: goString(c.orig_packet),
		// orig_packet_len (removed)

		Header: goString(c.header),
		Body:   goString(c.body),
		// body_len (removed)

		SrcCallsign: goString(c.src_callsign),
		DstCallsign: goString(c.dst_callsign),
		// path (set below)
		// path_len (removed)

		Latitude:      goFloat64Ptr(c.latitude),
		Longitude:     goFloat64Ptr(c.longitude),
		PosResolution: goFloat64Ptr(c.pos_resolution),
		PosAmbiguity:  goUnsignedIntPtr(c.pos_ambiguity),
		DaoDatumByte:  byte(c.dao_datum_byte), // 0x00 = undef
		Altitude:      goFloat64Ptr(c.altitude),
		Course:        goUnsignedIntPtr(c.course),
		Speed:         goFloat64Ptr(c.speed),

		SymbolTable: byte(c.symbol_table), // 0x00 = undef
		SymbolCode:  byte(c.symbol_code),  // 0x00 = undef

		Messaging:   goBoolPtr(c.messaging),
		Destination: goString(c.destination),
		Message:     goString(c.message),
		MessageAck:  goString(c.message_ack),
		MessageNack: goString(c.message_nack),
		MessageId:   goString(c.message_id),

		// comment (set below)
		// comment_len (removed)

		ObjectOrItemName: goString(c.object_or_item_name),
		Alive:            goBoolPtr(c.alive),
		GpsFixStatus:     goBoolPtr(c.gps_fix_status),
		RadioRange:       goUnsignedIntPtr(c.radio_range),
		Phg:              goString(c.phg),

		// timestamp (set below)
		NmeaChecksumOk: goBoolPtr(c.nmea_checksum_ok),

		// wx_report (TODO)
		// telemetry (TODO)

		Messagebits: goString(c.messagebits),

		// status (set below)

		// capabilities (TODO)
		// capabilities_len (removed)
	}

	if C.packet_type(c) != nil {
		packet.PacketType = uint(*C.packet_type(c))
	}
	if c.format != nil {
		packet.Format = uint(*c.format)
	}
	if c.status != nil {
		packet.Status = C.GoStringN(c.status, C.int(c.status_len))
	}
	if c.comment != nil {
		packet.Comment = C.GoStringN(c.comment, C.int(c.comment_len))
	}
	if c.timestamp != nil {
		packet.Timestamp = time.Unix(int64(*c.timestamp), 0)
	}

	// Get path
	packet.Path = make([]string, int(c.path_len))
	for i := 0; i < int(c.path_len); i++ {
		packet.Path[i] = goString(C.str_at(c.path, C.int(i)))
	}

	return &packet, err
}

func (c *_Ctype_fap_packet_t) error() error {
	if c.error_code == nil {
		return nil
	}

	buffer := C.new_c_str(64)
	defer C.free(unsafe.Pointer(buffer))

	C.fap_explain_error(*c.error_code, buffer)

	return errors.New(C.GoString(buffer))
}
