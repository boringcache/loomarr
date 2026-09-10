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
	if fmt.Sprintf("%x", sha256.Sum256(contents)) != "82edf37d366f7915e194408ac8ba8461bd33560a2cf89c5402cb03dbf8b17dc0" {
		return fmt.Errorf("cache workflow plan differs from its reviewed command and storage authority")
	}
	return nil
}
