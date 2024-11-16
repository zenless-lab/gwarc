package gwarc

import (
	"fmt"
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

// WarcInfoRecord represents metadata about the WARC file itself.
// This information is stored in records of type "warcinfo".
type WarcInfoRecord struct {
	// Version indicates the WARC format version
	Version WARCVariant

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

// MetadataRecord represents additional information about another record.
// This information is stored in records of type "metadata".
type MetadataRecord struct {
	// Version indicates the WARC format version
	Version WARCVariant

	// Via contains the URI where the archived URI was discovered
	Via string `warc:"via"`
	// HopsFromSeed describes the type of each hop from the seed URI to the current URI
	HopsFromSeed string `warc:"hopsFromSeed"`
	// FetchTimeMs indicates the time taken to collect the archived URI (in milliseconds)
	FetchTimeMs uint64 `warc:"fetchTimeMs"`
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
	if w.Version == "" {
		return fmt.Errorf("WARC version is required")
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
	if m.Version == "" {
		return fmt.Errorf("WARC version is required")
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
