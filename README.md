# Indiff

I(nternalizatio)n diff is tool for reviewing changes in unstructured localization files with help of [git](https://git-scm.com/) commits history.

## Usecase

For my project I need to provide documentation in multiple languages. Documentation is written in english as primary language and it is also translated to german. I use markdown files organized in following directory layout:

    doc
    ├── de
    │   ├── first.md
    │   └── section
    │       └── one.md
    └── en
        ├── first.md
        ├── second.md
        └── section
            ├── one.md
            └── two.md

I am preparing new release and I want to check if all documents in primary language were translated and which documents were changed since previous release tagged as `v1.0.0`.

So I run following command inside `doc` directory:

    indiff -f v1.0.0 en,de

And I get following output:

    de: missing translation of: en/second.md
    de: missing translation of: en/section/second.md
    de: modified only base: en/first.md: de/first.md
    de: modified base and translation: en/subone/first.md: de/subone/first.md
         
## Installation

Download archive for your platform from [releases page](https://github.com/unravela/indiff/releases/latest) and unpack it to some directory on your file system.

### How to install on Linux

When you are on Linux you can use following commands to download and install latest release:

    wget https://github.com/unravela/indiff/releases/download/v0.1.0/indiff_0.1.0_Linux_x86_64.tar.gz
    tar -xzvf ./indiff_0.1.0_Linux_x86_64.tar.gz -C /tmp/
    sudo mv /tmp/indiff /usr/local/bin


## Features

To see current list of all options and features just run:

    indiff --help

### Configurable directory layout

Indiff is not limited to single directory layout. You can use it for documents in any layout which can be described with [GLOB-like](http://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm) pattern passed as `-g` flag argument.

    indiff -g "**.%l.%e" en,de

Pattern must contain `%l` placeholder which will be replaced by specific language code. It can additionally contain `%e` placeholder which will be replaced by one or more supported file extensions.

>File extensions can be specified with `-e` flag, e.g. `-e md,rst,html`, and by default all extensions are matched.

Instead of pattern you can use one of the following predefined placeholders for common directory layouts:

**SUB** (`%l/**.%e`) for language specific files in separate root directories:

    doc
    ├── de
    │   ├── first.md
    │   └── section
    │       └── one.md
    └── en
        ├── first.md
        └── section
            └── one.md

**EXT** (`**.%l.%e`) for files in same directory but with language code in file extension:

    doc
    ├── first.de.md
    ├── first.en.md
    └── section
        ├── one.de.md
        └── one.en.md
    
### Multiple languages

You are not limited to work with only two languages. Indiff supports as many languages as you want. It just need to know which one is the primary (base) language. You define it with flag `-b`.

    indiff -b sk en,de,sk

If you ommit `-b` flag, first language is considered as the primary.

### Revision range

Indiff by default looks in uncommited (untracked + staged) files to figure out what was changed. You can specify exact revision range using flags `-f` for oldest revision and `-t` for newest revision.

To check changes between revision tagged with `v1.0.0` and latest commit run:

    indiff -f v1.0.0 -t HEAD en,de

> Instead of tag name you can use commit hash. Tildes and carets are supported too so you can use expressions like `HEAD^` or `v.1.0.0~2`.

### It works without git too

If you project is not versioned with git you can still use indiff to look for missing translation files.

    indiff --no-git en,de 

## Credits

Indiff use [go-git](https://github.com/go-git/go-git) for git repository manipulation.