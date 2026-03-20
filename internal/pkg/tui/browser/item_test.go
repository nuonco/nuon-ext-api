package browser

import (
	"testing"

	"github.com/nuonco/nuon-ext-api/internal/spec"
)

func TestRouteItemDescriptionDeprecatedSummary(t *testing.T) {
	item := routeItem{route: spec.Route{Summary: "list apps", Deprecated: true}}

	if got, want := item.Description(), "[deprecated] list apps"; got != want {
		t.Fatalf("unexpected description: got %q want %q", got, want)
	}
}

func TestRouteItemDescriptionDeprecatedOperationIDFallback(t *testing.T) {
	item := routeItem{route: spec.Route{OperationID: "ListApps", Deprecated: true}}

	if got, want := item.Description(), "[deprecated] ListApps"; got != want {
		t.Fatalf("unexpected description: got %q want %q", got, want)
	}
}

func TestRouteItemDescriptionNonDeprecated(t *testing.T) {
	item := routeItem{route: spec.Route{Summary: "list apps"}}

	if got, want := item.Description(), "list apps"; got != want {
		t.Fatalf("unexpected description: got %q want %q", got, want)
	}
}
