package hack

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

type SwaggerExporter struct {
	routeDefs RouteDefinitions
}

func NewSwaggerExporter(routeDefs RouteDefinitions) SwaggerExporter {
	return SwaggerExporter{routeDefs: routeDefs}
}

var paramRegex = regexp.MustCompile(":([^/]*)")

func formatRoute(route string) string {
	route = strings.ReplaceAll(route, "?", "")
	route = paramRegex.ReplaceAllString(route, "{$1}")
	return route
}

func formatMethod(method string) string {
	if method == "*" {
		method = "trace"  // TODO: Fix * method
	}
	method = strings.ToLower(method)
	return method
}

func (s SwaggerExporter) PrintYaml() {
	s.ExportToYaml(os.Stdout)
}

func writeSchemaWithoutName(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%stype: %s\n", strings.Repeat(" ", indent), s.Type)
	if len(s.Children) > 0 {
		fmt.Fprintf(w, "%sproperties:\n", strings.Repeat(" ", indent))
		for _, c := range s.Children {
			writeSchema(w, c, indent + 2)
		}
	}
	if s.Type == FieldArray {
		fmt.Fprintf(w, "%sitems:\n", strings.Repeat(" ", indent))
		writeSchemaWithoutName(w, *s.ArrayType, indent + 2)
	}
}

func writeSchema(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", indent), s.Name)
	writeSchemaWithoutName(w, s, indent + 2)
}

func writeRequestBody(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%srequestBody:\n", strings.Repeat(" ", indent))
	fmt.Fprintf(w, "%scontent:\n", strings.Repeat(" ", indent + 2))
	fmt.Fprintf(w, "%sapplication/json:\n", strings.Repeat(" ", indent + 4))
	fmt.Fprintf(w, "%sschema:\n", strings.Repeat(" ", indent + 6))
	writeSchemaWithoutName(w, s, indent + 8)
}

func writeQuerySchema(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%s- in: query\n", strings.Repeat(" ", indent))
	fmt.Fprintf(w, "%sname: %s\n", strings.Repeat(" ", indent + 2), s.Name)
	fmt.Fprintf(w, "%sschema:\n", strings.Repeat(" ", indent + 2))
	fmt.Fprintf(w, "%stype: %s\n", strings.Repeat(" ", indent + 4), s.Type)
	if s.Type == FieldArray {
		fmt.Fprintf(w, "%sitems:\n", strings.Repeat(" ", indent + 4))
		fmt.Fprintf(w, "%stype: %s\n", strings.Repeat(" ", indent + 6), s.ArrayType.Type)
	}
}

func writeParams(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%sparameters:\n", strings.Repeat(" ", indent))
	for _, c := range s.Children {
		writeQuerySchema(w, c, indent + 2)
	}

}

func (s SwaggerExporter) exportPaths(w io.Writer) {
	fmt.Fprintln(w, "paths:")
	lastRoute := ""

	// TODO: Fix mutation
	sort.Slice(s.routeDefs, func(i, j int) bool {
		return s.routeDefs[i].Route < s.routeDefs[j].Route
	})

	for _, r := range s.routeDefs {
		if r.Route != lastRoute {
			fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", 2), formatRoute(r.Route))
		}
		method := formatMethod(r.Method)
		fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", 4), method)
		if r.Schema.Type != FieldInvalid {
			switch method {
			case "post", "put", "patch":
				writeRequestBody(w, r.Schema, 6)
			default:
				writeParams(w, r.Schema, 6)
			}
		}

		fmt.Fprintf(w, "%sresponses:\n", strings.Repeat(" ", 6))
		fmt.Fprintf(w, "%s'200':\n", strings.Repeat(" ", 8))
		fmt.Fprintf(w, "%sdescription: OK\n", strings.Repeat(" ", 10))

		lastRoute = r.Route
	}
}

func (s SwaggerExporter) ExportToYaml(w io.Writer) {
	fmt.Fprintln(w, "openapi: 3.0.0")
	fmt.Fprintf(w, "info:\n  version: 1.0.0\n  title: HACK THE PLANET\n")
	s.exportPaths(w)
}
