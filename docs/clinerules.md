# .clinerules Template for Projects Using go-pkg

This is a template for `.clinerules` file to be used in projects that consume `go-pkg` or other utility libraries.

## Design Principles

- ✅ **Reference-based**: Points to documentation, not duplicates content
- ✅ **Multi-repository support**: Can reference multiple external libraries
- ✅ **Zero maintenance**: No need to update when functions are added
- ✅ **Model agnostic**: Works with Claude, GPT, Grok, and other AI models

---

## Template File

Copy this template to your project root as `.clinerules`:

```markdown
# AI Agent Development Rules

> Guidance for AI coding agents (Claude, GPT, Grok, etc.) working on this project.

## Core Principle

**ALWAYS check if functionality already exists in external dependencies before implementing new code.**

Before creating any utility function, validation, response handler, or helper:

1. **Search** this project's codebase
2. **Check** external dependencies documentation
3. **Reuse** existing implementations

---

## External Dependencies

This project uses the following utility libraries. Check their documentation before implementing common functionality:

### go-pkg (github.com/budimanlai/go-pkg)

**Purpose:** Common utilities for Go web applications

**Documentation:**
- **Primary:** `./vendor/github.com/budimanlai/go-pkg/docs/AI_AGENT_GUIDE.md` (if vendored)
- **Fallback:** https://github.com/budimanlai/go-pkg/blob/master/docs/AI_AGENT_GUIDE.md

**When to check:** Before implementing:
- Password hashing or authentication
- HTTP response formatting
- Input validation
- JSON/pointer/string utilities
- ID generation
- Internationalization
- Database connection management

**Rule:** If it's a common web utility, check go-pkg first.

---

### [another-lib] (github.com/owner/another-lib)

**Purpose:** [Brief description]

**Documentation:**
- **Primary:** `./vendor/github.com/owner/another-lib/docs/[GUIDE].md`
- **Fallback:** https://github.com/owner/another-lib/blob/master/docs/[GUIDE].md

**When to check:** Before implementing [specific functionality]

**Rule:** [Guidance on when to use this library]

---

### [third-lib] (github.com/owner/third-lib)

**Purpose:** [Brief description]

**Documentation:**
- **Primary:** `./vendor/github.com/owner/third-lib/docs/[GUIDE].md`
- **Fallback:** https://github.com/owner/third-lib/blob/master/docs/[GUIDE].md

**When to check:** Before implementing [specific functionality]

**Rule:** [Guidance on when to use this library]

---

## Quick Reference Commands

### Check if functionality exists in dependencies

```bash
# Search in vendor directory (if using vendoring)
grep -r "func.*[FunctionName]" ./vendor/

# Search in go.mod to see what's available
cat go.mod | grep "github.com"

# Search current project usage
grep -r "github.com/budimanlai/go-pkg" --include="*.go"
```

### View external documentation

```bash
# If vendored
cat ./vendor/github.com/budimanlai/go-pkg/docs/AI_AGENT_GUIDE.md

# If using go modules (download dependency docs)
go mod download
find $GOPATH/pkg/mod -name "AI_AGENT_GUIDE.md"
```

---

## Critical Rules

### ✅ ALWAYS DO:

1. **Check external dependencies** before implementing common utilities
2. **Read documentation** from vendor or online sources
3. **Search imports** to see what's already being used in the project
4. **Validate user input** using existing validator packages
5. **Handle errors properly** - check all error returns
6. **Use consistent patterns** shown in external library docs

### ❌ NEVER DO:

1. **Don't duplicate external library functions** - search first
2. **Don't assume functions exist** - verify in docs/source
3. **Don't skip documentation** - read guides before implementing
4. **Don't ignore go.mod** - it shows what's available
5. **Don't reinvent the wheel** - use provided utilities

---

## Development Workflow

### Starting New Feature

1. Understand requirements
2. **Check if similar functionality exists in:**
   - Current project codebase
   - External dependencies (check imports)
   - Utility libraries (read their docs)
3. Reuse or compose from existing code
4. Only implement if truly unique to this project

### Before Writing Utility Code

Ask yourself:
- Is this a common pattern? → **Check go-pkg**
- Is this authentication/validation? → **Check go-pkg**
- Is this response formatting? → **Check go-pkg**
- Is this domain-specific? → **Implement in this project**

---

## Finding Documentation

### Priority Order

1. **Vendor directory** (fastest, always available)
   ```bash
   ls -la ./vendor/github.com/budimanlai/go-pkg/docs/
   ```

2. **Online repository** (when vendor not available)
   - go-pkg: https://github.com/budimanlai/go-pkg/tree/master/docs

3. **Source code** (when docs unclear)
   ```bash
   cat ./vendor/github.com/budimanlai/go-pkg/[package]/[file].go
   ```

### Key Documentation Files

- `AI_AGENT_GUIDE.md` - Complete guide for AI agents
- `README.md` - Package overview
- `*_test.go` - Usage examples

---

## Project-Specific Guidelines

[Add any project-specific rules, patterns, or conventions here]

### Architecture Decisions

[Document why certain patterns are used in this project]

### Custom Packages

[List project-specific packages and their purposes]

---

## Version Information

- **Go Version:** [e.g., 1.21+]
- **go-pkg Version:** [check go.mod]
- **Other Dependencies:** [list critical versions]

---

## Getting Help

1. **External Lib Docs:** Check vendor or GitHub documentation
2. **Source Code:** Read actual implementation when unclear
3. **Tests:** Look at `*_test.go` for usage examples
4. **Imports:** Search project for import usage patterns

---

**Last Updated:** [Date]  
**Target AI Models:** Claude, GPT, Grok, and other code-capable LLMs  
**Maintainer:** [Your name/team]
```

---

## Usage Instructions

1. **Copy the template above** to your project root as `.clinerules`
2. **Update the External Dependencies section:**
   - Add all utility libraries your project uses
   - Update documentation paths
   - Specify when to check each library
3. **Fill in Project-Specific Guidelines section**
4. **Update Version Information**
5. **Commit to repository**

## Example for Multi-Repository Project

If your project uses 3 external libraries (go-pkg, go-utils, business-lib):

```markdown
## External Dependencies

### go-pkg (github.com/budimanlai/go-pkg)
**Documentation:** `./vendor/.../docs/AI_AGENT_GUIDE.md`
**Check before:** Common web utilities

### go-utils (github.com/company/go-utils)  
**Documentation:** `./vendor/.../docs/DEV_GUIDE.md`
**Check before:** Data processing utilities

### business-lib (github.com/company/business-lib)
**Documentation:** `./vendor/.../docs/API_GUIDE.md`
**Check before:** Business logic helpers
```

---

## Benefits

✅ **No maintenance** - Just point to external docs  
✅ **Always up-to-date** - AI reads latest documentation  
✅ **Multi-repo ready** - Support unlimited external libraries  
✅ **Model agnostic** - Works with any AI coding assistant  
✅ **Simple** - Only pointers, no content duplication

---

## Notes

- Keep `.clinerules` under 500 lines
- Focus on **where to look** not **what exists**
- Update only when adding/removing dependencies
- Let external repos manage their own function documentation
