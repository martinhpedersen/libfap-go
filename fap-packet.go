package fap

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

// FapPacket is the APRS packet type
type FapPacket struct {
	PacketType uint // See const
	OrigPacket string

	Header string
	Body   string

	SrcCallsign string
	DstCallsign string
	Path        []string

	Latitude  float64
	Longitude float64
	Format    uint // See const

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
func (p *FapPacket) HasLocation() bool {
	return (p.Latitude != 0 && p.Longitude != 0)
}

// MicEMessage returns the texual Mic-E message.
func (p *FapPacket) MicEMessage() string {
	if p.Messagebits == "" {
		return ""
	}

	return MicEMbitsToMessage(p.Messagebits)
}

// Distance returns the distance to the given packet b in km.
func (a *FapPacket) Distance(b *FapPacket) (float64, error) {
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

func (a *FapPacket) Direction(b *FapPacket) (float64, error) {
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

func (a *FapPacket) HumanReadableDirection(b *FapPacket) (string, error) {
	degrees, err := a.Direction(b)
	if err != nil {
		return "", err
	}

	return cardinals[int((degrees+22.5)/45.0)%8], err
}

func (p *FapPacket) String() string {
	buffer := bytes.NewBufferString("")

	if p.PacketType == OBJECT {
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
