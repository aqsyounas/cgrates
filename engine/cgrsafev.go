/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOev.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package engine

import (
	"time"

	"github.com/cgrates/cgrates/utils"
)

func NewCGRSafEventFromCGREvent(cgrEv *utils.CGREvent) *CGRSafEvent {
	return &CGRSafEvent{
		Tenant: cgrEv.Tenant,
		ID:     cgrEv.ID,
		Time:   cgrEv.Time,
		Event:  NewSafEvent(cgrEv.Event),
	}
}

// CGRSafEvent is a safe CGREvent
type CGRSafEvent struct {
	Tenant string
	ID     string
	Time   *time.Time // event time
	Event  *SafEvent
}

func (cgrSafEv *CGRSafEvent) AsCGREvent() *utils.CGREvent {
	return &utils.CGREvent{
		Tenant: cgrSafEv.Tenant,
		ID:     cgrSafEv.ID,
		Time:   cgrSafEv.Time,
		Event:  cgrSafEv.Event.AsMapInterface(),
	}
}