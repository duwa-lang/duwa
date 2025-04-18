package parser

import "github.com/tliron/kutil/problems"

type Context struct {
	Name     string
	Problems *problems.Problems
	Parser   *Parser
}