package profile

import (
	"context"
)

var (
	profilesIDs   map[uint16]string
	profilesCodes map[string]uint16
	profiles      = []*Profile{
		{1, "kufar"},
		{2, "onliner"},
		{3, "realt"},
	}
)

func GetByCode(ctx context.Context, code string) uint16 {
	if i, ok := profilesCodes[code]; ok {
		return i
	}
	fill(ctx)
	if i, ok := profilesCodes[code]; ok {
		return i
	}

	return 0
}

func GetByID(ctx context.Context, id uint16) string {
	if s, ok := profilesIDs[id]; ok {
		return s
	}
	fill(ctx)
	if s, ok := profilesIDs[id]; ok {
		return s
	}

	return ""
}

func fill(ctx context.Context) {
	_ = ctx

	profilesIDs = make(map[uint16]string, len(profiles))
	profilesCodes = make(map[string]uint16, len(profiles))

	for _, p := range profiles {
		profilesIDs[p.ID] = p.Code
		profilesCodes[p.Code] = p.ID
	}
}
