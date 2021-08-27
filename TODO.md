# TODO

## CLI

### MVP

- [ ] Core

  - [ ] Before action
  - [ ] After action

- [ ] Help generator

  - [ ] Format table

- [ ] Parser

  - [ ] Errors if flags and args registred twice
  - [ ] Global flags
  - [ ] Required flags
  - [ ] Optional args
  - [ ] `"--"` bypass
  - [ ] POSIX-style short flag combining (`-a -b -> -ab`)
  - [ ] Short-flag+parameter combining (`-a parm -> -aparm`)

- [ ] Documentation generator

  - [ ] Markdown

- [ ] Completion generator

  - [ ] Bash
  - [ ] ZSH

- [ ] Commands

  - [ ] Documentation generator
  - [ ] Completion generator

- [ ] Command flags

  - [ ] Help

### Dream release

- [ ] Core

  - [ ] Usager

    - [ ] Template

- [ ] Documentation generator

  - [ ] Man
  - [ ] Org Mode

- [ ] Completion generator

  - [ ] Fish
  - [ ] PowerShell

---

## Colors

### MVP

- [ ] Core

  - [ ] Supports colors checker
  - [ ] Setters for color mode and optimize `should*`
  - [ ] Optimize uint8 to string
  - [ ] Optimize uint8 with prefixes to string

- [ ] Tests

  - [ ] Cases when terminal does not support colors

---

## Misc

### MVP

- [ ] Remove and move `escape` package into other branch
