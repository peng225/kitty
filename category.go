package kitty

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const Identity = "identity"

type Object string
type MorphismID string

type Morphism struct {
	ID          MorphismID
	Source      Object
	Destination Object
}

type Category struct {
	Objects   map[Object]struct{}
	Morphisms map[MorphismID]*Morphism

	// key: f,g -> g∘f
	composeTable map[[2]MorphismID]MorphismID
}

func (o Object) GetIdentityID() MorphismID {
	return MorphismID("_id_" + o)
}

func (m *Morphism) IsIdentity() bool {
	return strings.HasPrefix(string(m.ID), "_id") &&
		m.Source == m.Destination
}

func (m *Morphism) Inverse() *Morphism {
	var inverseID MorphismID
	if strings.HasSuffix(string(m.ID), "^{-1}") {
		inverseID = MorphismID(strings.TrimSuffix(string(m.ID), "^{-1}"))
	} else {
		inverseID = m.ID + "^{-1}"
	}

	return &Morphism{
		ID:          inverseID,
		Source:      m.Destination,
		Destination: m.Source,
	}
}

func NewCategory(
	objs []Object, morphisms []*Morphism, compose map[[2]MorphismID]MorphismID,
) (*Category, error) {
	objects := make(map[Object]struct{})
	for _, o := range objs {
		objects[o] = struct{}{}
	}
	C := &Category{
		Objects:      objects,
		Morphisms:    make(map[MorphismID]*Morphism),
		composeTable: compose,
	}

	for _, m := range morphisms {
		if strings.Contains(string(m.ID), "_") {
			return nil, errors.New("Morphism names with '_' cannot be used")
		}
		C.Morphisms[m.ID] = m
	}

	for o := range C.Objects {
		id := o.GetIdentityID()
		m := &Morphism{
			ID:          id,
			Source:      o,
			Destination: o,
		}
		C.Morphisms[m.ID] = m
	}

	processedCompose := C.composeTable
	for k, v := range compose {
		if v == Identity {
			// Since k[1]◦k[0] is identity, its object should be the destination of k[1].
			processedCompose[k] = C.Morphisms[k[1]].Destination.GetIdentityID()
		}
	}

	for k, v := range C.composeTable {
		for _, w := range []MorphismID{k[0], k[1], v} {
			if _, ok := C.Morphisms[w]; !ok {
				return nil, fmt.Errorf("unknown morphism found in the composition definition: %s", w)
			}
		}
	}

	err := C.constructComposition()
	if err != nil {
		return nil, err
	}

	if err := C.validate(); err != nil {
		return nil, err
	}

	return C, nil
}

func (C *Category) constructComposition() error {
	initialMorphismCount := len(C.Morphisms)
	previousCTableCount := len(C.composeTable)
	for len(C.Morphisms) < initialMorphismCount*initialMorphismCount {
		toBeAddedMorphisms := make([]*Morphism, 0)
		for _, m1 := range C.Morphisms {
			for _, m2 := range C.Morphisms {
				if m1.Destination != m2.Source {
					continue
				}

				key := [2]MorphismID{m1.ID, m2.ID}
				if _, ok := C.composeTable[key]; ok {
					continue
				}
				switch {
				case m1.IsIdentity():
					// This case contains (m1.ID, m2.ID) = (identity, identity).
					composedID := m2.ID
					C.composeTable[key] = composedID
				case m2.IsIdentity():
					composedID := m1.ID
					C.composeTable[key] = composedID
				default:
					// Avoid the infinite morphism loop.
					// e.g. m1.ID = "g_f", m2.ID = "f"
					m1IDTokens := strings.Split(string(m1.ID), "_")
					m2IDTokens := strings.Split(string(m2.ID), "_")
					if containsAsSubsequence(m1IDTokens, m2IDTokens) ||
						containsAsSubsequence(m2IDTokens, m1IDTokens) {
						continue
					}
					// Eliminate the adjacent pair of inverse morphisms.
					// e.g. m1.ID = fg^{-1}", m2.ID = hgf^{-1} => m1.ID = id, m2.ID = h"
					for len(m1IDTokens) != 0 && len(m2IDTokens) != 0 {
						m1FirstMorphism := C.Morphisms[MorphismID(m1IDTokens[0])]
						m2LastMorphism := C.Morphisms[MorphismID(m2IDTokens[len(m2IDTokens)-1])]
						if m1FirstMorphism.Inverse().ID == m2LastMorphism.ID {
							m1IDTokens = m1IDTokens[1:]
							m2IDTokens = m2IDTokens[:len(m2IDTokens)-1]
							continue
						}
						break
					}
					var composedID MorphismID
					switch {
					case len(m1IDTokens) == 0 && len(m2IDTokens) == 0:
						composedID = m1.Source.GetIdentityID()
					case len(m1IDTokens) == 0:
						composedID = MorphismID(strings.Join(m2IDTokens, "_"))
					case len(m2IDTokens) == 0:
						composedID = MorphismID(strings.Join(m1IDTokens, "_"))
					default:
						composedID = MorphismID(fmt.Sprintf("%s_%s",
							strings.Join(m2IDTokens, "_"), strings.Join(m1IDTokens, "_")))
					}
					C.composeTable[key] = composedID
					toBeAddedMorphisms = append(toBeAddedMorphisms,
						&Morphism{
							ID:          composedID,
							Source:      m1.Source,
							Destination: m2.Destination,
						},
					)
				}
			}
		}
		for _, m := range toBeAddedMorphisms {
			C.Morphisms[m.ID] = m
		}
		if previousCTableCount == len(C.composeTable) {
			return nil
		}
		previousCTableCount = len(C.composeTable)
	}
	return errors.New("composition construction stuck detected")
}

func (C *Category) Compose(f, g MorphismID) (MorphismID, error) {
	key := [2]MorphismID{f, g}
	res, ok := C.composeTable[key]
	if !ok {
		return "", fmt.Errorf("composition not defined for %s and %s",
			f, g)
	}
	return res, nil
}

func containsAsSubsequence(tokens, subTokens []string) bool {
	if len(tokens) == 0 {
		return len(subTokens) == 0
	}
	if len(subTokens) == 0 {
		return true
	}
	i := slices.Index(tokens, subTokens[0])
	if i == -1 {
		return false
	}
	partialTokens := tokens[i:min(len(tokens), i+len(subTokens))]
	return slices.Compare(partialTokens, subTokens) == 0
}

func (C *Category) validate() error {
	// Type check
	for key, res := range C.composeTable {
		f, ok := C.Morphisms[key[0]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule",
				key[0])
		}
		g, ok := C.Morphisms[key[1]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule",
				key[1])
		}

		if f.Destination != g.Source {
			return errors.New("invalid composition domain")
		}

		r := C.Morphisms[res]
		if r.Source != f.Source || r.Destination != g.Destination {
			return errors.New("composition result mismatch")
		}
	}

	// Associative law
	for f := range C.Morphisms {
		for g := range C.Morphisms {
			for h := range C.Morphisms {
				gf, ok1 := C.Compose(f, g)
				hg, ok2 := C.Compose(g, h)
				if ok1 == nil && ok2 == nil {
					left, err := C.Compose(gf, h)
					if err != nil {
						continue
					}
					right, err := C.Compose(f, hg)
					if err != nil {
						continue
					}
					if left != right {
						return fmt.Errorf("associativity failed: left=%s, right=%s", left, right)
					}
				}
			}
		}
	}

	// Invalid source or destination
	for _, m := range C.Morphisms {
		if _, ok := C.Objects[m.Source]; !ok {
			return fmt.Errorf("%s has an invalid source: %s", m.ID, m.Source)
		}
		if _, ok := C.Objects[m.Destination]; !ok {
			return fmt.Errorf("%s has an invalid destination: %s", m.ID, m.Destination)
		}
	}

	return nil
}

func (C *Category) Hom(a, b Object) []MorphismID {
	var res []MorphismID

	for id, m := range C.Morphisms {
		if m.Source == a && m.Destination == b {
			res = append(res, id)
		}
	}

	return res
}

func (C *Category) IsIsomorphic(a, b Object) bool {
	for _, f := range C.Hom(a, b) {
		for _, g := range C.Hom(b, a) {
			gf, err := C.Compose(f, g)
			if err != nil {
				continue
			}
			fg, err := C.Compose(g, f)
			if err != nil {
				continue
			}

			if gf == a.GetIdentityID() &&
				fg == b.GetIdentityID() {
				return true
			}
		}
	}

	return false
}
