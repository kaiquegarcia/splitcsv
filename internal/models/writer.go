package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type (
	Writer interface {
		Split()
	}

	writer struct {
		inputFile         *os.File
		currentOutputFile *os.File
		csvWriter         *csv.Writer
		fileCounter       int
		input             *Input
		headers           []string
	}
)

func NewWriter(input *Input) Writer {
	return &writer{
		fileCounter: 1,
		input:       input,
	}
}

func (w *writer) close() {
	if w.inputFile != nil {
		if err := w.inputFile.Close(); err != nil {
			fmt.Printf("could not close input file stream: %s\n", err.Error())
		}
	}

	if w.currentOutputFile != nil {
		if err := w.currentOutputFile.Close(); err != nil {
			fmt.Printf("could not close current output file stream: %s\n", err.Error())
		}
	}
}

func (w *writer) createOutputFile() bool {
	if w.currentOutputFile != nil {
		w.currentOutputFile.Close()
	}

	outputFile := fmt.Sprintf(w.input.OutputFilePathPattern, w.fileCounter)
	fmt.Printf("creating new file: %s\n", outputFile)

	if f, err := os.Create(outputFile); err != nil {
		fmt.Printf("could not create new file: %s\n", err.Error())
		return false
	} else {
		w.currentOutputFile = f
	}

	w.csvWriter = csv.NewWriter(w.currentOutputFile)
	if err := w.csvWriter.Write(w.headers); err != nil {
		fmt.Printf("could not write headers: %s\n", err.Error())
		return false
	}

	return true
}

func (w *writer) Split() {
	stream, err := os.Open(w.input.InputFilePath)
	if err != nil {
		fmt.Printf("could not open input file: %s\n", err.Error())
		return
	}

	w.inputFile = stream
	defer w.close()

	reader := csv.NewReader(stream)
	headers, err := reader.Read()
	if err != nil {
		fmt.Printf("could not read input CSV headers: %s\n", err.Error())
		return
	}

	w.headers = headers
	if success := w.createOutputFile(); !success {
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("could not read CSV record: %s\n", err.Error())
			return
		}

		if err := w.csvWriter.Write(record); err != nil {
			fmt.Printf("could not write CSV record: %s\n", err.Error())
			return
		}

		w.csvWriter.Flush()

		info, err := w.currentOutputFile.Stat()
		if err != nil {
			fmt.Printf("could not get file info: %s\n", err.Error())
			return
		}

		if info.Size() >= int64(w.input.MaxSize) {
			w.fileCounter++
			if success := w.createOutputFile(); !success {
				return
			}
		}
	}
}
