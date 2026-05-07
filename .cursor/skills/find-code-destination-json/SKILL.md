---
name: find-code-destination-json
description: >-
  When the user asks to find or locate a class, function, constant, or named
  symbol in the codebase (by symbol name), run
  app_agent/explore_source_code_agent/base/find_code_destination_in_json.sh
  with --name only—do not read or load base/.../cb.json into
  context. Default script output is compact JSON (file_path, line_start,
  line_end, address, short line_code). Parse stdout then read the source file
  at file_path for that line range with a small margin, not the whole file.
  Use --full on the script only if docstring, constant_value, or call graph
  fields from the index are explicitly needed.
---

# Find code via `cb.json` (symbol lookup)

## Critical: use the shell script only for the index

- **Do not** use `read_file`, `cat`, or search tools to load **`cb.json`** into the conversation. It can be very large and wastes tokens.
- **Do** run **`find_code_destination_in_json.sh --name …`** (default **compact** output). The script reads the JSON on disk and prints a **small filtered array** to stdout; that output is the only index payload you need.
- After that, read **the source file** at `file_path` for **`line_start`–`line_end`** plus a **small margin** (see below)—not the whole file, not the whole index.

## When to use this skill

- The user wants to **find** or **locate** a **named symbol**: class, function, method, constant, variable name, etc.
- Prefer this workflow **before** using `grep` / ripgrep / semantic search for the same purpose. Assume the index exists if the script runs successfully; you do not preflight by opening `code_base.json`.

## When *not* to rely on it

- **Finding a file by path or folder name** — the script matches JSON entries by the symbol **`name`** field only, not file paths. For path-based search, use normal file tools or search.
- **`cb.json` or `jq` is missing** — the script exits with an error; fall back to grep / file search / read_file as usual.

## Command (run from repository root)

**Default (compact — prefer this):** small JSON objects; `line_code` is truncated (default max 120 chars; override with env `LINE_CODE_MAX` if needed).

```bash
bash app_agent/explore_source_code_agent/base/find_code_destination_in_json.sh --name "<SymbolName>"
```

**Full index rows (large — only when necessary):** e.g. user needs `docstring`, `constant_value`, `calls_to` / `called_by` from the index without opening the source file yet.

```bash
bash codebase/find_code_destination_in_json.sh --name "<SymbolName>" --full
```

Replace `<SymbolName>` with the exact identifier. Quote the name if it contains special shell characters.

## Interpreting the output

**Compact (default)** — each element includes:

| Field | Use |
|--------|-----|
| `name` | Symbol name |
| `kind` | e.g. `constant`, `class`, `function` |
| `file_path` | Path to the source file **relative to repo root** |
| `line_start`, `line_end` | Line range in that file |
| `address` | Stable id, e.g. `path::Symbol` |
| `line_code` | Short preview only (may end with `…`) — **do not** treat as full source |

There may be **multiple** elements (e.g. several `__init__`). Use `file_path`, `kind`, and `address` to disambiguate.

## Read source efficiently

1. Parse the JSON on stdout.
2. For each match, call **`read_file`** on `file_path` with:
   - **offset** ≈ `max(1, line_start - 10)` (or a smaller margin if the file is huge and context is not needed),
   - **limit** ≈ `(line_end - line_start + 1) + 20` so the span plus small context fits one read.
3. Expand the range only if the question needs imports above the block or code below `line_end`.

## What to do next (analysis)

1. Prefer **source** from step above over trusting only `line_code` from JSON.
2. Do **not** stop at the shell output alone if the user asked for behavior or callers — follow imports / call sites in the source tree as needed (grep or index `--full` for `called_by` if appropriate).

## Regeneration note

If lookups are empty or stale, regenerate `cb.json` via the repo pipeline; until then, use normal search. Still **do not** load the full JSON into context to “debug” — rerun the script or grep source.
