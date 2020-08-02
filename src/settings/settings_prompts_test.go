package settings

import "testing"

func TestExecuteIfNil(t *testing.T) {
	executed := false
	logic := func() error {
		executed = true
		return nil
	}
	updated := false
	err := executeIfNil(nil, logic, &updated)
	if err != nil {
		t.Errorf("should not have errored: %v", err)
	}
	if !executed {
		t.Error("logic should have executed if setting was nil")
	}
	if !updated {
		t.Error("updated should be true if it was false before and logic ran")
	}

	executed = false
	updated = false
	setting := "some setting"
	err = executeIfNil(&setting, logic, &updated)
	if err != nil {
		t.Errorf("should not have errored: %v", err)
	}
	if executed {
		t.Error("logic should not have executed if setting was not nil")
	}
	if updated {
		t.Error("updated should not be true if logic did not run")
	}

	executed = false
	updated = true
	err = executeIfNil(nil, logic, &updated)
	if err != nil {
		t.Errorf("should not have errored: %v", err)
	}
	if !executed {
		t.Error("logic should have executed if setting was nil")
	}
	if !updated {
		t.Error("updated should still be true after logic runs")
	}

	executed = false
	updated = true
	err = executeIfNil(&setting, logic, &updated)
	if err != nil {
		t.Errorf("should not have errored: %v", err)
	}
	if executed {
		t.Error("logic should not have executed if the setting was not nil")
	}
	if !updated {
		t.Error("updated should not be flipped from true to false if logic did not run")
	}
}
