package fap_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/martinhpedersen/libfap-go"
)

func ExampleParseAprs() {
	packet, _ := fap.ParseAprs("LA5NTA-9>V0QRR9,LD5BE*,WIDE2-2,qAR,LA2VSA-1:`{6qnfR>/]\"4M}WLNK-1=", false)
	fmt.Println(packet)
	// Output:
	//Location: LA5NTA-9
	//Path: ["LD5BE*" "WIDE2-2" "qAR" "LA2VSA-1"]
	//Pos: 60.204833,5.447500
	//Speed: 50km/h
	//Comment: ]WLNK-1=
	//Mic-E: in service
}

func TestTimestamp(t *testing.T) {
	// The APRS timestamp is parsed relative to runtime date, so we construct a packet with a timestamp for today.
	rawPacket := fmt.Sprintf("LA5NTA-9>V0QRR9,LD5BE*,WIDE2-2,qAR,LA2VSA-1:/%02d2045z4903.50N/07201.75W>Test1234", time.Now().Day())

	packet, _ := fap.ParseAprs(rawPacket, false)
	if packet.RawTimestamp == "" {
		t.Error("Expected non-empty raw timestamp.")
	}

	if packet.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp.")
	}

	packet.Timestamp = packet.Timestamp.UTC()

	if h := packet.Timestamp.Hour(); h != 20 {
		t.Errorf("Unexpected timestamp hour: %d.", h)
	}

	if m := packet.Timestamp.Minute(); m != 45 {
		t.Errorf("Unexpected timestamp minute: %d.", m)
	}

	if time.Now().UTC().Month() != packet.Timestamp.Month() {
		t.Errorf("Unexpected month, got '%s' expected '%s'.", packet.Timestamp.Month(), time.Now().UTC().Month())
	}

	if time.Now().UTC().Year() != packet.Timestamp.Year() {
		t.Errorf("Unexpected year, got '%d' expected '%d'.", packet.Timestamp.Year(), time.Now().UTC().Year())
	}
}
