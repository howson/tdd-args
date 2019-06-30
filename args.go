package main

import (
	"encoding/json"
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

func NewArgs(schemaInput string, flagInput string) (*Args, error) {
	args := new(Args)
	args.initSchemaParserMap(schemaInput)
	err := args.initFlagMap(flagInput)
	return args, err
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

func (args *Args) initFlagMap(flagInput string) error {
	args.FlagMap = make(map[string]string)
	flagInput = strings.Trim(flagInput, " ")
	if flagInput == "" {
		return nil
	}
	for {
		if flagInput != "" && flagInput[0:1] != "-" {
			return UnsupportedError("unsupported input format, initialization will be terminated.")
		}

		inputCursor, spaceCursor := findCursor(flagInput)

		if inputCursor == -1 {
			break
		}

		if spaceCursor <= inputCursor {
			spaceCursor = len(flagInput)
		}

		flagChar := flagInput[inputCursor+1 : spaceCursor]

		schemaDetail := args.containsFlagChar(flagChar)
		if schemaDetail == nil {
			flagInput = strings.Replace(flagInput, flagInput[inputCursor:spaceCursor], "", 1)
			return UnsupportedError("unsupported input format, initialization will be terminated.")
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

	return nil
}

func findParam(flagInput string, schemaDetail *SchemaDetail) (string, string, int) {
	nextInputCursor, nextSpaceCursor := findCursor(flagInput)

	if nextInputCursor == -1 || nextInputCursor+2 > len(flagInput) {
		value := strings.Trim(flagInput, " ")
		if value == "" {
			return schemaDetail.DefaultVal, flagInput, -1
		} else {
			return value, flagInput, -1
		}

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

	value := args.FlagMap[flagStr]
	if value == NoDefaultError {
		return nil, nil
	}

	if schemaDetail.SchemaType == "bool" {
		result, err := atob(value)
		if err != nil {
			return UnsupportedError(fmt.Sprintf("value %s can not be parsed to bool.", value)), nil
		}
		return nil, result
	}

	if schemaDetail.SchemaType == "int" {
		result, err := atoi(value)
		if err != nil {
			return UnsupportedError(fmt.Sprintf("value %s can not be parsed to int.", value)), nil
		}
		return nil, result
	}

	if schemaDetail.SchemaType == "string" {
		return nil, value
	}

	return nil, value
}

func atoi(str string) (int, error) {
	return strconv.Atoi(str)
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

// convert all the object to a json string
func MarshalObjToJson(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(result)
}
