package hack

import (
	"fmt"
	"regexp"
	"strings"
)

type RouteMatch int

const (
	NoMatch RouteMatch = iota
	StarMatch
	OverloadMatch
	ExactMatch
)

func (routes RouteDefinitions) MatchPath(method string, path string) RouteDefinition {
	bestMatch := NoMatch
	var bestRoute RouteDefinition

	for _, route := range routes {
		match := route.match(method, path)
		if match > bestMatch {
			bestMatch = match
			bestRoute = route
			if match == ExactMatch {
				break
			}
		}
	}

	return bestRoute
}

func (route RouteDefinition) match(matchMethod string, path string) RouteMatch {
	// add Any method matching support
	match := MatchMethod(route, matchMethod)
	if match == NoMatch {
		return match
	}
	regex := buildRegex(route.Route)
	matches := regex.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		return match
	}

	return NoMatch
}

func buildRegex(pattern string) *regexp.Regexp {
	pattern = routeReg1.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	var index int
	pattern = routeReg2.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	return regexp.MustCompile(pattern)
}

func MatchMethod(r RouteDefinition, method string) RouteMatch {
	method = strings.ToUpper(method)
	switch {
	case method == r.Method:
		return ExactMatch
	case method == "HEAD" && r.Method == "GET":
		return OverloadMatch
	case r.Method == "*":
		return StarMatch
	default:
		return NoMatch
	}
}
