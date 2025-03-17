Set Padding 4
Set Margin 0
Set FontSize 14
Set Width 690
Set Height 530
Set FontFamily "Monospace"

Output build/math.gif
Set Theme "GitHub Dark"

Hide
Type "export PS1='$ '"
Enter
Type "export SHELL=/bin/sh"
Enter
Type "cd examples"
Enter
Type "asciinema rec --overwrite ../docs/math.cast"
Enter
Sleep 2s
Show
Sleep 1s

Type "bat ./math.plugin.ts"
Enter
Sleep 10s

Type "bat math.js"
Enter
Sleep 10s

Type "../k6 run -q --no-summary math.js"
Enter
Sleep 5s

Type "K6_PLUGIN_RUNTIME=bun ../k6 run -q --no-summary math.js"
Enter

Sleep 10s

Hide

Ctrl+D
