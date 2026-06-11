/*
tag-monorepo: Per-module Monorepo Tagging
Copyright (C) 2026 Joseph Skubal

See the COPYING file for copyright information
*/

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alexflint/go-arg"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Args struct {
	Copyright      bool   `arg:"-c,--copyright" help:"display GPL copyright notice"`
	MaxDepth       int    `arg:"-d,--max-depth" help:"show all directories up to the given depth"`
	RepositoryPath string `arg:"positional" default:"." placeholder:"REPOSITORY_PATH" help:"path to the git repo to be searched"`
}

func (Args) Version() string {
	return "tag-monorepo 0.0.1"
}

func (Args) Epilogue() string {
	return `This program is free software released under the GNU GPLv3
Copyright (C) 2026 Joseph Skubal`
}

const gplCopyrightNotice = `tag-monorepo: Per-module Monorepo Tagging
Copyright (C) 2026 Joseph Skubal

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

func main() {
	var args Args
	arg.MustParse(&args)

	if args.Copyright {
		fmt.Println(gplCopyrightNotice)
		return
	}

	rows, err := getTags(args.RepositoryPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Apply max-depth, but preserve anything with tags
	rows = slices.DeleteFunc(rows, func(r Row) bool {
		return r.Version == nil && strings.Count(r.Module, "/") >= args.MaxDepth
	})

	p := tea.NewProgram(initialModel(rows))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	typedModel, ok := m.(*model)
	if !ok {
		return
	}

	var tags []string
	for _, row := range typedModel.Rows {
		if row.NeedsUpdate() {
			tags = append(tags, row.UpdateTagName())
		}
	}

	err = createTags(args.RepositoryPath, tags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Get the given tags
func getTags(repositoryPath string) ([]Row, error) {
	var modules []string
	err := fs.WalkDir(os.DirFS(repositoryPath), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != "." && !strings.HasPrefix(path, ".git") {
			modules = append(modules, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Inspect the git repo and analyze the tags
	repo, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return nil, err
	}

	tags, err := getMostRecentTagNames(repo)
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(modules))
	for _, module := range modules {
		version, ok := tags[module]
		if ok {

			rows = append(rows, Row{
				Module:        module,
				Version:       &version.Version,
				AppliedSuffix: version.Version.Suffix,
				Changed:       false, // reserved for future use
			})
		} else {
			rows = append(rows, Row{Module: module})
		}
	}

	return rows, err
}

type tagVersion struct {
	Module  string
	Version Version
	Hash    plumbing.Hash
}

func getMostRecentTagNames(repo *git.Repository) (map[string]tagVersion, error) {
	tags := make(map[string]tagVersion)
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	err = iter.ForEach(func(r *plumbing.Reference) error {
		tagName := r.Name().Short()
		moduleName := path.Dir(tagName)
		if moduleName != "." {
			version, ok := ParseVersion(path.Base(tagName))
			if ok {
				currentVersion, ok := tags[moduleName]
				if !ok || version.Compare(currentVersion.Version) > 0 {
					tags[moduleName] = tagVersion{
						Module:  moduleName,
						Version: version,
						Hash:    r.Hash(),
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Create the given tags on top of Head
func createTags(path string, tags []string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		fmt.Printf("Creating tag '%s'\n", tag)
		_, err = repo.CreateTag(tag, head.Hash(), nil)
		if err != nil {
			return err
		}
	}

	return nil
}
