# Loomarr cache validation

The fork follows the Adorsys integration pattern: a repository-owned
`.boringcache.toml`, released One, GitHub OIDC and the existing build/test targets.
The dedicated workspace is `boringcache/loomarr`, with a 20 GB cleanup budget.
No static cache tokens are configured. Pull requests and merge-queue jobs restore
only; trusted seed and rolling jobs publish. The warm dispatch is restore-only.

## Workloads

| Workload | Cache surface | Preserved checks |
| --- | --- | --- |
| Go contracts and three test shards | Native Go cache plus module/tool download archives | Lint, policy, race rules, package partition and real test execution |
| Rust contracts and image certification | Cargo target snapshots, registry/git archives, sccache 0.17.0 | Pinned Rust 1.93.0, fmt, Clippy, tests and real-codec certification |
| Frontend | pnpm download store | Two test shards, type checks, builds and generated-token checks |
| Android TV and mobile | Gradle task cache, ccache 4.14, dependency archives | Clean Expo prebuild, native compiler checks, four TV ABIs, unsigned bundle verification and mobile APK |
| Docker amd64 and arm64 | Layers, Go and sccache tool caches, pnpm and Cargo registry mounts | Native architectures, all Dockerfile proofs and packaged license/notice verification |

The Dockerfile's Rust build still uses the checked-in Rust toolchain. Go uses
1.27.0; Node uses 22.23.2; Java uses Temurin 21.0.12+1. The One distribution is
`404b744a2053da4cf963f13f615f7fafe94f3cf7` (v1.30.1); leave `cli_version` empty to
use its released default. An explicit value is a canary, not a comparison default.

## Replay protocol

The exact six upstream revisions are in `boringcache-replay.tsv`. They are one
seed and five consecutive first-parent commits ending at upstream main captured
on 10 September 2026. Keep this sequence fixed while collecting the series.

1. Apply the same reviewed integration to the seed on `boringcache-replay`.
2. Dispatch `BoringCache validation` with `phase=seed`. The new workspace and
   branch tags have no workload state restored. Record any setup failures as
   failed attempts; do not relabel retries as globally cold storage.
3. After a successful seed, dispatch `phase=warm` at the same commit on fresh
   runners. It restores without publishing and establishes same-source reuse.
4. Merge the next listed upstream commit into the replay branch, update
   `boringcache-upstream`, and dispatch `phase=rolling`. Repeat in order for all
   five commits using unchanged cache tags, workflow commands and toolchains.
   Do not add synthetic source changes or hidden prewarming. Retain failed runs
   and diagnose them before treating a later result as valid data.
5. Check that the replay's final source contains the captured upstream main and
   its integration matches `boringcache-validation`. Future upstream syncs can
   continue the same branch/cache history with one rolling run per real commit.

Record provider run/job IDs, full job duration, restore/build/save timing, CPU
model and runner image, source/lock hashes, native cache statistics, product
cache errors and retained artifacts. The workflow logs record source and runner
identity; One emits its own cache evidence. The Android job retains its existing
profile and exact unsigned bundle, and image certification retains its report.
Inspect BoringCache's run summaries alongside the native logs.

These five changes chiefly exercise Go. Report each revision's changed paths:
unchanged Rust or Android inputs establish reuse only. A fully cached Docker RUN
may make no tool-cache requests and perform no mount restore; count that as a
layer hit. A future actual Rust/native-input or lock change is required to prove
invalidation for those workloads. This series has no matched GitHub-cache
control, so it supports correctness and reuse findings, not a provider speedup
claim. Record elapsed time including cache transport and publication.

The validation entry builds and tests within the fork. It does not sign or
publish release images, upload to Play, deploy, or contact upstream maintainers.
