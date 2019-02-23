# chronicast -- Multicast alarm system
This is a fairly simple project for creating alarms and responding to them on a UDP multicast network.

From the CLI you simply issue an `alarm "My Alarm Name" some_command [arguments...]` command and await another instance -- or the same instance -- to receive a broadcast via the `ring "My Alarm Name" <Time/Offset>` command.
