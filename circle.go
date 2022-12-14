package main

type Circle struct {
	BaseProfile
}

func NewCircle(init Init) *Circle {
	return &Circle {
		BaseProfile: NewBaseProfile(init),
	}
}

func (c Circle) Radius() float64 {
	return c.Dim().X / 2
}

func (c Circle) RadiusSqr() float64 {
	return c.Radius() * c.Radius()
}

func (c Circle) Contains(point Vec2) ContainResults {
	results := c.BaseProfile.Contains(point)
	if results.contains {
		return results
	}

	selfResults := NewContainResults()
	pos := c.Pos()
	distX := pos.X - point.X
	distY := pos.Y - point.Y
	selfResults.contains = distX * distX + distY * distY <= c.RadiusSqr()
	
	results.Merge(selfResults)
	return results
}

func (c Circle) Intersects(line Line) IntersectResults {
	results := c.BaseProfile.Intersects(line)

	// TODO: circle intersects line
	return results
}

func (c Circle) OverlapProfile(profile Profile) CollideResult {
	result := c.BaseProfile.OverlapProfile(profile)

	switch other := profile.(type) {
	case *RotPoly:
		result.Merge(other.OverlapProfile(&c))
	case *Rec2:
		result.Merge(other.OverlapProfile(&c))
	case *Circle:
		radius := c.Radius() + other.Radius()
		if c.DistSqr(other) <= radius * radius {
			result.SetHit(true)
		}
	}

	result.SetPosAdjustment(c.PosAdjustment(profile))
	return result
}