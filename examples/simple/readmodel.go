// Copyright (c) 2014 - Max Ekman <max@looplab.se>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/looplab/eventhorizon"

	"github.com/looplab/eventhorizon/examples/domain"
)

// Invitation is a read model object for an invitation.
type Invitation struct {
	ID     eventhorizon.UUID
	Name   string
	Status string
}

// InvitationProjector is a projector that updates the invitations.
type InvitationProjector struct {
	repository eventhorizon.ReadRepository
}

// NewInvitationProjector creates a new InvitationProjector.
func NewInvitationProjector(repository eventhorizon.ReadRepository) *InvitationProjector {
	p := &InvitationProjector{
		repository: repository,
	}
	return p
}

// HandleEvent implements the HandleEvent method of the EventHandler interface.
func (p *InvitationProjector) HandleEvent(event eventhorizon.Event) {
	switch event := event.(type) {
	case *domain.InviteCreated:
		i := &Invitation{
			ID:   event.InvitationID,
			Name: event.Name,
		}
		p.repository.Save(i.ID, i)
	case *domain.InviteAccepted:
		m, _ := p.repository.Find(event.InvitationID)
		i := m.(*Invitation)
		i.Status = "accepted"
		p.repository.Save(i.ID, i)
	case *domain.InviteDeclined:
		m, _ := p.repository.Find(event.InvitationID)
		i := m.(*Invitation)
		i.Status = "declined"
		p.repository.Save(i.ID, i)
	}
}

// GuestList is a read model object for the guest list.
type GuestList struct {
	NumGuests   int
	NumAccepted int
	NumDeclined int
}

// GuestListProjector is a projector that updates the guest list.
type GuestListProjector struct {
	repository eventhorizon.ReadRepository
	eventID    eventhorizon.UUID
}

// NewGuestListProjector creates a new GuestListProjector.
func NewGuestListProjector(repository eventhorizon.ReadRepository, eventID eventhorizon.UUID) *GuestListProjector {
	p := &GuestListProjector{
		repository: repository,
		eventID:    eventID,
	}
	return p
}

// HandleEvent implements the HandleEvent method of the EventHandler interface.
func (p *GuestListProjector) HandleEvent(event eventhorizon.Event) {
	switch event.(type) {
	case *domain.InviteCreated:
		m, _ := p.repository.Find(p.eventID)
		if m == nil {
			m = &GuestList{}
		}
		g := m.(*GuestList)
		p.repository.Save(p.eventID, g)
	case *domain.InviteAccepted:
		m, _ := p.repository.Find(p.eventID)
		g := m.(*GuestList)
		g.NumAccepted++
		p.repository.Save(p.eventID, g)
	case *domain.InviteDeclined:
		m, _ := p.repository.Find(p.eventID)
		g := m.(*GuestList)
		g.NumDeclined++
		p.repository.Save(p.eventID, g)
	}
}
