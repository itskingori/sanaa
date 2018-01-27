#!/bin/bash

set -eu

version="$1"
major_version="$(echo ${version} | cut -d'.' -f1)"
minor_version="$(echo ${version} | cut -d'.' -f2)"
patch_version="$(echo ${version} | cut -d'.' -f3)"

binary_output_path="binaries"
package_path="github.com/itskingori/sanaa"
target_platforms=(
  'darwin/amd64'
  'linux/amd64'
)

ldflags="-X ${package_path}/service.majorVersion=${major_version} \
         -X ${package_path}/service.minorVersion=${minor_version} \
         -X ${package_path}/service.patchVersion=${patch_version}"

echo -e "Installing gox, Go's cross compilation tool"
go get -v github.com/mitchellh/gox

echo -e "Building binaries for targetted platforms"
rm -rf ${binary_output_path}/*
mkdir -p ${binary_output_path}
gox \
-osarch="${target_platforms[*]}" \
-ldflags="${ldflags}" \
-output="${binary_output_path}/sanaa-${version}-{{.OS}}-{{.Arch}}"

echo -e "\nCompressing built binaries:\n"
binaries=$(find "./${binary_output_path}" -name "*" -type f)
for binary in ${binaries[*]}; do
  uncompressed_file="${binary}"
  compressed_file="${uncompressed_file}.tar.gz"
  echo "--> ${uncompressed_file} ~> ${compressed_file}"
  tar -czf "${compressed_file}" "${uncompressed_file}"
  rm -rf ${uncompressed_file}
done

echo -e "\nCalculating SHA256-sums:\n"
cd ${binary_output_path}/
sha265sum_file="sanaa-${version}-shasum256.txt"
shasum -a 256 sanaa-${version}-* > "${sha265sum_file}"
cd ../
cat "${binary_output_path}/${sha265sum_file}"

echo -e "\nDone!"
