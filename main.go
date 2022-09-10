package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type Root struct {
	XMLName  xml.Name   `xml:"fluxASS"`
	Attrs    []xml.Attr `xml:",any,attr"`
	Accounts []Account  `xml:"compte"`
}

type Account struct {
	XMLName xml.Name   `xml:"compte"`
	Attrs   []xml.Attr `xml:",any,attr"`
	Data    []byte     `xml:",innerxml"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	if len(args) != 2 {
		return fmt.Errorf("program must take <threshold> and <file-path> as parameter")
	}

	threshold, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse threshold to a valid integer: %v", err)
	}

	fpath := args[1]
	file, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("failed to read XML file: %v", err)
	}
	defer file.Close()

	fmt.Println("parsing XML file...")

	var inputXML Root
	if err = xml.NewDecoder(file).Decode(&inputXML); err != nil {
		return fmt.Errorf("failed to parse XML input file: %v", err)
	}

	fmt.Println("generating output files...")

	fpart := 1
	var accounts = make([]Account, 0, threshold)
	for _, account := range inputXML.Accounts {
		accounts = append(accounts, account)
		if len(accounts) == cap(accounts) {
			err = writeXML(Root{
				XMLName:  inputXML.XMLName,
				Attrs:    inputXML.Attrs,
				Accounts: accounts,
			}, generateFilename(fpath, fpart))

			if err != nil {
				return err
			}

			accounts = make([]Account, 0, threshold)
			fpart += 1
		}
	}

	if len(accounts) > 0 {
		err = writeXML(Root{
			XMLName:  inputXML.XMLName,
			Attrs:    inputXML.Attrs,
			Accounts: accounts,
		}, generateFilename(fpath, fpart))
		if err != nil {
			return err
		}
	}

	return nil
}

func writeXML(r Root, fname string) error {
	if err := os.MkdirAll(path.Dir(fname), 0755); err != nil {
		return fmt.Errorf("cannot create subfolder: %v", err)
	}

	output, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot create output file: %v", err)
	}
	defer output.Close()

	encoder := xml.NewEncoder(output)
	encoder.Indent("", "  ")
	if err := encoder.Encode(r); err != nil {
		return fmt.Errorf("cannot generate XML file: %v", err)
	}

	fmt.Printf("file %s generated\n", fname)

	return nil
}

func generateFilename(fpath string, fpart int) string {
	folder := strings.TrimSuffix(fpath, path.Ext(fpath))
	return fmt.Sprintf("%s/part-%d.xml", folder, fpart)
}
