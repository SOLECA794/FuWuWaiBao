# -*- coding: utf-8 -*-
"""Normalize c.JSON(http.StatusOK/BadRequest/...) to pkg/apiresp in internal/handler/*.go."""
from __future__ import annotations

import re
import sys
from pathlib import Path

HANDLER_DIR = Path(__file__).resolve().parents[1] / "internal" / "handler"

LINE_SUBS: list[tuple[re.Pattern[str], str]] = [
    (
        re.compile(
            r'c\.JSON\(http\.StatusBadRequest,\s*gin\.H\{"code":\s*400,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.BadRequest(c, "\1", "")',
    ),
    (
        re.compile(
            r'c\.JSON\(http\.StatusInternalServerError,\s*gin\.H\{"code":\s*500,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.Internal(c, "\1", "")',
    ),
    (
        re.compile(
            r'c\.JSON\(http\.StatusNotFound,\s*gin\.H\{"code":\s*404,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.NotFound(c, "\1", "")',
    ),
    (
        re.compile(
            r'c\.JSON\(http\.StatusServiceUnavailable,\s*gin\.H\{"code":\s*503,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.ServiceUnavailable(c, "\1", "")',
    ),
    (
        re.compile(
            r'c\.JSON\(http\.StatusUnauthorized,\s*gin\.H\{"code":\s*401,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.Unauthorized(c, "\1", "")',
    ),
    (
        re.compile(
            r'c\.JSON\(http\.StatusForbidden,\s*gin\.H\{"code":\s*403,\s*"message":\s*"([^"]*)"\s*\}\)'
        ),
        r'apiresp.Forbidden(c, "\1", "")',
    ),
]


def ensure_apiresp_import(src: str) -> str:
    if "smart-teaching-backend/pkg/apiresp" in src:
        return src
    if "apiresp." not in src:
        return src
    needle = '\t"github.com/gin-gonic/gin"\n'
    ins = '\t"github.com/gin-gonic/gin"\n\n\t"smart-teaching-backend/pkg/apiresp"\n'
    if needle in src:
        return src.replace(needle, ins, 1)
    return src


def extract_quoted_string(s: str, key: str) -> str | None:
    m = re.search(r'"' + re.escape(key) + r'"\s*:\s*"([^"]*)"', s)
    return m.group(1) if m else None


def parse_expr_at(s: str, pos: int) -> tuple[str, int] | None:
    n = len(s)
    while pos < n and s[pos] in " \t\n":
        pos += 1
    if pos >= n:
        return None
    if s.startswith("gin.H{", pos):
        start = pos + len("gin.H{")
        d = 1
        p = start
        while p < n and d > 0:
            if s[p] == "{":
                d += 1
            elif s[p] == "}":
                d -= 1
            p += 1
        return s[pos:p], p
    if s[pos] == "[":
        d = 1
        p = pos + 1
        while p < n and d > 0:
            if s[p] == "[":
                d += 1
            elif s[p] == "]":
                d -= 1
            p += 1
        return s[pos:p], p
    m = re.match(r"([a-zA-Z_][a-zA-Z0-9_.]*)", s[pos:])
    if m:
        return m.group(1), pos + len(m.group(1))
    return None


def ok_inner_to_apiresp(inner: str) -> str | None:
    t = inner.strip()
    if not re.search(r'"code"\s*:\s*200', t):
        return None
    dm = re.search(r'"data"\s*:\s*', t)
    msg = "请求成功"
    if dm:
        head = t[: dm.start()]
        mm = re.search(r'"message"\s*:\s*"([^"]*)"', head)
        if mm:
            msg = mm.group(1)
    else:
        mm = re.search(r'"message"\s*:\s*"([^"]*)"', t)
        if mm:
            return f"apiresp.OKMessage(c, {repr(mm.group(1))})"
        return None
    pos = dm.end()
    parsed = parse_expr_at(t, pos)
    if not parsed:
        return None
    expr, end = parsed
    rest = t[end:].strip()
    if rest not in ("", ","):
        if rest.startswith(","):
            rest2 = rest[1:].strip()
            if rest2:
                return None
    return f"apiresp.OK(c, {repr(msg)}, {expr})"


def replace_status_ok_calls(src: str) -> str:
    out: list[str] = []
    i = 0
    n = len(src)
    needle = "c.JSON(http.StatusOK, gin.H{"
    while i < n:
        start = src.find(needle, i)
        if start == -1:
            out.append(src[i:])
            break
        out.append(src[i:start])
        j = start + len(needle)
        depth = 1
        k = j
        while k < n and depth > 0:
            if src[k] == "{":
                depth += 1
            elif src[k] == "}":
                depth -= 1
            k += 1
        if depth != 0:
            out.append(src[start : start + 20])
            i = start + 1
            continue
        inner = src[j : k - 1]
        rest = src[k:].lstrip()
        if not rest.startswith(")"):
            out.append(src[start])
            i = start + 1
            continue
        after = k + rest.find(")") + 1
        rep = ok_inner_to_apiresp(inner)
        if rep:
            out.append(rep)
            i = after
        else:
            out.append(src[start:after])
            i = after
    return "".join(out)


def process_file(path: Path) -> bool:
    if path.name == "compat_common.go":
        return False
    raw = path.read_text(encoding="utf-8")
    s = raw
    for pat, repl in LINE_SUBS:
        s = pat.sub(repl, s)
    s = replace_status_ok_calls(s)
    s = ensure_apiresp_import(s)
    if s != raw:
        path.write_text(s, encoding="utf-8")
        return True
    return False


def main() -> int:
    changed = []
    for path in sorted(HANDLER_DIR.glob("*.go")):
        if process_file(path):
            changed.append(path.name)
    print("updated:", ", ".join(changed) if changed else "(none)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
