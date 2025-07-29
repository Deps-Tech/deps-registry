# Launcher Dependencies

This repository serves as the central, trusted source for all dependencies used by the game launcher ecosystem. Every dependency is vetted and versioned to ensure stability and security for all mods and users.

## How It Works

This repository is managed via the `dep-tools` Go CLI utility. The tool standardizes the process of adding and updating dependencies. A CI/CD pipeline automatically validates all contributions to ensure they meet our strict quality and security criteria.

## Contribution Criteria

All dependencies submitted must adhere to the following rules without exception. Submissions that violate these rules will be rejected.

1.  **No Binaries:** Only Lua source code (`.lua`) and common asset files (e.g., `.ttf`, `.png`, `.wav`) are permitted. Any compiled code (e.g., `.luac`, `.dll`, `.so`) is strictly forbidden.
2.  **No Obfuscation:** Code must be human-readable and well-formatted. Any attempt to obfuscate or minify code will result in rejection.
3.  **Strict Sandboxing:** Dependencies must not attempt to access the network or the user's filesystem outside of the game's designated directories.
4.  **Clear Versioning:** All dependencies must follow Semantic Versioning (e.g., `1.0.0`).
5.  **Manifest Integrity:** Every version of every dependency must include a valid `dep.json` manifest file.

## Manifest Specification (`dep.json`)

The `dep.json` file is the manifest that describes a specific version of a dependency. It must be present in every version's directory.

**Structure:**
```json
{
  "id": "awesome-lib",
  "version": "1.2.0",
  "sourceUrl": "https://some-forum.com/t/awesome-lib-release/123",
  "files": [
    "awesome-lib.lua"
  ]
}
```

**Field Descriptions:**

| Field       | Type           | Description                                                                                             |
|-------------|----------------|---------------------------------------------------------------------------------------------------------|
| `id`        | String         | The unique, lowercase, dash-cased identifier for the dependency. This is immutable.                     |
| `version`   | String         | The Semantic Version of this specific release (e.g., "1.2.0").                                          |
| `sourceUrl` | String         | The original URL where this version of the dependency was found. Used for reference and verification.     |
| `files`     | Array of Strings | A complete list of all filenames included in this dependency version, relative to the manifest's location. |

## Usage

### Compile the Tool

```bash
go build -o dep-tools .
```

### Add a New Dependency

```bash
./dep-tools add
```

### Update an Existing Dependency

```bash
./dep-tools update