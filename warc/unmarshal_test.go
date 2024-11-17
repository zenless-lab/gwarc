package warc_test

import (
	"testing"
	"time"

	. "github.com/zenless-lab/gwarc/warc"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    WARCRecord
		wantErr bool
	}{
		{
			name: "basic record",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: 13
Content-Type: text/plain

Hello, World!`),
			want: WARCRecord{
				Version:       WARCVariant1_0,
				Type:          WARCTypeResponse,
				Date:          time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				RecordID:      "<urn:uuid:12345678-1234-1234-1234-123456789012>",
				ContentType:   "text/plain",
				ContentLength: 13,
				Content:       []byte("Hello, World!"),
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			input: []byte(`WARC/2.0
WARC-Type: response
Content-Length: 0

`),
			wantErr: true,
		},
		{
			name: "invalid header format",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: none
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: 0
WARC-IP-Address: 192.168.0.1
WARC-Target-URI: http://example.com
`),
			wantErr: true,
		},
		{
			name: "with optional fields",
			input: []byte(`WARC/1.1
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: 0
WARC-IP-Address: 192.168.1.1
WARC-Target-URI: http://example.com
Content-Type: text/html

`),
			want: WARCRecord{
				Version:     WARCVariant1_1,
				Type:        WARCTypeResponse,
				Date:        time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				RecordID:    "<urn:uuid:12345678-1234-1234-1234-123456789012>",
				IPAddress:   "192.168.1.1",
				TargetURI:   "http://example.com",
				ContentType: "text/html",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got WARCRecord
			err := Unmarshal(tt.input, &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Version != tt.want.Version {
					t.Errorf("Version = %v, want %v", got.Version, tt.want.Version)
				}
				if got.Type != tt.want.Type {
					t.Errorf("Type = %v, want %v", got.Type, tt.want.Type)
				}
				if !got.Date.Equal(tt.want.Date) {
					t.Errorf("Date = %v, want %v", got.Date, tt.want.Date)
				}
				if got.RecordID != tt.want.RecordID {
					t.Errorf("RecordID = %v, want %v", got.RecordID, tt.want.RecordID)
				}
				if got.ContentType != tt.want.ContentType {
					t.Errorf("ContentType = %v, want %v", got.ContentType, tt.want.ContentType)
				}
				if got.IPAddress != tt.want.IPAddress {
					t.Errorf("IPAddress = %v, want %v", got.IPAddress, tt.want.IPAddress)
				}
				if got.TargetURI != tt.want.TargetURI {
					t.Errorf("TargetURI = %v, want %v", got.TargetURI, tt.want.TargetURI)
				}
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name: "valid record",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: 13
Content-Type: text/plain

Hello, World!`),
			wantErr: false,
		},
		{
			name: "unsupported version",
			input: []byte(`WARC/2.0
WARC-Type: response
Content-Length: 0

`),
			wantErr: true,
		},
		{
			name: "missing Content-Length",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Type: text/plain

Hello, World!`),
			wantErr: true,
		},
		{
			name: "invalid Content-Length value",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: abc
Content-Type: text/plain

Hello, World!`),
			wantErr: true,
		},
		{
			name: "invalid header format",
			input: []byte(`WARC/1.0
WARC-Type: response
WARC-Date: 2024-01-01T10:00:00Z
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length 13
Content-Type: text/plain

Hello, World!`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Valid(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
