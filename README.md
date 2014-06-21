Gitql [![Build Status](https://travis-ci.org/cloudson/gitql.png)](https://travis-ci.org/cloudson/gitql)
===============

Gitql is a Git query language.  
In a repository path ...

![how to use](./howtouse.gif)

See more [here](https://asciinema.org/a/8863)

## Requirements 
- Go  
- cmake  

## Install
- `go get -u -d github.com/cloudson/gitql`
- `cd $GOPATH/src/github.com/cloudson/gitql`
- `make`
- `sudo make install`
- `export LD_LIBRARY_PATH=$PWD/libgit2/install/lib` on linux or `export DYLD_LIBRARY_PATH=$PWD/libgit2/install/lib`on Mac OS. 


## Examples 

`gitql "your query" `   
or   
`git ql "your query" `


Look the table of commits:

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

(see more tables [here](./tables.md))

You can do:   
* `select hash, author, message from commits limit 3`  
* `select hash, message from commits where 'hell' in full_message or 'Fuck' in full_message`  
* `select hash, message, author_email from commits where author = 'cloudson'`  
* `select date, message from commits where date < '2014-04-10' `  
* :warning: `select message from commits where 'hell' in message order by date asc`

## Questions? 

`gitql -h` or open an [issue](https://github.com/cloudson/gitql/issues)

Notes:   
* Gitql doesn't want kill `git log` :sweat_smile: . It was created just for science!!  
* It's  read-only. Nothing about delete, insert or update commits :stuck_out_tongue_closed_eyes:  
* The limit default is 10 rows  
* It's inspired by [textql](https://github.com/dinedal/textql)   
* But, why gitql is a compiler/interpreter instead of just read a sqlite database with all commits, tags and etc? Answer: Because we would need to sync the tables everytime before run sql and we would have sqlite bases for each repository. :neutral_face:
