package main

import (
	"time"
)

type ConnectionType uint8
const (
	unknownConnectionType ConnectionType = iota
	offsetConnectionType
	attractConnectionType
)

type Connection struct {
	connectionType ConnectionType
	offset Vec2
	attractFactor float64
}

func NewOffsetConnection(offset Vec2) Connection {
	return Connection {
		connectionType: offsetConnectionType,
		offset: offset,
	}
}

func NewAttractConnection(attract float64) Connection {
	return Connection {
		connectionType: attractConnectionType,
		attractFactor: attract,
	}
}

func (c Connection) GetType() ConnectionType {
	return c.connectionType
}

func (c *Connection) SetOffset(offset Vec2) {
	c.offset = offset
}

func (c Connection) GetOffset() Vec2 {
	return c.offset
}

func (c Connection) GetAttractFactor() float64 {
	return c.attractFactor
}

type Attachment struct {
	connections map[SpacedId]Connection
	sid SpacedId
}

func NewAttachment(sid SpacedId) Attachment {
	return Attachment {
		connections: make(map[SpacedId]Connection),
		sid: sid,
	}
}

func (a *Attachment) AddConnection(parent SpacedId, connection Connection) {
	a.connections[parent] = connection
}

func (a *Attachment) UpdateState(grid *Grid, now time.Time) bool {
	for parentId, connection := range(a.GetConnections()) {
		parent := grid.Get(parentId)
		child := grid.Get(a.sid)

		if parent == nil || child == nil {
			a.RemoveConnection(parentId)

			if len(a.GetConnections()) == 0 {
				child.RemoveAttribute(attachedAttribute)
			}
			continue
		}

		if connection.GetType() == attractConnectionType {
			force := parent.Pos()
			force.Add(connection.GetOffset(), 1.0)
			force.Sub(child.Pos(), 1.0)

			smoothing := Min(1, force.Len())
			force.Normalize()
			force.Scale(smoothing * connection.GetAttractFactor())
			child.SetVel(force)
		}

		if !child.HasAttribute(attachedAttribute) {
			child.AddAttribute(attachedAttribute)
		}
	}

	return false
}

func (a *Attachment) Postprocess(grid *Grid, now time.Time) {
	for parentId, connection := range(a.GetConnections()) {
		parent := grid.Get(parentId)
		child := grid.Get(a.sid)

		if parent == nil || child == nil {
			a.RemoveConnection(parentId)

			if len(a.GetConnections()) == 0 {
				child.RemoveAttribute(attachedAttribute)
			}
			continue
		}

		if connection.GetType() == offsetConnectionType {
			pos := parent.Pos()
			pos.Add(connection.GetOffset(), 1.0)
			child.SetPos(pos)

			child.SetVel(parent.Vel())
			child.SetExtVel(parent.ExtVel())
			child.SetAcc(parent.Acc())
			child.SetJerk(parent.Jerk())
			grid.Upsert(child)
		}

		if !child.HasAttribute(attachedAttribute) {
			child.AddAttribute(attachedAttribute)
		}
	}
}

func (a *Attachment) RemoveConnection(parentId SpacedId) {
	delete(a.connections, parentId)
}

func (a Attachment) GetConnections() map[SpacedId]Connection {
	return a.connections
}