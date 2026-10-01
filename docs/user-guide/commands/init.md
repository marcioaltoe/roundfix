### init

```bash
roundfix init [--scope user] [--force]
```

Creates Project Config (`<repo>/.roundfixrc.yml`) or, with `--scope user`,
User Config (`~/.roundfix/config.yml`). When `--scope` is omitted, Roundfix
asks and defaults to Project Config. `--force` overwrites an existing file.

