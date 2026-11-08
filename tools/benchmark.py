#!/usr/bin/env python3
"""Mesure un appel par processus ; génération des points exclue du chronométrage."""
import argparse
import json
import math
import platform
import statistics
import subprocess
from datetime import datetime, timezone
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument('--binary', required=True)
parser.add_argument('--output', required=True)
parser.add_argument('--trials', type=int, default=3)
parser.add_argument('--iterations', type=int, default=100)
args = parser.parse_args()
if args.trials < 1:
    parser.error('au moins une répétition')
measurements = []
# A→B : k fixé. A→C : k = n/100. Un nouveau processus mesure chaque pic RSS.
for n, k in [(100000, 1000), (1000000, 1000), (1000000, 10000)]:
    for trial in range(args.trials):
        command = [str(Path(args.binary).resolve()), '-n', str(n), '-k', str(k),
                   '-seed', str(42 + trial), '-threads', '1', '-iterations', str(args.iterations)]
        completed = subprocess.run(command, check=True, text=True, capture_output=True)
        value = json.loads(completed.stdout)
        measurements.append(value)
        print(json.dumps(value, ensure_ascii=False), flush=True)
summary = []
for n, k in [(100000, 1000), (1000000, 1000), (1000000, 10000)]:
    values = [v for v in measurements if v['n'] == n and v['k_requested'] == k]
    summary.append(dict(n=n, k=k,
        median_seconds=statistics.median(v['seconds'] for v in values),
        min_seconds=min(v['seconds'] for v in values),
        max_seconds=max(v['seconds'] for v in values),
        median_peak_rss_bytes=statistics.median(v['peak_rss_bytes'] for v in values),
        median_seconds_per_iteration=statistics.median(v['seconds']/v['iterations'] for v in values),
        iterations=[v['iterations'] for v in values],
        k_actual=[v['k_actual'] for v in values], converged=[v['converged'] for v in values]))
reference = summary[0]
ratios = [dict(case=f'k={v["k"]}', n_ratio=v['n']/reference['n'],
    time_ratio=v['median_seconds']/reference['median_seconds'],
    rss_ratio=v['median_peak_rss_bytes']/reference['median_peak_rss_bytes'],
    effective_time_exponent=math.log10(v['median_seconds']/reference['median_seconds']),
    per_iteration_time_ratio=v['median_seconds_per_iteration']/reference['median_seconds_per_iteration']) for v in summary[1:]]
result = dict(utc=datetime.now(timezone.utc).isoformat(), machine=platform.platform(),
    logical_cpus=__import__('os').cpu_count(), cpu=Path('/proc/cpuinfo').read_text().split('model name\t: ')[-1].splitlines()[0],
    trials=args.trials, iteration_limit=args.iterations, measurements=measurements, summary=summary, ratios=ratios,
    interpretation='Le quotient mesure cette distribution et ces convergences ; il ne démontre pas une complexité asymptotique. Le pic RSS inclut Go, C++, le runtime et les entrées/sorties.')
output = Path(args.output)
output.parent.mkdir(parents=True, exist_ok=True)
output.write_text(json.dumps(result, indent=2, ensure_ascii=False) + '\n')
print(json.dumps(result, ensure_ascii=False), flush=True)
