# Claude Code — Manual Installation

Clone the repo into `/opt/doc-torn` and symlink each skill:

```bash
sudo git clone https://github.com/Anhydrite/doc-torn /opt/doc-torn
ln -s /opt/doc-torn/skills/structured-documentation ~/.claude/skills/
ln -s /opt/doc-torn/skills/doc-driven-exploration ~/.claude/skills/
ln -s /opt/doc-torn/skills/documentation-consistency ~/.claude/skills/
```

Build the `doc-torn-scan` binary (required by the `init` mode):

```bash
cd /opt/doc-torn/tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

Claude Code auto-discovers skills from `~/.claude/skills/`. Load them by name when needed.
