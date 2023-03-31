package routing

import (
	"github.com/gopulse/pulse-router/constants"
	"strings"
)

type RouteParser struct {
	segments []*segment // segments of the route
	params   []string   // parameters of the route
}

type segment struct {
	isParam    bool   // true if the segment is a parameter
	paramName  string // name of the parameter
	isLast     bool   // true if the segment is the last segment
	isOptional bool   // true if the segment is optional
}

// Parse parses the route and returns a routeParser
func Parse(route string) RouteParser {
	var parser RouteParser
	var segmentStruct segment
	var lastSegment *segment

	// Split the route into segments
	segments := strings.Split(route, "/")

	// Iterate over the segments
	for i, s := range segments {
		// Check if the segmentStruct is a parameter
		if strings.HasPrefix(s, constants.ParamSign) {
			segmentStruct.isParam = true
			segmentStruct.paramName = s[1:]
			parser.params = append(parser.params, segmentStruct.paramName)
		} else if strings.HasPrefix(s, constants.OptionalSign) {
			segmentStruct.isParam = true
			segmentStruct.isOptional = true
			segmentStruct.paramName = s[1:]
			parser.params = append(parser.params, segmentStruct.paramName)
		} else if s == constants.WildcardSign {
			segmentStruct.isParam = true
			segmentStruct.paramName = constants.WildcardSign
			parser.params = append(parser.params, segmentStruct.paramName)
		} else {
			segmentStruct.isParam = false
		}

		// Check if the segmentStruct is the last segmentStruct
		if i == len(segments)-1 {
			segmentStruct.isLast = true
		}

		// Append the segmentStruct to the parser
		parser.segments = append(parser.segments, &segmentStruct)

		// Set the last segmentStruct
		lastSegment = &segmentStruct

		// Reset the segmentStruct
		segmentStruct = segment{}
	}

	// Check if the last segmentStruct is a parameter
	if lastSegment != nil && lastSegment.isParam {
		// Set the last segmentStruct as optional
		lastSegment.isOptional = true
	}

	return parser
}
