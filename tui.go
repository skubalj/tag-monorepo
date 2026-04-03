/*
tag-monorepo: Per-module Monorepo Tagging
Copyright (C) 2026 Joseph Skubal

See the COPYING file for copyright information
*/

package main

import (
	"cmp"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	Rows           []Row
	cursor         int
	viewportWidth  int
	viewportHeight int
	scroll         int
}

func initialModel(rows []Row) tea.Model {
	return model{
		Rows: rows,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewportHeight = msg.Width
		m.viewportHeight = msg.Height - 2 // For top/bottom padding
		if m.viewportHeight+m.scroll > len(m.Rows) {
			m.scroll = max(0, len(m.Rows)-m.viewportHeight)
		} else if m.cursor >= m.scroll+m.viewportHeight {
			m.scroll = m.cursor - m.viewportHeight + 1
		}

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "escape":
			m.Rows = nil
			return m, tea.Quit

		case "enter":
			return m, tea.Quit

		case "0", "`":
			m.Rows[m.cursor].AppliedChange = Update_None
		case "1":
			m.Rows[m.cursor].AppliedChange = Update_Major
		case "2":
			m.Rows[m.cursor].AppliedChange = Update_Minor
		case "3":
			m.Rows[m.cursor].AppliedChange = Update_Patch

		// The "up" and "k" keys move the cursor up
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll--
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down":
			if m.cursor < len(m.Rows)-1 {
				m.cursor++
				if m.cursor >= m.scroll+m.viewportHeight {
					m.scroll++
				}
				// if m.cursor > m.viewportHeight {
				// 	m.scroll++
				// }
			}

		case "right":
			m.Rows[m.cursor].AppliedChange = m.Rows[m.cursor].AppliedChange.Increment()

		case "left":
			m.Rows[m.cursor].AppliedChange = m.Rows[m.cursor].AppliedChange.Decrement()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true

	buf := new(strings.Builder)

	buf.WriteString(lipgloss.NewStyle().Underline(true).Render("Select Directories / Tags to Update"))
	buf.WriteRune('\n')

	rowStyle := lipgloss.NewStyle().Padding(0, 1)
	fromStyle := lipgloss.NewStyle().Inherit(rowStyle).Foreground(lipgloss.Red)
	toStyle := lipgloss.NewStyle().Inherit(rowStyle).Foreground(lipgloss.Green)

	// Subtract space for header and footer
	displayRows := m.Rows
	if int(m.viewportHeight) < len(m.Rows) {
		displayRows = m.Rows[m.scroll : m.scroll+m.viewportHeight]
	}

	for i, row := range displayRows {
		selected := i+m.scroll == m.cursor

		if selected {
			buf.WriteRune('>')
		} else {
			buf.WriteRune(' ')
		}

		var versionString strings.Builder
		if row.Version != nil {
			versionString.WriteRune('/')

			if row.AppliedChange != Update_None {
				versionString.WriteString(fromStyle.Render(row.Version.String()))
				versionString.WriteString(rowStyle.Render(" -> "))
				versionString.WriteString(toStyle.Render(row.UpdateVersion().String()))
				versionString.WriteString(rowStyle.Render(fmt.Sprintf(" (Update %s)", row.AppliedChange)))
			} else {
				versionString.WriteString(row.Version.String())
			}
		} else if row.AppliedChange != Update_None {
			versionString.WriteRune('/')
			versionString.WriteString(toStyle.Render(row.UpdateVersion().String()))
			versionString.WriteString(rowStyle.Render(fmt.Sprintf(" (Tag %s)", row.AppliedChange)))
		}

		buf.WriteString(rowStyle.Render(row.Module + versionString.String()))
		buf.WriteRune('\n')
	}

	out := lipgloss.NewStyle().Foreground(lipgloss.BrightBlack).Render(`0: Clear Update   1: Update Major   2: Update Minor   3: Update Patch   Enter: Accept     Q: Abort`)
	buf.WriteString(out)

	v.SetContent(buf.String())
	return v
}

type Row struct {
	Module        string
	Version       *Version
	AppliedChange Update
	AppliedSuffix string
	Changed       bool
}

func (r Row) NeedsUpdate() bool {
	return r.AppliedChange != Update_None || r.AppliedSuffix != r.Version.Suffix
}

func (r Row) UpdateTagName() string {
	return r.Module + "/" + r.UpdateVersion().String()
}

func (r Row) UpdateVersion() Version {
	var v Version
	if r.Version != nil {
		v = *r.Version
	}

	return v.Apply(r.AppliedChange, r.AppliedSuffix)
}

type Update int

const (
	Update_None Update = iota
	Update_Major
	Update_Minor
	Update_Patch

	update_count
)

func (u Update) Increment() Update {
	return (u + 1) % update_count
}

func (u Update) Decrement() Update {
	return (u + update_count - 1) % update_count
}

func (u Update) String() string {
	switch u {
	case Update_None:
		return "None"
	case Update_Major:
		return "Major"
	case Update_Minor:
		return "Minor"
	case Update_Patch:
		return "Patch"
	default:
		panic(fmt.Errorf("unknown variant '%d'", u))
	}
}

type Version struct {
	Major  int
	Minor  int
	Patch  int
	Suffix string
}

var versionRegex = regexp.MustCompile(`v(\d+)(?:\.(\d+)(?:\.(\d+))?)?(?:-([-\w\d]+))?`)

func ParseVersion(s string) (Version, bool) {
	matches := versionRegex.FindStringSubmatch(s)
	if matches == nil {
		return Version{}, false
	}

	var major, minor, patch int64
	var err1, err2, err3 error
	major, err1 = strconv.ParseInt(matches[1], 10, 0)
	if matches[2] != "" {
		minor, err2 = strconv.ParseInt(matches[2], 10, 0)
		if matches[3] != "" {
			patch, err3 = strconv.ParseInt(matches[3], 10, 0)
		}
	}
	if cmp.Or(err1, err2, err3) != nil {
		return Version{}, false
	}

	return Version{
		Major:  int(major),
		Minor:  int(minor),
		Patch:  int(patch),
		Suffix: matches[4],
	}, true
}

func (v Version) Apply(u Update, suffix string) Version {
	switch u {
	case Update_None:
	case Update_Major:
		v.Major += 1
		v.Minor = 0
		v.Patch = 0
	case Update_Minor:
		v.Minor += 1
		v.Patch = 0
	case Update_Patch:
		v.Patch += 1
	default:
		panic(fmt.Errorf("unknown variant '%d'", u))
	}

	v.Suffix = suffix
	return v
}

func (v Version) String() string {
	s := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Suffix != "" {
		s += "-" + v.Suffix
	}

	return s
}

func (v Version) Compare(other Version) int {
	return cmp.Or(v.Major-other.Major, v.Minor-other.Minor, v.Patch-other.Patch)
}
