package fap_test

import (
	"fmt"

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
