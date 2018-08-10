package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseTargetDependency(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawTargetDependency), &raw)
	require.NoError(t, err)

	targetDependency, err := parseTargetDependency("13E76E511F4AC94F0028096E", raw)
	require.NoError(t, err)
	// fmt.Printf("targetDependency:\n%s\n", pretty.Object(targetDependency))
	require.Equal(t, expectedTargetDependency, pretty.Object(targetDependency))
}

const rawTargetDependency = `
{
	13E76E511F4AC94F0028096E /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		target = 13E76E461F4AC94F0028096E /* share-extension */;
		targetProxy = 13E76E501F4AC94F0028096E /* PBXContainerItemProxy */;
	};

	13E76E461F4AC94F0028096E /* share-extension */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "share-extension" */;
		buildPhases = (
			13E76E431F4AC94F0028096E /* Sources */,
			13E76E441F4AC94F0028096E /* Frameworks */,
			13E76E451F4AC94F0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = (
		);
		name = "share-extension";
		productName = "share-extension";
		productReference = 13E76E471F4AC94F0028096E /* share-extension.appex */;
		productType = "com.apple.product-type.app-extension";
	};

	13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E0E1F4AC90A0028096E /* code-sign-test.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "code-sign-test.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E111F4AC90A0028096E /* AppDelegate.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = AppDelegate.h; sourceTree = "<group>"; };
	13E76E121F4AC90A0028096E /* AppDelegate.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = AppDelegate.m; sourceTree = "<group>"; };
	13E76E141F4AC90A0028096E /* ViewController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ViewController.h; sourceTree = "<group>"; };
	13E76E151F4AC90A0028096E /* ViewController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ViewController.m; sourceTree = "<group>"; };
	13E76E181F4AC90A0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Main.storyboard; sourceTree = "<group>"; };
	13E76E1A1F4AC90A0028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E1D1F4AC90A0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/LaunchScreen.storyboard; sourceTree = "<group>"; };
	13E76E1F1F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E201F4AC90A0028096E /* main.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = main.m; sourceTree = "<group>"; };
	13E76E261F4AC90A0028096E /* code-sign-testTests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "code-sign-testTests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E2A1F4AC90A0028096E /* code_sign_testTests.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = code_sign_testTests.m; sourceTree = "<group>"; };
	13E76E2C1F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E311F4AC90A0028096E /* code-sign-testUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "code-sign-testUITests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E351F4AC90A0028096E /* code_sign_testUITests.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = code_sign_testUITests.m; sourceTree = "<group>"; };
	13E76E371F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E471F4AC94F0028096E /* share-extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "share-extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E491F4AC94F0028096E /* ShareViewController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ShareViewController.h; sourceTree = "<group>"; };
	13E76E4A1F4AC94F0028096E /* ShareViewController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ShareViewController.m; sourceTree = "<group>"; };
	13E76E4D1F4AC94F0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/MainInterface.storyboard; sourceTree = "<group>"; };
	13E76E4F1F4AC94F0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E591F4AC9800028096E /* watchkit-app.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "watchkit-app.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E5C1F4AC9800028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Interface.storyboard; sourceTree = "<group>"; };
	13E76E5E1F4AC9800028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E601F4AC9800028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E651F4AC9800028096E /* watchkit-app Extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "watchkit-app Extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E6A1F4AC9800028096E /* InterfaceController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = InterfaceController.h; sourceTree = "<group>"; };
	13E76E6B1F4AC9800028096E /* InterfaceController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = InterfaceController.m; sourceTree = "<group>"; };
	13E76E6D1F4AC9800028096E /* ExtensionDelegate.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ExtensionDelegate.h; sourceTree = "<group>"; };
	13E76E6E1F4AC9800028096E /* ExtensionDelegate.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ExtensionDelegate.m; sourceTree = "<group>"; };
	13E76E701F4AC9800028096E /* NotificationController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = NotificationController.h; sourceTree = "<group>"; };
	13E76E711F4AC9800028096E /* NotificationController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = NotificationController.m; sourceTree = "<group>"; };
	13E76E731F4AC9800028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E751F4AC9800028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E761F4AC9800028096E /* PushNotificationPayload.apns */ = {isa = PBXFileReference; lastKnownFileType = text; path = PushNotificationPayload.apns; sourceTree = "<group>"; };
}`

const expectedTargetDependency = `{
	"ID": "13E76E511F4AC94F0028096E",
	"Target": {
		"Type": "PBXNativeTarget",
		"ID": "13E76E461F4AC94F0028096E",
		"Name": "share-extension",
		"BuildConfigurationList": {
			"ID": "13E76E3A1F4AC90A0028096E",
			"DefaultConfigurationName": "Release",
			"BuildConfigurations": null
		},
		"Dependencies": null,
		"ProductReference": {
			"Path": "share-extension.appex"
		}
	}
}`
