#!/usr/bin/env node
const { Command } = require('commander');
const path = require('path');
const program = new Command();
const { spawn } = require('child_process');

function startCompression(dirPath) {
    const engine = spawn('./bin/engine', [dirPath]);

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
    let totalSaved = 0;

    stats.forEach(img => {
        const saved = img.Original - img.Final;
        totalSaved += saved;
        console.log(`File: ${img.Name}`);
        console.log(`   ${img.Original} bytes -> ${img.Final} bytes (Saved: ${saved} bytes)`);
    });

    console.log('--------------------------');
    console.log(`TOTAL SPACE SAVED: ${(totalSaved / 1024).toFixed(2)} KB\n`);
}

program.name('img-size-compress').description('A CLI tool to compress image sizes').version('1.0.0')
    .argument('[dir]', 'Directory containing images to compress', process.cwd())
    .action((dir) => {
        const targetDir = path.resolve(dir);
        console.log(`🔍 Scanning directory: ${targetDir}`);

        startCompression(targetDir);
    });

program.parse()