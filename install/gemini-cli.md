# Gemini CLI — Manual Installation

Clone the repo into `/opt/doc-torn`:

```bash
sudo git clone https://github.com/Anhydrite/doc-torn /opt/doc-torn
```

Build the `doc-torn-scan` binary (required by the `init` mode):

```bash
cd /opt/doc-torn/tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

Gemini CLI does not auto-scan directories. Use `activate_skill` with the path to each SKILL.md:

```bash
activate_skill /opt/doc-torn/skills/structured-documentation/SKILL.md
activate_skill /opt/doc-torn/skills/doc-driven-exploration/SKILL.md
activate_skill /opt/doc-torn/skills/documentation-consistency/SKILL.md
```
