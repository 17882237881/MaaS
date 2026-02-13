from __future__ import annotations

from pathlib import Path
import sys

_proto_dir = Path(__file__).resolve().parent
if str(_proto_dir) not in sys.path:
    sys.path.append(str(_proto_dir))
