package browser

import (
	"fmt"

	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// routeItem wraps a spec.Route as a bubbles list.Item.
type routeItem struct {
	route spec.Route
}

func (i routeItem) Title() string {
	return fmt.Sprintf("%-6s %s", i.route.Method, i.route.Path)
}

func (i routeItem) Description() string {
	if i.route.Summary != "" {
		return i.route.Summary
	}
	return i.route.OperationID
}

func (i routeItem) FilterValue() string {
	return i.route.Method + " " + i.route.Path + " " + i.route.OperationID + " " + i.route.Summary + " " + i.route.Tag
}
