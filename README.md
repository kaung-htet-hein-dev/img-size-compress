# bimg-size-compress

A command-line tool to compress image files using bimg/libvips for high-performance compression.

## What it does

Compresses JPEG and PNG images in a directory to reduce file sizes while maintaining excellent quality. Uses libvips for 4-8x faster compression compared to traditional methods. Shows you how much space was saved.

## Prerequisites

⚠️ **libvips must be installed on your system before using this tool.**

### macOS
```bash
brew install vips
```

### Ubuntu/Debian
```bash
sudo apt-get install libvips-dev
```

### Windows
Download and install libvips from: https://github.com/libvips/libvips/releases

Or use vcpkg:
```bash
vcpkg install libvips
```

## Installation

```bash
npm install -g bimg-size-compress
```

## Usage

Compress images in current directory:
```bash
bimg-size-compress
```

Compress images in a specific directory:
```bash
bimg-size-compress /path/to/images
```

Or use with npx without installing:
```bash
npx bimg-size-compress
npx bimg-size-compress /path/to/images
```

The tool will display a table showing compression results for each file and total space saved.

## How it works

Uses a Go-based compression engine powered by bimg/libvips for high-performance image processing. The tool scans your specified directory, compresses each image using optimized algorithms, and generates a detailed report showing original size, final size, and percentage saved.

**Features:**
- 4-8x faster than traditional Go image libraries
- Low memory footprint with streaming processing
- JPEG quality: 75 (optimal balance of size/quality)
- PNG compression: Level 9 (best compression)

## Development

Build the Go engine:
```bash
npm run build-engine
```

## License

ISC
