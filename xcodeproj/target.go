package xcodeproj

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-tools/xcode-project/serialized"
)

// TargetType ...
type TargetType string

// TargetTypes
const (
	NativeTargetType    TargetType = "PBXNativeTarget"
	AggregateTargetType TargetType = "PBXAggregateTarget"
	LegacyTargetType    TargetType = "PBXLegacyTarget"
)

// Target ...
type Target struct {
	Type                   TargetType
	ID                     string
	Name                   string
	BuildConfigurationList ConfigurationList
	Dependencies           []TargetDependency
	ProductReference       ProductReference
}

// DependentTargets ...
func (t Target) DependentTargets() []Target {
	var targets []Target
	for _, targetDependency := range t.Dependencies {
		childTarget := targetDependency.Target
		targets = append(targets, childTarget)

		childDependentTargets := childTarget.DependentTargets()
		targets = append(targets, childDependentTargets...)
	}

	return targets
}

// DependentExecutableProductTargets ...
func (t Target) DependentExecutableProductTargets() []Target {
	var targets []Target
	for _, targetDependency := range t.Dependencies {
		childTarget := targetDependency.Target
		if !childTarget.IsExecutableProduct() {
			continue
		}

		targets = append(targets, childTarget)

		childDependentTargets := childTarget.DependentTargets()
		targets = append(targets, childDependentTargets...)
	}

	return targets
}

// IsAppProduct ...
func (t Target) IsAppProduct() bool {
	return filepath.Ext(t.ProductReference.Path) == ".app"
}

// IsAppExtensionProduct ...
func (t Target) IsAppExtensionProduct() bool {
	return filepath.Ext(t.ProductReference.Path) == ".appex"
}

// IsExecutableProduct ...
func (t Target) IsExecutableProduct() bool {
	return t.IsAppProduct() || t.IsAppExtensionProduct()
}

func parseTarget(id string, objects serialized.Object) (Target, error) {
	rawTarget, err := objects.Object(id)
	if err != nil {
		return Target{}, err
	}

	isa, err := rawTarget.String("isa")
	if err != nil {
		return Target{}, err
	}

	var targetType TargetType
	switch isa {
	case "PBXNativeTarget":
		targetType = NativeTargetType
	case "PBXAggregateTarget":
		targetType = AggregateTargetType
	case "PBXLegacyTarget":
		targetType = LegacyTargetType
	default:
		return Target{}, fmt.Errorf("unknown target type: %s", isa)
	}

	name, err := rawTarget.String("name")
	if err != nil {
		return Target{}, err
	}

	buildConfigurationListID, err := rawTarget.String("buildConfigurationList")
	if err != nil {
		return Target{}, err
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return Target{}, err
	}

	dependencyIDs, err := rawTarget.StringSlice("dependencies")
	if err != nil {
		return Target{}, err
	}

	var dependencies []TargetDependency
	for _, dependencyID := range dependencyIDs {
		dependency, err := parseTargetDependency(dependencyID, objects)
		if err != nil {
			// KeyNotFoundError can be only raised if the 'target' property not found on the raw target dependency object
			// we only care about target dependency, which points to a target
			if serialized.IsKeyNotFoundError(err) {
				continue
			} else {
				return Target{}, err
			}
		}

		dependencies = append(dependencies, dependency)
	}

	var productReference ProductReference
	productReferenceID, err := rawTarget.String("productReference")
	if err != nil {
		if !serialized.IsKeyNotFoundError(err) {
			return Target{}, err
		}
	} else {
		productReference, err = parseProductReference(productReferenceID, objects)
		if err != nil {
			return Target{}, err
		}
	}

	return Target{
		Type: targetType,
		ID:   id,
		Name: name,
		BuildConfigurationList: buildConfigurationList,
		Dependencies:           dependencies,
		ProductReference:       productReference,
	}, nil
}
