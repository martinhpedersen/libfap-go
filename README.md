# fap - A go wrapper for libfap - APRS parser

## Install

go get github.com/martinhpedersen/libfap-go

## libfap

Libfap is a C port of the Ham::APRS::FAP Finnish APRS Parser perl
module. As the original Perl code, libfap parses normal, mic-e and
compressed location packets, NMEA location packets, objects, items,
messages, telemetry and most weather packets. More information on
HAM::APRS::FAP is available at <http://search.cpan.org/dist/Ham-APRS-FAP/>

## Requirements

* latest libfap (<http://www.pakettiradio.net/libfap/>)

## Usage
Simple program that decodes the APRS-IS feed (omitting errors):

	func main() {
		defer fap.Cleanup()

		conn, _ := net.Dial("tcp", "rotate.aprs.net:23")
		fmt.Fprintf(conn, "user N0CALL pass -1 vers goAPRS 0.00\r\n")

		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}

			packet, _ := fap.ParseAprs(line, false)
			fmt.Println(packet)
		}
	}

## Missing features / Things you can help with

* wx_report
* telemetry
* capabilities
* tests
* ++

## Copyright and disclaimer

The parser (Ham::APRS::FAP) has been originally written
by Tapio Sokura, OH2KKU and Heikki Hannikainen, OH7JZB. It has
been ported to C by Tapio Aaltonen, OH2GVE and wrapped in go
by Martin Hebnes Pedersen, LA5NTA.

Libfap is free software; for more details see <http://www.pakettiradio.net/libfap/>.

libfap-go is also free software; 

Copyright (c) 2013, Contributors and Martin Hebnes Pedersen
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

