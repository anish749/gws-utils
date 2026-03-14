# gws_utils

Utilities and wrappers on top of [Google Workspace CLI (`gws`)](https://github.com/googleworkspace/cli).

`gws` provides raw API access to Google Workspace services. This project adds higher-level operations that combine multiple API calls — things like downloading a Google Doc with all its tabs as structured markdown files.

## Prerequisites

Requires `gws` to be installed and authenticated. If `gws` is not on your PATH, install it:

```bash
npm install -g @googleworkspace/cli
gws auth setup
gws auth login
```

See https://github.com/googleworkspace/cli for full setup instructions.

## Usage

### download

Download a Google Doc as markdown. Single-tab docs become a single `.md` file. Multi-tab docs become a folder with one numbered `.md` file per tab.

```bash
# By URL
gws_utils download -o ./output "https://docs.google.com/document/d/DOC_ID/edit"

# By document ID
gws_utils download -o ./output DOC_ID
```

**Flags:**
- `-o, --output` — output directory (default: `.`)
