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

## Configuration examples
### Kakoune:
```shell
define-command reple-eval %{
        evaluate-commands %sh{echo "$kak_selection" | reple eval}
}
map global normal <A-ret> ":reple-eval<ret>" # Alt + Enter to send selection to reple
```
### Helix:
```toml
[keys.normal]
"A-ret" = ":pipe-to reple eval"

[keys.select]
"A-ret" = ":pipe-to reple eval"
```
### Lite-XL:
I've experimented with the lite-xl and created little plugin for reple: [reple.lua](https://github.com/j3ka/litexl-reple/blob/master/reple.lua)

## Installation
```shell
go install github.com/j3ka/reple@latest
```

## Todo
- Add ability to spawn different repls in the same time
