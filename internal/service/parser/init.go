package parser

import (
	"github.com/sku4/ad-parser/internal/service/parser/kufar"
	"github.com/sku4/ad-parser/internal/service/parser/onliner"
	"github.com/sku4/ad-parser/internal/service/parser/realt"
)

var (
	codeProfiles = map[string]iProfile{}
)

func init() {
	profiles := []iProfile{
		kufar.New(),
		onliner.New(),
		realt.New(),
	}
	for _, profile := range profiles {
		codeProfiles[profile.GetCode()] = profile
	}
}
