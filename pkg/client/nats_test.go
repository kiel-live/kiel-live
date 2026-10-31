package client

import "testing"

func TestLegacyVehicleID(t *testing.T) {
	if got := legacyVehicleID("kvg-123"); got != "kvg-123" {
		t.Errorf("legacyVehicleID() = %q, want %q", got, "kvg-123")
	}

	if got := legacyVehicleID(""); got == "" {
		t.Error("legacyVehicleID() must not return an empty id")
	}
}
