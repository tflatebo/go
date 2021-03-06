// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"internal/syscall/windows/registry"
	"testing"
	. "time"
)

func testZoneAbbr(t *testing.T) {
	t1 := Now()
	// discard nsec
	t1 = Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, t1.Location())
	t2, err := Parse(RFC1123, t1.Format(RFC1123))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if t1 != t2 {
		t.Fatalf("t1 (%v) is not equal to t2 (%v)", t1, t2)
	}
}

func TestLocalZoneAbbr(t *testing.T) {
	ResetLocalOnceForTest() // reset the Once to trigger the race
	defer ForceUSPacificForTesting()
	testZoneAbbr(t)
}

func TestAusZoneAbbr(t *testing.T) {
	ForceAusForTesting()
	defer ForceUSPacificForTesting()
	testZoneAbbr(t)
}

func TestToEnglishName(t *testing.T) {
	const want = "Central Europe Standard Time"
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Time Zones\`+want, registry.READ)
	if err != nil {
		t.Fatalf("cannot open CEST time zone information from registry: %s", err)
	}
	defer k.Close()
	std, _, err := k.GetStringValue("Std")
	if err != nil {
		t.Fatalf("cannot read CEST Std registry key: %s", err)
	}
	dlt, _, err := k.GetStringValue("Dlt")
	if err != nil {
		t.Fatalf("cannot read CEST Dlt registry key: %s", err)
	}
	name, err := ToEnglishName(std, dlt)
	if err != nil {
		t.Fatalf("toEnglishName failed: %s", err)
	}
	if name != want {
		t.Fatalf("english name: %q, want: %q", name, want)
	}
}
