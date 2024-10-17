package renderer

import (
	"github.com/Jeffail/gabs/v2"
	"strings"
)

type Template struct {
	Segments []Segment `@@*`
}

func (template Template) Render(data *gabs.Container) (string, error) {
	var output strings.Builder

	for _, segment := range template.Segments {
		segmentString, err := segment.ToString(data)
		if err != nil {
			return "", err
		}
		output.WriteString(segmentString)
	}
	return output.String(), nil
}
