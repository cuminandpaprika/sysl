package roothandler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type folderStructure struct {
	folders, files []string
}

func buildFolderTest(folders, files []string) (fs afero.Fs, err error) {
	fs = afero.NewMemMapFs()
	var folder, file string

	for _, folder = range folders {
		folder, err = filepath.Abs(folder)
		if err != nil {
			return
		}

		err = fs.MkdirAll(folder, os.ModeTemporary)
		if err != nil {
			return
		}
	}

	for _, file = range files {
		file, err = filepath.Abs(file)
		if err != nil {
			return
		}

		_, err = fs.Create(file)
		if err != nil {
			return
		}
	}

	return
}

func TestRootHandler(t *testing.T) {

	successfulTest := folderStructure{
		folders: []string{
			"./SuccessfulTest/path/to/module",
			fmt.Sprintf("./SuccessfulTest/%s", syslRootMarker),
			"./SuccessfulTest/path/to/another/module",
			fmt.Sprintf("./SuccessfulTest/path/to/another/%s", syslRootMarker),
		},
		files: []string{
			"./SuccessfulTest/path/to/module/test.sysl",
			"./SuccessfulTest/test2.sysl",
			"./SuccessfulTest/path/to/another/module/test3.sysl",
		},
	}

	definedRootFlagUndefinedMarker := folderStructure{
		folders: []string{
			"./DefinedRootAndSyslRootUndefinedTest/path/to/module/",
		},
		files: []string{
			"./DefinedRootAndSyslRootUndefinedTest/path/to/module/test.sysl",
		},
	}

	definedRootFlagAndMarkerFound := folderStructure{
		folders: []string{
			"./DefinedRootAndSyslRootDefinedTest/path/to/module/",
			fmt.Sprintf("./DefinedRootAndSyslRootDefinedTest/path/%s", syslRootMarker),
		},
		files: []string{
			"./DefinedRootAndSyslRootDefinedTest/path/to/module/test.sysl",
		},
	}

	undefinedRoot := folderStructure{
		folders: []string{
			"./UndefinedRootAndUndefinedSyslRoot/",
		},
		files: []string{
			"./UndefinedRootAndUndefinedSyslRoot/test.sysl",
		},
	}

	currentDirectory, err := filepath.Abs(".")
	assert.NoError(t, err)

	absPathRelativeToCurrentDirectory := func(path string) string {
		absPath, err := filepath.Abs(path)
		assert.NoError(t, err)

		relPath, err := filepath.Rel(currentDirectory, absPath)
		assert.NoError(t, err)
		return relPath
	}

	tests := []struct {
		root, module, foundRoot, name string
		errAssert                     func(t *testing.T, err error)
		folders, files                []string
		rootIsFound                   bool
	}{
		{
			root:        "",
			name:        "Successful test: finding a root marker",
			module:      successfulTest.files[0],
			foundRoot:   "SuccessfulTest",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Successful test: finding a root marker in the same directory as the module",
			module:      successfulTest.files[1],
			foundRoot:   "SuccessfulTest",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Successful test: finding the closest root marker",
			module:      successfulTest.files[2],
			foundRoot:   "SuccessfulTest/path/to/another",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "DefinedRootAndSyslRootUndefinedTest/path/",
			name:        "Root flag is defined and root marker does not exist",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   "DefinedRootAndSyslRootUndefinedTest/path",
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        ".",
			name:        "Root with relative output",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   ".",
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        "/",
			name:        "Root with absolute output",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   absPathRelativeToCurrentDirectory("/"),
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        "~",
			name:        "Root flag and Root marker is defined",
			module:      definedRootFlagAndMarkerFound.files[0],
			foundRoot:   absPathRelativeToCurrentDirectory("~"),
			folders:     definedRootFlagAndMarkerFound.folders,
			files:       definedRootFlagAndMarkerFound.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Root is not defined",
			module:      undefinedRoot.files[0],
			foundRoot:   "UndefinedRootAndUndefinedSyslRoot",
			folders:     undefinedRoot.folders,
			files:       undefinedRoot.files,
			rootIsFound: false,
		},
	}

	for _, test := range tests {

		testRun := func(t *testing.T) {
			logger := logrus.StandardLogger()
			fs, err := buildFolderTest(test.folders, test.files)
			assert.NoError(t, err)

			rootHandler := NewRootHandler(test.root, test.module)
			err = rootHandler.HandleRoot(fs, logger)
			
			assert.NoError(t, err)
			assert.Equal(t, test.foundRoot, rootHandler.Root())
			assert.Equal(t, test.rootIsFound, rootHandler.RootIsFound())
		}

		t.Run(test.name, testRun)
	}
}