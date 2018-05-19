#!/bin/bash

set -eu

version="$1"
commit_version="$2"
major_version="$(echo "${version}" | cut -d'.' -f1)"
minor_version="$(echo "${version}" | cut -d'.' -f2)"
patch_version="$(echo "${version}" | cut -d'.' -f3)"

binary_output_path="bin"
package_path="github.com/itskingori/sanaa"
target_platforms=(
  'darwin/amd64'
  'linux/amd64'
)

ldflags="-X ${package_path}/service.majorVersion=${major_version} \
         -X ${package_path}/service.minorVersion=${minor_version} \
         -X ${package_path}/service.patchVersion=${patch_version} \
         -X ${package_path}/service.commitVersion=${commit_version}"

echo -e "Installing gox, Go's cross compilation tool"
go get -v github.com/mitchellh/gox

echo -e "Building binaries for targetted platforms"
rm -rf "${binary_output_path:?}"/*
mkdir -p ${binary_output_path}
gox \
-osarch="${target_platforms[*]}" \
-ldflags="${ldflags}" \
-output="${binary_output_path}/sanaa-${version}-{{.OS}}-{{.Arch}}"

echo -e "\nCompressing built binaries:\n"
binaries=$(find "./${binary_output_path}" -name "*" -type f)
for binary in ${binaries[*]}; do
  uncompressed_filepath="${binary}"
  uncompressed_filename=$(basename "${uncompressed_filepath}")
  compressed_filepath="${uncompressed_filepath}.tar.gz"
  compressed_filename=$(basename "${compressed_filepath}")

  echo "--> ${uncompressed_filepath} ~> ${compressed_filepath}"
  tar -czf "${compressed_filepath}" -C "${binary_output_path}/" "${uncompressed_filename}"
  rm -rf "${uncompressed_filepath:?}"

  cd ${binary_output_path}/
  shasum_256_file="${uncompressed_filename}-shasum-256.txt"
  shasum -a 256 "${compressed_filename}" > "${shasum_256_file}"
  cd ../
done

echo -e "\nDone!"
