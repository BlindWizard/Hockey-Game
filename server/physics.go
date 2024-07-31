package main

import "math"

type Vector struct {
	x float64
	y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{x: x, y: y}
}

func NewVectorFromPoints(a *Vector, b *Vector) *Vector {
	return SubstractVectors(b, a)
}

func Sqr(x float64) float64 {
	return x * x
}

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func VectorLength(a *Vector) float64 {
	return math.Sqrt(a.x*a.x + a.y*a.y)
}

func SetVectorLength(a *Vector, l float64) *Vector {
	scale := l / VectorLength(a)
	return NewVector(a.x*scale, a.y*scale)
}

func MultiplyVectors(a *Vector, b *Vector) float64 {
	return a.x*b.x + a.y*b.y
}

func MultiplyVectorNumber(a *Vector, num float64) *Vector {
	return NewVector(a.x*num, a.y*num)
}

func SumVectors(a *Vector, b *Vector) *Vector {
	return NewVector(a.x+b.x, a.y+b.y)
}

func SubstractVectors(a *Vector, b *Vector) *Vector {
	return NewVector(a.x-b.x, a.y-b.y)
}

func Lerp(a float64, b float64, t float64) float64 {
	return a + (b-a)*t
}

func LerpVector(a *Vector, b *Vector, t float64) *Vector {
	return NewVector(Lerp(a.x, b.x, t), Lerp(a.y, b.y, t))
}

func NormalizeVector(a *Vector) *Vector {
	length := VectorLength(a)

	return NewVector(a.x/length, a.y/length)
}

func DegreeToRad(r float64) float64 {
	return r * math.Pi / 180
}

func RadToDegree(d float64) float64 {
	return d * 180 / math.Pi
}

func DistanceBetweenPoints(a *Vector, b *Vector) float64 {
	return math.Sqrt(Sqr(a.x-b.x) + Sqr(a.y-b.y))
}

func PointOnLine(lineStart *Vector, lineEnd *Vector, distance float64) *Vector {
	lineLength := DistanceBetweenPoints(lineStart, lineEnd)
	t := distance / lineLength

	return LerpVector(lineStart, lineEnd, t)
}

func AngleBetweenLines(a1 *Vector, b1 *Vector, a2 *Vector, b2 *Vector) float64 {
	op1 := NewVector(b1.x-a1.x, b1.y-a1.y)
	op2 := NewVector(b2.x-a2.x, b2.y-a2.y)

	cos := MultiplyVectors(op1, op2) / (VectorLength(op1) * VectorLength(op2))
	angle := RadToDegree(math.Acos(cos))
	if angle > 90 {
		angle = 180 - angle
	}

	return DegreeToRad(angle)
}

func LineFromPoints(lineStart *Vector, lineEnd *Vector) (float64, float64, float64) {
	a := lineStart.y - lineEnd.y
	b := lineEnd.x - lineStart.x
	c := lineEnd.x*lineStart.y - lineStart.x*lineEnd.y

	return a, b, -c
}

func PointLineDistance(point *Vector, lineStart *Vector, lineEnd *Vector) float64 {
	a, b, c := LineFromPoints(lineStart, lineEnd)

	return math.Abs((a*point.x + b*point.y + c)) / (math.Sqrt(Sqr(a) + Sqr(b)))
}

func PointBelongsSegment(lineStart, lineEnd, target *Vector) bool {
	return DistanceBetweenPoints(lineStart, target)+DistanceBetweenPoints(target, lineEnd)-DistanceBetweenPoints(lineStart, lineEnd) < EPSILON
}

func ClosestPoint(a *Vector, b *Vector, target *Vector) *Vector {
	if DistanceBetweenPoints(a, target) < DistanceBetweenPoints(b, target) {
		return a
	} else {
		return b
	}
}

func CheckLineCircleIntercection(lineStart *Vector, lineEnd *Vector, circlePos *Vector, circleRadius float64) (bool, []*Vector) {
	if DistanceBetweenPoints(lineStart, lineEnd) < EPSILON {
		return false, nil
	}

	a, b, c := LineFromPoints(lineStart, lineEnd)
	c = -c

	A := Sqr(a) + Sqr(b)
	B := 2*a*b*circlePos.y - 2*a*c - 2*Sqr(b)*circlePos.x
	C := Sqr(b)*Sqr(circlePos.x) + Sqr(b)*Sqr(circlePos.y) - 2*b*c*circlePos.y + Sqr(c) - Sqr(b)*Sqr(circleRadius)

	Discr := Sqr(B) - 4*A*C

	if math.Abs(b) < EPSILON {
		x1 := c / a

		if math.Abs(circlePos.x-x1) > circleRadius {
			return false, nil
		}

		if math.Abs((x1-circleRadius)-circlePos.x) < EPSILON || math.Abs((x1+circleRadius)-circlePos.x) < EPSILON {
			return true, []*Vector{NewVector(x1, circlePos.y)}
		}

		dx := math.Abs(x1 - circlePos.x)
		dy := math.Sqrt(Sqr(circleRadius) - Sqr(dx))

		return true, []*Vector{NewVector(x1, circlePos.y+dy), NewVector(x1, circlePos.y-dy)}
	} else if math.Abs(Discr) < EPSILON {
		x1 := -B / (2 * A)
		y1 := (c - a*x1) / b

		return true, []*Vector{NewVector(x1, y1)}
	} else if Discr < 0 {
		return false, nil
	} else {
		Discr = math.Sqrt(Discr)

		x1 := (-B + Discr) / (2 * A)
		y1 := (c - a*x1) / b

		x2 := (-B - Discr) / (2 * A)
		y2 := (c - a*x2) / b

		return true, []*Vector{NewVector(x1, y1), NewVector(x1, y2)}
	}
}

func CheckSegmentCircleIntercection(lineStart *Vector, lineEnd *Vector, circlePos *Vector, circleRadius float64) (bool, []*Vector) {
	collision, hitPoints := CheckLineCircleIntercection(lineStart, lineEnd, circlePos, circleRadius)
	if collision {
		validHitPoints := make([]*Vector, 0, 2)
		for _, hitPoint := range hitPoints {
			if PointBelongsSegment(lineStart, lineEnd, hitPoint) {
				validHitPoints = append(validHitPoints, hitPoint)
			}
		}

		if len(validHitPoints) > 0 {
			return true, validHitPoints
		}
	}

	return false, nil
}

func CheckLineLineIntercection(a *Vector, b *Vector, c *Vector, d *Vector) (bool, *Vector) {
	if DistanceBetweenPoints(a, b) < EPSILON || DistanceBetweenPoints(c, d) < EPSILON {
		return false, nil
	}

	denom := ((d.y-c.y)*(b.x-a.x) - (d.x-c.x)*(b.y-a.y))

	if denom == 0 {
		return false, nil
	}

	ua := ((d.x-c.x)*(a.y-c.y) - (d.y-c.y)*(a.x-c.x)) / denom

	return true, NewVector(a.x+ua*(b.x-a.x), a.y+ua*(b.y-a.y))
}

func CheckSegmentSegmentIntercection(a *Vector, b *Vector, c *Vector, d *Vector) (bool, *Vector) {
	if DistanceBetweenPoints(a, b) < EPSILON || DistanceBetweenPoints(c, d) < EPSILON {
		return false, nil
	}

	line1 := NewVector(b.x-a.x, b.y-a.y)
	line2 := NewVector(d.x-c.x, d.y-c.y)

	denom := line1.x*line2.y - line2.x*line1.y

	if denom == 0 {
		return false, nil
	}

	denomPositive := denom > 0

	ua := a.x - c.x
	ub := a.y - c.y

	sn := line1.x*ub - line1.y*ua

	if (sn < 0) == denomPositive {
		return false, nil
	}

	tn := line2.x*ub - line2.y*ua
	if (tn < 0) == denomPositive {
		return false, nil
	}

	if sn > denom == denomPositive || tn > denom == denomPositive {
		return false, nil
	}

	t := tn / denom

	return true, NewVector(a.x+(t*line1.x), a.y+(t*line1.y))
}
