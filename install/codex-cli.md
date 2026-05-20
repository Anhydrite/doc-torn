# Codex CLI — Manual Installation

Clone the repo into `/opt/doc-torn` and symlink each skill:

```bash
sudo git clone https://github.com/Anhydrite/doc-torn /opt/doc-torn
ln -s /opt/doc-torn/skills/structured-documentation ~/.agents/skills/
ln -s /opt/doc-torn/skills/doc-driven-exploration ~/.agents/skills/
ln -s /opt/doc-torn/skills/documentation-consistency ~/.agents/skills/
```

Build the `doc-torn-scan` binary (required by the `init` mode):

```bash
cd /opt/doc-torn/tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

Codex CLI auto-discovers skills from `~/.agents/skills/`. No additional configuration needed.
