package releaseverify

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

// The Docker action delegates its command and output destinations to this plan.
// Bind those inputs to the reviewed source just as the action mappings are bound.
func verifyCacheWorkflowPlan(root string) error {
	contents, err := os.ReadFile(filepath.Join(root, ".boringcache.toml"))
	if err != nil {
		return fmt.Errorf("read cache workflow plan: %w", err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(contents)) != "cd97e9a5ed3c3f72c1f420425c8eb9946ab38ae7d8c7b4bfb0d7b2ab50afc6c2" {
		return fmt.Errorf("cache workflow plan differs from its reviewed command and storage authority")
	}
	return nil
}
