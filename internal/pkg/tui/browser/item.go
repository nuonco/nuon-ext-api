package browser

import (
	"fmt"
	"strings"

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
	desc := i.route.OperationID
	if i.route.Summary != "" {
		desc = i.route.Summary
	}

	if i.route.Deprecated {
		return "[deprecated] " + strings.TrimSpace(desc)
	}

	return desc
}

func (i routeItem) FilterValue() string {
	return i.route.Method + " " + i.route.Path + " " + i.route.OperationID + " " + i.route.Summary + " " + i.route.Tag
}
