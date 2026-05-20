# GitHub Copilot CLI — Manual Installation

Clone the repo into `/opt/doc-torn`:

```bash
sudo git clone https://github.com/Anhydrite/doc-torn /opt/doc-torn
```

Build the `doc-torn-scan` binary (required by the `init` mode):

```bash
cd /opt/doc-torn/tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

Copilot CLI does not auto-scan directories. Use the `skill` tool and reference the SKILL.md files by path when needed.
