# `tag-monorepo`: Per-module Monorepo Tagging

Working in monorepo codebases is great. There's no need to configure access to
private git modules; everything is all in one place! But what about when you
need to make one module in a monorepo accessible outside of the workspace? When
working in Go, this is a surprisingly common occurrence.

In Go, it is expected that each module in a monorepo is versioned separately
from the others (See
[the multi-module source doc](https://go.dev/doc/modules/managing-source#multiple-module-source)).
This program provides a simple text user interface (TUI) to simplify the
multi-package versioning process. It is not intended to solve all possible
cases, merely to streamline the most common cases. You will likely still need to
tag manually in more "interesting" cases.

## Features

- Current tag detection for each module
- Update by major, minor, or patch version
- Apply tag suffixes, like `beta` or `rc1`
- Handle modules that are not at the top level

## Depth

When working on monorepos, the modules are usually the top-level directories.
However, it is also valid to have a versioned module that is inside another
directory. In these cases, the tag is the full file path, ending with the
verison number. (for example, generated Go protobufs may be located in a module
called `proto/go/v0.0.1`)

By default, this program will only display directories that have existing tags.
This means that in the general case, you get a "focused" view that excludes
un-tagged folders like "docs" or "scripts". If you want to add a tag for a new
directory, run the program with the `-d` flag to specify the maximum directory
depth. The program will then display all directories at or above that depth in
addition to any directories with tags.

## Installation

```bash
go install github.com/skubalj/tag-monorepo@latest
```

## Acknowledgements

This program is made possible thanks to the generous open source contributions
of others.

| Dependency                    | License      |
| :---------------------------- | :----------- |
| `charm.land/bubbletea/v2`     | MIT          |
| `charm.land/bubbles/v2`       | MIT          |
| `charm.land/lipgloss/v2`      | MIT          |
| `github.com/alexflint/go-arg` | BSD-2-Clause |
| `github.com/go-git/go-git/v5` | Apache-2.0   |
| `github.com/stretchr/testify` | MIT          |

Special thanks are also due to the contributors and maintainers of Git. This
program (and let's face it, modern software development in general) wouldn't be
possible without their hard work.

## License

Copyright (C) 2026 Joseph Skubal

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <https://www.gnu.org/licenses/>.
