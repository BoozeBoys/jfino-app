package loc

type Box [2]Point

func NewBox(center, size Point) Box {
	var p0, p1 Point
	for i := range center {
		p0[i] = center[i] - size[i]
		p1[i] = center[i] + size[i]
	}
	return Box{p0, p1}
}

func (b *Box) Center() Point {
	var a Point
	for i := range a {
		a[i] = (b[0][i] + b[1][i]) / 2
	}

	return a
}

func (b *Box) Expand(scale float64) Box {
	var size Point
	c := b.Center()
	for i := range b[0] {
		size[i] = (b[0][i] - c[i]) * scale
	}
	return NewBox(c, size)
}

func (b *Box) transposeCoord(coord uint) Box {
	p0 := b[0]
	p1 := b[1]

	tmp := p0[coord]
	p0[coord] = p1[coord]
	p1[coord] = tmp

	return Box{p0, p1}
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
	return (b[0].IsEqual(b1[0]) && b[1].IsEqual(b1[1])) || (b[0].IsEqual(b1[1]) && b[1].IsEqual(b1[0]))
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
	p[0] = Point{b[0][0], b[0][1], b[0][2]}
	p[1] = Point{b[0][0], b[1][1], b[0][2]}
	p[2] = Point{b[1][0], b[1][1], b[0][2]}
	p[3] = Point{b[1][0], b[0][1], b[0][2]}
	p[4] = Point{b[0][0], b[0][1], b[1][2]}
	p[5] = Point{b[0][0], b[1][1], b[1][2]}
	p[6] = Point{b[1][0], b[1][1], b[1][2]}
	p[7] = Point{b[1][0], b[0][1], b[1][2]}

	c := b.Center()
	a := [8]Box{}
	for i := range a {
		a[i] = Box{p[i], c}
	}

	return a
}
