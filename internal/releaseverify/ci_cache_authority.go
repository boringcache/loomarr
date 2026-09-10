package releaseverify

import "strings"

// applyCacheWorkflowAuthorities registers the fork's reviewed cache setup. Native
// command authorities retain their acquisition and test-target constraints.
func applyCacheWorkflowAuthorities(catalog *workflowAuthorityRegistry) {
	{
		workflow := catalog.runs["ci-go.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["rustup toolchain install 1.93.0 --profile minimal"] = cacheWorkflowStep(originalSteps["rustup toolchain install 1.93.0 --profile minimal"], 3, "Install pinned Rust toolchain", "", "", map[string]string{})
		job.steps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"] = cacheWorkflowStep(originalSteps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"], 6, "Use the Cargo adapter in Make targets", "bash", "", map[string]string{"CACHE_ACCESS": "${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}"})
		job.steps["./scripts/ci-ffmpeg.sh metadata >> \"$GITHUB_OUTPUT\""] = cacheWorkflowStep(originalSteps["./scripts/ci-ffmpeg.sh metadata >> \"$GITHUB_OUTPUT\""], 9, "Resolve production FFmpeg pin", "", "", map[string]string{})
		job.steps["./scripts/ci-ffmpeg.sh download \"$HOME/.cache/loomarr-ffmpeg\""] = cacheWorkflowStep(originalSteps["./scripts/ci-ffmpeg.sh download \"$HOME/.cache/loomarr-ffmpeg\""], 10, "Download production FFmpeg", "", "", map[string]string{})
		job.steps["./scripts/ci-ffmpeg.sh install \"$HOME/.cache/loomarr-ffmpeg\" \"$RUNNER_TEMP/loomarr-ffmpeg\"\necho \"$RUNNER_TEMP/loomarr-ffmpeg\" >> \"$GITHUB_PATH\"\n"] = cacheWorkflowStep(originalSteps["./scripts/ci-ffmpeg.sh install \"$HOME/.cache/loomarr-ffmpeg\" \"$RUNNER_TEMP/loomarr-ffmpeg\"\necho \"$RUNNER_TEMP/loomarr-ffmpeg\" >> \"$GITHUB_PATH\"\n"], 11, "Install production FFmpeg", "", "", map[string]string{})
		job.steps["test \"$(command -v ffmpeg)\" = \"$RUNNER_TEMP/loomarr-ffmpeg/ffmpeg\"\ntest \"$(command -v ffprobe)\" = \"$RUNNER_TEMP/loomarr-ffmpeg/ffprobe\"\nffmpeg -version\nffprobe -version\n"] = cacheWorkflowStep(originalSteps["test \"$(command -v ffmpeg)\" = \"$RUNNER_TEMP/loomarr-ffmpeg/ffmpeg\"\ntest \"$(command -v ffprobe)\" = \"$RUNNER_TEMP/loomarr-ffmpeg/ffprobe\"\nffmpeg -version\nffprobe -version\n"], 12, "Verify production FFmpeg on PATH", "", "", map[string]string{})
		job.steps["make test GO_SHARD=${{ matrix.shard }}/${{ strategy.job-total }}"] = cacheWorkflowStep(originalSteps["make test GO_SHARD=${{ matrix.shard }}/${{ strategy.job-total }}"], 13, "", "", "", map[string]string{"GOFLAGS": "-p=1"})
		workflow.jobs["run"] = job
		catalog.runs["ci-go.yml"] = workflow
		catalog.topology["ci-go.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 14}}
		key := workflowJobContextKey{workflow: "ci-go.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-go.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e\",\"with\":{\"go-version\":\"${{ env.GO_VERSION }}\",\"cache\":false}}"
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 4}] = "{\"uses\":\"taiki-e/install-action@c44f6b046f1c29ae5918b1e0bfdbb2f1813836fd\",\"with\":{\"tool\":\"sccache@0.17.0\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 5}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"cargo-dependencies\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 7}] = "{\"name\":\"BoringCache Go dependencies and FFmpeg\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"go\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-go.yml", job: "run", step: 8}] = "{\"name\":\"BoringCache go\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"go\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
	}
	{
		workflow := catalog.runs["ci-go-contracts.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["make fmt shellcheck privacy-verify vet tags-verify vet-tags lint agent-harness-test compose-verify release-verify go-race-verify"] = cacheWorkflowStep(originalSteps["make fmt shellcheck privacy-verify vet tags-verify vet-tags lint agent-harness-test compose-verify release-verify go-race-verify"], 4, "Static analysis and repository contracts", "", "", map[string]string{})
		job.steps["make go-shard-verify SHARDS=3"] = cacheWorkflowStep(originalSteps["make go-shard-verify SHARDS=3"], 5, "The test shards cover every package", "", "", map[string]string{})
		job.steps["make openapi-verify"] = cacheWorkflowStep(originalSteps["make openapi-verify"], 6, "OpenAPI spec is committed and current", "", "", map[string]string{})
		job.steps["make config-docs-verify"] = cacheWorkflowStep(originalSteps["make config-docs-verify"], 7, "Config docs are committed and current", "", "", map[string]string{})
		job.steps["make arch-docs-verify"] = cacheWorkflowStep(originalSteps["make arch-docs-verify"], 8, "The §2 package map is committed and current", "", "", map[string]string{})
		job.steps["make dev-docs-verify"] = cacheWorkflowStep(originalSteps["make dev-docs-verify"], 9, "Command reference is committed and current", "", "", map[string]string{})
		job.steps["make retired-verify"] = cacheWorkflowStep(originalSteps["make retired-verify"], 10, "Retired identifiers are gone from prose", "", "", map[string]string{})
		job.steps["make ci-lint"] = cacheWorkflowStep(originalSteps["make ci-lint"], 11, "Workflows are valid", "", "", map[string]string{})
		job.steps["make observability-verify"] = cacheWorkflowStep(originalSteps["make observability-verify"], 12, "Observability artifacts are provisionable", "", "", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-go-contracts.yml"] = workflow
		catalog.topology["ci-go-contracts.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 13}}
		key := workflowJobContextKey{workflow: "ci-go-contracts.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-go-contracts.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-go-contracts.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-go-contracts.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e\",\"with\":{\"go-version\":\"${{ env.GO_VERSION }}\",\"cache\":false}}"
		catalog.actions[workflowActionKey{workflow: "ci-go-contracts.yml", job: "run", step: 3}] = "{\"name\":\"BoringCache archive\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"go-contracts\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
	}
	{
		workflow := catalog.runs["ci-rust-contracts.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["rustup toolchain install 1.93.0 --profile minimal"] = cacheWorkflowStep(originalSteps["rustup toolchain install 1.93.0 --profile minimal"], 2, "Install pinned Rust toolchain", "", "", map[string]string{})
		job.steps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"] = cacheWorkflowStep(originalSteps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"], 5, "Use the Cargo adapter in Make targets", "bash", "", map[string]string{"CACHE_ACCESS": "${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}"})
		job.steps["make rust-check"] = cacheWorkflowStep(originalSteps["make rust-check"], 6, "", "", "", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-rust-contracts.yml"] = workflow
		catalog.topology["ci-rust-contracts.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 7}}
		key := workflowJobContextKey{workflow: "ci-rust-contracts.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-rust-contracts.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-rust-contracts.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-rust-contracts.yml", job: "run", step: 3}] = "{\"uses\":\"taiki-e/install-action@c44f6b046f1c29ae5918b1e0bfdbb2f1813836fd\",\"with\":{\"tool\":\"sccache@0.17.0\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-rust-contracts.yml", job: "run", step: 4}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"cargo-dependencies\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
	}
	{
		workflow := catalog.runs["ci-image-certification.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["rustup toolchain install 1.93.0 --profile minimal"] = cacheWorkflowStep(originalSteps["rustup toolchain install 1.93.0 --profile minimal"], 3, "Install pinned Rust toolchain", "", "", map[string]string{})
		job.steps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"] = cacheWorkflowStep(originalSteps["if [[ \"$CACHE_ACCESS\" == restore ]]; then\n  echo 'CARGO=boringcache cargo --read-only' >> \"$GITHUB_ENV\"\nelse\n  echo 'CARGO=boringcache cargo --write' >> \"$GITHUB_ENV\"\nfi\n"], 6, "Use the Cargo adapter in Make targets", "bash", "", map[string]string{"CACHE_ACCESS": "${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}"})
		job.steps["IMAGE_CERT_REPORT=\"$RUNNER_TEMP/image-certification.json\" make image-cert"] = cacheWorkflowStep(originalSteps["IMAGE_CERT_REPORT=\"$RUNNER_TEMP/image-certification.json\" make image-cert"], 8, "Certify the release worker against the deterministic real-codec corpus", "", "", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-image-certification.yml"] = workflow
		catalog.topology["ci-image-certification.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 10}}
		key := workflowJobContextKey{workflow: "ci-image-certification.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-image-certification.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e\",\"with\":{\"go-version\":\"${{ env.GO_VERSION }}\",\"cache\":false}}"
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 4}] = "{\"uses\":\"taiki-e/install-action@c44f6b046f1c29ae5918b1e0bfdbb2f1813836fd\",\"with\":{\"tool\":\"sccache@0.17.0\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 5}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"cargo-dependencies\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 7}] = "{\"name\":\"BoringCache go\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"go\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-image-certification.yml", job: "run", step: 9}] = "{\"name\":\"Keep the Rust image certification report\",\"if\":\"always()\",\"uses\":\"actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a\",\"with\":{\"name\":\"rust-image-certification\",\"path\":\"${{ runner.temp }}/image-certification.json\",\"if-no-files-found\":\"ignore\",\"retention-days\":14}}"
	}
	{
		workflow := catalog.runs["ci-android.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["corepack enable"] = cacheWorkflowStep(originalSteps["corepack enable"], 5, "", "", "", map[string]string{})
		job.steps["set -euo pipefail\ncurl -fsSL https://github.com/ccache/ccache/releases/download/v4.14/ccache-4.14-linux-x86_64-glibc.tar.gz -o \"$RUNNER_TEMP/ccache.tar.gz\"\necho \"c64760b0b85ba86068f4cd162dc42e2dc39c6f46b0cb8c1990dfccbec7a1fed0  $RUNNER_TEMP/ccache.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-4.14-linux-x86_64-glibc/ccache\" /usr/local/bin/ccache\ncurl -fsSL https://github.com/ccache/ccache-storage-http-go/releases/download/v0.9/ccache-storage-http-go-0.9-linux-amd64.tar.gz -o \"$RUNNER_TEMP/ccache-http.tar.gz\"\necho \"875dbf6d575d06e4c4492f1ba639beb68530bc23382031a6bacc767cded9f463  $RUNNER_TEMP/ccache-http.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache-http.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-storage-http-go-0.9-linux-amd64/ccache-storage-http\" /usr/local/bin/ccache-storage-http\nccache --version\ntest -x /usr/local/bin/ccache-storage-http\n{\n  echo 'CCACHE_COMPILERCHECK=content'\n  echo 'CCACHE_SLOPPINESS='\n  echo 'CMAKE_C_COMPILER_LAUNCHER=ccache'\n  echo 'CMAKE_CXX_COMPILER_LAUNCHER=ccache'\n} >> \"$GITHUB_ENV\"\n"] = cacheWorkflowStep(originalSteps["set -euo pipefail\ncurl -fsSL https://github.com/ccache/ccache/releases/download/v4.14/ccache-4.14-linux-x86_64-glibc.tar.gz -o \"$RUNNER_TEMP/ccache.tar.gz\"\necho \"c64760b0b85ba86068f4cd162dc42e2dc39c6f46b0cb8c1990dfccbec7a1fed0  $RUNNER_TEMP/ccache.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-4.14-linux-x86_64-glibc/ccache\" /usr/local/bin/ccache\ncurl -fsSL https://github.com/ccache/ccache-storage-http-go/releases/download/v0.9/ccache-storage-http-go-0.9-linux-amd64.tar.gz -o \"$RUNNER_TEMP/ccache-http.tar.gz\"\necho \"875dbf6d575d06e4c4492f1ba639beb68530bc23382031a6bacc767cded9f463  $RUNNER_TEMP/ccache-http.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache-http.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-storage-http-go-0.9-linux-amd64/ccache-storage-http\" /usr/local/bin/ccache-storage-http\nccache --version\ntest -x /usr/local/bin/ccache-storage-http\n{\n  echo 'CCACHE_COMPILERCHECK=content'\n  echo 'CCACHE_SLOPPINESS='\n  echo 'CMAKE_C_COMPILER_LAUNCHER=ccache'\n  echo 'CMAKE_CXX_COMPILER_LAUNCHER=ccache'\n} >> \"$GITHUB_ENV\"\n"], 6, "Install native cache tools", "bash", "", map[string]string{})
		job.steps["make fe-install"] = cacheWorkflowStep(originalSteps["make fe-install"], 10, "", "", "", map[string]string{})
		job.steps["make fe-codegen"] = cacheWorkflowStep(originalSteps["make fe-codegen"], 11, "", "", "", map[string]string{})
		job.steps["make android-profile"] = cacheWorkflowStep(originalSteps["make android-profile"], 12, "Build and profile the verified four-ABI bundle", "", "", map[string]string{"ANDROID_CI_OUTPUT_DIR": "${{ github.workspace }}/.artifacts/android-ci"})
		workflow.jobs["run"] = job
		catalog.runs["ci-android.yml"] = workflow
		catalog.topology["ci-android.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 15}}
		key := workflowJobContextKey{workflow: "ci-android.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-android.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-java@dd06d9cba3e5552c54d9f8ea23572deb30010f7c\",\"with\":{\"distribution\":\"temurin\",\"java-version\":\"21.0.12+1\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 3}] = "{\"uses\":\"android-actions/setup-android@40fd30fb8d7440372e1316f5d1809ec01dcd3699\"}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 4}] = "{\"uses\":\"actions/setup-node@820762786026740c76f36085b0efc47a31fe5020\",\"with\":{\"node-version\":\"22.23.2\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 7}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"android\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 8}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"gradle\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 9}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"ccache\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 13}] = "{\"name\":\"Retain the exact unsigned merge-result bundle\",\"uses\":\"actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a\",\"with\":{\"name\":\"loomarr-android-unsigned-${{ github.sha }}\",\"path\":\".artifacts/android-ci/*\",\"if-no-files-found\":\"error\",\"retention-days\":30}}"
		catalog.actions[workflowActionKey{workflow: "ci-android.yml", job: "run", step: 14}] = "{\"name\":\"Retain Android build diagnostics separately from promotion artifacts\",\"if\":\"${{ always() && steps.android.outcome != 'skipped' }}\",\"uses\":\"actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a\",\"with\":{\"name\":\"loomarr-android-profile-${{ github.sha }}-${{ github.run_id }}-${{ github.run_attempt }}\",\"path\":\".artifacts/android-build-profile/\",\"if-no-files-found\":\"error\",\"retention-days\":30}}"
	}
	{
		workflow := catalog.runs["ci-expo-android-mobile.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["corepack enable"] = cacheWorkflowStep(originalSteps["corepack enable"], 5, "", "", "", map[string]string{})
		job.steps["set -euo pipefail\ncurl -fsSL https://github.com/ccache/ccache/releases/download/v4.14/ccache-4.14-linux-x86_64-glibc.tar.gz -o \"$RUNNER_TEMP/ccache.tar.gz\"\necho \"c64760b0b85ba86068f4cd162dc42e2dc39c6f46b0cb8c1990dfccbec7a1fed0  $RUNNER_TEMP/ccache.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-4.14-linux-x86_64-glibc/ccache\" /usr/local/bin/ccache\ncurl -fsSL https://github.com/ccache/ccache-storage-http-go/releases/download/v0.9/ccache-storage-http-go-0.9-linux-amd64.tar.gz -o \"$RUNNER_TEMP/ccache-http.tar.gz\"\necho \"875dbf6d575d06e4c4492f1ba639beb68530bc23382031a6bacc767cded9f463  $RUNNER_TEMP/ccache-http.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache-http.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-storage-http-go-0.9-linux-amd64/ccache-storage-http\" /usr/local/bin/ccache-storage-http\nccache --version\ntest -x /usr/local/bin/ccache-storage-http\n{\n  echo 'CCACHE_COMPILERCHECK=content'\n  echo 'CCACHE_SLOPPINESS='\n  echo 'CMAKE_C_COMPILER_LAUNCHER=ccache'\n  echo 'CMAKE_CXX_COMPILER_LAUNCHER=ccache'\n} >> \"$GITHUB_ENV\"\n"] = cacheWorkflowStep(originalSteps["set -euo pipefail\ncurl -fsSL https://github.com/ccache/ccache/releases/download/v4.14/ccache-4.14-linux-x86_64-glibc.tar.gz -o \"$RUNNER_TEMP/ccache.tar.gz\"\necho \"c64760b0b85ba86068f4cd162dc42e2dc39c6f46b0cb8c1990dfccbec7a1fed0  $RUNNER_TEMP/ccache.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-4.14-linux-x86_64-glibc/ccache\" /usr/local/bin/ccache\ncurl -fsSL https://github.com/ccache/ccache-storage-http-go/releases/download/v0.9/ccache-storage-http-go-0.9-linux-amd64.tar.gz -o \"$RUNNER_TEMP/ccache-http.tar.gz\"\necho \"875dbf6d575d06e4c4492f1ba639beb68530bc23382031a6bacc767cded9f463  $RUNNER_TEMP/ccache-http.tar.gz\" | sha256sum --check -\ntar -xzf \"$RUNNER_TEMP/ccache-http.tar.gz\" -C \"$RUNNER_TEMP\"\nsudo install -m 755 \"$RUNNER_TEMP/ccache-storage-http-go-0.9-linux-amd64/ccache-storage-http\" /usr/local/bin/ccache-storage-http\nccache --version\ntest -x /usr/local/bin/ccache-storage-http\n{\n  echo 'CCACHE_COMPILERCHECK=content'\n  echo 'CCACHE_SLOPPINESS='\n  echo 'CMAKE_C_COMPILER_LAUNCHER=ccache'\n  echo 'CMAKE_CXX_COMPILER_LAUNCHER=ccache'\n} >> \"$GITHUB_ENV\"\n"], 6, "Install native cache tools", "bash", "", map[string]string{})
		job.steps["make fe-install"] = cacheWorkflowStep(originalSteps["make fe-install"], 10, "", "", "", map[string]string{})
		job.steps["make client-android-debug CLIENT_APP=mobile"] = cacheWorkflowStep(originalSteps["make client-android-debug CLIENT_APP=mobile"], 11, "Generate and build the standalone mobile APK", "", "", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-expo-android-mobile.yml"] = workflow
		catalog.topology["ci-expo-android-mobile.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 13}}
		key := workflowJobContextKey{workflow: "ci-expo-android-mobile.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-expo-android-mobile.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-java@dd06d9cba3e5552c54d9f8ea23572deb30010f7c\",\"with\":{\"distribution\":\"temurin\",\"java-version\":\"21.0.12+1\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 3}] = "{\"uses\":\"android-actions/setup-android@40fd30fb8d7440372e1316f5d1809ec01dcd3699\"}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 4}] = "{\"uses\":\"actions/setup-node@820762786026740c76f36085b0efc47a31fe5020\",\"with\":{\"node-version\":\"22.23.2\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 7}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"android\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 8}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"gradle\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 9}] = "{\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"ccache\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
		catalog.actions[workflowActionKey{workflow: "ci-expo-android-mobile.yml", job: "run", step: 12}] = "{\"name\":\"Keep the standalone mobile APK\",\"if\":\"always()\",\"uses\":\"actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a\",\"with\":{\"name\":\"expo-android-mobile\",\"path\":\"web/apps/mobile/android/app/build/outputs/apk/debug/app-debug.apk\",\"if-no-files-found\":\"error\",\"retention-days\":7}}"
	}
	{
		workflow := catalog.runs["ci-image.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["set -euo pipefail\nimage=loomarr-ci:verify\ntest \"$(docker image inspect \"$image\" --format '{{ index .Config.Labels \"org.opencontainers.image.documentation\" }}')\" = \\\n  \"https://github.com/loomarr/loomarr/blob/main/THIRD_PARTY_NOTICES.md\"\ntest \"$(docker image inspect \"$image\" --format '{{ index .Config.Labels \"org.opencontainers.image.licenses\" }}')\" = \\\n  \"MIT AND GPL-3.0-or-later\"\ncontainer=$(docker create \"$image\")\ntrap 'docker rm -f \"$container\" >/dev/null' EXIT\ndocker cp \"$container:/usr/share/doc/loomarr/LICENSE\" \"$RUNNER_TEMP/LICENSE\"\ndocker cp \"$container:/usr/share/doc/loomarr/THIRD_PARTY_NOTICES.md\" \"$RUNNER_TEMP/THIRD_PARTY_NOTICES.md\"\ncmp LICENSE \"$RUNNER_TEMP/LICENSE\"\ncmp THIRD_PARTY_NOTICES.md \"$RUNNER_TEMP/THIRD_PARTY_NOTICES.md\"\n"] = cacheWorkflowStep(originalSteps["set -euo pipefail\nimage=loomarr-ci:verify\ntest \"$(docker image inspect \"$image\" --format '{{ index .Config.Labels \"org.opencontainers.image.documentation\" }}')\" = \\\n  \"https://github.com/loomarr/loomarr/blob/main/THIRD_PARTY_NOTICES.md\"\ntest \"$(docker image inspect \"$image\" --format '{{ index .Config.Labels \"org.opencontainers.image.licenses\" }}')\" = \\\n  \"MIT AND GPL-3.0-or-later\"\ncontainer=$(docker create \"$image\")\ntrap 'docker rm -f \"$container\" >/dev/null' EXIT\ndocker cp \"$container:/usr/share/doc/loomarr/LICENSE\" \"$RUNNER_TEMP/LICENSE\"\ndocker cp \"$container:/usr/share/doc/loomarr/THIRD_PARTY_NOTICES.md\" \"$RUNNER_TEMP/THIRD_PARTY_NOTICES.md\"\ncmp LICENSE \"$RUNNER_TEMP/LICENSE\"\ncmp THIRD_PARTY_NOTICES.md \"$RUNNER_TEMP/THIRD_PARTY_NOTICES.md\"\n"], 3, "Inspect packaged license and notice metadata", "bash", "", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-image.yml"] = workflow
		catalog.topology["ci-image.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 4}}
		key := workflowJobContextKey{workflow: "ci-image.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "${{ matrix.runner }}"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-image.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-image.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-image.yml", job: "run", step: 2}] = "{\"name\":\"BoringCache docker\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"docker\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
	}
	{
		workflow := catalog.runs["ci-frontend.yml"]
		workflow.environment = map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"}
		workflow.permissions = map[string]string{"contents": "read", "id-token": "write"}
		job := workflow.jobs["run"]
		originalSteps := job.steps
		job.steps = map[string]workflowStepAuthority{}
		job.steps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"] = cacheWorkflowStep(originalSteps["git rev-parse HEAD\ncat .github/boringcache-upstream\ngit status --short\nlscpu\nsha256sum go.sum Cargo.lock rust-toolchain.toml web/pnpm-lock.yaml\n"], 1, "Record source and runner identity", "bash", "", map[string]string{})
		job.steps["corepack enable"] = cacheWorkflowStep(originalSteps["corepack enable"], 3, "Enable pnpm (pinned via web/package.json packageManager)", "", "", map[string]string{})
		job.steps["make fe-install"] = cacheWorkflowStep(originalSteps["make fe-install"], 5, "", "", "", map[string]string{})
		job.steps["make fe FE_SHARD=${{ matrix.shard }}/${{ strategy.job-total }}"] = cacheWorkflowStep(originalSteps["make fe FE_SHARD=${{ matrix.shard }}/${{ strategy.job-total }}"], 6, "", "", "", map[string]string{})
		job.steps["make fe-tokens-verify"] = cacheWorkflowStep(originalSteps["make fe-tokens-verify"], 7, "Token artifacts are committed and current", "", "matrix.shard == 1", map[string]string{})
		workflow.jobs["run"] = job
		catalog.runs["ci-frontend.yml"] = workflow
		catalog.topology["ci-frontend.yml"] = workflowTopologyAuthority{jobs: map[string]int{"run": 8}}
		key := workflowJobContextKey{workflow: "ci-frontend.yml", job: "run"}
		context := catalog.jobContexts[key]
		context.runsOn = "ubuntu-24.04"
		catalog.jobContexts[key] = context
		for key := range catalog.actions {
			if key.workflow == "ci-frontend.yml" {
				delete(catalog.actions, key)
			}
		}
		catalog.actions[workflowActionKey{workflow: "ci-frontend.yml", job: "run", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\"}"
		catalog.actions[workflowActionKey{workflow: "ci-frontend.yml", job: "run", step: 2}] = "{\"uses\":\"actions/setup-node@820762786026740c76f36085b0efc47a31fe5020\",\"with\":{\"node-version\":\"${{ env.NODE_VERSION }}\"}}"
		catalog.actions[workflowActionKey{workflow: "ci-frontend.yml", job: "run", step: 4}] = "{\"name\":\"BoringCache archive\",\"uses\":\"boringcache/one@404b744a2053da4cf963f13f615f7fafe94f3cf7\",\"with\":{\"cli-version\":\"${{ inputs.cli_version }}\",\"mode\":\"archive\",\"cache-profiles\":\"frontend\",\"trust-policy\":\"${{ (github.event_name == 'pull_request' || github.event_name == 'merge_group' || inputs.cache_access == 'restore') && 'restore' || 'auto' }}\",\"fail-on-cache-error\":true}}"
	}

	root := catalog.runs["ci.yml"]
	root.permissions = map[string]string{"contents": "read", "id-token": "write"}
	catalog.runs["ci.yml"] = root

	const replay = "boringcache-validation.yml"
	catalog.topology[replay] = workflowTopologyAuthority{jobs: map[string]int{"connect": 3}}
	catalog.runs[replay] = workflowAuthority{
		environment: map[string]string{"GO_VERSION": "1.27.0", "NODE_VERSION": "22.23.2"},
		permissions: map[string]string{"contents": "read", "id-token": "write"},
		jobs:        map[string]workflowJobAuthority{},
	}
	catalog.jobContexts[workflowJobContextKey{workflow: replay, job: "connect"}] = workflowJobContextAuthority{runsOn: "ubuntu-24.04", timeoutMinutes: 20}
	catalog.runs[replay].jobs["connect"] = workflowJobAuthority{condition: "github.repository == 'boringcache/loomarr' && inputs.phase == 'connect'", steps: map[string]workflowStepAuthority{
		"curl -fsSL https://install.boringcache.com/install.sh | sh\necho \"$HOME/.local/bin\" >> \"$GITHUB_PATH\"\n": exactWorkflowStep(1, "Install the released CLI for enrollment", workflowStepAuthority{allowsAcquisition: true}),
		"boringcache ci connect --oidc-provider github-actions":                                                       exactWorkflowStep(2, "Connect the fork with GitHub OIDC", workflowStepAuthority{allowsAcquisition: true}),
	}}
	catalog.actions[workflowActionKey{workflow: replay, job: "connect", step: 0}] = "{\"uses\":\"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\",\"with\":{\"persist-credentials\":false}}"
	registerCacheReplayJob(catalog, "ci-go.yml", "cache_go")
	registerCacheReplayJob(catalog, "ci-go-contracts.yml", "cache_go_contracts")
	registerCacheReplayJob(catalog, "ci-rust-contracts.yml", "cache_rust_contracts")
	registerCacheReplayJob(catalog, "ci-image-certification.yml", "cache_image_certification")
	registerCacheReplayJob(catalog, "ci-android.yml", "cache_android")
	registerCacheReplayJob(catalog, "ci-expo-android-mobile.yml", "cache_expo_android_mobile")
	registerCacheReplayJob(catalog, "ci-image.yml", "cache_image")
	registerCacheReplayJob(catalog, "ci-frontend.yml", "cache_frontend")
}

func cacheWorkflowStep(original workflowStepAuthority, index int, name, shell, condition string, environment map[string]string) workflowStepAuthority {
	original.shell = shell
	original.condition = condition
	original.environment = environment
	return exactWorkflowStep(index, name, original)
}

func cacheReplayPolicy(value string) string {
	return strings.ReplaceAll(value, "inputs.cache_access == 'restore'", "inputs.phase == 'warm'")
}

func registerCacheReplayJob(catalog *workflowAuthorityRegistry, source, name string) {
	const replay = "boringcache-validation.yml"
	catalog.topology[replay].jobs[name] = catalog.topology[source].jobs["run"]
	catalog.jobContexts[workflowJobContextKey{workflow: replay, job: name}] = catalog.jobContexts[workflowJobContextKey{workflow: source, job: "run"}]
	sourceJob := catalog.runs[source].jobs["run"]
	job := sourceJob
	job.condition = "github.repository == 'boringcache/loomarr' && github.event_name == 'workflow_dispatch' && inputs.phase != 'connect'"
	job.steps = make(map[string]workflowStepAuthority, len(sourceJob.steps))
	for command, step := range sourceJob.steps {
		step.environment = make(map[string]string, len(step.environment))
		for key, value := range sourceJob.steps[command].environment {
			step.environment[key] = cacheReplayPolicy(value)
		}
		job.steps[command] = step
	}
	catalog.runs[replay].jobs[name] = job
	for key, authority := range catalog.actions {
		if key.workflow == source {
			catalog.actions[workflowActionKey{workflow: replay, job: name, step: key.step}] = cacheReplayPolicy(authority)
		}
	}
}
