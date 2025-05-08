package lsp

import (
	"go.lsp.dev/protocol"
)

type Keyword struct {
	Text        string                      // The keyword text
	Kind        protocol.CompletionItemKind // The kind of completion item
	Detail      string                      // Short label for the kind of item
	Description string                      // Longer description of the keyword
}

func (se *Server) getKeywords() []Keyword {
	return []Keyword{
		{
			Text:        "nambala",
			Kind:        protocol.CompletionItemKindTypeParameter,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "mawu",
			Kind:        protocol.CompletionItemKindTypeParameter,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "zoona",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "bodza",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "ngati",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "kapena",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "bweza",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "ndondomeko",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "za",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "pamene",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "mgwirizano",
			Kind:        protocol.CompletionItemKindTypeParameter,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "kalasi",
			Kind:        protocol.CompletionItemKindTypeParameter,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "siya",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "pitirizani",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "palibe",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "tenga",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
		{
			Text:        "kuchokera",
			Kind:        protocol.CompletionItemKindKeyword,
			Detail:      "",
			Description: "",
		},
	}
}
