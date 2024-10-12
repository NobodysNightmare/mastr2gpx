package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"nur-jan.de/go/mastr2gpx/xmlstream"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Gpx struct {
	Metadata  GpxMetadata
	Waypoints []GpxWaypoint
}

type GpxMetadata struct {
	Name        string `xml:"name,omitempty"`
	Description string `xml:"desc,omitempty"`
}

type GpxWaypoint struct {
	Name        string  `xml:"name,omitempty"`
	Description string  `xml:"desc,omitempty"`
	Lat         float64 `xml:"lat,attr"`
	Lon         float64 `xml:"lon,attr"`
}

type EinheitSolar struct {
	ID         string  `xml:"EinheitMastrNummer"`
	Name       string  `xml:"NameStromerzeugungseinheit"`
	PostalCode string  `xml:"Postleitzahl"`
	NetPower   float64 `xml:"Nettonennleistung"`
	Lat        float64 `xml:"Breitengrad"`
	Lng        float64 `xml:"Laengengrad"`
	Modules    int     `xml:"AnzahlModule"`
}

func main() {
	directory := "/home/jan/Downloads/Gesamtdatenexport_20241012_24.2/"
	outputFile := "generators.gpx"
	postalCode := "12345"

	generators, err := findAllGenerators(directory, postalCode)
	if err != nil {
		fmt.Println("Could not get all generators:", err)
		os.Exit(1)
	}

	gpx := Gpx{Metadata: GpxMetadata{Name: "Generator List", Description: "A list of generator waypoints, extracted from a Marktstammdatenregister data export."}}

	for _, gen := range generators {
		name := gen.Name
		if name == "" {
			name = "Generator"
		}
		gpx.Waypoints = append(gpx.Waypoints, GpxWaypoint{
			Name:        fmt.Sprintf("%s (%s)", name, gen.ID),
			Description: fmt.Sprintf("Net power: %f kW (%d modules)", gen.NetPower, gen.Modules),
			Lat:         gen.Lat,
			Lon:         gen.Lng,
		})
	}

	gpxData, err := xml.Marshal(gpx)
	if err != nil {
		fmt.Println("Could not generate GPX:", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, gpxData, 0777)
	if err != nil {
		fmt.Println("Could not write GPX file:", err)
		os.Exit(1)
	}

	fmt.Println("GPX-Export finished! Wrote", len(gpx.Waypoints), "waypoints to the file.")
}

func findAllGenerators(directory, postalCode string) ([]EinheitSolar, error) {
	fileEntries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("Could not read directory: %w", err)
	}

	generators := []EinheitSolar{}
	for _, fileEntry := range fileEntries {
		if fileEntry.IsDir() || !strings.HasPrefix(fileEntry.Name(), "EinheitenSolar_") {
			continue
		}

		fmt.Println("Reading file", fileEntry.Name())

		f, err := os.Open(fmt.Sprintf("%s/%s", directory, fileEntry.Name()))
		if err != nil {
			return nil, fmt.Errorf("Could not open file: %w", err)
		}

		defer f.Close()

		newGenerators, err := findGenerators(f, postalCode)
		if err != nil {
			return nil, fmt.Errorf("Could not read generators: %w", err)
		}

		generators = append(generators, newGenerators...)
	}

	return generators, nil
}

func findGenerators(f *os.File, postalCode string) ([]EinheitSolar, error) {
	utf16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(utf16be.NewDecoder())
	unicodeReader := transform.NewReader(f, utf16bom)

	scanner := xmlstream.NewScanner(unicodeReader, new(EinheitSolar))
	result := []EinheitSolar{}

	for scanner.Scan() {
		tag := scanner.Element()
		switch el := tag.(type) {
		case *EinheitSolar:
			generator := *el
			if generator.Lat == 0 && generator.Lng == 0 {
				continue
			}
			if generator.PostalCode != postalCode {
				continue
			}

			result = append(result, generator)
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return result, nil
}

func (gpx Gpx) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	metaElement := xml.StartElement{Name: xml.Name{Local: "metadata"}}
	waypointElement := xml.StartElement{Name: xml.Name{Local: "wpt"}}

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "gpx", Space: "http://www.topografix.com/GPX/1/1"}, Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "version"}, Value: "1.1"}}})
	e.EncodeElement(gpx.Metadata, metaElement)
	for _, wp := range gpx.Waypoints {
		e.EncodeElement(wp, waypointElement)
	}
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "gpx", Space: "http://www.topografix.com/GPX/1/1"}})
	return nil
}
