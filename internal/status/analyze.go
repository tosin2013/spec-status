package status

import (
	"errors"
	"path"
	"sort"
	"strings"

	"github.com/tosin2013/spec-status/internal/safeio"
)

func build(rel string, items []rawItem, fs *safeio.FS) Report {
	rep := Report{
		Outcome:     "report",
		Complete:    true,
		Roadmap:     rel,
		InputFormat: inputFormat,
		Items:       []Item{},
	}
	byID := map[string][]int{}
	for i, it := range items {
		byID[it.ID] = append(byID[it.ID], i)
	}
	supers := make([]Supersession, len(items))
	superCode := make([]string, len(items))
	for i, it := range items {
		supers[i], superCode[i] = classifySupersession(it, items)
	}
	valid := make([]bool, len(items))
	for i := range items {
		if !items[i].HasSuper {
			continue
		}
		if superCode[i] == "" {
			cyc, bad := supersessionWalk(i, items)
			if cyc {
				superCode[i] = "SS_SUPERSESSION_CYCLE"
			} else if bad {
				superCode[i] = "SS_SUPERSESSION_CHAIN_INVALID"
			}
		}
		ok := superCode[i] == ""
		valid[i] = ok
		flag := ok
		supers[i].Valid = &flag
		if ok {
			supers[i].TerminalID = terminalID(i, items)
		}
	}

	diags := make([]Diagnostic, 0)
	for i, it := range items {
		if len(byID[it.ID]) > 1 {
			diags = append(diags, Diagnostic{
				Code: "SS_DUPLICATE_ID", Severity: "error", ItemID: it.ID, Path: rel, Line: it.Line,
				Message: "Roadmap item ID is duplicated.",
			})
		}
		if !knownStatus[it.Status] {
			diags = append(diags, Diagnostic{
				Code: "SS_STATUS_UNKNOWN", Severity: "warning", ItemID: it.ID, Path: rel, Line: it.Line,
				Message: "Status is not a known lifecycle token; the raw value was kept.",
			})
		}
		if !it.HasSpec && stageLink(it.Status) != "" {
			sev := stageLink(it.Status)
			msg := "Spec link is required for this recorded status."
			if sev == "info" {
				msg = "Planned item has no spec link."
			}
			diags = append(diags, Diagnostic{
				Code: "SS_SPEC_REFERENCE_MISSING", Severity: sev, ItemID: it.ID, Path: rel, Line: it.Line, Message: msg,
			})
		}
		if superCode[i] != "" {
			line := it.SuperLine
			if line == 0 {
				line = it.Line
			}
			diags = append(diags, Diagnostic{
				Code: superCode[i], Severity: "error", ItemID: it.ID, Path: rel, Line: line,
				Message: supersessionMessage(superCode[i]),
			})
		}
	}

	cycles := dependencyCycles(items, byID)
	depStates := make([][]Dependency, len(items))
	for i, it := range items {
		if it.DepsBad {
			diags = append(diags, Diagnostic{
				Code: "SS_DEPENDENCY_MALFORMED", Severity: "error", ItemID: it.ID, Path: rel, Line: lineOr(it.DepsLine, it.Line),
				Message: "Dependency list is not a comma-separated list of item IDs.",
			})
			continue
		}
		for _, dep := range it.Deps {
			d, diag := classifyDep(i, it, dep, rel, byID, items, valid, cycles)
			depStates[i] = append(depStates[i], d)
			if diag != nil {
				diags = append(diags, *diag)
			}
			if it.Status == "verified" && !d.Satisfied {
				diags = append(diags, Diagnostic{
					Code: "SS_VERIFIED_DEPENDENCY_CONTRADICTION", Severity: "error", ItemID: it.ID, Path: rel, Line: lineOr(it.DepsLine, it.Line),
					Message: "Recorded verified item has an unsatisfied dependency.",
				})
			}
		}
	}

	taskReads := map[string]taskFile{}
	taskMissing := map[string]struct{}{}
	taskSkipped := map[string]struct{}{}
	taskIDs := map[string]string{}
	outItems := make([]Item, len(items))
	for i, it := range items {
		item := Item{
			ID: it.ID, Title: it.Title, Status: it.Status, StatusKnown: knownStatus[it.Status],
			Path: rel, Line: it.Line,
			Active:      it.Status != "deferred" && it.Status != "abandoned" && !valid[i],
			Description: it.Description, Outcome: it.Outcome, Scope: it.Scope,
			Supersession: supers[i], Dependencies: depStates[i],
		}
		if item.Dependencies == nil {
			item.Dependencies = []Dependency{}
		}
		if it.HasSuper && !it.SuperBad {
			item.SupersededBy = it.SuperID
		}
		if it.HasSuper {
			item.Supersession.Line = it.SuperLine
		}
		item.Tasks.State = "not_linked"
		if it.HasSpec {
			itemDiags, tasks, specDir := inspectSpec(fs, rel, it, &rep, taskReads, taskMissing, taskSkipped, taskIDs)
			if rep.Outcome == "fatal" {
				return rep
			}
			diags = append(diags, itemDiags...)
			item.Tasks = tasks
			item.SpecDir = specDir
		}
		outItems[i] = item
	}

	checked, total := 0, 0
	for _, tf := range taskReads {
		checked += tf.Checked
		total += tf.Total
	}
	sum := &Summary{
		TotalItems:   len(items),
		StatusCounts: map[string]int{},
		Tasks: TaskSummary{
			Checked: checked, Total: total, Percent: percentOf(checked, total),
			Scope: taskScope, FilesRead: len(taskReads), FilesMissing: len(taskMissing), FilesSkipped: len(taskSkipped),
		},
	}
	for _, it := range outItems {
		sum.StatusCounts[it.Status]++
		if it.Supersession.Supplied && it.Supersession.Valid != nil && *it.Supersession.Valid {
			sum.SupersededItems++
		}
		if it.Supersession.Supplied && it.Supersession.Valid != nil && !*it.Supersession.Valid {
			sum.InvalidSupersessionItems++
		}
		if !it.Active {
			continue
		}
		sum.ActiveItems++
		if it.StatusKnown && it.Status == "verified" {
			sum.VerifiedActiveItems++
		}
		blocked, review := false, false
		for _, d := range it.Dependencies {
			if d.State == "human_review" {
				review = true
			}
			if !d.Satisfied {
				blocked = true
			}
		}
		if blocked {
			sum.BlockedCount++
		}
		if review {
			sum.HumanReviewCount++
		}
	}
	sum.VerifiedPercent = percentOf(sum.VerifiedActiveItems, sum.ActiveItems)
	rep.Items = outItems
	rep.Summary = sum
	rep.Diagnostics = diags
	sort.SliceStable(rep.Diagnostics, func(a, b int) bool { return diagLess(rep.Diagnostics[a], rep.Diagnostics[b]) })
	return rep
}

func classifySupersession(it rawItem, items []rawItem) (Supersession, string) {
	if !it.HasSuper {
		return Supersession{}, ""
	}
	s := Supersession{Supplied: true}
	if it.SuperBad {
		s.RawReference = sanitizeRaw(it.SuperRaw)
		s.RawSet = s.RawReference != ""
		return s, "SS_SUPERSESSION_MALFORMED"
	}
	s.Parsed = it.SuperID
	s.RawReference = it.SuperID
	s.RawSet = true
	if it.SuperID == it.ID {
		return s, "SS_SUPERSESSION_SELF"
	}
	n := 0
	for i := range items {
		if items[i].ID == it.SuperID {
			n++
		}
	}
	switch {
	case n == 0:
		return s, "SS_SUPERSESSION_MISSING"
	case n > 1:
		return s, "SS_SUPERSESSION_AMBIGUOUS"
	default:
		return s, ""
	}
}

func supersessionWalk(start int, items []rawItem) (cycle bool, bad bool) {
	index := map[int]int{}
	order := make([]int, 0, 8)
	cur := start
	for {
		it := items[cur]
		if !it.HasSuper {
			return false, false
		}
		if it.SuperBad || it.SuperID == "" || it.SuperID == it.ID || unique(it.SuperID, items) < 0 {
			return false, true
		}
		if _, seen := index[cur]; seen {
			return true, true
		}
		index[cur] = len(order)
		order = append(order, cur)
		next := unique(it.SuperID, items)
		if pos, seen := index[next]; seen {
			for _, n := range order[pos:] {
				if n == start {
					return true, true
				}
			}
			if next == start {
				return true, true
			}
			return false, true
		}
		if !items[next].HasSuper {
			return false, false
		}
		cur = next
	}
}

func unique(id string, items []rawItem) int {
	found := -1
	for i := range items {
		if items[i].ID == id {
			if found >= 0 {
				return -1
			}
			found = i
		}
	}
	return found
}

func terminalID(i int, items []rawItem) string {
	seen := map[int]bool{}
	for {
		if seen[i] || !items[i].HasSuper || items[i].SuperBad {
			return ""
		}
		seen[i] = true
		next := unique(items[i].SuperID, items)
		if next < 0 {
			return ""
		}
		if !items[next].HasSuper {
			return items[next].ID
		}
		i = next
	}
}

func supersessionMessage(code string) string {
	switch code {
	case "SS_SUPERSESSION_MALFORMED":
		return "Supersession value must be one item ID."
	case "SS_SUPERSESSION_MISSING":
		return "Supersession target does not exist."
	case "SS_SUPERSESSION_AMBIGUOUS":
		return "Supersession target ID is duplicated."
	case "SS_SUPERSESSION_SELF":
		return "Item supersedes itself."
	case "SS_SUPERSESSION_CYCLE":
		return "Supersession chain is cyclic."
	default:
		return "Supersession leads into an invalid chain."
	}
}

// dependencyCycles returns each item's strongly connected component in the
// dependency graph. Two items share a component with more than one member
// exactly when each can reach the other, so an edge is part of a cycle only
// when both of its ends are in the same multi-member component.
func dependencyCycles(items []rawItem, byID map[string][]int) cycleInfo {
	next := make([][]int, len(items))
	for i, it := range items {
		if it.DepsBad {
			continue
		}
		for _, dep := range it.Deps {
			if dep == it.ID || len(byID[dep]) != 1 {
				continue
			}
			next[i] = append(next[i], byID[dep][0])
		}
	}
	info := cycleInfo{comp: make([]int, len(items)), size: map[int]int{}}
	index := make([]int, len(items))
	low := make([]int, len(items))
	onStack := make([]bool, len(items))
	for i := range index {
		index[i] = -1
	}
	var stack []int
	counter, comps := 0, 0
	var strong func(int)
	strong = func(v int) {
		index[v], low[v] = counter, counter
		counter++
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range next[v] {
			if index[w] < 0 {
				strong(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] == index[v] {
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				info.comp[w] = comps
				info.size[comps]++
				if w == v {
					break
				}
			}
			comps++
		}
	}
	for i := range items {
		if index[i] < 0 {
			strong(i)
		}
	}
	return info
}

type cycleInfo struct {
	comp []int
	size map[int]int
}

// edge reports whether the edge from src to dst lies on a dependency cycle.
func (c cycleInfo) edge(src, dst int) bool {
	return src != dst && c.comp[src] == c.comp[dst] && c.size[c.comp[src]] > 1
}

// member reports whether item i is on any dependency cycle.
func (c cycleInfo) member(i int) bool {
	return c.size[c.comp[i]] > 1
}

func classifyDep(src int, it rawItem, dep, rel string, byID map[string][]int, items []rawItem, valid []bool, cycles cycleInfo) (Dependency, *Diagnostic) {
	d := Dependency{ID: dep, Path: rel, Line: lineOr(it.DepsLine, it.Line), State: "unfinished"}
	base := Diagnostic{ItemID: it.ID, Path: rel, Line: d.Line}
	switch {
	case dep == it.ID:
		d.State = "self"
		base.Code, base.Severity, base.Message = "SS_DEPENDENCY_SELF", "error", "Item depends on itself."
		return d, &base
	case len(byID[dep]) == 0:
		d.State = "missing"
		base.Code, base.Severity, base.Message = "SS_DEPENDENCY_MISSING", "error", "Dependency ID does not exist."
		return d, &base
	case len(byID[dep]) > 1:
		d.State = "ambiguous"
		base.Code, base.Severity, base.Message = "SS_DEPENDENCY_AMBIGUOUS", "error", "Dependency ID is duplicated."
		return d, &base
	}
	target := byID[dep][0]
	tgt := items[target]
	if cycles.edge(src, target) {
		d.State = "cycle"
		base.Code, base.Severity, base.Message = "SS_DEPENDENCY_CYCLE", "error", "Dependency cycle includes this item."
		return d, &base
	}
	if valid[target] || tgt.HasSuper {
		// Any supplied superseded-by line, valid or not, requires human review
		// of its dependents (owner decision, 7 November 2026). An invalid line
		// still leaves the target active for metrics, per D-02A.
		d.State = "human_review"
		d.ReplacementID = tgt.SuperID
		msg := "Dependency on a superseded item requires human review."
		if !valid[target] {
			msg = "Dependency on an item with an invalid superseded-by line requires human review."
		}
		base.Code, base.Severity, base.Message = "SS_DEPENDENCY_HUMAN_REVIEW", "warning", msg
		base.Related = []Ref{{Path: rel, Line: tgt.Line}}
		return d, &base
	}
	if knownStatus[tgt.Status] && tgt.Status == "verified" {
		d.State = "satisfied"
		d.Satisfied = true
		return d, nil
	}
	base.Code, base.Severity, base.Message = "SS_DEPENDENCY_BLOCKED", "warning", "Dependency is not recorded verified."
	return d, &base
}

// ids maps a file identity (device and inode) to the first path it was read
// under, so one tasks file reached through two spellings, for example on a
// case-insensitive filesystem, is read and counted once.
func inspectSpec(fs *safeio.FS, roadmap string, it rawItem, rep *Report, reads map[string]taskFile, missing, skipped map[string]struct{}, ids map[string]string) ([]Diagnostic, ItemTasks, string) {
	var diags []Diagnostic
	tasks := ItemTasks{State: "not_linked"}
	specDir, err := safeio.NormalizeRel(it.SpecDir)
	if err == nil && (strings.Contains(it.SpecDir, ":") || strings.HasPrefix(it.SpecDir, "/") || strings.HasPrefix(it.SpecDir, "\\")) {
		// A URL or absolute path is not a root-relative spec directory.
		err = errUnsupportedSpecDir
	}
	if err != nil || safeio.Sensitive(it.SpecDir) || (err == nil && safeio.Sensitive(specDir)) {
		tasks.State = "skipped"
		if safeio.Sensitive(it.SpecDir) || (err == nil && safeio.Sensitive(specDir)) {
			skipped["sensitive:"+it.ID] = struct{}{}
			diags = append(diags, Diagnostic{
				Code: "SS_REFERENCE_SENSITIVE_SKIPPED", Severity: "error", ItemID: it.ID, Path: roadmap, Line: lineOr(it.SpecLine, it.Line),
				Message: "Spec reference was skipped because it matched the sensitive-name matcher.",
			})
		} else {
			diags = append(diags, Diagnostic{
				Code: "SS_PATH_UNSAFE", Severity: "error", ItemID: it.ID, Path: roadmap, Line: lineOr(it.SpecLine, it.Line),
				Message: "Spec directory path is outside the selected root.",
			})
		}
		return diags, tasks, ""
	}
	ent, err := fs.StatRel(specDir)
	if err != nil {
		if fatal, d := fsDiag(err, roadmap, it.SpecLine, it.ID, false); fatal != nil {
			*rep = *fatal
			return nil, tasks, ""
		} else if d != nil {
			diags = append(diags, *d)
			if d.Code == "SS_REFERENCE_SENSITIVE_SKIPPED" {
				tasks.State = "skipped"
				skipped[specDir] = struct{}{}
			}
		}
		return diags, tasks, ""
	}
	if !ent.Exists || !ent.Dir {
		diags = append(diags, Diagnostic{
			Code: "SS_SPEC_DIRECTORY_MISSING", Severity: "error", ItemID: it.ID, Path: roadmap, Line: lineOr(it.SpecLine, it.Line),
			Message: "Linked spec directory is absent.",
		})
		tasks.State = "missing"
		missing[specDir+"/tasks.md"] = struct{}{}
		return diags, tasks, specDir
	}
	need := stageFiles(it.Status)
	for _, f := range []struct {
		name string
		need bool
	}{{"spec.md", need.spec}, {"plan.md", need.plan}, {"tasks.md", need.tasks}} {
		p := specDir + "/" + f.name
		st, statErr := fs.StatRel(p)
		if statErr != nil {
			if fatal, d := fsDiag(statErr, roadmap, it.SpecLine, it.ID, false); fatal != nil {
				*rep = *fatal
				return nil, tasks, ""
			} else if d != nil {
				diags = append(diags, *d)
				if f.name == "tasks.md" {
					tasks.State = "skipped"
					tasks.Source = ""
					skipped[p] = struct{}{}
				}
			}
			continue
		}
		if f.name != "tasks.md" {
			if f.need && (!st.Exists || !st.Regular) {
				diags = append(diags, Diagnostic{
					Code: "SS_STAGE_FILE_MISSING", Severity: "error", ItemID: it.ID, Path: roadmap, Line: lineOr(it.SpecLine, it.Line),
					Message: "Stage-required file " + f.name + " is absent.",
				})
			}
			continue
		}
		if st.ID != "" {
			if first, seen := ids[st.ID]; seen {
				p = first
			} else {
				ids[st.ID] = p
			}
		}
		tasks.Source = p
		if !st.Exists || !st.Regular {
			tasks.State = "missing"
			missing[p] = struct{}{}
			if f.need {
				diags = append(diags, Diagnostic{
					Code: "SS_STAGE_FILE_MISSING", Severity: "error", ItemID: it.ID, Path: roadmap, Line: lineOr(it.SpecLine, it.Line),
					Message: "Stage-required file tasks.md is absent.",
				})
			}
			continue
		}
		if _, ok := reads[p]; !ok {
			data, readErr := fs.ReadRel(p)
			if readErr != nil {
				if fatal, d := fsDiag(readErr, roadmap, it.SpecLine, it.ID, true); fatal != nil {
					*rep = *fatal
					rep.Roadmap = roadmap
					return nil, tasks, ""
				} else if d != nil {
					diags = append(diags, *d)
				}
				tasks.State = "skipped"
				skipped[p] = struct{}{}
				continue
			}
			tf, ferr := parseTasks(p, data)
			if ferr != nil {
				*rep = fatalReport(ferr.Code, ferr.Message, roadmap, "")
				return nil, tasks, ""
			}
			reads[p] = tf
			total := 0
			for _, read := range reads {
				total += read.Total
			}
			if total > maxTasks {
				*rep = fatalReport("SS_INPUT_LIMIT", "Task checkbox limit exceeded across task files.", roadmap, "")
				return nil, tasks, ""
			}
			for _, link := range tf.Links {
				diags = append(diags, checkTaskLink(fs, p, it.ID, link)...)
			}
		}
		tf := reads[p]
		tasks.State = "read"
		tasks.HasNums = true
		tasks.Checked = tf.Checked
		tasks.Total = tf.Total
		tasks.Percent = percentOf(tf.Checked, tf.Total)
		tasks.Checks = append([]TaskCheck(nil), tf.Checks...)
	}
	return diags, tasks, specDir
}

func stageLink(status string) string {
	switch status {
	case "planned":
		return "info"
	case "specced", "in-progress", "implemented", "verified":
		return "error"
	default:
		return ""
	}
}

func stageFiles(status string) (n struct{ spec, plan, tasks bool }) {
	// D-03B implementation proposal, not an owner-approved lifecycle law.
	switch status {
	case "specced":
		n.spec = true
	case "in-progress", "implemented", "verified":
		n.spec, n.plan, n.tasks = true, true, true
	}
	return n
}

func checkTaskLink(fs *safeio.FS, tasksPath, itemID string, link taskLink) []Diagnostic {
	if link.Kind == "unsupported" {
		return []Diagnostic{{
			Code: "SS_REFERENCE_UNSUPPORTED", Severity: "warning", ItemID: itemID, Path: tasksPath, Line: link.Line,
			Message: "Task reference is not a supported local path and was not opened.",
		}}
	}
	dir := path.Dir(tasksPath)
	if dir == "." {
		dir = ""
	}
	combined := link.Path
	if dir != "" {
		combined = dir + "/" + link.Path
	}
	clean, err := safeio.NormalizeRel(combined)
	sensitive := safeio.Sensitive(link.Path) || safeio.Sensitive(combined) || (err == nil && safeio.Sensitive(clean))
	if err != nil || sensitive {
		code, msg := "SS_PATH_UNSAFE", "Task reference escapes the selected root."
		if sensitive {
			code, msg = "SS_REFERENCE_SENSITIVE_SKIPPED", "Task reference was skipped because it matched the sensitive-name matcher."
		}
		return []Diagnostic{{Code: code, Severity: "error", ItemID: itemID, Path: tasksPath, Line: link.Line, Message: msg}}
	}
	ent, statErr := fs.StatRel(clean)
	if statErr != nil {
		if _, d := fsDiag(statErr, tasksPath, link.Line, itemID, false); d != nil {
			return []Diagnostic{*d}
		}
		return nil
	}
	if ent.Exists && ent.Regular {
		return nil
	}
	sev, msg := "info", "Unchecked task names a file that is not present."
	if link.Checked {
		sev, msg = "error", "Checked task names a file that is not present."
	}
	return []Diagnostic{{Code: "SS_TASK_TARGET_MISSING", Severity: sev, ItemID: itemID, Path: tasksPath, Line: link.Line, Message: msg}}
}

func fsDiag(err error, filePath string, line int, itemID string, required bool) (*Report, *Diagnostic) {
	code := safeio.CodeUnavailable
	if se, ok := err.(*safeio.Error); ok {
		code = se.Code
	}
	d := Diagnostic{Severity: "error", Path: filePath, Line: line, ItemID: itemID}
	switch code {
	case safeio.CodeNotFound:
		if required {
			r := fatalReport("SS_ROADMAP_UNAVAILABLE", "Required input is unavailable.", "", "")
			return &r, nil
		}
		return nil, nil
	case safeio.CodeSensitive:
		d.Code, d.Message = "SS_REFERENCE_SENSITIVE_SKIPPED", "Reference was skipped because it matched the sensitive-name matcher."
	case safeio.CodeLimit:
		r := fatalReport("SS_INPUT_LIMIT", "Input budget exceeded.", "", "")
		return &r, nil
	case safeio.CodeChanged:
		r := fatalReport("SS_SOURCE_CHANGED", "Input changed while it was being read.", "", "")
		return &r, nil
	case safeio.CodeUnsafe:
		d.Code, d.Message = "SS_PATH_UNSAFE", "Path is outside the selected root or uses a denied symlink."
	default:
		d.Code, d.Message = "SS_PATH_UNSAFE", "Path could not be read under the confinement rules."
	}
	if required && (code == safeio.CodeUnsafe || code == safeio.CodeSensitive || code == safeio.CodeUnavailable) {
		r := fatalReport(d.Code, d.Message, "", "")
		return &r, nil
	}
	if filePath != "" && safeio.Sensitive(filePath) {
		d.Path = ""
	}
	return nil, &d
}

func diagLess(a, b Diagnostic) bool {
	if (a.Path == "") != (b.Path == "") {
		return a.Path == ""
	}
	if a.Path != b.Path {
		return a.Path < b.Path
	}
	if (a.Line == 0) != (b.Line == 0) {
		return a.Line == 0
	}
	if a.Line != b.Line {
		return a.Line < b.Line
	}
	if a.Code != b.Code {
		return a.Code < b.Code
	}
	if (a.ItemID == "") != (b.ItemID == "") {
		return a.ItemID == ""
	}
	if a.ItemID != b.ItemID {
		return a.ItemID < b.ItemID
	}
	return a.Message < b.Message
}

func lineOr(line, fallback int) int {
	if line == 0 {
		return fallback
	}
	return line
}

func sanitizeRaw(s string) string {
	s = trimSpace(s)
	if s == "" || len(s) > 80 || containsAny(s, "\r\n/") || safeio.Sensitive(s) {
		return ""
	}
	for _, r := range s {
		if r < 32 || r == 127 {
			return ""
		}
	}
	return s
}

func trimSpace(s string) string {
	i, j := 0, len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t') {
		j--
	}
	return s[i:j]
}

func containsAny(s, chars string) bool {
	for _, r := range s {
		for _, c := range chars {
			if r == c {
				return true
			}
		}
	}
	return false
}

var errUnsupportedSpecDir = errors.New("spec dir is not a root-relative path")
