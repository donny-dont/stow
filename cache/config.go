package cache

import (
	"errors"
	"net/url"

	"github.com/graymeta/stow"
)

// Kind is the kind of Location this package provides
const Kind = "cache"

const (
	// ConfigOverlayKind is the Location type for the overlay.
	ConfigOverlayKind = "overlay_type"
	// ConfigOverlay is the configuration of the overlay.
	ConfigOverlay = "overlay"
	// ConfigBaseKind is the Location type for the base.
	ConfigBaseKind = "base_type"
	// ConfigBase is the configuration of the base.
	ConfigBase = "base"
)

func init() {
	validatefn := func(config stow.Config) error {
		overlayKind, ok := config.Config(ConfigOverlayKind)
		if !ok || overlayKind == "" {
			return errors.New("no overlay_type")
		}
		overlayConfig, ok := config.NestedConfig(ConfigOverlay)
		if !ok {
			return errors.New("no overlay config")
		}
		err := stow.Validate(overlayKind, overlayConfig)
		if err != nil {
			return err
		}

		baseKind, ok := config.Config(ConfigBaseKind)
		if !ok || baseKind == "" {
			return errors.New("no base_type")
		}
		baseConfig, ok := config.NestedConfig(ConfigBase)
		if !ok {
			return errors.New("no base config")
		}
		return stow.Validate(baseKind, baseConfig)
	}
	makefn := func(config stow.Config) (stow.Location, error) {
		overlayType, ok := config.Config(ConfigOverlayKind)
		if !ok || overlayType == "" {
			return nil, errors.New("no overlay_type")
		}
		overlayConfig, ok := config.NestedConfig(ConfigOverlay)
		if !ok {
			return nil, errors.New("no overlay config")
		}

		baseType, ok := config.Config(ConfigBaseKind)
		if !ok || overlayType == "" {
			return nil, errors.New("no base_type")
		}
		baseConfig, ok := config.NestedConfig(ConfigBase)
		if !ok {
			return nil, errors.New("no base config")
		}

		o, err := stow.Dial(overlayType, overlayConfig)
		if err != nil {
			return nil, err
		}
		b, err := stow.Dial(baseType, baseConfig)
		if err != nil {
			return nil, err
		}

		return &location{
			overlay: o,
			base:    b,
		}, nil
	}
	kindfn := func(u *url.URL) bool {
		return u.Scheme == Kind
	}

	stow.Register(Kind, makefn, kindfn, validatefn)
}
