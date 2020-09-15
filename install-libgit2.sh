#!/bin/bash

set -ex

git clone --depth 1 --branch v1.0.1 git://github.com/libgit2/libgit2 libgit2

cd libgit2
# Those files are temporary, should not be in git. To minimize 
# impact, we'll delete them now
rm -rf CMakeCache.txt CMakeFiles 

# There is no official way to specify a generator
# inside CMakeLists.txt, see 
# https://stackoverflow.com/questions/11269833/cmake-selecting-a-generator-within-cmakelists-txt
# So we deriving the environment from uname and setting generator
# in the command line, 
# https://stackoverflow.com/questions/3466166/how-to-check-if-running-in-cygwin-mac-or-linux
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac

# Passing an argument with a space is a trouble, so we switch from sh to bash
# https://stackoverflow.com/a/2249967/9469533
EXTRA_ARGS=( )
if [[ $machine == MinGw ]]
then
    extra_args=( -G "MSYS Makefiles" )
fi

cmake -DTHREADSAFE=ON -DBUILD_CLAR=OFF -DCMAKE_INSTALL_PREFIX=$PWD/install "${extra_args[@]}" .

cmake --build .

cmake --build . --target install
