# TODO

## CLI

### MVP

- [x] Core

  - [?] Before action
  - [?] After action
  - [x] Human friendly errors
  - [x] Rest args

- [x] Help generator

  - [x] Format table

- [x] Parser

  - [x] Errors if flags and args registred twice
    - [?] Show line and file where it was registred
  - [?] Global flags
  - [x] Required flags
  - [x] Optional args
  - [x] `"--"` bypass
  - [x] POSIX-style short flag combining (`-a -b -> -ab`)
  - [x] Short-flag+parameter combining (`-a parm -> -aparm`)
  - [?] Default value
  - [x] Better `isBoolValue`

- [x] Completion generator

  - [x] ZSH

- [x] Commands

  - [x] Help
  - [x] Completion generator

- [x] Command flags

  - [x] Help

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
  - [ ] Long description for commands, flags and args

- [ ] Documentation generator

  - [ ] Markdown
  - [ ] Man
  - [ ] Org Mode
  - [ ] reStructuredText

- [ ] Completion generator

  - [ ] Bash
  - [ ] Fish
  - [ ] PowerShell

- [ ] Commands

  - [ ] Documentation generator

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
  - [x] Test runner
  - [x] Linter runner
  - [ ] Build packages
  - [ ] Commit messages

- [x] Codecov

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

- [x] Logo
- [x] Github tags
- [x] Github thumbnail
- [ ] Article on dev.to
- [ ] Article on habr.com
- [ ] Quick introduction on YouTube
- [x] Subreddit
- [ ] Twitter posts
- [x] Custom domain

### Dream release

- [ ] Article on medium
- [ ] Who use Nice
- [ ] Discourse forum

---

## Misc

### MVP

- [x] Remove and move `escape` package into other branch

### Dream release

- [ ] Packages for OSes

  - [ ] NixOS
  - [ ] Arch
  - [ ] Ubunta / Debian
  - [ ] Brew
  - [ ] Mac Ports
  - [ ] Windows
