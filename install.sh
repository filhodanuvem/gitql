# Don't run this one directly, use «make install» to run it
set -ex

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac

if [[ $machine == MinGw ]]
then
 # Everything here is ugly. But MinGW seem to have no ldconfig, 
 # no respect to /usr/local/bin, /usr/local/lib
 cp ./libgit2/install/bin/libgit2.dll  /usr/bin/
 cp ./gitql.exe /usr/bin/gitql.exe
 ln -s -f /usr/bin/gitql.exe /usr/bin/git-ql
 echo "Gitql.exe is in /usr/bin/gitql.exe"
 echo "You can also use: git ql 'query here'"
else
 cp ./libgit2/install/lib/lib*  /usr/local/lib/
 ldconfig /usr/local/lib >/dev/null 2>&1 || echo "ldconfig not found">/dev/null
 cp ./gitql /usr/local/bin/gitql
 ln -s -f /usr/local/bin/gitql /usr/local/bin/git-ql
 echo "Git is in /usr/local/bin/gitql"
 echo "You can also use: git ql 'query here'"
fi


