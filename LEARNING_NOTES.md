# project-metagross Learning Reference

Topics I've learned while building project-metagross. Local only (listed in `.gitignore`).

---

## To learn

### Interfaces
- Implicit satisfaction, method sets
- Value vs pointer receivers and which type satisfies the interface
- `any` / empty interface
- Type assertions, comma-ok form, type switches
- nil interface vs interface holding a nil pointer
- Small stdlib interfaces: `io.Reader`, `io.Writer`, `io.Closer`, `fmt.Stringer`, `error`
- Interface embedding (`io.ReadWriter`)
- "Accept interfaces, return structs"
- Define interfaces at the consumer, not the producer
- Compile-time check `var _ I = (*T)(nil)`
- Interfaces for testing (fakes/mocks)
- Interfaces as generic constraints

### Error handling (deep)
- `errors.As`, custom error types
- `Unwrap() error`, `Unwrap() []error`, custom `Is` / `As` methods
- `errors.Join`, multiple `%w`
- Typed nil error gotcha
- Errors as part of a package's public API
- Errors from `defer` (`Close`), named results
- Errors in goroutines, `errgroup`
- `context.Canceled`, `context.DeadlineExceeded`, `context.Cause`
- Retryable vs permanent errors
- Testing errors (`errors.Is` in tests)

### Logging
- `log/slog` handlers, `HandlerOptions`, `AddSource`, `ReplaceAttr`
- Dynamic level with `slog.LevelVar`
- Groups, `WithGroup`, `LogValuer`
- Logger passed via dependency vs `slog.Default`
- Request-scoped logger, request ID middleware
- Logging HTTP requests (status, duration) in middleware
- Redacting secrets, what never to log
- Level policy: what goes in Debug / Info / Warn / Error
- Logs in tests

### File operations
- `os.Open`, `os.Create`, `os.OpenFile` flags and perms
- `defer f.Close()` and checking `Close` errors on writes
- `os.ReadFile`, `os.WriteFile`
- `io.Reader` / `io.Writer`, `io.Copy`, `io.ReadAll`
- `bufio.Scanner`, `bufio.Reader`, `bufio.Writer` + `Flush`
- `os.Stat`, `fs.FileInfo`, `errors.Is(err, fs.ErrNotExist)`
- `os.ReadDir`, `filepath.WalkDir`, `filepath.Glob`
- `path/filepath` vs `path`; `Join`, `Base`, `Dir`, `Ext`, `Clean`, `Rel`
- `os.Remove`, `os.RemoveAll`, `os.Rename`
- Atomic writes: `os.CreateTemp` + rename
- `os.Root` / path traversal protection
- `embed`, `io/fs`, `fs.FS`, `testing/fstest.MapFS`

### `time` package
- `time.Time`, `time.Duration`, `time.Now`, `time.Since`, `time.Until`
- `Add`, `Sub`, `AddDate`, `Before`, `After`, `Equal` vs `==`
- Layout reference date `2006-01-02 15:04:05`, `Format`, `Parse`, `time.RFC3339`
- Time zones: `UTC`, `Local`, `time.LoadLocation`, `In`
- Monotonic clock reading
- Zero value, `IsZero`
- `Truncate`, `Round`, `Unix`, `UnixMilli`
- `time.Timer`, `time.Ticker`, `time.After`, `time.AfterFunc`, `Stop`
- `time.Sleep` vs `select` with context
- `context.WithTimeout` / `WithDeadline`
- Durations in config (`time.ParseDuration`)
- Testing time: injecting a clock, `testing/synctest`

### Observability and tracing
- Three pillars: logs, metrics, traces
- `net/http/pprof`: CPU, heap, goroutine, mutex, block profiles
- `go tool pprof`, flame graphs
- Benchmarks with `-cpuprofile`, `-memprofile`
- `runtime/trace`, `go tool trace`, flight recorder
- `runtime/metrics`, `expvar`
- `GODEBUG=gctrace=1`
- Prometheus `client_golang`: counters, gauges, histograms, `/metrics`
- OpenTelemetry Go: tracer, spans, attributes, exporters, `otelhttp`
- Context propagation across calls and services
- Trace/request IDs in logs
- Health vs readiness endpoints

---

## Go language

### Pointers and functions
- `*T`, `&x`, `*p`, `nil`
- `&T{}`, `new(T)`, `new(expr)` (Go 1.26)
- No pointer arithmetic
- Arguments are always passed by value (copied)
- When to take a `*T` parameter: mutate, large struct, no-copy types (`sync.Mutex`), optional/nil
- Slices, maps, channels, funcs: copied header, shared underlying data
- `append` inside a function doesn't change the caller's slice length; return the new slice
- Returning `&local` is safe: escape analysis, heap allocation, `go build -gcflags=-m`
- Value vs pointer receivers; keep receivers consistent
- Method sets: `T` has only value-receiver methods, `*T` has both; interface satisfaction
- Functions as values, passing functions as arguments
- Closures capture variables, not values
- Per-iteration loop variables (Go 1.22)
- Nil pointer dereference panics; nil-receiver method calls

### Method signatures and interfaces
- Reading `func (Receiver) Name(params) (results)`
- Unnamed receivers and parameters
- Can't declare methods on interface types or on types from other packages
- Interfaces are satisfied implicitly (no `implements`)
- `http.ResponseWriter` satisfies `io.Writer`
- Interface values vs pointers to structs (`http.ResponseWriter` vs `*http.Request`)

### Type conversions and `[]byte(s)`
- `T(x)` conversion syntax
- No automatic conversions, even between numeric types
- `int(3.9)` truncates
- `string` (read-only) vs `[]byte` (modifiable); `[]byte(s)` and `string(b)` copy
- `len` counts bytes, not characters; `[]rune(s)` gives code points
- Compiler skips the copy in some cases (map lookups, range)
- Conversion `T(x)` vs type assertion `x.(T)`
- `string(65)` is `"A"`; use `strconv.Itoa`; `go vet` stringintconv check
- `io.WriteString`, `fmt.Fprint` as alternatives to converting yourself

### Slices
- Array `[N]T` (fixed, value type) vs slice `[]T` (view: pointer, `len`, `cap` over a backing array)
- Create: literal `[]T{...}`, `make([]T, len, cap)`, slicing; pre-allocate with `make([]T, 0, n)`
- `nil` vs empty slice: both `len` 0; JSON `null` vs `[]`
- `s[low:high]` shares the backing array (aliasing); `s[low:high:max]` limits `cap`
- `s = append(s, x)`: always use the result; reallocates past `cap`
- `copy(dst, src)` copies `min(len(dst), len(src))`
- Passed by value: header copied, elements shared
- `for i, v := range s`: `v` is a copy
- Index out of range panics; compare contents with `slices.Equal`, not `==`
- Small slice of a big array keeps it alive; `slices.Clone`
- `slices`: `Contains`, `Index`, `Sort`, `Clone`, `Equal`, `Grow`, `Delete`, `Insert`; `clear(s)`
- 2D slices `[][]T`

### Function parameters and types
- Every parameter needs a type; no type inference for parameters
- Shared type grouping: `a, b string`
- Unnamed parameters (types only), e.g. in func types / interface methods
- Error: `mixed named and unnamed parameters`
- Lone identifier in a parameter list is read as a type
- Cross-package types: `pkg.Type` + import
- Accept only what you need (narrow parameters, fewer dependencies)
- Passing structs by value (copy) vs pointer
- Unused import is a compile error

### Source files and packages
- Every `.go` file starts with `package` (only comments/build tags before); error `expected 'package', found '<token>'`
- Directory = package (not file = module like JS/Python); all files share one package name
- Error `found packages main (main.go) and parser (parser.go)` / `MismatchedPkgName`; `_test.go` may use `package foo_test`
- Files in one package see each other's identifiers, unexported too, with no import; each file has its own `import` block
- Duplicate declarations across files: `redeclared in this block`
- Import path = module path + directory; package name = last directory element (convention)
- Naming: short, lowercase, no underscores, no stdlib clashes (`go/parser`), no stutter (`parser.Config`)
- Cross-package use needs exported names and fields; exported func returning an unexported type is a lint warning
- `package main` + `func main()` = executable; can't be imported
- `go run ./cmd/migrate` vs `go run cmd/migrate/main.go` (`undefined:` with file args)
- Editor auto-import can pick the wrong package
- Split into a package for reuse, a clear API, isolated tests; otherwise YAGNI. Cost: API surface, exports, import path
- `cmd/<binary>/`: one thin `package main` per binary
- `internal/`: importable only under its parent directory; root `internal/` (module-wide) vs `cmd/<binary>/internal/` (one binary)
- `cmd/migrate/<pkg>/` is legal but not protected
- Libraries return errors; `Must*` / panics belong in `main`
- Official "Organizing a Go module" doc vs community `project-layout` repo (not official)

---

## Standard library: `net/http`

### ServeMux and DefaultServeMux
- `http.Handler` interface, `http.HandlerFunc` adapter; `ServeMux` is a router and itself a `Handler`
- `http.NewServeMux`, `Handle`, `HandleFunc`
- Avoid `DefaultServeMux` (`http.HandleFunc`, `ListenAndServe(addr, nil)`): global state, `net/http/pprof` side effects, testing
- Patterns (Go 1.22+): `METHOD HOST/PATH`; `GET` also matches `HEAD`; 405 with `Allow`
- Subtree `/path/`, exact `/{$}`, catch-all `/`; wildcards `{id}`, `{path...}`, `r.PathValue` (full segments only)
- Most specific wins; conflicting or duplicate patterns panic
- 301 to subtree, path cleaning, 404 on no match
- Middleware wraps the mux; `GODEBUG=httpmuxgo121=1` for old routing

### `http.ResponseWriter`
- `Header()`, `Write([]byte)`, `WriteHeader(code)`
- Order: headers, then status, then body
- Implicit 200 on first `Write`
- "superfluous response.WriteHeader call"
- Content-Type sniffing (`http.DetectContentType`)
- `Write` can fail (client disconnect); log it or ignore it explicitly

---

## Standard library: `strings`

### `HasPrefix`, `TrimLeft`, `TrimPrefix`, `Cut`
- `strings.HasPrefix(s, prefix) bool`: empty prefix always true, case-sensitive
- `strings.TrimLeft(s, cutset)`: cutset is a set of characters, not a prefix; removes repeatedly
- `TrimLeft` vs `TrimPrefix` gotcha
- `strings.TrimPrefix(s, prefix)`: removes once, unchanged if absent
- `strings.Cut(s, sep) (before, after, found)`: first occurrence only
- `Cut` not found: `before = s`, `after = ""`
- `strings.CutPrefix` / `CutSuffix` (Go 1.20+)
- `TrimSpace`, `TrimFunc`, `TrimRight`, `TrimSuffix`

---

## Standard library: `errors`

### Sentinel errors, `errors.New`, `errors.Is`
- `error` interface: `Error() string`
- `errors.New` makes a distinct value each call; compared by identity, not message
- Sentinel: package-level `var ErrX = errors.New("pkg: ...")`; lowercase, no punctuation
- Wrap with `fmt.Errorf("...: %w", err)`; `%v` doesn't wrap; `==` breaks once wrapped
- `errors.Is(err, target)` walks the `Unwrap` chain; `errors.Is(nil, nil)` true
- `errors.Unwrap`, multiple `%w` / `errors.Join` (Go 1.20+), custom `Is(target error) bool`
- `errors.As` for typed errors; sentinel vs typed vs opaque
- Sentinels are public API; never compare `err.Error()` strings
- Stdlib: `io.EOF`, `sql.ErrNoRows`, `os.ErrNotExist`, `http.ErrServerClosed`, `context.Canceled`, `context.DeadlineExceeded`

---

## Standard library: `os`

### `os.MkdirAll` and Unix permissions
- `os.MkdirAll(path, perm)`: like `mkdir -p`; existing dir → `nil`; existing file → error (`ENOTDIR`)
- `os.Mkdir`: no parents, fails if it exists (`errors.Is(err, fs.ErrExist)`)
- `perm` applies only to directories it creates
- `fs.FileMode`; octal `0o755` (legacy `0755`); `755` is decimal (wrong)
- Owner / group / others; `r=4 w=2 x=1`; `rwxr-xr-x` = 755, `rw-r--r--` = 644
- On dirs: `r` list, `w` create/delete entries, `x` enter; dirs need `x`, files usually don't
- umask: final = `perm &^ umask` (022 vs 002); `os.Chmod` sets the exact mode
- Common: 755/644 normal, 700/600 private; `ls -l`, `stat`, `chmod`
- `os.OpenFile` flags `O_CREATE`, `O_EXCL`, `O_WRONLY`; perm only on create

---

## Output, logging and errors

### Output channels
- stdout (program output) vs stderr (diagnostics) vs exit code
- Piping/redirection only captures stdout (`>`, `|`, `2>`)
- Exit codes: 0 success, 1 runtime failure, 2 usage error
- `return` from `main` = exit 0

### `fmt` print family
- `Print` / `Println` / `Printf`: to stdout
- `Sprint*`: return a string
- `Fprint*`: to any `io.Writer` (`os.Stderr`, files, `http.ResponseWriter`)
- `Errorf`: builds an error
- `Print` spacing rule (only between non-string operands) vs `Println` (always spaces + newline)
- `Printf` needs explicit `\n`
- Verbs: `%v`, `%+v`, `%#v`, `%T`, `%s`, `%q`, `%d`, `%w`
- `go vet` printf checks

### `log` package
- Default logger: stderr, date/time prefix, adds newline
- `log.Print*`, `log.Fatal*` (log + `os.Exit(1)`), `log.Panic*` (log + `panic`)
- `log.SetFlags` (`LstdFlags`, `Lshortfile`, `Lmicroseconds`, `LUTC`), `SetPrefix`, `SetOutput`, `log.New`
- `log.Fatal` skips deferred calls

### `log/slog` (Go 1.21+)
- Structured logging: message + key/value attributes
- Levels: `Debug`, `Info`, `Warn`, `Error`; default level Info
- Handlers: default (via `log`), `TextHandler`, `JSONHandler`; `HandlerOptions{Level}`
- `slog.New`, `slog.SetDefault`, `logger.With`
- `slog.String`, `slog.Int`, `slog.Any` typed attrs; `"err", err`
- No Fatal: log `Error` then `os.Exit(1)`
- `sloglint` / vet check for odd key-value pairs

### `panic`, `recover`, `os.Exit`
- `panic`: unwinds stack, runs defers, stack trace, exit code 2
- For programmer bugs / impossible states, not expected failures
- `recover` only works inside a deferred function
- `net/http` recovers handler panics per request
- `Must*` convention: panic only at startup / with constant input
- `os.Exit`: immediate, skips defers, no stack trace

### Error handling strategy
- Errors are values; return them, don't print them in libraries
- Library packages: no printing, logging, `os.Exit`, `log.Fatal`
- `main` decides: message, stream, exit code
- Wrap once per layer with context (`%w`); `%v` at the final print
- Handle an error once: log OR return, not both
- Message style: lowercase, no punctuation, no "failed to" chains
- `errors.New` vs `fmt.Errorf` without verbs
- Never discard errors silently (`_ =` only on purpose)
- `http.ErrServerClosed` after `Shutdown` is not a failure
- HTTP handlers: `http.Error` for the client + log for operators; don't leak internals

---

## Tooling: database migrations

### `//go:embed` and `//go:` directives
- `//go:name` with no space after `//` is a directive; `// go:embed` is a plain comment
- `//go:embed` sits directly above a package-level `var` (`string`, `[]byte`, `embed.FS`); needs `import "embed"` (`_ "embed"` for `string` / `[]byte`)
- Resolved at build time: files must exist and the pattern must match (`no matching files found`)
- Embedded files are a snapshot; rebuild after adding files (`go run` hides this)
- `var FS embed.FS` without the directive is an empty filesystem
- Patterns can't use `.` / `..`
- Other directives: `//go:build`, `//go:generate`, `//go:noinline`, `//go:linkname`
- `go doc embed`, `go help buildconstraint`, `go help generate`

### golang-migrate basics
- Sources: `iofs` (embedded) vs `file://` (disk at runtime); dev CLI reads disk, app binary embeds
- `iofs` skips file names that aren't migrations
- `migrations/` with `embed.go` is a package: commit it, not generated output
- `schema_migrations` table: `version`, `dirty`; empty = no version
- `m.Up()`, `m.Down()`, `m.Steps(±n)` relative, `m.Migrate(uint)` absolute, `m.Force(int)`
- `force` runs no SQL; only rewrites the version and clears `dirty`
- `force -1` = no version; `force V` = "V fully applied"; `< -1` rejected (`ErrInvalidVersion`)
- Forcing a version with no migration file breaks later commands
- `iofs.New` takes any `fs.FS`: `embed.FS` (build time) or `os.DirFS(dir)` (run time)
- Accept the interface (`fs.FS`) as a parameter, let `main` choose the implementation
- A flag that only half the commands honour is a bug; remove it or make every command use it
- Relative paths resolve against the current working directory, not the binary or repo

### Writing migrations
- Never run `up` on empty files: the version is recorded anyway
- Never edit an applied migration; add a new one
- Down file undoes exactly the up file; `IF EXISTS` hides mistakes
- Round trip `up` → `down` → `up` proves the down file
- Missing `;` after `BEGIN` merges statements into one syntax error
- Postgres DDL is transactional; a failed file leaves nothing behind
- Schema basics: identity vs UUID keys, `NOT NULL`, `CHECK`, money as cents / `numeric`, `timestamptz`
- `psql`: `\d <table>`

### Dirty state and recovery
- Failed migration → version marked `dirty`; tool refuses to continue
- Check the real DB state first, then fix the SQL
- Make record and schema agree: `force` to the last truly applied version, or fix the schema by hand (local only)
- Lying with `force` → next `up` fails (`relation already exists`) → dirty again
- Shared/prod DBs: every change goes through a migration

### CLI design lessons
- `create` needs no DB: return right after generating files
- Every dispatcher branch ends in `return`; no fall-through
- One `--dir` value drives `create` and `up`/`down`
- Dangerous commands require their argument; no defaults
- Check presence (`len(args)`), not a sentinel value that is also valid input
- Validation range matches the domain (`down`: `> 0`, `force`: `>= -1`)
- Check what the library already validates before adding your own checks
- Copy-pasted guards carry the wrong assumptions
- `defer` right after a successful open; `os.Exit(realMain())` lets defers run
- `strconv.Atoi` returns `int`; `uint(v)` for `m.Migrate`
- Library code doesn't print; printing an `*os.File` shows only a pointer
- Read error output in order: the last line before the error tells you where it failed
- Importing a missing package breaks the whole binary build
- Test every documented command form

---

## Tooling: Make

### Makefile basics
- Rule: `target: prerequisites` + recipe lines
- Recipe lines start with a TAB, not spaces (`missing separator`)
- `make` runs the first target (default goal); `make <target>` runs one
- Prerequisites run first (`run: build`)
- `.PHONY`: target is a command, not a file; list every non-file target
- Make rebuilds by file timestamps; Go's build cache does the real caching
- `@` hides the command echo; `-` ignores a failing line
- Each recipe line runs in its own shell (`cd` doesn't carry over)
- Variables: `NAME ?= default`, `$(NAME)`, override with `make target NAME=value`
- `$(error ...)` to require a variable
- `go run` (rebuild each time) vs built binary (embedded snapshot)
- Run from the repo root → stable relative paths; `make -C <dir>`
- `make -n` dry run; `help` target convention

---

## Learning approach

### When to study the references
- Learn things when you need them, not up front
- Skim the topic soon; open the docs when stuck, before a design decision, or when something works and you don't know why
- Weekly or milestone deep read, then review your own code
- Doc sources by effort: Tour < Effective Go / FAQ / pkg docs / blog < spec
- Put off performance internals until you're measuring performance

