package renderer

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"io"
)

var (
	templateLexer = lexer.MustStateful(lexer.Rules{
		"Root": {
			{`String`, `[^{\\]+`, nil},
			{`InterpolationStart`, `{{`, lexer.Push("Interpolation")},
			{`SingleBrace`, `{|\\{`, nil},
			{"Whitespace", `\s+`, nil},
		},
		"Interpolation": {
			{`InterpolationEnd`, `}`, lexer.Pop()},
			{`PropertyLookup`, `p\.`, lexer.Push("PropertyLookup")},
			{`SpecLookup`, `spec\.`, lexer.Push("SpecLookup")},
		},
		"PropertyLookup": {
			{`InterpolationEnd`, `}`, lexer.Pop()},
			{`PropertyLookupIdentifier`, `[^}]+`, nil},
		},
		"SpecLookup": {
			{`InterpolationEnd`, `}`, lexer.Pop()},
			{`SpecLookupIdentifier`, `[^}]+`, nil},
		},
	})

	templateParser = participle.MustBuild[Template](
		participle.Lexer(templateLexer),
		participle.Union[Segment](StringSegment{}, PropertyInterpolationSegment{}, SpecInterpolationSegment{}, SingleBraceSegment{}, WhitespaceSegment{}),
		participle.Elide("InterpolationStart", "InterpolationEnd", "PropertyLookup", "SpecLookup"),
	)
)

type Segment interface {
	ToString(data *gabs.Container) (string, error)
}

type StringSegment struct {
	Body string `@String`
}

func (segment StringSegment) ToString(*gabs.Container) (string, error) {
	return segment.Body, nil
}

type WhitespaceSegment struct {
	Body string `@Whitespace`
}

func (segment WhitespaceSegment) ToString(*gabs.Container) (string, error) {
	return segment.Body, nil
}

type PropertyInterpolationSegment struct {
	InterpolationString string `@PropertyLookupIdentifier`
}

func (segment PropertyInterpolationSegment) ToString(data *gabs.Container) (string, error) {
	currentData := data.Search("properties").Path(segment.InterpolationString)
	if currentData == nil {
		return "", errors.New(fmt.Sprintf("p.%s did not match any provided properties", segment.InterpolationString))
	}
	rawData := currentData.Data()
	formatted, ok := rawData.(string)
	if ok {
		return formatted, nil
	}
	return currentData.String(), nil
}

type SpecInterpolationSegment struct {
	InterpolationString string `@SpecLookupIdentifier`
}

func (segment SpecInterpolationSegment) ToString(data *gabs.Container) (string, error) {
	currentData := data.Search("spec").Path(segment.InterpolationString)
	if currentData == nil {
		return "", errors.New(fmt.Sprintf("spec.%s did not match any provided properties", segment.InterpolationString))
	}
	rawData := currentData.Data()
	formatted, ok := rawData.(string)
	if ok {
		return formatted, nil
	}
	return currentData.String(), nil
}

type SingleBraceSegment struct {
	Body string `@SingleBrace`
}

func (segment SingleBraceSegment) ToString(*gabs.Container) (string, error) {
	return segment.Body[len(segment.Body)-1:], nil
}

func Parse(r io.Reader) (*Template, error) {
	template, err := templateParser.Parse("", r)
	if err != nil {
		return nil, err
	}
	return template, nil
}
