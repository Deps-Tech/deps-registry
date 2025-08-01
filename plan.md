# Dependency Consolidation & Tooling Upgrade Plan

This plan outlines the two-phase process to fix the dependency structure and upgrade the tooling to prevent future issues.

### Phase 1: Consolidate Existing Dependencies (The "Fix")

The goal is to programmatically merge all `*-lib` packages into their primary counterparts.

**Todo:**
1.  **Create `tools/cmd/merge_deps.go`:** A temporary tool to handle the migration.
2.  **Implement Merge Logic:**
    - Scan `deps/` for directories ending in `-lib`.
    - For each `*-lib` directory, find the corresponding main directory.
    - Move contents from `[name]-lib/[version]/` to `[name]/[version]/`.
    - Read `dep.json` from both directories.
    - Merge the `files` array from the `-lib` manifest into the main manifest.
    - Save the updated `dep.json` in the main directory.
3.  **Register and Run:** Add the new command to `tools/main.go` and execute it.
4.  **Cleanup:** Remove the temporary `merge_deps.go` tool, the command from `tools/main.go`, and all the now-empty `*-lib` directories.

### Phase 2: Upgrade The Tooling (The "Future-Proofing")

The goal is to upgrade the `dep add` command to handle complex libraries with multiple source paths correctly.

**Todo:**
1.  **Refactor `copyFiles` in `tools/cmd/utils.go`:** Modify the function to support recursive copying of directories.
2.  **Upgrade `addDependency` in `tools/cmd/dep_add.go`:**
    - Change the `--source-path` flag to accept multiple values.
    - Update the command logic to iterate over all source paths and use the new `copyFiles` function.
3.  **Update Documentation:** Ensure the command's help text reflects the new capabilities.

### Mermaid Diagram: Structural Change

```mermaid
graph TD
    subgraph Before
        A["deps/socket/1.0.0/"] --> A1["socket.lua"]
        A --> A2["dep.json (files: [socket.lua])"]
        B["deps/socket-lib/1.0.0/"] --> B1["http.lua"]
        B --> B2["dep.json (files: [http.lua])"]
    end

    subgraph After
        C["deps/socket/1.0.0/"] --> C1["socket.lua"]
        C --> C2["http.lua"]
        C --> C3["dep.json (files: [socket.lua, http.lua])"]
    end

    Before -->|Migration & Tooling Upgrade| After