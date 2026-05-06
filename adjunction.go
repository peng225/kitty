package kitty

import (
	"errors"
	"fmt"
)

var (
	ErrorTriangleIdentityViolation = errors.New("triangle identity violation.")
)

func CheckAdjunction(F, G *Functor, eta, epsilon *NaturalTransformation) error {
	// --- triangle 1: εF ∘ Fη = id_F ---
	FEta, err := eta.ApplyFunctorLeft(F)
	if err != nil {
		return err
	}
	epsilonF, err := epsilon.ApplyFunctorRight(F)
	if err != nil {
		return err
	}
	composedNat, err := FEta.ComposeWith(epsilonF)
	if err != nil {
		return err
	}

	idF, err := IdentityNaturalTransformation(F)
	if err != nil {
		return err
	}

	for obj := range composedNat.Components {
		if composedNat.Components[obj] != idF.Components[obj] {
			return fmt.Errorf("triangle identity 1: %w", ErrorTriangleIdentityViolation)
		}
	}

	// --- triangle 2: Gε ∘ ηG = id_G ---
	GEpsilon, err := epsilon.ApplyFunctorLeft(G)
	if err != nil {
		return err
	}
	etaG, err := eta.ApplyFunctorRight(G)
	if err != nil {
		return err
	}
	composedNat, err = etaG.ComposeWith(GEpsilon)
	if err != nil {
		return err
	}

	idG, err := IdentityNaturalTransformation(G)
	if err != nil {
		return err
	}

	for obj := range composedNat.Components {
		if composedNat.Components[obj] != idG.Components[obj] {
			return fmt.Errorf("triangle identity 2: %w", ErrorTriangleIdentityViolation)
		}
	}

	return nil
}
