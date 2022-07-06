package main

type ProfileMath interface {
	Contains(point Vec2) ContainResults
	Intersects(line Line) IntersectResults
	// TODO: this should take object
	Overlap(profile Profile) OverlapResults
	Snap(colliders ObjectHeap) SnapResults

	getIgnored() map[SpacedId]bool
	resetIgnored()
	addIgnored(sid SpacedId) 
}

type ContainResults struct {
	contains bool
	ignored bool
}

func NewContainResults() ContainResults {
	return ContainResults {
		contains: false,
		ignored: false,
	}
}

func (cr *ContainResults) Merge(other ContainResults) {
	cr.contains = cr.contains || other.contains
}

type OverlapResults struct {
	overlap bool

	// TODO: use a map?
	amount Vec2
}

func NewOverlapResults() OverlapResults {
	return OverlapResults {
		overlap: false,
		amount: NewVec2(0, 0),
	}
}

func (or *OverlapResults) Merge(other OverlapResults) {
	or.overlap = or.overlap || other.overlap
	if (other.amount.Area() > or.amount.Area()) {
		or.amount = other.amount
	}
}

type IntersectResults struct {
	hit bool
	ignored bool
	t float64
}

func NewIntersectResults() IntersectResults {
	return IntersectResults {
		hit: false,
		ignored: false,
		t: 1.0,
	}
}

func (ir *IntersectResults) Merge(other IntersectResults) {
	ir.hit = ir.hit || other.hit
	ir.t = Min(ir.t, other.t)
}

type SnapResults struct {
	snap bool
	ignored bool
	posAdj Vec2
	extVel Vec2

	// TODO: need IDs
}

func NewSnapResults() SnapResults {
	return SnapResults {
		snap: false,
		ignored: false,
		posAdj: NewVec2(0, 0),
		extVel: NewVec2(0, 0),
	}
}

func (sr *SnapResults) Merge(other SnapResults) {
	sr.snap = sr.snap || other.snap

	pos := &sr.posAdj
	otherPos := other.posAdj
	pos.X = AbsMax(pos.X, otherPos.X)
	pos.Y = AbsMax(pos.Y, otherPos.Y)

	extVel := &sr.extVel
	extVel.X = AbsMax(extVel.X, other.extVel.X)
	extVel.Y = AbsMax(extVel.Y, other.extVel.Y)
}