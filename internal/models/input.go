package models

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFilePathID              = "--name="
	outputFilePathPatternID      = "--out="
	maxSizeID                    = "--size="
	defaultInputFilePath         = "./input.csv"
	defaultOutputFilePathPattern = "./output_%d.csv"
	defaultMaxSize               = 10 * GB
	stringDelimiter              = "\""
)

var regexNotLetters = regexp.MustCompile("[^a-z]")

type Input struct {
	// InputFilePath is the path of the CSV to be splitted.
	//
	// Identified as --name="{path}".
	//
	// The {path} can be relative or absolute.
	//
	// Default: "./input.csv"
	InputFilePath string

	// OutputFilePathPattern is the pattern of the path of the outputs.
	//
	// Identified as --out="{path_pattern}".
	//
	// The {path_pattern} can be relative or absolute. The pattern will be handled by fmt.Sprintf()
	//
	// Default: "./output_%d.csv"
	OutputFilePathPattern string

	// MaxSize is the byte limit to split the CSV.
	//
	// Identified as --size="{size}".
	//
	// If the {size} has only numbers (being int or float), it will be interpreted as bytes.
	// You can also add a type abbreviation to the number as suffix:
	//    "b" for bytes
	//    "kb" for kilobytes
	//    "mb" for megabytes
	//    "gb" for gigabytes.
	//
	// Default: 10gb
	MaxSize FileSize `param:"size"`
}

func NewInput() *Input {
	input := &Input{
		InputFilePath: defaultInputFilePath,
		MaxSize:       defaultMaxSize,
	}

	if len(os.Args) <= 1 {
		return input
	}

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, inputFilePathID) {
			input.InputFilePath = parseArg(arg, inputFilePathID)
			continue
		}

		if strings.HasPrefix(arg, outputFilePathPatternID) {
			input.OutputFilePathPattern = parseArg(arg, outputFilePathPatternID)
			continue
		}

		if strings.HasPrefix(arg, maxSizeID) {
			arg = strings.ToLower(parseArg(arg, maxSizeID))
			unit := B
			if strings.HasSuffix(arg, "b") {
				unitID := regexNotLetters.ReplaceAllString(arg, "")
				switch unitID {
				case "b":
					unit = B
				case "kb":
					unit = KB
				case "mb":
					unit = MB
				case "gb":
					unit = GB
				default:
					panic(fmt.Sprintf("invalid size unit '%s'", unitID))
				}

				arg = strings.Replace(arg, unitID, "", 1)
			}

			size, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				panic(err)
			}

			size *= float64(unit)
			input.MaxSize = FileSize(int64(math.Ceil(size)))
			continue
		}

		panic(fmt.Sprintf("unexpected argument '%s'", arg))
	}

	return input
}

func parseArg(arg string, prefix string) string {
	arg = strings.Replace(arg, prefix, "", 1)
	if strings.HasPrefix(arg, stringDelimiter) && strings.HasSuffix(arg, stringDelimiter) {
		arg, _ = strings.CutPrefix(arg, stringDelimiter)
		arg, _ = strings.CutSuffix(arg, stringDelimiter)
	}

	return arg
}
