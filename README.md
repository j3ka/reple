# Reple

## Description
Simple wrapper for any repl. If your editor can send selected text as stdin for shell command, you can use repl-driven-development in your workflow.

First spawn the shell:

```shell
reple spawn 'bash --noprofile --norc -i'
```

Next send to shell peace of code from your text editor as stdin, or from another shell:
```shell
echo 'ls -la' | reple eval

# or

reple eval < my_script.sh
```

## Installation
```shell
go install github.com/j3ka/reple@latest
```

## Todo
- Add ability to spawn different repls in the same time