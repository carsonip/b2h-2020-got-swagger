package hack

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type SwaggerExporter struct {
	routeDefs RouteDefinitions
}

func NewSwaggerExporter(routeDefs RouteDefinitions) SwaggerExporter {
	return SwaggerExporter{routeDefs: routeDefs}
}

func formatRoute(route string) string {
	return route
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
}

func writeSchema(w io.Writer, s Schema, indent int) {
	fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", indent), s.Name)
	writeSchemaWithoutName(w, s, indent + 2)
}

func (s SwaggerExporter) ExportToYaml(w io.Writer) {
	fmt.Fprintln(w, "openapi: 3.0.0")
	fmt.Fprintf(w, "info:\n  version: 1.0.0\n  title: HACK THE PLANET\n")
	fmt.Fprintln(w, "paths:")
	for _, r := range s.routeDefs {
		fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", 2), formatRoute(r.Route))
		fmt.Fprintf(w, "%s%s:\n", strings.Repeat(" ", 4), strings.ToLower(r.Method))
		if r.Schema.Name != "" {
			fmt.Fprintf(w, "%srequestBody:\n", strings.Repeat(" ", 6))
			fmt.Fprintf(w, "%scontent:\n", strings.Repeat(" ", 8))
			fmt.Fprintf(w, "%sapplication/json:\n", strings.Repeat(" ", 10))
			fmt.Fprintf(w, "%sschema:\n", strings.Repeat(" ", 12))
			writeSchemaWithoutName(w, r.Schema, 14)
		}

		fmt.Fprintf(w, "%sresponses:\n", strings.Repeat(" ", 6))
		fmt.Fprintf(w, "%s'200':\n", strings.Repeat(" ", 8))
		fmt.Fprintf(w, "%sdescription: OK\n", strings.Repeat(" ", 10))
	}
}
