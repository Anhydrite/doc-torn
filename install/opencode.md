# OpenCode — Manual Installation

Clone the repo into `/opt/doc-torn` and symlink each skill:

```bash
sudo git clone https://github.com/Anhydrite/doc-torn /opt/doc-torn
ln -s /opt/doc-torn/skills/structured-documentation ~/.config/opencode/skills/
ln -s /opt/doc-torn/skills/doc-driven-exploration ~/.config/opencode/skills/
ln -s /opt/doc-torn/skills/documentation-consistency ~/.config/opencode/skills/
```

Build the `doc-torn-scan` binary (required by the `init` mode):

```bash
cd /opt/doc-torn/tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

OpenCode natively scans `~/.config/opencode/skills/`, `~/.claude/skills/`, and `~/.agents/skills/`. No plugin needed.
