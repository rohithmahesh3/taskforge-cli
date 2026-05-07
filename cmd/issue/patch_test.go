package issue

import (
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPatchFromFlagsReplaceRequiresActual(t *testing.T) {
	t.Cleanup(resetPatchFlags)
	patchOp = "replace"
	patchField = "description"
	patchValue = "new text"

	_, err := buildPatchFromFlags()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--actual is required")
}

func TestBuildPatchFromFlagsMapsServerFields(t *testing.T) {
	t.Cleanup(resetPatchFlags)
	patchOp = "replace"
	patchField = "description"
	patchValue = "new text"
	patchActual = "old text"

	ops, err := buildPatchFromFlags()
	require.NoError(t, err)
	require.Len(t, ops, 1)
	assert.Equal(t, "old text", ops[0].Old)
	assert.Equal(t, "new text", ops[0].New)
}

func TestNormalizePatchesLegacyFields(t *testing.T) {
	in := []taskforge.PatchOp{
		{Op: "replace", Field: "description", Value: "v2", Old: "v1"},
		{Op: "insert", Field: "description", Value: " plus", After: "start"},
		{Op: "delete", Field: "description", Before: "remove-me"},
		{Op: "diff", Field: "description", Diff: "@@ -1 +1 @@\n-a\n+b"},
	}
	out := normalizePatches(in)
	require.Len(t, out, 4)
	assert.Equal(t, "v2", out[0].New)
	assert.Equal(t, " plus", out[1].Content)
	assert.Equal(t, "remove-me", out[2].Old)
	assert.Equal(t, "@@ -1 +1 @@\n-a\n+b", out[3].Unified)
}

func resetPatchFlags() {
	patchOp = ""
	patchField = ""
	patchValue = ""
	patchDiff = ""
	patchAfter = ""
	patchBefore = ""
	patchFile = ""
	patchActual = ""
}
