Set Padding 4
Set Margin 0
Set FontSize 14
Set Width 760
Set Height 512
Set FontFamily "Monospace"

Output build/dump.gif
Set Theme "GitHub Dark"

Hide
Type "export PS1='$ '"
Enter
Type "export SHELL=/bin/sh"
Enter
Type "cd examples"
Enter
Type "asciinema rec --overwrite ../docs/dump.cast"
Enter
Sleep 2s
Show
Sleep 1s

Type "bat ./dump.plugin.js"
Enter
Sleep 10s

Type "bat dump.js"
Enter
Sleep 10s

Type "clear"
Enter
Type "rm -rf dump"
Enter
Sleep 1

Type "../k6 run -q --no-summary dump.js"
Enter
Sleep 5

Type "ls dump"
Enter
Sleep 2

Type "bat dump/*"
Enter

Sleep 10s

Hide

Ctrl+D
