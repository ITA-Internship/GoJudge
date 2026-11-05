package report

import (
	"strings"
	"testing"
)

func TestParseSuccess(t *testing.T) {
	input := `{"Time":"2026-11-05T10:00:00Z","Action":"run","Package":"trusted","Test":"TestAddSuccess"}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Test":"TestAddSuccess","Output":"=== RUN   TestAddSuccess\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"pass","Package":"trusted","Test":"TestAddSuccess","Elapsed":0.01}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Output":"PASS\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"pass","Package":"trusted","Elapsed":0.02}
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TestsRun != 1 {
		t.Errorf("expected TestsRun=1, got %d", s.TestsRun)
	}
	if s.TestsPassed != 1 {
		t.Errorf("expected TestsPassed=1, got %d", s.TestsPassed)
	}
	if s.PackageVerdict != "pass" {
		t.Errorf("expected PackageVerdict=pass, got %s", s.PackageVerdict)
	}
}

func TestParseNoTestsToRunTrap(t *testing.T) {
	// When -run has a typo, go test produces warning and package pass, but no test 'run' events
	input := `{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Output":"testing: warning: no tests to run\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Output":"PASS\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"pass","Package":"trusted","Elapsed":0.00}
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TestsRun != 0 {
		t.Errorf("expected TestsRun=0 for empty run, got %d", s.TestsRun)
	}
	if s.PackageVerdict != "pass" {
		t.Errorf("expected PackageVerdict=pass from go test, got %s", s.PackageVerdict)
	}
}

func TestParseFail(t *testing.T) {
	input := `{"Time":"2026-11-05T10:00:00Z","Action":"run","Package":"trusted","Test":"TestFailAssertion"}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Test":"TestFailAssertion","Output":"=== RUN   TestFailAssertion\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Test":"TestFailAssertion","Output":"    calc_test.go:21: intentional failure\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"fail","Package":"trusted","Test":"TestFailAssertion","Elapsed":0.00}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Output":"FAIL\n"}
{"Time":"2026-11-05T10:00:00Z","Action":"fail","Package":"trusted","Elapsed":0.01}
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TestsRun != 1 || s.TestsFailed != 1 || s.TestsPassed != 0 {
		t.Errorf("unexpected counts: run=%d passed=%d failed=%d", s.TestsRun, s.TestsPassed, s.TestsFailed)
	}
	if s.PackageVerdict != "fail" {
		t.Errorf("expected PackageVerdict=fail, got %s", s.PackageVerdict)
	}
}

func TestParseTruncatedStream(t *testing.T) {
	// Stream killed mid-execution by DeadlineExceeded or OOMKilled
	input := `{"Time":"2026-11-05T10:00:00Z","Action":"start","Package":"trusted"}
{"Time":"2026-11-05T10:00:00Z","Action":"run","Package":"trusted","Test":"TestInfiniteLoop"}
{"Time":"2026-11-05T10:00:00Z","Action":"output","Package":"trusted","Test":"TestInfiniteLoop","Output":"=== RUN   TestInfiniteLoop\n"}
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TestsRun != 1 {
		t.Errorf("expected TestsRun=1, got %d", s.TestsRun)
	}
	if s.PackageVerdict != "unknown" {
		t.Errorf("expected PackageVerdict=unknown on truncated stream, got %s", s.PackageVerdict)
	}
}

func TestParseNonJSONNoiseBetweenLines(t *testing.T) {
	input := `some container startup banner
{"Time":"2026-11-05T10:00:00Z","Action":"run","Package":"trusted","Test":"TestAddSuccess"}
non-json stderr line
{"Time":"2026-11-05T10:00:00Z","Action":"pass","Package":"trusted","Test":"TestAddSuccess","Elapsed":0.01}
{"Time":"2026-11-05T10:00:00Z","Action":"pass","Package":"trusted","Elapsed":0.02}
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TestsRun != 1 || s.TestsPassed != 1 || s.PackageVerdict != "pass" {
		t.Errorf("unexpected summary after noise: %+v", s)
	}
}

