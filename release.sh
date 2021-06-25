#!/bin/bash
name=$(basename $PWD)
zip -ju ${name}.zip bin/* *.yml README.md
