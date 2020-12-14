package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tbellembois/gochimitheque/globals"
)

// oneGroupAtomCount returns a count of the atoms of the f formula as a map.
// f must be a formula like (XYZ) (XYZ)n or [XYZ] [XYZ]n.
// example:
// (CH3)2 will return "C":2, "H":6
// CH3CH(NO2)CH3 will return "N":1 "O":2
// CH3CH(NO2)(CH3)2 will return "N":1 "O":2 - process only the first match
func oneGroupAtomCount(f string) map[string]int {
	var (
		// the result map
		c = make(map[string]int)
		r = globals.OneGroupMolRe
	)
	// Looking for non matching molecules.
	if !r.MatchString(f) {
		return nil
	}

	// sl is a list of 3 elements like
	// [[(CH3Na6CCl5H)2 CH3Na6CCl5H 2]]
	sl := r.FindAllStringSubmatch(f, -1)
	basicMol := sl[0][1]
	multiplier, _ := strconv.Atoi(sl[0][2])

	// if there is no multiplier
	if multiplier == 0 {
		multiplier = 1
	}

	// counting the atoms
	aCount := basicAtomCount(basicMol)
	for at, nb := range aCount {
		c[at] = nb * multiplier
	}

	return c
}

// basicAtomCount returns a count of the atoms of the f formula as a map.
// f must be a basic formula with only atoms and numbers.
// example:
// C6H5COC6H4CO2H will return "C1":4, "H":10, "O":3
// CH3CH(NO2)CH3 will return Nil, parenthesis are not allowed
func basicAtomCount(f string) map[string]int {
	var (
		// the result map
		c   = make(map[string]int)
		r   = globals.BasicMolRe
		err error
	)
	// Looking for non matching molecules.
	if !r.MatchString(f) {
		return nil
	}

	// sl is a slice like [[Na Na ] [Cl Cl ] [C2 C 2] [Cl3 Cl 3]]
	// for f = NaClC2Cl3
	// [ matchingString capture1 capture2 ]
	// capture1 is the atom
	// capture2 is the its number
	sl := r.FindAllStringSubmatch(f, -1)
	for _, i := range sl {
		atom := i[1]
		var nbAtom int
		if i[2] != "" {
			nbAtom, err = strconv.Atoi(i[2])
			if err != nil {
				return nil
			}
		} else {
			nbAtom = 1
		}
		if _, ok := c[atom]; ok {
			c[atom] = c[atom] + nbAtom
		} else {
			c[atom] = nbAtom
		}
	}
	return c
}

// LinearToEmpiricalFormula returns the empirical formula from the linear formula f.
// example: [(CH3)2SiH]2NH
//          (CH3)2C[C6H2(Br)2OH]2
func LinearToEmpiricalFormula(f string) string {
	var ef string

	s := "-"
	nf := ""

	// Finding the first (XYZ)n match
	reg := globals.OneGroupMolRe

	for s != "" {
		s = reg.FindString(f)

		// Counting the atoms and rebuilding the molecule string
		m := oneGroupAtomCount(s)
		ms := "" // molecule string
		for k, v := range m {
			ms += k
			if v != 1 {
				ms += fmt.Sprintf("%d", v)
			}
		}

		// Then replacing the match with the molecule string - nf is for "new f"
		nf = strings.Replace(f, s, ms, 1)
		f = nf
	}

	// Counting the atoms
	bAc := basicAtomCount(nf)

	// Sorting the atoms
	// C, H and then in alphabetical order
	var ats []string // atoms
	hasC := false    // C atom present
	hasH := false    // H atom present

	for k := range bAc {
		switch k {
		case "C":
			hasC = true
		case "H":
			hasH = true
		default:
			ats = append(ats, k)
		}
	}
	sort.Strings(ats)

	if hasH {
		ats = append([]string{"H"}, ats...)
	}
	if hasC {
		ats = append([]string{"C"}, ats...)
	}

	for _, at := range ats {
		ef += at
		nb := bAc[at]
		if nb != 1 {
			ef += fmt.Sprintf("%d", nb)
		}
	}

	return ef
}
