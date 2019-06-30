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

	var schemaType string
	var schemaDefaultVal string
	if len(flaglist) > 2 {
		schemaType = flaglist[1]
		schemaDefaultVal = flaglist[2]
	} else if len(flaglist) == 2 {
		schemaType = flaglist[1]
		schemaDefaultVal = NoDefaultError
	} else {
		return
	}

	args.SchemaParserMap[flaglist[0]] = &SchemaDetail{schemaType, schemaDefaultVal}

}

func (args *Args) initFlagMap(flagInput string) {
	args.FlagMap = make(map[string]string)
	for {
		inputCursor, spaceCursor := findCursor(flagInput)

		if inputCursor == -1 {
			break
		}

		//		fmt.Printf("flagInput:%s, inputCursor:%d, spaceCursor:%d\n", flagInput, inputCursor, spaceCursor)
		flagChar := flagInput[inputCursor+1 : spaceCursor]

		schemaDetail := args.containsFlagChar(flagChar)
		if schemaDetail == nil {
			flagInput = strings.Replace(flagInput, flagInput[inputCursor:spaceCursor], "", 1)
			continue
		}

		flagInput = strings.Trim(flagInput[spaceCursor:len(flagInput)], " ")
		value, changeInput, nextSpaceCursor := findParam(flagInput, schemaDetail)
		flagInput = changeInput
		args.FlagMap[flagChar] = value
		if nextSpaceCursor < 0 {
			break
		}
		if nextSpaceCursor == 0 {
			continue
		}
		flagInput = strings.Trim(flagInput[nextSpaceCursor+1:len(flagInput)], " ")

	}
}

func findParam(flagInput string, schemaDetail *SchemaDetail) (string, string, int) {
	nextInputCursor, nextSpaceCursor := findCursor(flagInput)

	if nextInputCursor == -1 || nextInputCursor+2 > len(flagInput) {
		return schemaDetail.DefaultVal, flagInput, -1
	}
	_, err := strconv.Atoi(flagInput[nextInputCursor+1 : nextInputCursor+2])
	if err != nil {
		if nextInputCursor > nextSpaceCursor {
			return flagInput[0:nextSpaceCursor], flagInput, nextSpaceCursor
		} else {
			return schemaDetail.DefaultVal, flagInput, 0
		}
	} else {
		value := flagInput[0:nextSpaceCursor]
		return value, strings.Replace(flagInput, value, "", 1), nextSpaceCursor
	}
}

func findCursor(flagInput string) (int, int) {
	inputCursor := strings.Index(flagInput, "-")
	spaceCursor := strings.Index(flagInput, " ")
	return inputCursor, spaceCursor
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
		value = args.FlagMap[flagStr]

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
