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

	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "")

	c.Assert(len(args.SchemaParserMap), Equals, 3)
}

// test whether the schema can accept the input schema format and parse to cache.
// check if it can be parse correctly
func (s *ParserSuite) TestSchemaParseCorrectly(c *C) {

	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "")
	schemaDetail := args.SchemaParserMap["f"]

	c.Assert(schemaDetail.SchemaType, Equals, "string")
	c.Assert(schemaDetail.DefaultVal, Equals, ".")
}

// test if the args can accept input and parse to map with correct number of flag.
func (s *ParserSuite) TestArgsInputNum(c *C) {

	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 9231 -f /hh/oo")

	c.Assert(len(args.FlagMap), Equals, 3)
}

// test if the args can parse bool value with no input, and it will use default value which i set in schema
func (s *ParserSuite) TestArgsInputBoolDefault(c *C) {

	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l -d 9231 -f /hh/oo")
	_, input := args.GetValue("l")
	c.Assert(input, Equals, false)

	args, _ = NewArgs("l:bool:true,f:string:.,d:int:0", "-l -d 9231 -f /hh/oo")
	_, input = args.GetValue("l")
	c.Assert(input, Equals, true)
}

// test if the args can parse bool value with input. legal input should be considered
func (s *ParserSuite) TestArgsInputBoolLegalParam(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 9231 -f /hh/oo")
	_, input := args.GetValue("l")
	c.Assert(input, Equals, true)
}

// test if the args can parse bool value with input. illegal input should be considered
func (s *ParserSuite) TestArgsInputBoolIllegalParam(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l TRue -d 9231 -f /hh/oo")
	err, _ := args.GetValue("l")
	c.Assert(err, NotNil)
}

// test if the args can parse integer value. legal value should be considered
func (s *ParserSuite) TestArgsInputIntParam(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 9231 -f /hh/oo")
	_, input := args.GetValue("d")
	c.Assert(input, Equals, 9231)
}

// test if the args can parse integer value. legal minus value should be considered
func (s *ParserSuite) TestArgsInputIntMinusParam(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d -1021 -f /hh/oo")
	_, input := args.GetValue("d")
	c.Assert(input, Equals, -1021)
}

// test if the args can parse integer value. default value should be considered
func (s *ParserSuite) TestArgsInputIntDefault(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d -f /hh/oo")
	_, input := args.GetValue("d")
	c.Assert(input, Equals, 0)
}

// test if the args can parse integer value. default value should be considered
func (s *ParserSuite) TestArgsInputIntIllegalParam(c *C) {
	args, _ := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 5abc -f /hh/oo")
	err, _ := args.GetValue("d")
	c.Assert(err, NotNil)
}

// test if the args can parse integer value. default value should be considered
func (s *ParserSuite) TestArgsInputIntIllegalParam2(c *C) {
	_, err := NewArgs("l:bool:false,f:string:.,d:int:0", "-l true -d 10 21 -f /hh/oo")
	c.Assert(err, NotNil)
}
