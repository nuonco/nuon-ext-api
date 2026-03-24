package dispatch

import (
	"testing"

	"github.com/nuonco/nuon-ext-api/internal/config"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

func TestResolvePreservesConcreteSegmentsInMixedPath(t *testing.T) {
	api := &spec.API{Routes: []spec.Route{
		{Path: "/v1/installs/{install_id}/actions/{action_id}", Method: "GET", OperationID: "GetInstallAction"},
	}}

	cfg := &config.Config{InstallID: "ins_123"}
	inputPath := "/v1/installs/{install_id}/actions/iawag6pbgfzvlkyqdiy2a1xw6j"

	req, err := Resolve(api, inputPath, "", "", cfg, nil)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	wantPath := "/v1/installs/ins_123/actions/iawag6pbgfzvlkyqdiy2a1xw6j"
	if req.Path != wantPath {
		t.Fatalf("expected resolved path %q, got %q", wantPath, req.Path)
	}
}

func TestResolveUsesTemplateParamNamesForUnresolvedSegments(t *testing.T) {
	api := &spec.API{Routes: []spec.Route{
		{Path: "/v1/installs/{install_id}/actions/{action_id}", Method: "GET", OperationID: "GetInstallAction"},
	}}

	cfg := &config.Config{InstallID: "ins_123"}
	inputPath := "/v1/installs/{foo}/actions/iawag6pbgfzvlkyqdiy2a1xw6j"

	req, err := Resolve(api, inputPath, "", "", cfg, nil)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	wantPath := "/v1/installs/ins_123/actions/iawag6pbgfzvlkyqdiy2a1xw6j"
	if req.Path != wantPath {
		t.Fatalf("expected resolved path %q, got %q", wantPath, req.Path)
	}
}
