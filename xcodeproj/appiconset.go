package xcodeproj

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/xcode-project/serialized"
)

// TargetsToAppIconSets maps target names to an array of asset catalog paths
type TargetsToAppIconSets map[string][]string

// AppIconSetPaths parses an Xcode project and returns targets mapped to app icon set paths
func AppIconSetPaths(projectPath string) (TargetsToAppIconSets, error) {
	absPth, err := pathutil.AbsPath(projectPath)
	if err != nil {
		return TargetsToAppIconSets{}, err
	}

	objects, projectID, err := open(absPth)
	proj, err := parseProj(projectID, objects)
	if err != nil {
		return TargetsToAppIconSets{}, err
	}

	return appIconSetPaths(proj, projectPath, objects)
}

func appIconSetPaths(project Proj, projectPath string, objects serialized.Object) (TargetsToAppIconSets, error) {
	type iconTarget struct {
		target          Target
		appIconSetNames []string
	}

	targetsWithAppIconSetName := []iconTarget{}
	for _, target := range project.Targets {
		appIconSetNames, err := getAppIconSetNames(target)
		if err != nil {
			return nil, fmt.Errorf("app icon set name not found in project, error: %s", err)
		} else if len(appIconSetNames) == 0 {
			continue
		}
		targetsWithAppIconSetName = append(targetsWithAppIconSetName, iconTarget{
			target:          target,
			appIconSetNames: appIconSetNames,
		})
	}

	targetToAppIcons := map[string][]string{}
	for _, iconTarget := range targetsWithAppIconSetName {
		assetCatalogs, err := assetCatalogs(iconTarget.target, project.ID, objects)
		if err != nil {
			return nil, err
		} else if len(assetCatalogs) == 0 {
			continue
		}

		appIcons := []string{}
		for _, appIconSetName := range iconTarget.appIconSetNames {
			appIconSetPaths, err := lookupAppIconPaths(projectPath, assetCatalogs, appIconSetName, project.ID, objects)
			if err != nil {
				return nil, err
			} else if len(appIconSetPaths) == 0 {
				// continue
				return nil, fmt.Errorf("not found app icon set (%s) on paths: %s", appIconSetName, assetCatalogs)
			}
			appIcons = append(appIcons, appIconSetPaths...)
		}
		targetToAppIcons[iconTarget.target.ID] = appIcons
	}

	return targetToAppIcons, nil
}

func lookupAppIconPaths(projectPath string, assetCatalogs []fileReference, appIconSetName string, projectID string, objects serialized.Object) ([]string, error) {
	for _, fileReference := range assetCatalogs {
		resolvedPath, err := resolveObjectAbsolutePath(fileReference.id, projectID, projectPath, objects)
		if err != nil {
			return nil, err
		} else if resolvedPath == "" {
			return nil, fmt.Errorf("could not resolve path")
		}

		wildcharAppIconSetName := appIconSetName
		baseIconName := strings.Split(appIconSetName, "${")
		if len(baseIconName) > 1 {
			wildcharAppIconSetName = baseIconName[0] + "*"
		}

		matches, err := filepath.Glob(path.Join(regexp.QuoteMeta(resolvedPath), wildcharAppIconSetName+".appiconset"))
		if err != nil {
			return nil, err
		}
		return matches, nil
	}
	return nil, nil
}

func assetCatalogs(target Target, projectID string, objects serialized.Object) ([]fileReference, error) {
	if target.Type == NativeTargetType { // Ignoring PBXAggregateTarget and PBXLegacyTarget as may not contain buildPhases key
		resourcesBuildPhase, err := filterResourcesBuildPhase(target.buildPhaseIDs, objects)
		if err != nil {
			return nil, fmt.Errorf("getting resource build phases failed, error: %s", err)
		}
		assetCatalogs, err := filterAssetCatalogs(resourcesBuildPhase, projectID, objects)
		if err != nil {
			return nil, err
		}
		return assetCatalogs, nil
	}
	return nil, nil
}

func filterResourcesBuildPhase(buildPhases []string, objects serialized.Object) (resourcesBuildPhase, error) {
	for _, buildPhaseUUID := range buildPhases {
		rawBuildPhase, err := objects.Object(buildPhaseUUID)
		if err != nil {
			return resourcesBuildPhase{}, err
		}
		if isResourceBuildPhase(rawBuildPhase) {
			buildPhase, err := parseResourcesBuildPhase(buildPhaseUUID, objects)
			if err != nil {
				return resourcesBuildPhase{}, fmt.Errorf("failed to parse ResourcesBuildPhase, error: %s", err)
			}
			return buildPhase, nil
		}
	}
	return resourcesBuildPhase{}, fmt.Errorf("resource build phase not found")
}

func filterAssetCatalogs(buildPhase resourcesBuildPhase, projectID string, objects serialized.Object) ([]fileReference, error) {
	assetCatalogs := []fileReference{}
	for _, fileUUID := range buildPhase.files {
		buildFile, err := parseBuildFile(fileUUID, objects)
		if err != nil {
			// ignore:
			// D0177B971F26869C0044446D /* (null) in Resources */ = {isa = PBXBuildFile; };
			continue
		}

		// can be PBXVariantGroup or PBXFileReference
		rawElement, err := objects.Object(buildFile.fileRef)
		if err != nil {
			return nil, err
		}
		if ok, err := isFileReference(rawElement); err != nil {
			return nil, err
		} else if !ok {
			// ignore PBXVariantGroup
			continue
		}

		fileReference, err := parseFileReference(buildFile.fileRef, objects)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(fileReference.path, ".xcassets") {
			assetCatalogs = append(assetCatalogs, fileReference)
		}
	}
	return assetCatalogs, nil
}

func getAppIconSetNames(target Target) ([]string, error) {
	const appIconSetNameKey = "ASSETCATALOG_COMPILER_APPICON_NAME"

	found, defaultConfiguration := defaultConfiguration(target)
	if !found {
		return nil, fmt.Errorf("default configuration not found for target: %s", target)
	}

	appIconSetNameRaw, ok := defaultConfiguration.BuildSettings[appIconSetNameKey]
	if !ok {
		return nil, nil
	}

	appIconSetName, ok := appIconSetNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for value of key %s", appIconSetNameKey)
	}

	return []string{appIconSetName}, nil
}

func defaultConfiguration(target Target) (bool, BuildConfiguration) {
	defaultConfigurationName := target.BuildConfigurationList.DefaultConfigurationName
	for _, configuration := range target.BuildConfigurationList.BuildConfigurations {
		if configuration.Name == defaultConfigurationName {
			return true, configuration
		}
	}
	return false, BuildConfiguration{}
}