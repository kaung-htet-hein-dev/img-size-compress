# img-size-compress

A command-line tool to compress image files and reduce their size.

## What it does

Compresses JPEG and PNG images in a directory to reduce file sizes while maintaining acceptable quality. Shows you how much space was saved.

## How it works

Uses a Go-based compression engine with the imaging library to process and optimize image files. The tool scans your specified directory, compresses each image, and generates a report showing original size, final size, and percentage saved.

## Installation

```bash
npm install -g img-size-compress
```

## Usage

Compress images in current directory:
```bash
img-size-compress
```

Compress images in a specific directory:
```bash
img-size-compress /path/to/images
```

Or use with npx without installing:
```bash
npx img-size-compress
npx img-size-compress /path/to/images
```

The tool will display a table showing compression results for each file and total space saved.

## Development

Build the Go engine:
```bash
npm run build-engine
```

## License

ISC
