#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');

const packageJson = require('../package.json');
const VERSION = packageJson.version;
const REPO = 'Faithful001/simflux';

// Map Node platform & arch to Go release targets
const PLATFORM_MAP = {
  win32: 'windows',
  darwin: 'darwin',
  linux: 'linux',
};

const ARCH_MAP = {
  x64: 'amd64',
  arm64: 'arm64',
};

function getBinaryName() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    throw new Error(`Unsupported platform/architecture: ${process.platform} ${process.arch}`);
  }

  const ext = platform === 'windows' ? '.exe' : '';
  return `simflux-${platform}-${arch}${ext}`;
}

function getDestinationPath() {
  const binDir = path.join(__dirname, '..', 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const ext = process.platform === 'win32' ? '.exe' : '';
  return path.join(binDir, `simflux-bin${ext}`);
}

function download(url, dest, cb) {
  const request = https.get(url, (response) => {
    // Handle HTTP redirects (GitHub releases redirect to S3/CDN)
    if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
      download(response.headers.location, dest, cb);
      return;
    }

    if (response.statusCode !== 200) {
      cb(new Error(`Download failed with HTTP status ${response.statusCode}`));
      return;
    }

    const file = fs.createWriteStream(dest, { mode: 0o755 });
    response.pipe(file);

    file.on('finish', () => {
      file.close(() => {
        try {
          fs.chmodSync(dest, 0o755);
        } catch (_) {}
        cb(null);
      });
    });

    file.on('error', (err) => {
      fs.unlink(dest, () => cb(err));
    });
  });

  request.on('error', (err) => {
    cb(err);
  });
}

function install() {
  try {
    const binaryName = getBinaryName();
    const dest = getDestinationPath();
    const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;

    console.log(`[simflux] Downloading binary for ${process.platform}-${process.arch} from:`);
    console.log(`[simflux] ${url}`);

    download(url, dest, (err) => {
      if (err) {
        console.warn(`[simflux] Warning: Could not download pre-built binary: ${err.message}`);
        console.warn(`[simflux] The binary will be downloaded automatically when running if available.`);
        process.exit(0);
      }
      console.log(`[simflux] Successfully installed binary to ${dest}`);
    });
  } catch (err) {
    console.warn(`[simflux] Warning: ${err.message}`);
    process.exit(0);
  }
}

if (require.main === module) {
  install();
}

module.exports = { getBinaryName, getDestinationPath, download, install };

