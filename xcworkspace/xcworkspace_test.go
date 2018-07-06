package xcworkspace

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	workspaceContentsPth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	workspacePth := filepath.Dir(workspaceContentsPth)
	parentDir := filepath.Dir(workspacePth)

	workspace, err := Open(workspacePth)
	require.NoError(t, err)

	require.Equal(t, filepath.Base(workspacePth), workspace.Name)
	require.Equal(t, workspacePth, workspace.Path)

	{
		require.Equal(t, 1, len(workspace.FileRefs))
		require.Equal(t, "group:XcodeProj.xcodeproj", workspace.FileRefs[0].Location)

		pth, err := workspace.FileRefs[0].AbsPath(parentDir)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(parentDir, "XcodeProj.xcodeproj"), pth)
	}

	{
		require.Equal(t, 1, len(workspace.Groups))

		group := workspace.Groups[0]
		require.Equal(t, "group:Group", group.Location)

		groupPth, err := group.AbsPath(parentDir)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(parentDir, "Group"), groupPth)

		require.Equal(t, 2, len(group.FileRefs))
		require.Equal(t, "group:../XcodeProj/AppDelegate.swift", group.FileRefs[0].Location)
		require.Equal(t, "group:XcodeProj.xcodeproj", group.FileRefs[1].Location)

		groupFileRef := group.FileRefs[0]
		groupFileRefPth, err := groupFileRef.AbsPath(groupPth)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(parentDir, "XcodeProj/AppDelegate.swift"), groupFileRefPth)
	}
}

const workspaceContentsContent = `<?xml version="1.0" encoding="UTF-8"?>
<Workspace
   version = "1.0">
   <Group
      location = "group:Group"
      name = "Group">
      <FileRef
         location = "group:../XcodeProj/AppDelegate.swift">
      </FileRef>
      <FileRef
         location = "group:XcodeProj.xcodeproj">
      </FileRef>
   </Group>
   <FileRef
      location = "group:XcodeProj.xcodeproj">
   </FileRef>
</Workspace>
`