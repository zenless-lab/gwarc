package cdx

import (
    "strings"
    "testing"
    "time"
)

func TestMarshal(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    []string
        wantErr bool
    }{
        {
            name: "basic CDX9 format",
            input: &CDXFile{
                Header: CDXHeader{
                    Format:    CDX9,
                    Delimiter: ' ',
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
            want: []string{
                "CDX N b a m s k r V g",
                "http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - 12345 example.warc.gz",
            },
            wantErr: false,
        },
        {
            name: "basic CDX11 format",
            input: &CDXFile{
                Header: CDXHeader{
                    Format:    CDX11,
                    Delimiter: ' ',
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
                        MetaTags:            "TagA",
                        CompressedSize:      123,
                        CompressedArcOffset: 12345,
                        Filename:            "example.warc.gz",
                    },
                },
            },
            want: []string{
                "CDX N b a m s k r M S V g",
                "http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - TagA 123 12345 example.warc.gz",
            },
            wantErr: false,
        },
        {
            name: "multiple records",
            input: &CDXFile{
                Header: CDXHeader{
                    Format:    CDX9,
                    Delimiter: ' ',
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
                    {
                        MassagedURL:         "http://example1.com/",
                        Date:                time.Date(2001, 4, 24, 21, 5, 51, 0, time.UTC),
                        OriginalURL:         "http://example1.com/",
                        MIMEType:            "text/plain",
                        StatusCode:          404,
                        NewChecksum:         "PRUJLL24DV7QYOOZMSA5TNJUKKRYAIM5",
                        Redirect:            "-",
                        CompressedArcOffset: 54321,
                        Filename:            "example1.warc.gz",
                    },
                },
            },
            want: []string{
                "CDX N b a m s k r V g",
                "http://example.com/ 20010424210312 http://example.com/ text/html 200 ZMSA5TNJUKKRYAIM5PRUJLL24DV7QYOO - 12345 example.warc.gz",
                "http://example1.com/ 20010424210551 http://example1.com/ text/plain 404 PRUJLL24DV7QYOOZMSA5TNJUKKRYAIM5 - 54321 example1.warc.gz",
            },
            wantErr: false,
        },
        {
            name:    "nil input",
            input:   nil,
            want:    nil,
            wantErr: true,
        },
        {
            name:    "non-struct input",
            input:   "invalid",
            want:    nil,
            wantErr: true,
        },
        {
            name:    "non-CDXFile input",
            input:   &struct{}{},
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Marshal(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr {
                gotStr := string(got)
                gotLines := strings.Split(strings.TrimSpace(gotStr), "\n")

                if len(gotLines) != len(tt.want) {
                    t.Errorf("Marshal() got %d lines, want %d", len(gotLines), len(tt.want))
                    return
                }

                for i, want := range tt.want {
                    if gotLines[i] != want {
                        t.Errorf("Marshal() line %d = %v, want %v", i, gotLines[i], want)
                    }
                }
            }
        })
    }
}