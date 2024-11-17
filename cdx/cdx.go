package cdx

import (
	"strings"
	"time"
)

// CDXField if a single character that represents a field in a CDX record.
type CDXField rune

const (
	// URL

	FieldCanonizedURL      CDXField = 'A' // canonized url
	FieldNewsGroup         CDXField = 'B' // news group
	FieldRulespaceCategory CDXField = 'C' // rulespace category(future use)
	FieldCanonizedFrame    CDXField = 'F' // canonized frame
	FieldCanonizedHost     CDXField = 'H' // canonized host
	FieldCanonizedImage    CDXField = 'I' // canonized image
	FieldCanonizedJump     CDXField = 'J' // canonized jump point
	FieldCanonizedLink     CDXField = 'L' // canonized link
	FieldMassagedURL       CDXField = 'N' // massaged url
	FieldCanonizedPath     CDXField = 'P' // canonized path
	FieldCanonizedRedirect CDXField = 'R' // canonized redirect
	FieldOriginalURL       CDXField = 'a' // original url (in alexa-made dat file meta-data line)
	FieldOriginalHost      CDXField = 'h' // original host
	FieldOriginalPath      CDXField = 'p' // original path

	// Metadata

	FieldDate     CDXField = 'b' // date (in alexa-made dat file meta-data line)
	FieldIP       CDXField = 'e' // IP (in alexa-made dat file meta-data line)
	FieldLanguage CDXField = 'Q' // language string
	FieldPort     CDXField = 'o' // port
	FieldTitle    CDXField = 't' // title (in alexa-made dat file)
	FieldMetaTags CDXField = 'M' // meta tags (AIF) (in alexa-made dat file)

	// Checksum and size

	FieldOldChecksum    CDXField = 'c' // old style checksum (in alexa-made dat file)
	FieldNewChecksum    CDXField = 'k' // new style checksum (in alexa-made dat file)
	FieldCompressedSize CDXField = 'S' // compressed record size
	FieldArcDocLength   CDXField = 'n' // arc document length (in alexa-made dat file)

	// Offset information

	FieldCompressedDatOffset   CDXField = 'D' // compressed dat file offset
	FieldUncompressedDatOffset CDXField = 'd' // uncompressed dat file offset
	FieldCompressedArcOffset   CDXField = 'V' // compressed arc file offset
	FieldUncompressedArcOffset CDXField = 'v' // uncompressed arc file offset (in alexa-made dat file)

	// Resource reference

	FieldFrame        CDXField = 'f' // frame (in alexa-made dat file)
	FieldImage        CDXField = 'i' // image (in alexa-made dat file)
	FieldJumpPoint    CDXField = 'j' // original jump point
	FieldLink         CDXField = 'l' // link (in alexa-made dat file)
	FieldURLsInHref   CDXField = 'x' // url in other href tags (in alexa-made dat file)
	FieldURLsInSrc    CDXField = 'y' // url in other src tags (in alexa-made dat file)
	FieldURLsInScript CDXField = 'z' // url found in script (in alexa-made dat file)

	// Response information

	FieldMIMEType   CDXField = 'm' // mime type of original document (in alexa-made dat file)
	FieldStatusCode CDXField = 's' // response code (in alexa-made dat file)
	FieldRedirect   CDXField = 'r' // redirect (in alexa-made dat file)

	// Other

	FieldFilename     CDXField = 'g' // file name
	FieldFBIS         CDXField = 'K' // Some weird FBIS what's changed
	FieldUniqueness   CDXField = 'U' // uniqueness (future use)
	FieldLanguageDesc CDXField = 'G' // multi-column language description (in alexa-made dat files, soon)
)

// CDXFormat is a list of fields that make up a CDX record.
type CDXFormat []CDXField

var (
	// Default CDX formats

	CDX9 = CDXFormat{
		FieldMassagedURL,         // N
		FieldDate,                // b
		FieldOriginalURL,         // a
		FieldMIMEType,            // m
		FieldStatusCode,          // s
		FieldNewChecksum,         // k
		FieldRedirect,            // r
		FieldCompressedArcOffset, // V
		FieldFilename,            // g
	}

	CDX11 = CDXFormat{
		FieldMassagedURL,         // N
		FieldDate,                // b
		FieldOriginalURL,         // a
		FieldMIMEType,            // m
		FieldStatusCode,          // s
		FieldNewChecksum,         // k
		FieldRedirect,            // r
		FieldMetaTags,            // M
		FieldCompressedSize,      // S
		FieldCompressedArcOffset, // V
		FieldFilename,            // g
	}
)

// String returns a string representation of the CDX format
func (f CDXFormat) String() string {
	fields := make([]string, len(f)+1)
	fields[0] = "CDX"
	for i, field := range f {
		fields[i+1] = string(field)
	}
	return strings.Join(fields, " ")
}

// CDXRecord is a single record in a CDX file.
type CDXRecord struct {
	// URL field group

	CanonizedURL      string `cdx:"A" json:"canonized_url"`      // A canonized url
	NewsGroup         string `cdx:"B" json:"news_group"`         // B news group
	RulespaceCategory string `cdx:"C" json:"rulespace_category"` // C rulespace category
	CanonizedFrame    string `cdx:"F" json:"canonized_frame"`    // F canonized frame
	CanonizedHost     string `cdx:"H" json:"canonized_host"`     // H canonized host
	CanonizedImage    string `cdx:"I" json:"canonized_image"`    // I canonized image
	CanonizedJump     string `cdx:"J" json:"canonized_jump"`     // J canonized jump point
	CanonizedLink     string `cdx:"L" json:"canonized_link"`     // L canonized link
	MassagedURL       string `cdx:"N" json:"massaged_url"`       // N massaged url
	CanonizedPath     string `cdx:"P" json:"canonized_path"`     // P canonized path
	CanonizedRedirect string `cdx:"R" json:"canonized_redirect"` // R canonized redirect
	OriginalURL       string `cdx:"a" json:"original_url"`       // a original url
	OriginalHost      string `cdx:"h" json:"original_host"`      // h original host
	OriginalPath      string `cdx:"p" json:"original_path"`      // p original path

	// Metadata field group

	Date     time.Time `cdx:"b" json:"date"`      // b date
	IP       string    `cdx:"e" json:"ip"`        // e IP
	Language string    `cdx:"Q" json:"language"`  // Q language string
	Port     int       `cdx:"o" json:"port"`      // o port
	Title    string    `cdx:"t" json:"title"`     // t title
	MetaTags string    `cdx:"M" json:"meta_tags"` // M meta tags (AIF)

	// Checksum and size field group

	OldChecksum       string `cdx:"c" json:"old_checksum"`        // c old style checksum
	NewChecksum       string `cdx:"k" json:"new_checksum"`        // k new style checksum
	CompressedSize    int64  `cdx:"S" json:"compressed_size"`     // S compressed record size
	ArcDocumentLength int64  `cdx:"n" json:"arc_document_length"` // n arc document length

	// Offset information field group

	CompressedDatOffset   int64 `cdx:"D" json:"compressed_dat_offset"`   // D compressed dat file offset
	UncompressedDatOffset int64 `cdx:"d" json:"uncompressed_dat_offset"` // d uncompressed dat file offset
	CompressedArcOffset   int64 `cdx:"V" json:"compressed_arc_offset"`   // V compressed arc file offset
	UncompressedArcOffset int64 `cdx:"v" json:"uncompressed_arc_offset"` // v uncompressed arc file offset

	// Resource reference field group

	Frame        string `cdx:"f" json:"frame"`          // f frame
	Image        string `cdx:"i" json:"image"`          // i image
	JumpPoint    string `cdx:"j" json:"jump_point"`     // j original jump point
	Link         string `cdx:"l" json:"link"`           // l link
	URLsInHref   string `cdx:"x" json:"urls_in_href"`   // x url in other href tags
	URLsInSrc    string `cdx:"y" json:"urls_in_src"`    // y url in other src tags
	URLsInScript string `cdx:"z" json:"urls_in_script"` // z url found in script

	// Response information field group

	MIMEType   string `cdx:"m" json:"mime_type"`   // m mime type of original document
	StatusCode int    `cdx:"s" json:"status_code"` // s response code
	Redirect   string `cdx:"r" json:"redirect"`    // r redirect

	// Other field group

	Filename   string `cdx:"g" json:"filename"`   // g file name
	FBIS       string `cdx:"K" json:"fbis"`       // K Some weird FBIS what's changed
	Uniqueness string `cdx:"U" json:"uniqueness"` // U uniqueness
}

// CDXHeader contains the format and field definitions for a CDX file.
type CDXHeader struct {
	Format    CDXFormat
	Fields    []CDXField
	Delimiter rune
}

// CDXFile is a collection of CDX records
type CDXFile struct {
	Header  CDXHeader
	Records []CDXRecord
}

// NewCDXFile initializes a new CDX file with the given format
func NewCDXFile(format CDXFormat) *CDXFile {
	return &CDXFile{
		Header: CDXHeader{
			Format:    format,
			Delimiter: ' ',
			Fields:    parseFormat(format),
		},
		Records: make([]CDXRecord, 0),
	}
}

// parseFormat parses a CDX format string into a list of fields
func parseFormat(format CDXFormat) []CDXField {
	parts := strings.Fields(string(format))
	if len(parts) > 1 {
		fields := make([]CDXField, len(parts)-1)
		for i, part := range parts[1:] {
			fields[i] = CDXField(part[0])
		}
		return fields
	}
	return nil
}
