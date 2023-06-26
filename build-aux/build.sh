#!/bin/bash

set -e

if [[ -z "${1}" ]]; then
    echo "error: missing output dir"
    exit 1
fi

OS_ARCH_LIST=(
    linux-386
    linux-amd64
    linux-arm
    linux-arm64
    windows-386
    windows-amd64
    darwin-amd64
    darwin-arm64
)

MYDIR="$(realpath "$(dirname "${0}")")"
ROOTDIR="$(realpath "${MYDIR}/../")"
PREFIXDIR="${1}"

for os_arch in "${OS_ARCH_LIST[@]}"; do
    os="$(echo ${os_arch} | cut -d- -f1)"
    arch="$(echo ${os_arch} | cut -d- -f2)"
    dir="${PREFIXDIR}/${os_arch}"

    mkdir -p "${dir}"

    pushd "${dir}" > /dev/null
    GOOS="${os}" GOARCH="${arch}" go build -v "${ROOTDIR}"
    popd > /dev/null
done

VERSION="$(date -u +'%Y%m%d%H%M')"
echo "${VERSION}" > "${PREFIXDIR}/VERSION"

if [[ -n "${GITHUB_OUTPUT}" ]]; then
    echo "version=$(cat "${PREFIXDIR}/VERSION")" >> "${GITHUB_OUTPUT}"
fi

pushd "${PREFIXDIR}" > /dev/null

cp "${ROOTDIR}/LICENSE" .
zip -r9 "synth-datagen-all-${VERSION}.zip" LICENSE *-*
sha512sum "synth-datagen-all-${VERSION}.zip" > "synth-datagen-all-${VERSION}.zip.sha512"
rm LICENSE

for dir in {linux,windows,darwin}-*; do
    f="synth-datagen-${dir}-${VERSION}"
    mv "${dir}" "${f}"
    cp "${ROOTDIR}/LICENSE" "${f}/"
    zip -r9 "${f}.zip" "${f}"
    sha512sum "${f}.zip" > "${f}.zip.sha512"
    rm -rf "${f}"
done

popd > /dev/null
