package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Args struct {
	SchemaParserMap map[string]*SchemaDetail
	FlagMap         map[string]string
}

const (
	NoDefaultError = "NoDefaultError"
)

func NewArgs(schemaInput string, flagInput string) *Args {
	args := new(Args)
	args.initSchemaParserMap(schemaInput)
	args.initFlagMap(flagInput)
	return args
}

func (args *Args) initSchemaParserMap(schemaInput string) {
	args.SchemaParserMap = make(map[string]*SchemaDetail)

	list := strings.Split(schemaInput, ",")

	for _, val := range list {
		args.newSchemaDetail(val)
	}
}

func (args *Args) newSchemaDetail(flagstr string) {
	flaglist := strings.Split(flagstr, ":")

	var schType string
	var schDefaultVal string
	if len(flaglist) > 2 {
		schType = flaglist[1]
		schDefaultVal = flaglist[2]
	} else if len(flaglist) == 2 {
		schType = flaglist[1]
		schDefaultVal = NoDefaultError
	} else {
		return
	}

	args.SchemaParserMap[flaglist[0]] = &SchemaDetail{schType, schDefaultVal}

}

func (args *Args) initFlagMap(flagInput string) {
	args.FlagMap = make(map[string]string)
	for {
		inputCursor := strings.Index(flagInput, "-")
		spaceCursor := strings.Index(flagInput, " ")

		if inputCursor == -1 {
			break
		}

		flagChar := flagInput[inputCursor+1 : spaceCursor]

		if schemaDetail := args.containsFlagChar(flagChar); schemaDetail != nil {
			args.FlagMap[flagChar] = "input"
		}

		flagInput = strings.Trim(flagInput[spaceCursor:len(flagInput)], " ")

		nextSpaceCursor := strings.Index(flagInput, " ")
		nextInputCursor := strings.Index(flagInput, "-")

		if nextInputCursor == -1 || nextInputCursor+2 > len(flagInput) {
			break
		}
		_, err := strconv.Atoi(flagInput[nextInputCursor : nextInputCursor+2])
		if err != nil {
			if nextInputCursor > nextSpaceCursor {
				flagInput = strings.Trim(flagInput[nextSpaceCursor+1:len(flagInput)], " ")
			}
		}

	}

}

func (args *Args) containsFlagChar(flagChar string) *SchemaDetail {
	for key, val := range args.SchemaParserMap {
		if key == flagChar {
			return val
		}
	}
	return nil
}

func (args *Args) GetValue(flagStr string) (error, interface{}) {
	schemaDetail := args.SchemaParserMap[flagStr]
	if schemaDetail == nil {
		return UnsupportedError("the input flag is not supported"), nil
	}

	var value string
	if schemaDetail.SchemaType == "bool" {
		value = schemaDetail.DefaultVal
		if value == NoDefaultError {
			return UnsupportedError(fmt.Sprintf("param %s does not support empty input.", flagStr)), nil
		}

		result, err := atob(value)
		if err != nil {
			return UnsupportedError(fmt.Sprintf("value %s can not be parsed.", value)), nil
		}
		return nil, result
	}
	return nil, nil
}

func atob(str string) (bool, error) {
	if str == "true" || str == "TRUE" {
		return true, nil
	}

	if str == "false" || str == "FALSE" {
		return false, nil
	}
	return false, errors.New("format error: can not parse the string to bool.")
}
