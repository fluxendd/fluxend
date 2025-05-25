package commands

import (
	"fluxend/internal/app"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// routesCmd represents the command to list all routes
var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List all registered API routes",
	Run: func(cmd *cobra.Command, args []string) {
		printRoutes()
	},
}

func printRoutes() {
	e := setupServer(app.InitializeContainer())

	registeredRoutes := e.Routes()
	routesGrouped := make(map[string][]string)

	for _, route := range registeredRoutes {
		if route.Method == "echo_route_not_found" {
			continue
		}

		routePrefix := getRoutePrefix(route.Path)
		routesGrouped[routePrefix] = append(routesGrouped[routePrefix], fmt.Sprintf("Method: %-6s  Path: %-30s", route.Method, route.Path))
	}

	for prefix, matchedRoutes := range routesGrouped {
		fmt.Printf("\nRoutes for %s:\n", prefix)
		for _, route := range matchedRoutes {
			fmt.Println(route)
		}
	}
}

func getRoutePrefix(path string) string {
	segments := strings.Split(path, "/")
	if len(segments) > 1 {
		return "/" + segments[2]
	}
	return path
}
