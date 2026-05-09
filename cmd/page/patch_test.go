package page

import (
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPagePatchFromFlagsReplaceRequiresActual(t *testing.T) {
	t.Cleanup(resetPagePatchFlags)
	pagePatchOp = "replace"
	pagePatchField = "content"
	pagePatchValue = "new text"

	_, err := buildPagePatchFromFlags()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--actual is required")
}

func TestBuildPagePatchFromFlagsMapsServerFields(t *testing.T) {
	t.Cleanup(resetPagePatchFlags)
	pagePatchOp = "replace"
	pagePatchField = "content"
	pagePatchValue = "new text"
	pagePatchActual = "old text"

	ops, err := buildPagePatchFromFlags()
	require.NoError(t, err)
	require.Len(t, ops, 1)
	assert.Equal(t, "old text", ops[0].Old)
	assert.Equal(t, "new text", ops[0].New)
}

func TestNormalizePagePatchesLegacyFields(t *testing.T) {
	in := []taskforge.PatchOp{
		{Op: "replace", Field: "content", Value: "v2", Old: "v1"},
		{Op: "insert", Field: "content", Value: " plus", After: "start"},
		{Op: "delete", Field: "content", Before: "remove-me"},
		{Op: "diff", Field: "content", Diff: "@@ -1 +1 @@\n-a\n+b"},
	}
	out := normalizePagePatches(in)
	require.Len(t, out, 4)
	assert.Equal(t, "v2", out[0].New)
	assert.Equal(t, " plus", out[1].Content)
	assert.Equal(t, "remove-me", out[2].Old)
	assert.Equal(t, "@@ -1 +1 @@\n-a\n+b", out[3].Unified)
}

func resetPagePatchFlags() {
	pagePatchOp = ""
	pagePatchField = ""
	pagePatchValue = ""
	pagePatchDiff = ""
	pagePatchAfter = ""
	pagePatchBefore = ""
	pagePatchFile = ""
	pagePatchActual = ""
}
