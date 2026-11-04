#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

function resolveBinary() {
  const ext = process.platform === 'win32' ? '.exe' : '';
  const installedBin = path.join(__dirname, `simflux-bin${ext}`);

  if (fs.existsSync(installedBin)) {
    return installedBin;
  }

  // Fallback for local development inside the repository
  const localDevBin = path.join(__dirname, '..', `simflux${ext}`);
  if (fs.existsSync(localDevBin)) {
    return localDevBin;
  }

  return null;
}

function run() {
  let binaryPath = resolveBinary();

  if (!binaryPath) {
    console.log('[simflux] Binary not found. Attempting to download...');
    const { install } = require('../scripts/install');
    // If not found yet, inform the user
    console.error('[simflux] Error: Simflux binary is not installed yet.');
    console.error('[simflux] Please run `npm install` or build locally using `go build -o simflux ./cmd/simflux`.');
    process.exit(1);
  }

  const args = process.argv.slice(2);
  const result = spawnSync(binaryPath, args, { stdio: 'inherit' });

  if (result.error) {
    console.error(`[simflux] Failed to execute binary: ${result.error.message}`);
    process.exit(1);
  }

  if (result.signal) {
    process.kill(process.pid, result.signal);
  } else {
    process.exit(result.status !== null ? result.status : 0);
  }
}

run();

