package fap

/*
#cgo LDFLAGS: -lfap
#include <stdlib.h>
#include <string.h>
#include <fap.h>

char* str_at(char **lst, int idx) {
    return lst[idx];
}
char* new_c_str(unsigned int size) {
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

type fap_packet_t C.fap_packet_t

func init() {
	C.fap_init()
}

// Cleanup should be called when done using this package.
func Cleanup() {
	C.fap_cleanup()
}

// ParseAprs is the main parser method. Parses content of input
// string. Setting isAX25 to true source callsign and path
// elements are checked to be strictly compatible with AX.25
// specs so that theycan be sent into AX.25 network. Destination
// callsign is always checked this way.
func ParseAprs(input string, isAX25 bool) (*Packet, error) {
	c_input := C.CString(input)
	defer C.free(unsafe.Pointer(c_input))

	c_len := C.uint(C.strlen(c_input))

	var c_isAX25 C.short
	if isAX25 {
		c_isAX25 = 1
	}

	c_packet := C.fap_parseaprs(c_input, c_len, c_isAX25)
	defer C.fap_free(c_packet)

	if c_packet == nil {
		log.Fatal("fap_parseaprs returned nil. Is libfap initialized?")
	}

	packet, err := (*fap_packet_t)(c_packet).goPacket()

	return packet, err
}

// Calculates distance between given locations,
// returning the the distance in kilometers.
func Distance(lon0, lat0, lon1, lat1 float64) float64 {
	c_dist := C.fap_distance(
		C.double(lon0), C.double(lat0),
		C.double(lon1), C.double(lat1),
	)

	return float64(c_dist)
}

// Calculates direction from first to second location.
func Direction(lon0, lat0, lon1, lat1 float64) float64 {
	c_dir := C.fap_direction(
		C.double(lon0), C.double(lat0),
		C.double(lon1), C.double(lat1),
	)

	return float64(c_dir)
}

// MicEMbitsToMessage converts mic-e message bits (three numbers 0-2)
// to a textual message.
func MicEMbitsToMessage(mbits string) string {
	if mbits == "" {
		log.Fatal("MicEMbitsToMessage() called with empty string")
	}

	buffer := C.new_c_str(60)
	defer C.free(unsafe.Pointer(buffer))

	C.fap_mice_mbits_to_message(C.CString(mbits), buffer)

	return C.GoString(buffer)
}

func (c *fap_packet_t) goPacket() (*Packet, error) {
	err := c.error()

	packet := Packet{
		// error_code (removed)
		Type: UNKNOWN, // set below

		OrigPacket: goString(c.orig_packet),
		// orig_packet_len (removed)

		Header: goString(c.header),
		Body:   goString(c.body),
		// body_len (removed)

		SrcCallsign: goString(c.src_callsign),
		DstCallsign: goString(c.dst_callsign),
		// path (set below)
		// path_len (removed)

		Latitude:      goFloat64(c.latitude),
		Longitude:     goFloat64(c.longitude),
		PosResolution: goFloat64(c.pos_resolution),
		PosAmbiguity:  goUnsignedInt(c.pos_ambiguity),
		PosFormat: POS_UNKNOWN, // set below
		DaoDatumByte:  byte(c.dao_datum_byte), // 0x00 = undef
		Altitude:      goFloat64(c.altitude),
		Course:        goUnsignedInt(c.course),
		Speed:         goFloat64(c.speed),

		SymbolTable: byte(c.symbol_table), // 0x00 = undef
		SymbolCode:  byte(c.symbol_code),  // 0x00 = undef

		Messaging:   goBool(c.messaging),
		Destination: goString(c.destination),
		Message:     goString(c.message),
		MessageAck:  goString(c.message_ack),
		MessageNack: goString(c.message_nack),
		MessageId:   goString(c.message_id),

		// comment (set below)
		// comment_len (removed)

		ObjectOrItemName: goString(c.object_or_item_name),
		Alive:            goBool(c.alive),
		GpsFixStatus:     goBool(c.gps_fix_status),
		RadioRange:       goUnsignedInt(c.radio_range),
		Phg:              goString(c.phg),

		// timestamp (set below)
		NmeaChecksumOk: goBool(c.nmea_checksum_ok),

		// wx_report (TODO)
		// telemetry (TODO)

		Messagebits: goString(c.messagebits),

		// status (set below)

		// capabilities (TODO)
		// capabilities_len (removed)
	}

	if t := C.packet_type((*C.fap_packet_t)(c)); t != nil {
		packet.Type = PacketType(*t)
	}
	if c.format != nil {
		packet.PosFormat = PositionFormat(*c.format)
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

func (c *fap_packet_t) error() error {
	if c.error_code == nil {
		return nil
	}

	buffer := C.new_c_str(64)
	defer C.free(unsafe.Pointer(buffer))

	C.fap_explain_error(*c.error_code, buffer)

	return errors.New(C.GoString(buffer))
}
