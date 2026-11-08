#!/usr/bin/env python3
"""Vérifie l’invalidation cgo après modification d’un en-tête C++ imbriqué."""
import argparse
import json
import re
import subprocess
import tempfile
import time
from pathlib import Path

p = argparse.ArgumentParser()
p.add_argument('--native', required=True, help='checkout modifiable de milvi/mlpack')
p.add_argument('--output', required=True)
a = p.parse_args()
root = Path(__file__).resolve().parents[1]
header = Path(a.native).resolve() / 'src/mlpack/methods/kmeans/dual_tree_kmeans.hpp'
original = header.read_bytes()
results = []
with tempfile.TemporaryDirectory(prefix='mlpack-cache-') as tmp:
    def build(name):
        start = time.monotonic()
        value = subprocess.run(['go', 'build', '-p', '2', '-x', '-o', str(Path(tmp)/'kmeans'), './cmd/kmeans'],
            cwd=root, capture_output=True, text=True)
        if value.returncode:
            raise RuntimeError(value.stderr[-12000:])
        units = re.findall(r' -c (?:\S*/)?(native_[^\s]+\.cpp)(?:\s|$)', value.stderr)
        result = dict(name=name, seconds=time.monotonic()-start, native_compilations=len(units))
        results.append(result)
        print(json.dumps(result), flush=True)
    build('préparation')
    build('sans modification')
    header.touch()
    build('horodatage seul modifié')
    try:
        header.write_bytes(original + b'\n// MILVI_CACHE_PROBE_HEADER_MODIFIE\n')
        build('en-tête imbriqué modifié')
    finally:
        header.write_bytes(original)
    build('source restaurée')
    binary = (Path(tmp)/'kmeans').read_bytes()
    marker = 'Les matrices restent dans std::any ; seuls les modèles pointés sont concernés.'.encode()
    native_header = Path(a.native).resolve()/'src/mlpack/bindings/go/model_memory.hpp'
    if marker not in native_header.read_bytes():
        raise AssertionError('témoin source absent du test')
    if marker in binary:
        raise AssertionError('les sources C++ embarquées sont présentes dans le binaire')
    if (results[1]['native_compilations'] != 0 or
            results[2]['native_compilations'] != 0 or
            results[3]['native_compilations'] != 49 or
            results[4]['native_compilations'] != 0):
        raise AssertionError(f'cache inattendu : {results}')
    report = dict(builds=results, source_marker_in_binary=False, binary_bytes=len(binary),
        limitation='Les en-têtes et bibliothèques APT restent extérieurs au cache Go ; forcer go build -a après leur mise à jour.')
Path(a.output).write_text(json.dumps(report, indent=2, ensure_ascii=False)+'\n')
