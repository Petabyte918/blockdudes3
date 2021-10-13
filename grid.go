package main

const (
	gridUnitLength int = 4
	gridUnitHeight int = 4
)

type GridCoord struct {
	x int
	y int
}

func (gc *GridCoord) advance(x int, y int) {
	gc.x += x * gridUnitLength
	gc.y += y * gridUnitHeight
}

type Grid struct {
	players map[int]*Player
	objects map[int]*Object

	grid map[GridCoord]map[int]*Object
	reverseGrid map[int][]GridCoord
}

func NewGrid() *Grid {
	return &Grid {
		players: make(map[int]*Player, 0),
		objects: make(map[int]*Object, 0),

		grid: make(map[GridCoord]map[int]*Object, 0),
		reverseGrid: make(map[int][]GridCoord, 0),
	}
}

func (g *Grid) addPlayer(id int, initData PlayerInitData) {
	g.players[id] = NewPlayer(id, initData)
}

func (g *Grid) hasPlayer(id int) bool {
	_, ok := g.players[id]
	return ok
}

func (g *Grid) deletePlayer(id int) {
	delete(g.players, id)
}

func (g *Grid) setPlayerData(id int, data PlayerData) {
	prof := g.players[id].Profile

	prof.SetPos(data.Pos)
	prof.SetVel(data.Vel)
	prof.SetAcc(data.Acc)
}

func (g* Grid) addObject(id int, object *Object) {
	g.objects[id] = object

	coords := g.getCoords(object.Profile)
	for _, coord := range(coords) {
		if _, ok := g.grid[coord]; !ok {
			g.grid[coord] = make(map[int]*Object)
		}
		g.grid[coord][id] = object
	}
	g.reverseGrid[id] = coords
}

func (g* Grid) updateObject(id int, object *Object) {
	if _, ok := g.objects[id]; !ok {
		return
	}

	if coords, ok := g.reverseGrid[id]; ok {
		for _, coord := range(coords) {
			delete(g.grid[coord], id)
		}
		delete(g.reverseGrid, id)
	}

	coords := g.getCoords(object.Profile)
	for _, coord := range(coords) {
		if _, ok := g.grid[coord]; !ok {
			g.grid[coord] = make(map[int]*Object)
		}
		g.grid[coord][id] = object
	}
	g.reverseGrid[id] = coords
}

func (g* Grid) setObjects(objects map[int]*Object) {
	g.objects = make(map[int]*Object, len(objects))
	g.grid = make(map[GridCoord]map[int]*Object, 0)
	g.reverseGrid = make(map[int][]GridCoord, len(objects))

	for id, object := range(objects) {
		g.addObject(id, object)
	}
}

func (g *Grid) getColliders(prof Profile) ObjectHeap {
	objectHeap := make(ObjectHeap, 0)

	for id, object := range(g.getObjectsNearProfile(prof)) {
		if prof.Overlap(object.Profile) {
			item := &ObjectItem {
				id: id,
				object: object,
			}
			objectHeap.Push(item)
			objectHeap.priority(item, prof.OverlapX(object.Profile) * prof.OverlapY(object.Profile))
		}
	}
	return objectHeap
}

func (g *Grid) getLineCollider(line Line) (bool, int, *Object, float64) {
	var collider *Object
	colliderId := -1
	closest := 1.0

	coord := g.getCoord(line.O)
	for {
		for id, object := range(g.grid[coord]) {
			collision, t := object.Profile.Intersects(line)
			if collision && t < closest {
				colliderId = id
				collider = object 
				closest = t
			}
		}

		if collider != nil {
			return true, colliderId, collider, closest
		}

		xstart := NewVec2(float64(coord.x), float64(coord.y))
		if Sign(line.R.X) > 0 {
			xstart.X += float64(gridUnitLength)
		}
		xline := NewLine(xstart, NewVec2(0, float64(gridUnitHeight)))

		ystart := NewVec2(float64(coord.x), float64(coord.y))
		if Sign(line.R.Y) > 0 {
			ystart.Y += float64(gridUnitHeight)
		}
		yline := NewLine(ystart, NewVec2(float64(gridUnitLength), 0))

		xcollide, xt := line.Intersects(xline)
		ycollide, yt := line.Intersects(yline)

		if !xcollide && !ycollide {
			break
		}
		if xcollide && (xt <= yt || !ycollide) {
			coord.advance(int(Sign(line.R.X)), 0)
		}
		if ycollide && (yt <= xt || !xcollide) {
			coord.advance(0, int(Sign(line.R.Y)))
		}
	}

	return false, colliderId, collider, closest
}

func (g* Grid) getCoord(point Vec2) GridCoord {
	cx := IntLeft(point.X)
	cy := IntLeft(point.Y)

	return GridCoord {
		x: cx - Mod(cx, gridUnitLength),
		y: cy - Mod(cy, gridUnitHeight),
	}
}

func (g* Grid) getCoords(prof Profile) []GridCoord {
	pos := prof.Pos()
	dim := prof.Dim()

	coords := make([]GridCoord, 0)

	xmin := pos.X - dim.X / 2
	xmax := pos.X + dim.X / 2
	ymin := pos.Y - dim.Y / 2
	ymax := pos.Y + dim.Y / 2

	cxmin := IntOut(xmin) - Mod(IntOut(xmin), gridUnitLength)
	cxmax := IntOut(xmax) - Mod(IntOut(xmax), gridUnitLength)
	cymin := IntOut(ymin) - Mod(IntOut(ymin), gridUnitHeight)
	cymax := IntOut(ymax) - Mod(IntOut(ymax), gridUnitHeight)

	for x := cxmin; x <= cxmax; x += gridUnitLength {
		for y := cymin; y <= cymax; y += gridUnitHeight {
			coords = append(coords, GridCoord{x : x, y: y})
		}
	}
	return coords
}

func (g *Grid) getObjectsNearProfile(prof Profile) map[int]*Object {
	objects := make(map[int]*Object)

	for _, coord := range(g.getCoords(prof)) {
		for id, object := range(g.grid[coord]) {
			objects[id] = object
		}
	}
	return objects
}