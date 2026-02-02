#!/usr/bin/env node
const { Command } = require('commander');
const path = require('path');
const program = new Command();
const { spawn } = require('child_process');

function startCompression(dirPath) {
    const binaryPath = path.join(__dirname, 'bin', 'engine');
    const engine = spawn(binaryPath, [dirPath]);
    let resultData = '';

    engine.stdout.on('data', (data) => {
        resultData += data;
    });

    engine.stderr.on('data', (data) => {
        console.error(`${data.toString()}`);
    });

    engine.on('close', (code) => {
        if (code === 0) {
            const stats = JSON.parse(resultData);
            renderTable(stats);
        }
    });
}

function renderTable(stats) {
    console.log('\n--- Compression Report ---');

    // Normalize fields and compute saved per file
    const rows = stats.map(img => {
        const name = img.name || 'unknown';
        const original = Number(img.original ?? 0);
        const final = Number(img.final ?? 0);
        const saved = original - final;
        return { File: name, Original: original, Final: final, Saved: saved };
    });

    // Helper: human-readable bytes
    function formatBytes(bytes) {
        if (bytes === 0) return '0 B';
        const abs = Math.abs(bytes);
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log10(abs) / 3);
        const value = bytes / Math.pow(1024, i);
        return `${value.toFixed(2)} ${units[i]}`;
    }

    // Build display rows with human-readable sizes and percent
    const display = rows.map(r => {
        const percent = r.Original > 0 ? (r.Saved / r.Original) * 100 : 0;
        return {
            File: r.File,
            Original: formatBytes(r.Original),
            Final: formatBytes(r.Final),
            Saved: `${percent.toFixed(2)}%`
        };
    });

    // Print a nice table
    console.table(display);

    // Totals (numeric)
    const totalSaved = rows.reduce((acc, r) => acc + (r.Saved || 0), 0);
    console.log('--------------------------');
    console.log(`TOTAL SPACE SAVED: ${formatBytes(totalSaved)} (${(totalSaved / 1024).toFixed(2)} KB)\n`);
}

program.name('img-size-compress').description('A CLI tool to compress image sizes').version('1.0.0')
    .argument('[dir]', 'Directory containing images to compress', process.cwd())
    .action((dir) => {
        const targetDir = path.resolve(dir);
        console.log(`🔍 Scanning directory: ${targetDir}`);

        startCompression(targetDir);
    });

program.parse()