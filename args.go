package main

import (
	"strings"
)

type Args struct {
	SchemaParserMap map[string]SchemaDetail
}

const (
	NoDefaultError = "NoDefaultError"
)

func NewArgs(schemaInput string) *Args {
	args := new(Args)
	args.initSchemaParserMap(schemaInput)
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
