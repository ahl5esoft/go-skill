package apisvc

import (
	"bytes"
	"testing"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("level = 1", func(t *testing.T) {
		mockIOFactory := contract.NewMockIIOFactory(ctrl)
		mockIOPath := contract.NewMockIIOPath(ctrl)
		workspace := "github.com"

		rootPath := "root-path"
		mockIOPath.EXPECT().GetRoot().Return(rootPath)

		mockApiDir := contract.NewMockIIODirectory(ctrl)
		mockIOFactory.EXPECT().BuildDirectory(rootPath, "api").Return(mockApiDir)

		mockApiDir.EXPECT().FindFiles()

		mockEndpointDir := contract.NewMockIIODirectory(ctrl)
		mockApiDir.EXPECT().FindDirectories().Return([]contract.IIODirectory{mockEndpointDir})

		mockApiFile := contract.NewMockIIOFile(ctrl)
		mockMetadataFile := contract.NewMockIIOFile(ctrl)
		mockTestFile := contract.NewMockIIOFile(ctrl)
		mockEndpointDir.EXPECT().FindFiles().Return([]contract.IIOFile{
			mockApiFile,
			mockMetadataFile,
			mockTestFile,
		})

		mockApiFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockApiFile.EXPECT().GetName().Return("api.go").AnyTimes()

		mockMetadataFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockMetadataFile.EXPECT().GetName().Return(metadataFilename).AnyTimes()

		mockTestFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockTestFile.EXPECT().GetName().Return("api_test.go").AnyTimes()

		mockApiFile.EXPECT().Read(
			gomock.Any(),
		).SetArg(0, `type TestApi struct`)

		mockEndpointDir.EXPECT().FindDirectories()

		mockEndpointDir.EXPECT().GetName().Return("endpoint").AnyTimes()

		mockEndpointDir.EXPECT().GetParent().Return(mockApiDir)

		mockApiDir.EXPECT().GetName().Return("api").AnyTimes()

		mockWorkspaceDir := contract.NewMockIIODirectory(ctrl)
		mockApiDir.EXPECT().GetParent().Return(mockWorkspaceDir)

		mockWorkspaceDir.EXPECT().GetName().Return(workspace)

		apiDirPath := "api-path"
		mockApiDir.EXPECT().GetPath().Return(apiDirPath)

		mockIOFactory.EXPECT().BuildFile(apiDirPath, metadataFilename).Return(mockMetadataFile)

		bf := bytes.NewBufferString(`package api

import (
    "github.com/ahl5esoft/go-skill/hide-ctx/contract"
    endpoint "github.com/api/endpoint"
)

func Register(apiFactory contract.IApiFactory) {
    apiFactory.Register("endpoint", "api", endpoint.TestApi{})
}`)
		mockMetadataFile.EXPECT().Write(*bf)

		err := GenerateMetadata(mockIOFactory, mockIOPath, workspace)
		assert.NoError(t, err)
	})

	t.Run("level > 1", func(t *testing.T) {
		mockIOFactory := contract.NewMockIIOFactory(ctrl)
		mockIOPath := contract.NewMockIIOPath(ctrl)
		workspace := "github.com"

		rootPath := "root-path"
		mockIOPath.EXPECT().GetRoot().Return(rootPath)

		mockApiDir := contract.NewMockIIODirectory(ctrl)
		mockIOFactory.EXPECT().BuildDirectory(rootPath, "api").Return(mockApiDir)

		mockApiDir.EXPECT().FindFiles()

		mockEndpointDir := contract.NewMockIIODirectory(ctrl)
		mockApiDir.EXPECT().FindDirectories().Return([]contract.IIODirectory{mockEndpointDir})

		mockApiFile := contract.NewMockIIOFile(ctrl)
		mockMetadataFile := contract.NewMockIIOFile(ctrl)
		mockTestFile := contract.NewMockIIOFile(ctrl)
		mockEndpointDir.EXPECT().FindFiles().Return([]contract.IIOFile{
			mockApiFile,
			mockMetadataFile,
			mockTestFile,
		})

		mockApiFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockApiFile.EXPECT().GetName().Return("api.go").AnyTimes()

		mockMetadataFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockMetadataFile.EXPECT().GetName().Return(metadataFilename).AnyTimes()

		mockTestFile.EXPECT().GetExt().Return(".go").AnyTimes()
		mockTestFile.EXPECT().GetName().Return("api_test.go").AnyTimes()

		mockApiFile.EXPECT().Read(
			gomock.Any(),
		).SetArg(0, `type TestApi struct`)

		mockEndpointDir.EXPECT().FindDirectories()

		mockEndpointDir.EXPECT().GetName().Return("endpoint").AnyTimes()

		mockModuleDir := contract.NewMockIIODirectory(ctrl)
		mockEndpointDir.EXPECT().GetParent().Return(mockModuleDir)

		mockModuleDir.EXPECT().GetName().Return("module").AnyTimes()

		mockModuleDir.EXPECT().GetParent().Return(mockApiDir)

		mockApiDir.EXPECT().GetName().Return("api").AnyTimes()

		mockWorkspaceDir := contract.NewMockIIODirectory(ctrl)
		mockApiDir.EXPECT().GetParent().Return(mockWorkspaceDir)

		mockWorkspaceDir.EXPECT().GetName().Return(workspace)

		apiDirPath := "api-path"
		mockApiDir.EXPECT().GetPath().Return(apiDirPath)

		mockIOFactory.EXPECT().BuildFile(apiDirPath, metadataFilename).Return(mockMetadataFile)

		bf := bytes.NewBufferString(`package api

import (
    "github.com/ahl5esoft/go-skill/hide-ctx/contract"
    moduleendpoint "github.com/api/module/endpoint"
)

func Register(apiFactory contract.IApiFactory) {
    apiFactory.Register("module/endpoint", "api", moduleendpoint.TestApi{})
}`)
		mockMetadataFile.EXPECT().Write(*bf)

		err := GenerateMetadata(mockIOFactory, mockIOPath, workspace)
		assert.NoError(t, err)
	})
}
