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
	if fmt.Sprintf("%x", sha256.Sum256(contents)) != "4a9d56f17f936d793cb62da183b6043d37f87fc55e5f5faca66d3c1886358939" {
		return fmt.Errorf("cache workflow plan differs from its reviewed command and storage authority")
	}
	return nil
}
