package loc

type Box struct {
	P0, P1 Point
}

func (b *Box) Center() Point {
	x := (b.P0.X + b.P1.X) / 2
	y := (b.P0.Y + b.P1.Y) / 2
	z := (b.P0.Z + b.P1.Z) / 2
	return Point{X: x, Y: y, Z: z}
}

/*Slice slices the box in 8 sub-boxes, cutting it in the center */
func (b *Box) Slice() [8]Box {
	p0 := Point{X: b.P0.X, Y: b.P0.Y, Z: b.P0.Z}
	p1 := Point{X: b.P0.X, Y: b.P1.Y, Z: b.P0.Z}
	p2 := Point{X: b.P1.X, Y: b.P1.Y, Z: b.P0.Z}
	p3 := Point{X: b.P1.X, Y: b.P0.Y, Z: b.P0.Z}

	p4 := Point{X: b.P0.X, Y: b.P0.Y, Z: b.P1.Z}
	p5 := Point{X: b.P0.X, Y: b.P1.Y, Z: b.P1.Z}
	p6 := Point{X: b.P1.X, Y: b.P1.Y, Z: b.P1.Z}
	p7 := Point{X: b.P1.X, Y: b.P0.Y, Z: b.P1.Z}

	c := b.Center()
	a := [8]Box{}
	a[0] = Box{p0, c}
	a[1] = Box{p1, c}
	a[2] = Box{p2, c}
	a[3] = Box{p3, c}
	a[4] = Box{p4, c}
	a[5] = Box{p5, c}
	a[6] = Box{p6, c}
	a[7] = Box{p7, c}
	return a
}
