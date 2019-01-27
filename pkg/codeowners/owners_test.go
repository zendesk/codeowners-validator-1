package codeowners_test

import (
	"fmt"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

const sampleCodeownerFile = `
# Sample codeowner file
*	@everyone

src/**	@org/hakuna-matata @pico-bello
pkg/github.com/**	@myk

`

func TestParseCodeownersSuccess(t *testing.T) {
	// given
	givenCodeownerPath := "workspace/go/repo-name"
	expEntries := []codeowners.Entry{
		{3, "*", []string{"@everyone"}},
		{5, "src/**", []string{"@org/hakuna-matata", "@pico-bello"}},
		{6, "pkg/github.com/**", []string{"@myk"}},
	}

	tFS := afero.NewMemMapFs()
	revert := codeowners.SetFS(tFS)
	defer revert()

	f, _ := tFS.Create(path.Join(givenCodeownerPath, "CODEOWNERS"))
	f.WriteString(sampleCodeownerFile)

	// when
	entries, err := codeowners.NewFromPath(givenCodeownerPath)

	// then
	require.NoError(t, err)
	assert.Len(t, entries, len(expEntries))
	for _, expEntry := range expEntries {
		assert.Contains(t, entries, expEntry)
	}
}

func TestFindCodeownersFileSuccess(t *testing.T) {
	tests := map[string]struct {
		basePath string
	}{
		"Should find the CODEOWNERS at root": {
			basePath: "/workspace/go/repo-name1/",
		},
		"Should find the CODEOWNERS in docs/": {
			basePath: "/workspace/go/repo-name2/docs/",
		},
		"Should find the CODEOWNERS IN .github": {
			basePath: "/workspace/go/repo-name3/.github/",
		},
		"Should manage situation without trailing slash": {
			basePath: "/workspace/go/repo-name3/.github",
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			tFS := afero.NewMemMapFs()
			revert := codeowners.SetFS(tFS)
			defer revert()

			tFS.Create(path.Join(tc.basePath, "CODEOWNERS"))
			// when
			_, err := codeowners.NewFromPath(tc.basePath)

			// then
			require.NoError(t, err)
		})
	}
}

func TestFindCodeownersFileFailure(t *testing.T) {
	// given
	tFS := afero.NewMemMapFs()
	revert := codeowners.SetFS(tFS)
	defer revert()

	givenRepoPath := "/workspace/go/repo-without-codeowners/"
	expErrMsg := fmt.Sprintf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository %s", givenRepoPath)

	// when
	entries, err := codeowners.NewFromPath(givenRepoPath)

	// then
	assert.EqualError(t, err, expErrMsg)
	assert.Nil(t, entries)
}
