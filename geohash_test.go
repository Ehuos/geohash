package geohash

import (
	"math/rand"
	"testing"
)

var wrappingtests = []struct {
	lat1, lon1 float64
	lat2, lon2 float64
}{
	{0, -190, 0, 170},
	{-100, 0, 80, 0},
	{-460, 0, 80, 0},
	{440, 0, 80, 0},
}

func TestWrapping(t *testing.T) {
	//Tests that ranges outside latitude: (-90,90) and longitude (-180,180) generate the same geohash
	for i, tst := range wrappingtests {

		h1 := New(tst.lat1, tst.lon1, LatMaxPrecision, LonMaxPrecision)
		h2 := New(tst.lat2, tst.lon2, LatMaxPrecision, LonMaxPrecision)

		if h1 != h2 {
			t.Errorf("Fail:Test:%d h1=%v h2=%v", i, h1, h2)
		} else {
			t.Logf("Pass:Test:%d h1=%v h2=%v", i, h1, h2)
		}
	}
}

func TestRoundTrip(t *testing.T) {

	tolerance := []float64{0.1, 0.01, 0.001}
	for _, tol := range tolerance {
		for i := 0; i < 10; i++ {
			lat := -90.0 + 180.0*rand.Float64()
			lon := -180.0 + 360.0*rand.Float64()
			h := New(lat, lon, tol, tol)
			latMin, latDelta, lonMin, lonDelta := h.Decode()

			t.Logf("Coord:%g:%g:%g\n", lat, lon, tol)
			t.Logf("h:%s,%g,%g,%g,%g\n", h, latMin, lonMin, latDelta, lonDelta)

			if latDelta > tol {
				t.Errorf("Latitude tol (%g) > tolerance (%g)", latDelta, tol)
			}

			if lonDelta > tol {
				t.Errorf("Longitude tol (%g) > tolerance (%g)", lonDelta, tol)
			}
		}
	}
}

func TestNbrRoundTrip(t *testing.T) {

	for i := 0; i < 1; i++ {
		lat := -90.0 + 180.0*rand.Float64()
		lon := -180.0 + 360.0*rand.Float64()
		dir := Direction(rand.Int31n(4))
		rev := Direction((int(dir) + 2) % 4)
		tol := 0.01

		gh := New(lat, lon, tol, tol)

		nbr, err := gh.Nbr(dir)
		if err != nil {
			t.Errorf("Error generating nbr:%s", err)
		}

		nbrnbr, err := nbr.Nbr(rev)
		if err != nil {
			t.Errorf("Error generating nbr:%s", err)
		}

		if gh != nbrnbr {
			t.Errorf("Neighbour of neighbour not neighbour: Hash:%s Nbr:%s NbrNbr:%s", gh, nbr, nbrnbr)
		}
		t.Logf("Neighbour of neighbours: Hash:%s Nbr:%s NbrNbr:%s", gh, nbr, nbrnbr)

	}
}
