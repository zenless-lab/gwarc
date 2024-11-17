package cdx

import (
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *CDXFile
		wantErr bool
	}{
		{
			name: "basic CDX9 format",
			input: `CDX N b a m s k r V g
http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - 12345 example.warc.gz`,
			want: &CDXFile{
				Header: CDXHeader{
					Format:    CDX9,
					Delimiter: ' ',
					Fields:    parseFormat(CDX9),
				},
				Records: []CDXRecord{
					{
						MassagedURL:         "http://example.com/",
						Date:                time.Date(2001, 4, 24, 21, 3, 12, 0, time.UTC),
						OriginalURL:         "http://example.com/",
						MIMEType:            "text/html",
						StatusCode:          200,
						NewChecksum:         "ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO",
						Redirect:            "-",
						CompressedArcOffset: 12345,
						Filename:            "example.warc.gz",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "basic CDX11 format",
			input: `CDX N b a m s k r M S V g
http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - TagA 123 12345 example.warc.gz`,
			want: &CDXFile{
				Header: CDXHeader{
					Format:    CDX11,
					Delimiter: ' ',
					Fields:    parseFormat(CDX11),
				},
				Records: []CDXRecord{
					{
						MassagedURL:         "http://example.com/",
						Date:                time.Date(2001, 4, 24, 21, 3, 12, 0, time.UTC),
						OriginalURL:         "http://example.com/",
						MIMEType:            "text/html",
						StatusCode:          200,
						Redirect:            "-",
						CompressedArcOffset: 12345,
						Filename:            "example.warc.gz",
						NewChecksum:         "ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO",
						MetaTags:            "TagA",
						CompressedSize:      123,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple records",
			input: `CDX N b a m s k r V g
http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - 12345 example.warc.gz
http://example1.com/ 20010424210551 http://example1.com/ text/plain 404 PRUJLL24DV7QYOOZMSA5TNJUKKRYAIM5 - 54321 example1.warc.gz`,
			want: &CDXFile{
				Header: CDXHeader{
					Format:    CDX9,
					Delimiter: ' ',
					Fields:    parseFormat(CDX9),
				},
				Records: []CDXRecord{
					{
						MassagedURL:         "http://example.com/",
						Date:                time.Date(2001, 4, 24, 21, 3, 12, 0, time.UTC),
						OriginalURL:         "http://example.com/",
						MIMEType:            "text/html",
						StatusCode:          200,
						Redirect:            "-",
						CompressedArcOffset: 12345,
						Filename:            "example.warc.gz",
						NewChecksum:         "ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO",
					},
					{
						MassagedURL:         "http://example1.com/",
						Date:                time.Date(2001, 4, 24, 21, 5, 51, 0, time.UTC),
						OriginalURL:         "http://example1.com/",
						MIMEType:            "text/plain",
						StatusCode:          404,
						Redirect:            "-",
						CompressedArcOffset: 54321,
						Filename:            "example1.warc.gz",
						NewChecksum:         "PRUJLL24DV7QYOOZMSA5TNJUKKRYAIM5",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid header",
			input:   "INVALID N b a m s k r V g",
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid field count",
			input: `CDX N b a m s k r V g
http://example.com/ 20010424210312 http://example.com/`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got CDXFile
			err := Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.Records) != len(tt.want.Records) {
					t.Errorf("Unmarshal() got %d records, want %d", len(got.Records), len(tt.want.Records))
					return
				}

				// Compare first record
				gotRecord := got.Records[0]
				wantRecord := tt.want.Records[0]

				if gotRecord.MassagedURL != wantRecord.MassagedURL {
					t.Errorf("Record.MassagedURL = %v, want %v", gotRecord.MassagedURL, wantRecord.MassagedURL)
				}
				if !gotRecord.Date.Equal(wantRecord.Date) {
					t.Errorf("Record.Date = %v, want %v", gotRecord.Date, wantRecord.Date)
				}
				if gotRecord.StatusCode != wantRecord.StatusCode {
					t.Errorf("Record.StatusCode = %v, want %v", gotRecord.StatusCode, wantRecord.StatusCode)
				}
				if gotRecord.CompressedArcOffset != wantRecord.CompressedArcOffset {
					t.Errorf("Record.CompressedArcOffset = %v, want %v", gotRecord.CompressedArcOffset, wantRecord.CompressedArcOffset)
				}
				if gotRecord.Filename != wantRecord.Filename {
					t.Errorf("Record.Filename = %v, want %v", gotRecord.Filename, wantRecord.Filename)
				}
				if gotRecord.NewChecksum != wantRecord.NewChecksum {
					t.Errorf("Record.NewChecksum = %v, want %v", gotRecord.NewChecksum, wantRecord.NewChecksum)
				}
				if gotRecord.Redirect != wantRecord.Redirect {
					t.Errorf("Record.Redirect = %v, want %v", gotRecord.Redirect, wantRecord.Redirect)
				}
				if gotRecord.MIMEType != wantRecord.MIMEType {
					t.Errorf("Record.MIMEType = %v, want %v", gotRecord.MIMEType, wantRecord.MIMEType)
				}
				if gotRecord.MetaTags != wantRecord.MetaTags {
					t.Errorf("Record.MetaTags = %v, want %v", gotRecord.MetaTags, wantRecord.MetaTags)
				}
				if gotRecord.CompressedSize != wantRecord.CompressedSize {
					t.Errorf("Record.CompressedSize = %v, want %v", gotRecord.CompressedSize, wantRecord.CompressedSize)
				}
			}
		})
	}
}
