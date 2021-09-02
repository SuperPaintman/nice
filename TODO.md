# TODO

## CLI

### MVP

- [ ] Core

  - [ ] Before action
  - [ ] After action
  - [ ] Human friendly errors
  - [ ] Long description for commands, flags and args
  - [ ] Rest args

- [ ] Help generator

  - [ ] Format table

- [ ] Parser

  - [ ] Errors if flags and args registred twice
    - [ ] Show line and file where it was registred
  - [ ] Global flags
  - [ ] Required flags
  - [ ] Optional args
  - [ ] `"--"` bypass
  - [ ] POSIX-style short flag combining (`-a -b -> -ab`)
  - [ ] Short-flag+parameter combining (`-a parm -> -aparm`)
  - [ ] Default value

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

- [ ] Parser

  - [ ] "Did you mean?" for unknown flags and commands
  - [ ] Value from ENV
  - [ ] Hidden commands and flags
  - [ ] Negative bool flags (`--no-*`)
  - [ ] Array types (`-i 1 -i 2 -i 3`)
  - [ ] Object types (`--a.b.c 2`)
  - [ ] "deprecated" option

- [ ] Documentation generator

  - [ ] Man
  - [ ] Org Mode
  - [ ] reStructuredText

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

## CI

### MVP

- [ ] Github CI

  - [ ] Check generated files (regenerate them)
  - [ ] Test runner
  - [ ] Linter runner
  - [ ] Build packages
  - [ ] Commit messages

- [ ] Coveralls

### Dream release

- [ ] Github CI

  - [ ] License updater
  - [ ] Python linter
  - [ ] Markdown linter
  - [ ] Test Go in Markdown comments
  - [ ] Lint Go in Markdown comments

- [ ] Code Climate

---

## Documentation

### MVP

- [ ] Update `README.md`
- [ ] Examples
- [ ] Go comments

### Dream release

- [ ] Table with simulator projects
- [ ] Tutorial

---

## Community

### MVP

- [ ] Logo
- [ ] Github tags
- [ ] Github thumbnail
- [ ] Article on dev.to
- [ ] Article on habr.com
- [ ] Quick introduction on YouTube
- [ ] Subreddit
- [ ] Twitter posts

### Dream release

- [ ] Article on medium
- [ ] Who use Nice

---

## Misc

### MVP

- [ ] Remove and move `escape` package into other branch

- [ ] Packages for OSes

  - [ ] NixOS
  - [ ] Arch
  - [ ] Ubunta / Debian
  - [ ] Brew
  - [ ] Mac Ports
  - [ ] Windows
