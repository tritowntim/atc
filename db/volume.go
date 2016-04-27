package db

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/concourse/atc"
)

type Volume struct {
	Handle     string
	WorkerName string
	TTL        time.Duration
	Identifier VolumeIdentifier
}

// pls gib algebraic data types
type VolumeIdentifier struct {
	ResourceCache *ResourceCacheIdentifier
	COW           *COWIdentifier
	Output        *OutputIdentifier
	Import        *ImportIdentifier
}

func (i VolumeIdentifier) Type() string {
	switch {
	case i.ResourceCache != nil:
		return "cache"
	case i.COW != nil:
		return "copy"
	case i.Output != nil:
		return "output"
	case i.Import != nil:
		return "import"
	default:
		return ""
	}
}

func (i VolumeIdentifier) String() string {
	switch {
	case i.ResourceCache != nil:
		return i.ResourceCache.String()
	case i.COW != nil:
		return i.COW.String()
	case i.Output != nil:
		return i.Output.String()
	case i.Import != nil:
		return i.Import.String()
	default:
		return ""
	}
}

type ResourceCacheIdentifier struct {
	ResourceVersion atc.Version
	ResourceHash    string
}

func (i ResourceCacheIdentifier) String() string {
	pairs := []string{}
	for k, v := range i.ResourceVersion {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k, v))
	}

	sort.Sort(sort.StringSlice(pairs))

	return strings.Join(pairs, ",")
}

type COWIdentifier struct {
	ParentVolumeHandle string
}

func (i COWIdentifier) String() string {
	return i.ParentVolumeHandle
}

type OutputIdentifier struct {
	Name string
}

func (i OutputIdentifier) String() string {
	return i.Name
}

type ImportIdentifier struct {
	WorkerName string
	Path       string
	Version    *string
}

func (i ImportIdentifier) String() string {
	return fmt.Sprintf("%s@%s", i.Path, *i.Version)
}

type SavedVolume struct {
	Volume

	ID        int
	ExpiresIn time.Duration
}