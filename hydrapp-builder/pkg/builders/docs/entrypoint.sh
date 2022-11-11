#!/bin/bash

set -e

pandoc --to markdown --shift-heading-level-by=-1 --resource-path=docs --standalone "${GOMAIN}/INSTALLATION.md" | pandoc --to html5 --citeproc --listings --shift-heading-level-by=1 --number-sections --resource-path=docs --toc --katex --self-contained --number-offset=1 -o '/dst/INSTALLATION.html'
