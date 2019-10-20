Gitql ![](https://github.com/cloudson/gitql/workflows/CI/badge.svg)
===============

Gitql is a Git query language.

In a repository path...

![how to use](howtouse.gif)

See more [here](https://asciinema.org/a/97094)

## Requirements 
- Go  
- cmake  
- pkg-config  

## How to install

We support static compiling for linux and windows platform (amd64), so you can access the [releases page](https://github.com/cloudson/gitql/releases) and just grab the binary. If you want to compile itself follow the instructions below: 

### linux/amd64 

Read the dockerfile to understand the whole process. 

### darwin/amd64

We **do not** support yet static compiling. You need to have pkg-config as dependencies, so after install that, run

```bash
# Inside this repository folder
export PKG_CONFIG_PATH=${PWD}/libgit2/install/lib/pkgconfig:/usr/local/opt/openssl/lib/pkgconfig
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
make build-dynamic
```

### windows/amd64

You need a C compiler, Cmake and Ninja installed. Using chocolately it should be easy 

```bash
choco install cmake ninja vcredist2017
set PATH=%HOMEDRIVE%\mingw64\bin;%PATH%
make build
```

You can always take a look in our [github actions file](./.github/workflows/ci.yml) to understand
how we build it in the ci server. If even after try [the binaries](https://github.com/cloudson/gitql/releases) or either compile yourself you couldn't use that. Open an issue. 

## Examples 

`gitql "your query" `   
or   
`git ql "your query" `


As an example, this is the `commits` table:

| commits |
| ---------|
| author |
| author_email |
| committer |
| committer_email |
| hash |
| date |
| message |
| full_message |

(see more tables [here](tables.md))

## Example Commands
* `select hash, author, message from commits limit 3`  
* `select hash, message from commits where 'hell' in full_message or 'Fuck' in full_message`  
* `select hash, message, author_email from commits where author = 'cloudson'`  
* `select date, message from commits where date < '2014-04-10' `  
* `select message from commits where 'hell' in message order by date asc`

## Questions?

`gitql` or open an [issue](https://github.com/cloudson/gitql/issues)

Notes:
* Gitql doesn't want to _kill_ `git log` - it was created just for science! :sweat_smile:
* It's read-only - no deleting, inserting, or updating tables or commits. :stuck_out_tongue_closed_eyes:
* The default limit is 10 rows.
* It's inspired by [textql](https://github.com/dinedal/textql).
* Gitql is a compiler/interpreter instead of just read a sqlite database with all commits, tags, etc. because we would need to sync the tables every time before run sql and we would have sqlite bases for each repository. :neutral_face:
