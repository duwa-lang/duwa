package lsp

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/kutil/problems"
)

func createDiagnostics(problems *problems.Problems, content string) []protocol.Diagnostic {
	diagnostics := make([]protocol.Diagnostic, len(problems.Problems))

	for index, problem := range problems.Problems {
		row := problem.Row - 1
		column := problem.Column - 1
		if row < 0 {
			row = 0
		}
		if column < 0 {
			column = 0
		}

		position := protocol.Position{
			Line:      protocol.UInteger(row),
			Character: protocol.UInteger(column),
		}

		severity := protocol.DiagnosticSeverityError
		diagnostics[index] = protocol.Diagnostic{
			Range: protocol.Range{
				Start: position,
				End:   position.EndOfLineIn(content),
			},
			Severity: &severity,
			Code:     &protocol.IntegerOrString{Value: problem.Item},
			//Source:   "TOSCA",
			Message: problem.Message,
		}
	}

	return diagnostics
}
