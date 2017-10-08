nanolist
========

nanolist is a lightweight mailing list manager written in Go. It's easy to
deploy, and easy to manage. It was written as an antithesis of the experience
of setting up other mailing list software.

Installation
------------

First, you'll need to build and install the nanolist binary:
`go get github.com/eXeC64/nanolist`

Second, you'll need to write a config to either `/etc/nanolist.ini`
or `/usr/local/etc/nanolist.ini` as follows:

You can also specify a custom config file location by invoking nanolist
with the `-config` flag: `-config=/path/to/config.ini`

```ini
# File for event and error logging. nanolist does not rotate its logs
# automatically. Recommended path is /var/log/mail/nanolist
# You'll need to set permissions on it depending on which account your MTA
# runs nanolist as.
log = /path/to/logfile

# An sqlite3 database is used for storing the email addresses subscribed to
# each mailing list. Recommended location is /var/db/nanolist.db
# You'll need to set permissions on it depending on which account your MTA
# runs nanolist as.
database = /path/to/sqlite/database

# Address nanolist should receive user commands on
command_address = lists@example.com

# SMTP details for sending mail
smtp_hostname = "smtp.example.com"
smtp_port = 25
smtp_username = "nanolist"
smtp_password = "hunter2"

# Create a [list.id] section for each mailing list.
# The 'list.' prefix tells nanolist you're creating a mailing list. The rest
# is the id of the mailing list.

[list.golang]
# Address this list should receieve mail on
address = golang@example.com
# Information to show in the list of mailing lists
name = "Go programming"
description = "General discussion of Go programming"
# bcc all posts to the listed addresses for archival
bcc = archive@example.com, datahoarder@example.com

[list.announcements]
address = announce@example.com
name = "Announcements"
description = "Important announcements"
# List of email addresses that are permitted to post to this list
posters = admin@example.com, moderator@example.com

[list.fight-club]
address = robertpaulson99@example.com
# Don't tell users this list exists
hidden = true
# Only let subscribed users post to this list
subscribers_only = true
```

Lastly, you need to hook the desired incoming addresses to nanolist:

In `/etc/aliases`:
```
nanolist: "| /path/to/bin/nanolist message"
```

And run `newaliases` for the change to take effect.

This creates an alias that pipes messages sent to the `nanolist` alias to the
nanolist command.

The final step is telling your preferred MTA to route mail to this address
when needed.

For postfix edit `/etc/postfix/aliases` and add:
```
lists@example.com nanolist
golang@example.com nanolist
announce@example.com nanolist
robertpaulson99@example.com nanolist
```
and restart postfix.

Congratulations, you've now set up 3 mailing lists of your own!

Commands
--------

Commands are sent to nanolist by emailing the command address (`command_address`
in the configuration file), with the command in the subject. The body of
messages sent to the command address is ignored.

The following commands are available:

* `help` - Reply with a list of valid commands
* `lists` - Reply with a list of available mailing lists
* `subscribe list-id` - Subscribe to receive mail sent to the given list
* `unsubscribe list-id` - Unsubscribe from receiving mail sent to the given list


License
-------

nanolist is made available under the BSD-3-Clause license.
