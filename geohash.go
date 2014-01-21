// Package geohash provides an implementation of geohashes
// as described in: http://en.wikipedia.com/wiki/Geohash
package geohash

import (
	"fmt"
	"math"
	"strings"
)

// A GeoHash is a base32 encoded string.
// Allowed characters are: 0123456789bcdefghjkmnpqrstuvwxyz
type GeoHash string

// Direction is used to specify the neighbours of each cell
// as: NORTH, SOUTH, EAST or WEST.
type Direction int

// LatMaxPrecision and LonMaxPrecision are the minimum size of a
// latitude and longitude cell used when calculating a GeoHash.
// The constant values correspond to 6.3 cm at the equator,
// assuming an earth radius of 6371 km. Should be close enough for
// government work.
const (
	base32          = "0123456789bcdefghjkmnpqrstuvwxyz"
	nbrsNE          = "p0r21436x8zb9dcf5h7kjnmqesgutwvy"
	nbrsEE          = "bc01fg45238967deuvhjyznpkmstqrwx"
	nbrsSE          = "14365h7k9dcfesgujnmqp0r2twvyx8zb"
	nbrsWE          = "238967debc01fg45kmstqrwxuvhjyznp"
	nbrsNO          = "bc01fg45238967deuvhjyznpkmstqrwx"
	nbrsEO          = "p0r21436x8zb9dcf5h7kjnmqesgutwvy"
	nbrsSO          = "238967debc01fg45kmstqrwxuvhjyznp"
	nbrsWO          = "14365h7k9dcfesgujnmqp0r2twvyx8zb"
	bdrsWE          = "0145hjnp"
	bdrsEE          = "bcfguvyz"
	bdrsNE          = "prxz"
	bdrsSE          = "028b"
	bdrsWO          = "028b"
	bdrsEO          = "prxz"
	bdrsNO          = "bcfguvyz"
	bdrsSO          = "0145hjnp"
	LatMaxPrecision = 0.00000001
	LonMaxPrecision = 0.00000001
)

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

// Creates a GeoHash from a (lat)itude, a (lon)gitude.
// The GeoHash will define a region that is at least as small as
// latTol and lonTol, unless latTol or lonTol are smaller than
// LatMaxPrecision or LonMaxPrecision respectively in which case these values
// are used instead. lat and lon are projected into (-90,90) and (-180,180).
func New(lat, lon, latTol, lonTol float64) GeoHash {

	if latTol < LatMaxPrecision {
		latTol = LatMaxPrecision
	}
	if lonTol < LonMaxPrecision {
		lonTol = LonMaxPrecision
	}

	// Initial boundary values of geohash
	latMin := -90.0
	latDelta := 90.0
	lonMin := -180.0
	lonDelta := 180.0

	// Reduce lat and lon into (-90,90) and (-180,180) respectively
	lat = math.Mod(math.Mod(lat-90.0, 180.0)+180, 180) - 90.0
	lon = math.Mod(math.Mod(lon-180.0, 360.0)+360, 360) - 180.0

	var hash []byte
	ind := 0
	bit := 0x10

	// Decompose coordinate until reach required tolerance
	for {
		if lon > lonMin+lonDelta {
			ind += bit
			lonMin += lonDelta
		}
		bit >>= 1

		// Check if at end of character and produce character
		if bit == 0 {
			// Produce character
			hash = append(hash, base32[ind])

			// reset counters
			ind = 0
			bit = 0x10

			// Check for end conditions
			if latDelta <= latTol && lonDelta <= lonTol {
				break
			}
		}

		if lat > latMin+latDelta {
			ind += bit
			latMin += latDelta
		}
		bit >>= 1

		// Check if at end of character and produce character
		if bit == 0 {
			// Produce character
			hash = append(hash, base32[ind])

			// reset counters
			ind = 0
			bit = 0x10

			// Check for end conditions
			if latDelta <= latTol && lonDelta <= lonTol {
				break
			}
		}

		// Reduce tolerance
		lonDelta *= 0.5
		latDelta *= 0.5
	}
	return GeoHash(string(hash))
}

// Creates a GeoHash from a string.
// An error returned if the input string contains characters
// outside the legitimate set.
func NewFromString(hash string) (GeoHash, error) {
	var err error
	var gh GeoHash
	if isValid(hash) {
		gh = GeoHash(hash)
	} else {
		err = fmt.Errorf("Hash contains illegal characters")
	}
	return gh, err
}

// Returns the bounding box defined by the GeoHash. latMin and lonMin are the lower latitude
// and longitude bounds, latDelta and lonDelta the size of the box.
func (gh GeoHash) Decode() (latMin, latDelta, lonMin, lonDelta float64) {

	latMin = -90.0
	latDelta = 90.0
	lonMin = -180.0
	lonDelta = 180.0
	hash := string(gh)

	for i := 0; i < len(hash); i++ {

		if ind := strings.IndexByte(base32, hash[i]); ind >= 0 {

			if i%2 == 0 {
				// Even character - bits are (lon,lat,lon,lat,lon)
				if ind&0x10 == 0x10 {
					lonMin += lonDelta
				}
				lonDelta *= 0.5

				if ind&0x8 == 0x8 {
					latMin += latDelta
				}
				latDelta *= 0.5

				if ind&0x4 == 0x4 {
					lonMin += lonDelta
				}
				lonDelta *= 0.5

				if ind&0x2 == 0x2 {
					latMin += latDelta
				}
				latDelta *= 0.5

				if ind&0x1 == 0x1 {
					lonMin += lonDelta
				}
				lonDelta *= 0.5
			} else {
				// Even character - bits are (lat,lon,lat,lon,lat)
				if ind&0x10 == 0x10 {
					latMin += latDelta
				}
				latDelta *= 0.5

				if ind&0x8 == 0x8 {
					lonMin += lonDelta
				}
				lonDelta *= 0.5

				if ind&0x4 == 0x4 {
					latMin += latDelta
				}
				latDelta *= 0.5

				if ind&0x2 == 0x2 {
					lonMin += lonDelta
				}
				lonDelta *= 0.5

				if ind&0x1 == 0x1 {
					latMin += latDelta
				}
				latDelta *= 0.5
			}
		} else {
			// Geohash at this point should always be valid, hence this branch should never be executed
			// if geohash has been properly created
			break
		}
	}
	return latMin, latDelta, lonMin, lonDelta
}

func isValid(hash string) bool {

	for i := 0; i < len(hash); i++ {
		if strings.IndexByte(base32, hash[i]) == -1 {
			return false
		}
	}
	return true
}

// Nbr returns the GeoHash in the direction specified by the input
// Direction dir.
// Returns an error if the GeoHash contains invalid parameters
// or if an the Direction is invalid.
func (gh GeoHash) Nbr(dir Direction) (GeoHash, error) {

	var newgh GeoHash
	var err error
	if dir >= NORTH && dir <= WEST {
		if isValid(string(gh)) {
			newgh = GeoHash(string(nbr([]byte(string(gh)), dir)))
		} else {
			err = fmt.Errorf("Hash contains illegal characters")
		}
	} else {
		err = fmt.Errorf("Illegal input direction")
	}
	return newgh, err
}

func nbr(hash []byte, dir Direction) []byte {

	base := hash[0 : len(hash)-1]
	last := hash[len(hash)-1]
	even := len(hash)%2 == 0

	if even {
		if dir == NORTH {
			if strings.IndexByte(bdrsNE, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsNE, last)])
		} else if dir == EAST {
			if strings.IndexByte(bdrsEE, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsEE, last)])
		} else if dir == SOUTH {
			if strings.IndexByte(bdrsSE, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsSE, last)])
		} else if dir == WEST {
			if strings.IndexByte(bdrsWE, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsWE, last)])
		}
	} else {
		if dir == NORTH {
			if strings.IndexByte(bdrsNO, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsNO, last)])
		} else if dir == EAST {
			if strings.IndexByte(bdrsEO, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsEO, last)])
		} else if dir == SOUTH {
			if strings.IndexByte(bdrsSO, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsSO, last)])
		} else if dir == WEST {
			if strings.IndexByte(bdrsWO, last) != -1 {
				base = nbr(base, dir)
			}
			return append(base, base32[strings.IndexByte(nbrsWO, last)])
		}
	}
	return make([]byte, 0)
}
