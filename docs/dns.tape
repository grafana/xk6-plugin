Set Padding 4
Set Margin 0
Set FontSize 14
Set Width 690
Set Height 380
Set FontFamily "Monospace"

Output build/dns.gif
Set Theme "GitHub Dark"

Hide
Type "export PS1='$ '"
Enter
Type "export SHELL=/bin/sh"
Enter
Type "cd examples"
Enter
Type "asciinema rec --overwrite ../docs/dns.cast"
Enter
Sleep 2s
Show
Sleep 1s

Type "bat ./dns.plugin.js"
Enter
Sleep 10s

Type "bat dns.js"
Enter
Sleep 10s

Type "../k6 run -q --no-summary dns.js"
Enter

Sleep 10s

Hide

Ctrl+D
