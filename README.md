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

## Installation

```bash
go install github.com/skubalj/tag-monorepo
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
