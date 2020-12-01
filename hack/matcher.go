package hack

import "net/http"

type RouteMatch int

const (
	NoMatch RouteMatch = iota
	StarMatch
	OverloadMatch
	ExactMatch
)

func (routes RouteDefinitions) MatchPath(path string) RouteDefinition{
	bestMatch := NoMatch
	var bestVals map[string]string
	var bestRoute *RouteDefinition
	for _, route := range routes {
		match, vals := route.Match(req.Method, req.URL.Path)
		if match.BetterThan(bestMatch) {
			bestMatch = match
			bestVals = vals
			bestRoute = route
			if match == ExactMatch {
				break
			}
		}
	}
	if bestMatch != NoMatch {
		params := Params(bestVals)
		context.Map(params)
		bestRoute.Handle(context, res)
		return
	}

	// no routes exist, 404
	c := &routeContext{context, 0, r.notFounds}
	context.MapTo(c, (*Context)(nil))
	c.run()

	return routes[0]
}

func Match(route RouteDefinition, matchMethod string) (RouteMatch, map[string]string) {
	// add Any method matching support
	match := MatchMethod(route, matchMethod)
	if match == NoMatch {
		return match, nil
	}

	//matches := r.regex.FindStringSubmatch(path)
	//if len(matches) > 0 && matches[0] == path {
	//	params := make(map[string]string)
	//	for i, name := range r.regex.SubexpNames() {
	//		if len(name) > 0 {
	//			params[name] = matches[i]
	//		}
	//	}
	//	return match, params
	//}
	return NoMatch, nil
}

func MatchMethod(r RouteDefinition, method string) RouteMatch {
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

