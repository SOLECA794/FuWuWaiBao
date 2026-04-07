# -*- coding: utf-8 -*-
import re
import sys
from pathlib import Path

HANDLER_DIR = Path(__file__).resolve().parents[1] / "internal" / "handler"


def main() -> int:
    for p in sorted(HANDLER_DIR.glob("*.go")):
        s = p.read_text(encoding="utf-8")
        s2 = re.sub(r"apiresp\.OK\(c, '([^']*)', ", r'apiresp.OK(c, "\1", ', s)
        s2 = re.sub(r"apiresp\.OKMessage\(c, '([^']*)'\)", r'apiresp.OKMessage(c, "\1")', s2)
        if s2 != s:
            p.write_text(s2, encoding="utf-8")
            print("fixed", p.name)
    return 0


if __name__ == "__main__":
    sys.exit(main())
