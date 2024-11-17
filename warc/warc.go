package warc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type WARCVariant string

const (
	// WARCVariant1_0 represents the WARC 1.0 format
	WARCVariant1_0 WARCVariant = "1.0"
	// WARCVariant1_1 represents the WARC 1.1 format
	WARCVariant1_1 WARCVariant = "1.1"
)

// WARCRecordType represents the type of a WARC record.
// It defines the nature and intended use of the record content.
type WARCRecordType string

const (
	// WARCTypeWarcinfo indicates a record that describes the records that follow
	WARCTypeWarcinfo WARCRecordType = "warcinfo"
	// WARCTypeResponse represents a response to a request
	WARCTypeResponse WARCRecordType = "response"
	// WARCTypeResource represents a resource captured independently of protocol
	WARCTypeResource WARCRecordType = "resource"
	// WARCTypeRequest represents a documented web request
	WARCTypeRequest WARCRecordType = "request"
	// WARCTypeMetadata holds content created to further describe another record
	WARCTypeMetadata WARCRecordType = "metadata"
	// WARCTypeRevisit indicates a subsequent visit to content previously archived
	WARCTypeRevisit WARCRecordType = "revisit"
	// WARCTypeConversion contains alternative versions of other records' content
	WARCTypeConversion WARCRecordType = "conversion"
	// WARCTypeContinuation holds the continuation of another record that exceeded size limits
	WARCTypeContinuation WARCRecordType = "continuation"
)

// TruncatedReason represents the reason why a record's content block was truncated
type TruncatedReason string

const (
	// TruncatedLength indicates truncation due to exceeding a size limit
	TruncatedLength TruncatedReason = "length"

	// TruncatedTime indicates truncation due to exceeding a time limit
	TruncatedTime TruncatedReason = "time"

	// TruncatedDisconnect indicates truncation due to network disconnect
	TruncatedDisconnect TruncatedReason = "disconnect"

	// TruncatedUnspecified indicates truncation for unspecified reasons
	TruncatedUnspecified TruncatedReason = "unspecified"
)

// WARCRecord represents a single WARC (Web ARChive) record.
// It contains both the record headers and content according to the WARC format specification.
type WARCRecord struct {
	// WARC version
	Version WARCVariant

	// Required fields

	// RecordID is a globally unique identifier for the record
	RecordID string `warc:"WARC-Record-ID"`
	// ContentLength specifies the length of the record's content block in bytes
	ContentLength uint64 `warc:"Content-Length"`
	// Date represents the record creation time in UTC
	Date time.Time `warc:"WARC-Date"`
	// Type indicates the type of WARC record
	Type WARCRecordType `warc:"WARC-Type"`

	// Optional fields

	// ContentType specifies the MIME type of the content block
	ContentType string `warc:"Content-Type,omitempty"`
	// ConcurrentTo lists record IDs created from the same capture event
	ConcurrentTo []string `warc:"WARC-Concurrent-To,omitempty"`
	// BlockDigest contains the digest value of the entire content block
	BlockDigest string `warc:"WARC-Block-Digest,omitempty"`
	// PayloadDigest contains the digest value of the record's payload
	PayloadDigest string `warc:"WARC-Payload-Digest,omitempty"`
	// IPAddress records the IP address of the server providing the archived content
	IPAddress string `warc:"WARC-IP-Address,omitempty"`
	// RefersTo contains the ID of the record this record references
	RefersTo string `warc:"WARC-Refers-To,omitempty"`
	// RefersToTargetURI contains the original URI of the referenced record
	RefersToTargetURI string `warc:"WARC-Refers-To-Target-URI,omitempty"`
	// RefersToDate contains the date of the referenced record
	RefersToDate time.Time `warc:"WARC-Refers-To-Date,omitempty"`
	// TargetURI contains the original URI of the archived content
	TargetURI string `warc:"WARC-Target-URI,omitempty"`
	// Truncated indicates if and why the content block was truncated
	Truncated TruncatedReason `warc:"WARC-Truncated,omitempty"`
	// WarcinfoID references the warcinfo record describing this record's capture
	WarcinfoID string `warc:"WARC-Warcinfo-ID,omitempty"`
	// Filename contains the name of the WARC file (used in warcinfo records)
	Filename string `warc:"WARC-Filename,omitempty"`
	// Profile specifies the URI of the profile used for revisit records
	Profile string `warc:"WARC-Profile,omitempty"`
	// IdentifiedPayloadType contains the independently determined content type
	IdentifiedPayloadType string `warc:"WARC-Identified-Payload-Type,omitempty"`
	// SegmentNumber indicates the sequence number in a segmented record
	SegmentNumber int `warc:"WARC-Segment-Number,omitempty"`
	// SegmentOriginID references the first record in a segmented series
	SegmentOriginID string `warc:"WARC-Segment-Origin-ID,omitempty"`
	// SegmentTotalLength specifies the total length of all segments
	SegmentTotalLength uint64 `warc:"WARC-Segment-Total-Length,omitempty"`

	// Body

	// Content holds the actual content block of the WARC record
	Content []byte
}

type WARCRecordMarshaler interface {
	MarshalWARCRecord() ([]byte, error)
}

type WARCRecordUnmarshaler interface {
	UnmarshalWARCRecord([]byte) error
}

type warcInfoExtraRecord struct {
	// Operator contains contact information for the WARC creator
	Operator string `warc:"operator"`
	// Software identifies the software used to create the WARC
	Software string `warc:"software"`
	// Robots describes the robots policy followed during crawling
	Robots string `warc:"robots"`
	// Hostname identifies the machine that created the WARC
	Hostname string `warc:"hostname"`
	// IP contains the IP address of the machine that created the WARC
	IP string `warc:"ip"`
	// UserAgent contains the HTTP user-agent header used during crawling
	UserAgent string `warc:"http-header-user-agent"`
	// From contains the HTTP from header used during crawling
	From string `warc:"http-header-from"`
}

// WarcInfoRecord represents metadata about the WARC file itself.
// This information is stored in records of type "warcinfo".
type WarcInfoRecord struct {
	WARCRecord

	// Extra fields
	warcInfoExtraRecord
}

func (w *WarcInfoRecord) MarshalWARCRecord() ([]byte, error) {
	warcRecordResult, err := Marshal(w.WARCRecord)
	if err != nil {
		return nil, err
	}

	var extraRecordResult []byte
	if w.Operator != "" {
		extraRecordResult = append(extraRecordResult, []byte("operator: "+w.Operator+"\n")...)
	}
	if w.Software != "" {
		extraRecordResult = append(extraRecordResult, []byte("software: "+w.Software+"\n")...)
	}
	if w.Robots != "" {
		extraRecordResult = append(extraRecordResult, []byte("robots: "+w.Robots+"\n")...)
	}
	if w.Hostname != "" {
		extraRecordResult = append(extraRecordResult, []byte("hostname: "+w.Hostname+"\n")...)
	}
	if w.IP != "" {
		extraRecordResult = append(extraRecordResult, []byte("ip: "+w.IP+"\n")...)
	}
	if w.UserAgent != "" {
		extraRecordResult = append(extraRecordResult, []byte("http-header-user-agent: "+w.UserAgent+"\n")...)
	}
	if w.From != "" {
		extraRecordResult = append(extraRecordResult, []byte("http-header-from: "+w.From+"\n")...)
	}

	return append(warcRecordResult, extraRecordResult...), nil

}

func (w *WarcInfoRecord) UnmarshalWARCRecord(data []byte) (err error) {
	err = Unmarshal(data, &w.WARCRecord)
	if err != nil {
		return
	}

	content := w.WARCRecord.Content

	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		switch key {
		case "robots":
			w.Robots = value
		case "hostname":
			w.Hostname = value
		case "software":
			w.Software = value
		case "operator":
			w.Operator = value
		}
	}

	return
}

type metadataExtraRecord struct {
	// Via contains the URI where the archived URI was discovered
	Via string `warc:"via"`
	// HopsFromSeed describes the type of each hop from the seed URI to the current URI
	HopsFromSeed string `warc:"hopsFromSeed"`
	// FetchTimeMs indicates the time taken to collect the archived URI (in milliseconds)
	FetchTimeMs uint64 `warc:"fetchTimeMs"`
}

// MetadataRecord represents additional information about another record.
// This information is stored in records of type "metadata".
type MetadataRecord struct {
	WARCRecord

	// Extra fields
	metadataExtraRecord
}

func (m *MetadataRecord) MarshalWARCRecord() ([]byte, error) {
	warcRecordResult, err := Marshal(m.WARCRecord)
	if err != nil {
		return nil, err
	}

	var extraRecordResult []byte
	if m.Via != "" {
		extraRecordResult = append(extraRecordResult, []byte("via: "+m.Via+"\n")...)
	}
	if m.HopsFromSeed != "" {
		extraRecordResult = append(extraRecordResult, []byte("hopsFromSeed: "+m.HopsFromSeed+"\n")...)
	}
	if m.FetchTimeMs != 0 {
		extraRecordResult = append(extraRecordResult, []byte("fetchTimeMs: "+strconv.FormatUint(m.FetchTimeMs, 10)+"\n")...)
	}

	return append(warcRecordResult, extraRecordResult...), nil
}

func (m *MetadataRecord) UnmarshalWARCRecord(data []byte) (err error) {
	err = Unmarshal(data, &m.WARCRecord)
	if err != nil {
		return
	}

	content := m.WARCRecord.Content

	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		switch key {
		case "via":
			m.Via = value
		case "hopsFromSeed":
			m.HopsFromSeed = value
		case "fetchTimeMs":
			m.FetchTimeMs, err = strconv.ParseUint(value, 10, 64)
			if err != nil {
				return
			}
		}
	}

	return
}

type WARC struct {
	scanner *bufio.Scanner
}

func NewWARC(r io.Reader) *WARC {
	scanner := bufio.NewScanner(r)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Look for "WARC/x.x" pattern
		if i := bytes.Index(data, []byte("WARC/")); i >= 0 {
			// Return the data before "WARC/" if we're not at the start
			if i > 0 {
				return i, data[0:i], nil
			}
			// Find the end of this block (next "WARC/" or EOF)
			if j := bytes.Index(data[i+5:], []byte("WARC/")); j >= 0 {
				return i + j + 5, data[i : i+j+5], nil
			}
			// If we're at EOF, return the rest
			if atEOF {
				return len(data), data, nil
			}
		}

		// Request more data
		return 0, nil, nil
	}
	scanner.Split(split)

	return &WARC{
		scanner: scanner,
	}
}

func NewWARCFromFile(path string) (*WARC, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewWARC(file), nil
}

func NewWARCFromString(s string) *WARC {
	return NewWARC(bytes.NewBufferString(s))
}

func NewWARCFromBytes(b []byte) *WARC {
	return NewWARC(bytes.NewBuffer(b))
}

func (w *WARC) NextChunk() (*[]byte, error) {
	if !w.scanner.Scan() {
		return nil, io.EOF
	}
	chunk := w.scanner.Bytes()
	return &chunk, nil
}

func (w *WARC) Next() (record any, kind WARCRecordType, err error) {
	chunk, err := w.NextChunk()

	if err != nil {
		return
	}

	var loadedRecord WARCRecord
	if err = Unmarshal(*chunk, &loadedRecord); err != nil {
		return
	}

	if loadedRecord.Type == "metadata" {
		var metadata MetadataRecord
		if err = Unmarshal(*chunk, &metadata); err != nil {
			return
		}
		record = metadata
		kind = WARCTypeMetadata
		return
	} else if loadedRecord.Type == "warcinfo" {
		var warcinfo WarcInfoRecord
		if err = Unmarshal(*chunk, &warcinfo); err != nil {
			return
		}
		record = warcinfo
		kind = WARCTypeWarcinfo
		return
	} else {
		record = loadedRecord
		kind = loadedRecord.Type
		return
	}
}

// Validate checks if all required fields are present and valid
func (w *WARCRecord) Validate() error {
	if w.Version == "" {
		return fmt.Errorf("WARC version is required")
	}
	if w.RecordID == "" {
		return fmt.Errorf("WARC-Record-ID is required")
	}
	if w.Date.IsZero() {
		return fmt.Errorf("WARC-Date is required")
	}
	if w.Type == "" {
		return fmt.Errorf("WARC-Type is required")
	}
	return nil
}

// Validate checks if all required fields are present and valid
func (w *WarcInfoRecord) Validate() error {
	if w.WARCRecord.Validate() != nil {
		return fmt.Errorf("WARC record is invalid")
	}
	if w.Operator == "" {
		return fmt.Errorf("operator is required")
	}
	if w.Software == "" {
		return fmt.Errorf("software is required")
	}
	if w.Robots == "" {
		return fmt.Errorf("robots is required")
	}
	if w.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if w.IP == "" {
		return fmt.Errorf("IP is required")
	}
	if w.UserAgent == "" {
		return fmt.Errorf("http-header-user-agent is required")
	}
	if w.From == "" {
		return fmt.Errorf("http-header-from is required")
	}
	return nil
}

// Validate checks if all required fields are present and valid
func (m *MetadataRecord) Validate() error {
	if m.WARCRecord.Validate() != nil {
		return fmt.Errorf("WARC record is invalid")
	}
	if m.Via == "" {
		return fmt.Errorf("via is required")
	}
	if m.HopsFromSeed == "" {
		return fmt.Errorf("hopsFromSeed is required")
	}
	if m.FetchTimeMs == 0 {
		return fmt.Errorf("fetchTimeMs is required")
	}
	return nil
}
