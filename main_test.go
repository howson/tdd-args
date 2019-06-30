package main

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

// preparement for all tests
func (s *ParserSuite) SetUpSuite(c *C) {}

// test whether the schema can accept the input schema format and parse to cache.
// only check the length
func (s *ParserSuite) TestSchemaParseLength(c *C) {

	args := NewArgs("l:bool:false,f:string:.,d:int:0", "")

	c.Assert(len(args.SchemaParserMap), Equals, 3)
}

// test whether the schema can accept the input schema format and parse to cache.
// check if it can be parse correctly
func (s *ParserSuite) TestSchemaParseCorrectly(c *C) {

	args := NewArgs("l:bool:false,f:string:.,d:int:0", "")
	schemaDetail := args.SchemaParserMap["f"]

	c.Assert(schemaDetail.SchemaType, Equals, "string")
	c.Assert(schemaDetail.DefaultVal, Equals, ".")
}

// test if the args can accept input and parse to map with correct number of flag.
func (s *ParserSuite) TestArgsInputNum(c *C) {

	args := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 9231 -f /hh/oo")

	c.Assert(len(args.FlagMap), Equals, 3)
}

// test if the args can parse bool value with no input, and it will use default value which i set in schema
func (s *ParserSuite) TestArgsInputBoolDefault(c *C) {

	args := NewArgs("l:bool:false,f:string:.,d:int:0", "-l -d 9231 -f /hh/oo")
	_, input := args.GetValue("l")
	c.Assert(input, Equals, false)

	args = NewArgs("l:bool:true,f:string:.,d:int:0", "-l -d 9231 -f /hh/oo")
	_, input = args.GetValue("l")
	c.Assert(input, Equals, true)
}
