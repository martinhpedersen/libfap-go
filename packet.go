package fap

//#include <fap.h>
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (t PacketType) String() string {
	switch (t) {
		case UNKNOWN:
			return "Unknown"
		case LOCATION:
			return "Location"
		case OBJECT:
			return "Object"
		case ITEM:
			return "Item"
		case MICE:
			return "Mic-E"
		case NMEA:
			return "NMEA"
		case WX:
			return "WX"
		case MESSAGE:
			return "Message"
		case CAPABILITIES:
			return "Capabilities"
		case STATUS:
			return "Status"
		case TELEMETRY:
			return "Telemetry"
		case TELEMETRYMESSAGE:
			return "Telemetry Message"
		case DXSPOT:
			return "DX Spot"
		case EXPERIMENTAL:
			return "Experimental"
		default:
			panic(fmt.Sprintf("Missing Stringer for %d", t))
	}
	return ""
}

const (
	UNKNOWN          PacketType = -1
	LOCATION         PacketType = C.fapLOCATION
	OBJECT           PacketType = C.fapOBJECT
	ITEM             PacketType = C.fapITEM
	MICE             PacketType = C.fapMICE
	NMEA             PacketType = C.fapNMEA
	WX               PacketType = C.fapWX
	MESSAGE          PacketType = C.fapMESSAGE
	CAPABILITIES     PacketType = C.fapCAPABILITIES
	STATUS           PacketType = C.fapSTATUS
	TELEMETRY        PacketType = C.fapTELEMETRY
	TELEMETRYMESSAGE PacketType = C.fapTELEMETRY_MESSAGE
	DXSPOT           PacketType = C.fapDX_SPOT
	EXPERIMENTAL     PacketType = C.fapEXPERIMENTAL
)

const (
	POS_UNKNOWN      PositionFormat = -1
	POS_COMPRESSED   PositionFormat = C.fapPOS_COMPRESSED
	POS_UNCOMPRESSED PositionFormat = C.fapPOS_UNCOMPRESSED
	POS_MICE         PositionFormat = C.fapPOS_MICE
	POS_NMEA         PositionFormat = C.fapPOS_NMEA
)

type PacketType int
type PositionFormat int

// Packet is the APRS packet type
type Packet struct {
	Type       PacketType
	OrigPacket string

	Header string
	Body   string

	SrcCallsign string
	DstCallsign string
	Path        []string

	Latitude      float64
	Longitude     float64
	PosFormat     PositionFormat
	PosResolution float64
	PosAmbiguity  uint
	DaoDatumByte  byte
	Altitude      float64
	Course        uint
	Speed         float64

	SymbolTable byte
	SymbolCode  byte

	Messaging   bool
	Destination string
	Message     string
	MessageAck  string
	MessageNack string
	MessageId   string

	Comment          string
	ObjectOrItemName string
	Alive            bool
	GpsFixStatus     bool
	RadioRange       uint
	Phg              string
	Timestamp        time.Time
	NmeaChecksumOk   bool

	WxReport  string
	Telemetry string

	Messagebits  string
	Status       string
	Capabilities []string
}

// HasLocation returns true if packet has location
// data.
func (p *Packet) HasLocation() bool {
	return p.PosFormat != POS_UNKNOWN
}

// MicEMessage returns the texual Mic-E message.
func (p *Packet) MicEMessage() string {
	if p.Messagebits == "" {
		return ""
	}

	return MicEMbitsToMessage(p.Messagebits)
}

// Distance returns the distance to the given packet b in km.
func (a *Packet) Distance(b *Packet) (float64, error) {
	if b == nil {
		return 0, errors.New("Distance between A and nil is undefined")
	}
	if a.Latitude == 0 || a.Longitude == 0 ||
		b.Longitude == 0 || b.Latitude == 0 {
		return 0, errors.New("One or more components is nil when calculating distance")
	}

	return Distance(
		a.Longitude, a.Latitude,
		b.Longitude, b.Latitude,
	), nil
}

func (a *Packet) Direction(b *Packet) (float64, error) {
	if b == nil {
		return 0, errors.New("Direction between A and nil is undefined")
	}
	if a.Latitude == 0 || a.Longitude == 0 ||
		b.Longitude == 0 || b.Latitude == 0 {
		return 0, errors.New("One or more components is nil when calculating direction")
	}

	return Direction(
		a.Longitude, a.Latitude,
		b.Longitude, b.Latitude,
	), nil
}

var cardinals = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

func (a *Packet) IntercardinalDirection(b *Packet) (string, error) {
	degrees, err := a.Direction(b)
	if err != nil {
		return "", err
	}

	return cardinals[int((degrees+22.5)/45.0)%8], err
}

func (p *Packet) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("%s: ", p.Type))

	if p.Type == OBJECT {
		fmt.Fprintf(buffer, "%s (via %s)\n", strings.TrimSpace(p.ObjectOrItemName), p.SrcCallsign)
	} else {
		fmt.Fprintf(buffer, "%s\n", p.SrcCallsign)
	}

	if !p.Timestamp.IsZero() {
		fmt.Fprintf(buffer, "Time: %s\n", p.Timestamp)
	}

	if len(p.Path) > 0 {
		fmt.Fprintf(buffer, "Path: %q\n", p.Path)
	}

	if p.HasLocation() {
		fmt.Fprintf(buffer, "Pos: %f,%f\n", p.Latitude, p.Longitude)
	}

	fmt.Fprintf(buffer, "Speed: %.0fkm/h\n", p.Speed)

	if p.Comment != "" {
		fmt.Fprintf(buffer, "Comment: %s\n", strings.TrimSpace(p.Comment))
	}

	if p.Status != "" {
		fmt.Fprintf(buffer, "Status: %s\n", strings.TrimSpace(p.Status))
	}

	if p.Messagebits != "" {
		fmt.Fprintf(buffer, "Mic-E: %s\n", p.MicEMessage())
	}

	return buffer.String()
}
