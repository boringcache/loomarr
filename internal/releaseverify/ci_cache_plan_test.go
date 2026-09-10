package releaseverify

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCacheWorkflowPlanRejectsCommandAndWorkspaceChanges(t *testing.T) {
	for name, mutation := range map[string][2]string{
		"published image":     {"\"--load\"", "\"--push\""},
		"different workspace": {"boringcache/loomarr", "other/workspace"},
		"secret directory":    {"~/.cache/loomarr-ffmpeg", "~/.ssh"},
		"compiler checks":     {"fail-on-cache-error = true", "fail-on-cache-error = false"},
	} {
		t.Run(name, func(t *testing.T) {
			root := writeCIContainerDownloadsFixture(t)
			if err := verifyCacheWorkflowPlan(root); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, ".boringcache.toml")
			original := readFixtureFile(t, path)
			changed := strings.Replace(original, mutation[0], mutation[1], 1)
			if changed == original {
				t.Fatal("mutation did not change the plan")
			}
			writeFixtureFile(t, path, changed)
			if err := verifyCacheWorkflowPlan(root); err == nil {
				t.Fatal("accepted changed cache authority")
			}
		})
	}
}
