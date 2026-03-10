#!/bin/bash
set -euo pipefail  # Migliore gestione errori

# Colori per output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funzioni di utilità
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verifica parametri
TARGET="$1"
if [ -z "$TARGET" ]; then
    print_error "Target della compilazione non fornito"
    echo "Possibili target:"
    echo "  DEBUG   - Build di debug con simboli"
    echo "  RELEASE - Build ottimizzata per produzione"
    echo "  CLEAN   - Pulisce cache e file temporanei"
    echo "  TEST    - Esegue i test"
    exit 1
fi

# Verifica dipendenze necessarie
check_dependencies() {
    local deps=("jq" "git" "go" "rsync")

    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            print_error "$dep non trovato. Installalo prima di continuare."
            exit 1
        fi
    done
}

check_dependencies

START_TIME=$SECONDS

# Directory di lavoro
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# File di versione
VERSION_FILE="./version.json"
TEMP_FILE="./version.tmp"

if [ ! -f "$VERSION_FILE" ]; then
    print_error "File $VERSION_FILE non trovato"
    exit 1
fi

# Leggi e incrementa versione
major=$(jq -r '.version.major' "$VERSION_FILE")
minor=$(jq -r '.version.minor' "$VERSION_FILE")
patch=$(jq -r '.version.patch' "$VERSION_FILE")

# Incrementa patch e ogni 50 build pulisci cache
((patch++))
if [ $((patch % 50)) -eq 0 ]; then
    print_info "Pulizia cache Go..."
    go clean -cache
fi

# Salva nuova versione
jq --argjson a $patch '.version.patch = $a' "$VERSION_FILE" > "$TEMP_FILE" && mv "$TEMP_FILE" "$VERSION_FILE"

# Genera informazioni versione
appver="$major.$minor.$patch"
githash=$(git rev-parse --short HEAD)
build_date=$(date -u '+%Y-%m-%d_%H:%M:%S')
go_version=$(go version | cut -d' ' -f3)

appname="combina"
if [ "$TARGET" == "DEBUG" ]; then
    appname="${appname}_debug"
else
    appname="${appname}_release"
fi

print_info "Versione: $appver"
print_info "Git hash: $githash"
print_info "Build date: $build_date"
print_info "Go version: $go_version"

# Pulizia build precedenti
print_info "Pulizia build precedenti..."
rm -f "./$appname" "./$appname.exe" 2>/dev/null || true

# Formatta codice
print_info "Formattazione codice..."
if command -v betteralign &> /dev/null; then
    betteralign -apply ./...
else
    go fmt ./...
    go mod tidy
fi

# Esegui test se richiesto
if [ "$TARGET" == "TEST" ]; then
    print_info "Esecuzione test..."
    go test -v ./...
    exit 0
fi

# Compilazione
print_info "Compilazione..."

# Flags comuni
LD_FLAGS="-X main.version=$appver -X main.commit=$githash -X main.date=$build_date"
if [ "$TARGET" == "DEBUG" ]; then
    print_info "Build DEBUG"
    go build -o "./$appname" -ldflags="$LD_FLAGS" .
else
    print_info "Build RELEASE"

    # Aggiungi flags di ottimizzazione
    LD_FLAGS="$LD_FLAGS -s -w"  # Strips debug symbols

    # Build per Linux
    GOOS=linux GOARCH=amd64 go build -ldflags="$LD_FLAGS" -o "./${appname}_linux_amd64" .

    # Build per Windows
    # GOOS=windows GOARCH=amd64 go build -ldflags="$LD_FLAGS" -o "./${appname}_windows_amd64.exe" ./src

	# Build per ARM64 (es. Raspberry Pi)
	GOOS=linux GOARCH=arm64  go build -ldflags="$LD_FLAGS" -o "./${appname}_linux_arm64" .

    # Opzionale: build per altri OS
    # GOOS=darwin GOARCH=amd64 go build -ldflags="$LD_FLAGS" -o "./${appname}_darwin_amd64" ./src
fi

# Rendi eseguibile
if [ -f "./$appname" ]; then
    chmod +x "./$appname"
fi

# Calcola dimensioni
if [ -f "./${appname}_linux_amd64" ]; then
    size=$(du -h "./${appname}_linux_amd64" | cut -f1)
    print_info "Dimensione eseguibile Linux: $size"
fi

# Crea checksum
print_info "Generazione checksums..."
if [ "$TARGET" == "RELEASE" ]; then
    sha256sum ${appname}_* > "checksums_${appver}.txt"
fi

# Durata build
DURATION=$((SECONDS - START_TIME))
print_success "Build completata in $((DURATION / 60)) minuti e $((DURATION % 60)) secondi"

# Avvio applicazione se richiesto
if [ "$TARGET" == "DEBUG" ]; then
    print_info "Avvio applicazione in modalità debug..."
    if [ -f "./$appname" ]; then
        ./$appname -debug
    else
        print_error "Eseguibile non trovato: $appname"
        exit 1
    fi
fi

# Crea un file con le informazioni di build
# cat > "build_info_${appver}.txt" << EOF
# Build Information
# ================
# Application: combina
# Version: $appver
# Git Hash: $githash
# Build Date: $build_date
# Go Version: $go_version
# Target OS/Arch: $os/$arch
# Build Target: $TARGET

# Files:
# $(ls -la ${appname}_* 2>/dev/null || echo "No binaries found")
# EOF

# print_success "Build info salvata in build_info_${appver}.txt"
