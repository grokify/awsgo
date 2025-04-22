package inspector2util

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/govex/poam"
	"github.com/grokify/mogo/pointer"
)

type Package types.VulnerablePackage

func (p Package) NameAndFilepathAtVersion() string {
	name := p.NameString()
	fp := p.FilepathString()
	if name != "" && fp != "" {
		return fmt.Sprintf("%s (%s)", p.NameAtVersion(), p.FilepathAtVersion())
	} else if name != "" {
		return p.NameAtVersion()
	} else if fp != "" {
		return p.FilepathAtVersion()
	} else {
		return ""
	}
}

func (p Package) NameAndFilepathAtVersionFixed() string {
	name := p.NameString()
	fp := p.FilepathString()
	if name != "" && fp != "" {
		return fmt.Sprintf("%s (%s)", p.NameAtVersionFixed(), p.FilepathAtVersionFixed())
	} else if name != "" {
		return p.NameAtVersionFixed()
	} else if fp != "" {
		return p.FilepathAtVersionFixed()
	} else {
		return ""
	}
}

func (p Package) FilepathAtVersion() string {
	fp := p.FilepathString()
	ver := p.VersionString()
	if fp == "" && ver == "" {
		return ""
	} else {
		return fp + sepFilepathVersion + ver
	}
}

func (p Package) FilepathAtVersionFixed() string {
	fp := p.FilepathString()
	ver := p.VersionFixedString()
	if fp == "" && ver == "" {
		return ""
	} else {
		return fp + sepFilepathVersion + ver
	}
}

func (p Package) FilepathString() string {
	if p.FilePath == nil {
		return ""
	} else {
		return strings.TrimSpace(*p.FilePath)
	}
}

func (p Package) NameAtVersion() string {
	name := p.NameString()
	ver := p.VersionString()
	if name == "" && ver == "" {
		return ""
	} else {
		return name + sepFilepathVersion + ver
	}
}

func (p Package) NameAtVersionFixed() string {
	fp := p.NameString()
	ver := p.VersionFixedString()
	if fp == "" && ver == "" {
		return ""
	} else {
		return fp + sepFilepathVersion + ver
	}
}

func (p Package) NameString() string {
	if p.Name == nil {
		return ""
	} else {
		return strings.TrimSpace(*p.Name)
	}
}

func (p Package) POAMItem() poam.POAMItemUpgradeRemedationPackage {
	return poam.POAMItemUpgradeRemedationPackage{
		Name:           pointer.Dereference(p.Name),
		CurVersion:     pointer.Dereference(p.Version),
		FixVersion:     pointer.Dereference(p.FixedInVersion),
		PackageManager: string(p.PackageManager)}
}

func (p Package) VersionString() string {
	if p.Version == nil {
		return ""
	} else {
		return strings.TrimSpace(*p.Version)
	}
}

func (p Package) VersionFixedString() string {
	if p.FixedInVersion == nil {
		return ""
	} else {
		return strings.TrimSpace(*p.FixedInVersion)
	}
}
