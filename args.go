package main

import (
	"fmt"
	"strings"
)

type Args struct {
	SchemaParserMap map[string]SchemaDetail
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
	args.SchemaParserMap = make(map[string]SchemaDetail)

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

	args.SchemaParserMap[flaglist[0]] = SchemaDetail{schType, schDefaultVal}

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
		flagInput = strings.Trim(flagInput[nextSpaceCursor+1:len(flagInput)], " ")
	}

}

func (args *Args) containsFlagChar(flagChar string) *SchemaDetail {
	for key, val := range args.SchemaParserMap {
		if key == flagChar {
			return &val
		}
	}
	return nil
}
