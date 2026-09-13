#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';

const mw = process.argv[2];
const outDir = process.argv[3];
const meta = JSON.parse(process.argv[4] || '{}');

function git(args) {
  try {
    return execFileSync('git', ['-C', mw, ...args], { encoding: 'utf8', maxBuffer: 1024 * 1024 * 256 });
  } catch (e) {
    return String((e && e.stdout) || '') + String((e && e.stderr) || '');
  }
}

const unmerged = git(['diff', '--name-only', '--diff-filter=U']).split('\n').map((s) => s.trim()).filter(Boolean);
const status = git(['status', '--porcelain=v1']).split('\n').filter(Boolean);

const conflicts = [];
for (const p of unmerged) {
  let text = '';
  try { text = fs.readFileSync(path.join(mw, p), 'utf8'); } catch (e) { text = ''; }
  let hunks = 0;
  let oursLines = 0;
  let theirsLines = 0;
  let side = 0;
  for (const ln of text.split('\n')) {
    if (ln.startsWith('<<<<<<< ')) { hunks += 1; side = 1; continue; }
    if (ln.startsWith('=======')) { side = 2; continue; }
    if (ln.startsWith('>>>>>>> ')) { side = 0; continue; }
    if (side === 1) oursLines += 1;
    else if (side === 2) theirsLines += 1;
  }
  conflicts.push({ path: p, hunks, oursLines, theirsLines, bytes: Buffer.byteLength(text) });
}
conflicts.sort((a, b) => a.hunks - b.hunks || a.path.localeCompare(b.path));

const report = {
  generatedAt: new Date().toISOString(),
  ...meta,
  mergeExit: Number(process.env.MERGE_EXIT || -1),
  changedFiles: status.length,
  conflictCount: conflicts.length,
  totalHunks: conflicts.reduce((s, c) => s + c.hunks, 0),
  conflicts,
};

fs.mkdirSync(outDir, { recursive: true });
fs.writeFileSync(path.join(outDir, 'report.json'), JSON.stringify(report, null, 2));
fs.writeFileSync(
  path.join(outDir, 'conflicts.txt'),
  conflicts.map((c) => String(c.hunks).padStart(4) + '  ' + c.path).join('\n') + '\n',
);
fs.writeFileSync(path.join(outDir, 'status.txt'), status.join('\n') + '\n');
console.log(JSON.stringify({ mergeExit: report.mergeExit, changedFiles: report.changedFiles, conflictCount: report.conflictCount, totalHunks: report.totalHunks }));
