package mimetype

import (
	"github.com/h2non/filetype"
	"os"
	"strings"
)

type FileMIME struct {
	Types map[FileExtension]string
}

type FileExtension int

const (
	JPG FileExtension = iota
	PNG
	GIF
	WEBP
	CR2
	TIF
	BMP
	HEIF
	JXR
	PSD
	ICO
	DWG
	AVIF
	MP4
	M4V
	MKV
	WEBM
	MOV
	AVI
	WMV
	MPG
	FLV
	THREEGP
	MID
	MP3
	M4A
	OGG
	FLAC
	WAV
	AMR
	AAC
	AIFF
	EPUB
	ZIP
	TAR
	RAR
	GZ
	BZ2
	SEVENZ
	XZ
	ZSTD
	PDF
	EXE
	SWF
	RTF
	ISO
	EOT
	PS
	SQLITE
	NES
	CRX
	CAB
	DEB
	AR
	Z
	LZ
	RPM
	ELF
	DCM
	DOC
	DOCX
	XLS
	XLSX
	PPT
	PPTX
	WOFF
	WOFF2
	TTF
	OTF
	WASM
	DEX
	DEY
)

var mimeTypes = map[FileExtension]string{
	JPG:     "image/jpeg",
	PNG:     "image/png",
	GIF:     "image/gif",
	WEBP:    "image/webp",
	CR2:     "image/x-canon-cr2",
	TIF:     "image/tiff",
	BMP:     "image/bmp",
	HEIF:    "image/heif",
	JXR:     "image/vnd.ms-photo",
	PSD:     "image/vnd.adobe.photoshop",
	ICO:     "image/vnd.microsoft.icon",
	DWG:     "image/vnd.dwg",
	AVIF:    "image/avif",
	MP4:     "video/mp4",
	M4V:     "video/x-m4v",
	MKV:     "video/x-matroska",
	WEBM:    "video/webm",
	MOV:     "video/quicktime",
	AVI:     "video/x-msvideo",
	WMV:     "video/x-ms-wmv",
	MPG:     "video/mpeg",
	FLV:     "video/x-flv",
	THREEGP: "video/3gpp",
	MID:     "audio/midi",
	MP3:     "audio/mpeg",
	M4A:     "audio/mp4",
	OGG:     "audio/ogg",
	FLAC:    "audio/x-flac",
	WAV:     "audio/x-wav",
	AMR:     "audio/amr",
	AAC:     "audio/aac",
	AIFF:    "audio/x-aiff",
	EPUB:    "application/epub+zip",
	ZIP:     "application/zip",
	TAR:     "application/x-tar",
	RAR:     "application/vnd.rar",
	GZ:      "application/gzip",
	BZ2:     "application/x-bzip2",
	SEVENZ:  "application/x-7z-compressed",
	XZ:      "application/x-xz",
	ZSTD:    "application/zstd",
	PDF:     "application/pdf",
	EXE:     "application/vnd.microsoft.portable-executable",
	SWF:     "application/x-shockwave-flash",
	RTF:     "application/rtf",
	ISO:     "application/x-iso9660-image",
	EOT:     "application/octet-stream",
	PS:      "application/postscript",
	SQLITE:  "application/vnd.sqlite3",
	NES:     "application/x-nintendo-nes-rom",
	CRX:     "application/x-google-chrome-extension",
	CAB:     "application/vnd.ms-cab-compressed",
	DEB:     "application/vnd.debian.binary-package",
	AR:      "application/x-unix-archive",
	Z:       "application/x-compress",
	LZ:      "application/x-lzip",
	RPM:     "application/x-rpm",
	ELF:     "application/x-executable",
	DCM:     "application/dicom",
	DOC:     "application/msword",
	DOCX:    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	XLS:     "application/vnd.ms-excel",
	XLSX:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	PPT:     "application/vnd.ms-powerpoint",
	PPTX:    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	WOFF:    "application/font-woff",
	WOFF2:   "application/font-woff",
	TTF:     "application/font-sfnt",
	OTF:     "application/font-sfnt",
	WASM:    "application/wasm",
	DEX:     "application/vnd.android.dex",
	DEY:     "application/vnd.android.dey",
}

func MIMEMap() *FileMIME {
	return &FileMIME{mimeTypes}
}

func (fm *FileMIME) GetMIME(extension FileExtension) string {
	return fm.Types[extension]
}

func (fm *FileMIME) CheckFileType(filePath string) (string, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return "", err
	}

	return kind.MIME.Value, nil
}

func (fm *FileMIME) GetExtensionFromMIME(mimeType string) string {
	// We split the MIME type at the forward slash (/)
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return ""
	}

	// Grab the extension which is split in two based on the seperator
	extension := parts[1]

	// account of jpeg
	if extension == "jpeg" {
		extension = "jpg"
	}

	// account of mpg
	if extension == "mpeg" {
		extension = "mpg"
	}

	// prepend the dot to form the file extension
	return "." + extension
}
