package main

import "testing"

func TestParseClusterResourceID(t *testing.T) {
	got, ok := parseClusterResourceID("/clusters/11111111-2222-3333-4444-555555555555")
	if !ok || got != "11111111-2222-3333-4444-555555555555" {
		t.Fatalf("%s %v", got, ok)
	}
	got, ok = parseClusterResourceID("11111111-2222-3333-4444-555555555555")
	if !ok || got == "" {
		t.Fatal("prefix is optional")
	}
	for _, bad := range []string{
		"",
		"/clusters/",
		"/clusters/../etc/passwd",
		"../release",
		"id/with/slash",
		"id:colon",
		"id\nnewline",
	} {
		if _, ok := parseClusterResourceID(bad); ok {
			t.Fatalf("must reject %q", bad)
		}
	}
}
