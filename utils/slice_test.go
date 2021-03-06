/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestIsSliceMember(t *testing.T) {
	if !IsSliceMember([]string{"1", "2", "3", "4", "5"}, "5") {
		t.Error("Expecting: true, received: false")
	}
	if IsSliceMember([]string{"1", "2", "3", "4", "5"}, "6") {
		t.Error("Expecting: true, received: false")
	}
}

func TestSliceHasMember(t *testing.T) {
	if !SliceHasMember([]string{"1", "2", "3", "4", "5"}, "5") {
		t.Error("Expecting: true, received: false")
	}
	if SliceHasMember([]string{"1", "2", "3", "4", "5"}, "6") {
		t.Error("Expecting: true, received: false")
	}
}

func TestSliceWithoutMember(t *testing.T) {
	rcv := SliceWithoutMember([]string{"1", "2", "3", "4", "5"}, "5")
	sort.Strings(rcv)
	eOut := []string{"1", "2", "3", "4"}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	rcv = SliceWithoutMember([]string{"1", "2", "3", "4", "5"}, "6")
	sort.Strings(rcv)
	eOut = []string{"1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestSliceMemberHasPrefix(t *testing.T) {
	if !SliceMemberHasPrefix([]string{"1", "*2", "3", "4", "5"}, "*") {
		t.Error("Expecting: true, received: false")
	}
	if SliceMemberHasPrefix([]string{"1", "2", "3", "4", "5"}, "*") {
		t.Error("Expecting: true, received: false")
	}
}

func TestAvg(t *testing.T) {
	if rcv := Avg([]float64{}); rcv != 0 {
		t.Errorf("Expecting: 0, received: %+v", rcv)
	}
	if rcv := Avg([]float64{1, 2, 3}); rcv != 2 {
		t.Errorf("Expecting: 2, received: %+v", rcv)
	}
	if rcv := Avg([]float64{1.5, 2.75, 3.25}); rcv != 2.5 {
		t.Errorf("Expecting: 2.5, received: %+v", rcv)
	}
}

func TestAvgNegative(t *testing.T) {
	if rcv := AvgNegative([]float64{}); rcv != -1 {
		t.Errorf("Expecting: -1, received: %+v", rcv)
	}
	if rcv := AvgNegative([]float64{1, 2, 3}); rcv != 2 {
		t.Errorf("Expecting: 2, received: %+v", rcv)
	}
	if rcv := Avg([]float64{1.5, 2.75, 3.25}); rcv != 2.5 {
		t.Errorf("Expecting: 2.5, received: %+v", rcv)
	}
}

func TestPrefixSliceItems(t *testing.T) {
	rcv := PrefixSliceItems([]string{"1", "2", "3", "4", "5"}, "*")
	sort.Strings(rcv)
	eOut := []string{"*1", "*2", "*3", "*4", "*5"}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestStripSlicePrefix(t *testing.T) {
	eSlc := make([]string, 0)
	if retSlc := StripSlicePrefix([]string{}, 2); !reflect.DeepEqual(eSlc, retSlc) {
		t.Errorf("expecting: %+v, received: %+v", eSlc, retSlc)
	}
	eSlc = []string{"1", "2"}
	if retSlc := StripSlicePrefix([]string{"0", "1", "2"}, 1); !reflect.DeepEqual(eSlc, retSlc) {
		t.Errorf("expecting: %+v, received: %+v", eSlc, retSlc)
	}
	eSlc = []string{}
	if retSlc := StripSlicePrefix([]string{"0", "1", "2"}, 3); !reflect.DeepEqual(eSlc, retSlc) {
		t.Errorf("expecting: %+v, received: %+v", eSlc, retSlc)
	}
}

func TestSliceStringToIface(t *testing.T) {
	exp := []interface{}{"*default", "ToR", "*voice"}
	if rply := SliceStringToIface([]string{"*default", "ToR", "*voice"}); !reflect.DeepEqual(exp, rply) {
		t.Errorf("Expected: %s ,received: %s", ToJSON(exp), ToJSON(rply))
	}
}
