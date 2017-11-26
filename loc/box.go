package loc

type Box struct {
	P0, P1 Point
}

func (b *Box) Center() Point {
	var a Point
	for i := range a {
		a[i] = (b.P0[i] + b.P1[i]) / 2
	}

	return a
}

func (b *Box) transposeCoord(coord uint) Box {
	p0 := b.P0
	p1 := b.P1

	tmp := p0[coord]
	p0[coord] = p1[coord]
	p1[coord] = tmp

	return Box{P0: p0, P1: p1}
}

func (b *Box) TransposeX() Box {
	return b.transposeCoord(0)
}

func (b *Box) TransposeY() Box {
	return b.transposeCoord(1)
}

func (b *Box) TransposeZ() Box {
	return b.transposeCoord(2)
}

func (b *Box) issame(b1 Box) bool {
	return (b.P0.IsEqual(b1.P0) && b.P1.IsEqual(b1.P1)) || (b.P0.IsEqual(b1.P1) && b.P1.IsEqual(b1.P0))
}

/*IsEqual returns true is 2 boxes are equivalent */
func (b *Box) IsEqual(b1 Box) bool {
	if b.issame(b1) || b.issame(b1.TransposeX()) || b.issame(b1.TransposeY()) || b.issame(b1.TransposeZ()) {
		return true
	}

	return false
}

/*Bisect slices the box in 8 sub-boxes, cutting it in the center */
func (b *Box) Bisect() [8]Box {
	var p [8]Point
	p[0] = Point{b.P0[0], b.P0[1], b.P0[2]}
	p[1] = Point{b.P0[0], b.P1[1], b.P0[2]}
	p[2] = Point{b.P1[0], b.P1[1], b.P0[2]}
	p[3] = Point{b.P1[0], b.P0[1], b.P0[2]}
	p[4] = Point{b.P0[0], b.P0[1], b.P1[2]}
	p[5] = Point{b.P0[0], b.P1[1], b.P1[2]}
	p[6] = Point{b.P1[0], b.P1[1], b.P1[2]}
	p[7] = Point{b.P1[0], b.P0[1], b.P1[2]}

	c := b.Center()
	a := [8]Box{}
	for i := range a {
		a[i] = Box{p[i], c}
	}

	return a
}
