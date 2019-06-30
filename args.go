package main

import (
	"strings"
)

type Args struct {
	SchemaParserMap map[string]SchemaDetail
}

func NewArgs(schemaInput string) *Args {
	args := new(Args)
	args.initSchemaParserMap(schemaInput)
	return args
}

func (args *Args) initSchemaParserMap(schemaInput string) {
	args.SchemaParserMap = make(map[string]SchemaDetail)

	list := strings.Split(schemaInput, ",")

	for _, val := range list {
		flaglist := strings.Split(val, ":")
		args.SchemaParserMap[flaglist[0]] = SchemaDetail{}
	}
}
