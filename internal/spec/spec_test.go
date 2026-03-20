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
