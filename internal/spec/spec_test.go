package spec

import "testing"

func TestParseDeprecatedRoutesAreLast(t *testing.T) {
	api, err := Parse()
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	seenDeprecated := false
	for _, route := range api.Routes {
		if route.Deprecated {
			seenDeprecated = true
			continue
		}

		if seenDeprecated {
			t.Fatalf("found non-deprecated route %q after deprecated routes", route.DisplayName())
		}
	}
}

func TestParseReadsDeprecatedFlag(t *testing.T) {
	api, err := Parse()
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	var deprecatedCount int
	for _, route := range api.Routes {
		if route.Deprecated {
			deprecatedCount++
		}
	}

	if deprecatedCount == 0 {
		t.Fatal("expected at least one deprecated route from embedded spec")
	}
}

func TestListRoutesHidesDeprecatedByDefault(t *testing.T) {
	api := &API{Routes: []Route{
		{Method: "GET", Path: "/v1/apps"},
		{Method: "GET", Path: "/v1/legacy", Deprecated: true},
	}}

	routes := api.ListRoutes(false)
	if len(routes) != 1 {
		t.Fatalf("expected 1 non-deprecated route, got %d", len(routes))
	}

	if routes[0].Deprecated {
		t.Fatal("did not expect deprecated route when includeDeprecated=false")
	}
}

func TestListRoutesCanIncludeDeprecated(t *testing.T) {
	api := &API{Routes: []Route{
		{Method: "GET", Path: "/v1/apps"},
		{Method: "GET", Path: "/v1/legacy", Deprecated: true},
	}}

	routes := api.ListRoutes(true)
	if len(routes) != 2 {
		t.Fatalf("expected all routes when includeDeprecated=true, got %d", len(routes))
	}

	if !routes[1].Deprecated {
		t.Fatal("expected deprecated route to be present when includeDeprecated=true")
	}
}
