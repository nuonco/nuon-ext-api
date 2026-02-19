package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nuonco/nuon-ext-api/internal/client"
)

// Print writes an API response to stdout.
// If raw is true, prints the body as-is. Otherwise pretty-prints JSON.
func Print(resp *client.Response, raw bool) error {
	if resp.StatusCode >= 400 {
		return printError(resp, raw)
	}

	if raw {
		_, err := os.Stdout.Write(resp.Body)
		if err == nil {
			fmt.Println()
		}
		return err
	}

	return prettyPrint(resp.Body)
}

func printError(resp *client.Response, raw bool) error {
	if raw {
		os.Stderr.Write(resp.Body)
		fmt.Fprintln(os.Stderr)
	} else {
		pretty, err := formatJSON(resp.Body)
		if err != nil {
			// Not JSON — print raw
			fmt.Fprintln(os.Stderr, string(resp.Body))
		} else {
			fmt.Fprintln(os.Stderr, pretty)
		}
	}
	return fmt.Errorf("HTTP %d", resp.StatusCode)
}

func prettyPrint(data []byte) error {
	pretty, err := formatJSON(data)
	if err != nil {
		// Not JSON — print as-is
		_, err := os.Stdout.Write(data)
		if err == nil {
			fmt.Println()
		}
		return err
	}

	fmt.Println(pretty)
	return nil
}

func formatJSON(data []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}
