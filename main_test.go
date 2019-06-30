package main

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

// preparement for all tests
func (s *ParserSuite) SetUpSuite(c *C) {}
