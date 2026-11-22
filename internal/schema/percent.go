package schema

import (
	"bytes"
	"errors"
	"strconv"
)

// Percent is a half-up percentage to two decimal places, or null when the
// denominator is zero. In JSON it is a bare number such as 50, 41.67 or 0.05.
type Percent struct {
	ok         bool
	hundredths int
}

// PercentOf returns num/den as a Percent; den <= 0 or num < 0 gives null.
func PercentOf(num, den int) Percent {
	if den <= 0 || num < 0 {
		return Percent{}
	}
	scaled := num * 10000
	q := scaled / den
	r := scaled % den
	if r*2 >= den {
		q++
	}
	return Percent{ok: true, hundredths: q}
}

// Valid reports whether the percentage is defined (not null).
func (p Percent) Valid() bool { return p.ok }

// Hundredths returns the value in hundredths of a percent (4167 for 41.67).
// It is 0 when the percentage is null.
func (p Percent) Hundredths() int { return p.hundredths }

// Text returns the number as written in JSON, or "N/A" when null.
func (p Percent) Text() string {
	if !p.ok {
		return "N/A"
	}
	whole := p.hundredths / 100
	frac := p.hundredths % 100
	if frac == 0 {
		return strconv.Itoa(whole)
	}
	if frac%10 == 0 {
		return strconv.Itoa(whole) + "." + strconv.Itoa(frac/10)
	}
	if frac < 10 {
		return strconv.Itoa(whole) + ".0" + strconv.Itoa(frac)
	}
	return strconv.Itoa(whole) + "." + strconv.Itoa(frac/10) + strconv.Itoa(frac%10)
}

// MarshalJSON writes the bare number, or null.
func (p Percent) MarshalJSON() ([]byte, error) {
	if !p.ok {
		return []byte("null"), nil
	}
	return []byte(p.Text()), nil
}

var errPercent = errors.New("schema: percent must be null or a non-negative number with at most two decimals")

// UnmarshalJSON reads null or a non-negative number with at most two
// decimals, so that decoding then encoding reproduces the bytes.
func (p *Percent) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if string(b) == "null" {
		*p = Percent{}
		return nil
	}
	s := string(b)
	whole, frac := s, ""
	if i := bytes.IndexByte(b, '.'); i >= 0 {
		whole, frac = s[:i], s[i+1:]
		if frac == "" || len(frac) > 2 {
			return errPercent
		}
	}
	if whole == "" || (len(whole) > 1 && whole[0] == '0') || len(whole) > 7 {
		return errPercent
	}
	for _, c := range whole + frac {
		if c < '0' || c > '9' {
			return errPercent
		}
	}
	w, _ := strconv.Atoi(whole)
	h := w * 100
	switch len(frac) {
	case 1:
		d, _ := strconv.Atoi(frac)
		h += d * 10
	case 2:
		d, _ := strconv.Atoi(frac)
		h += d
	}
	np := Percent{ok: true, hundredths: h}
	if np.Text() != s {
		// Non-canonical forms such as 41.10 or 41.0 would not round-trip.
		return errPercent
	}
	*p = np
	return nil
}
