# simple opsgenie CLI

provides the following commands:

| command                      | description |
|---|---|
|ops alert list                | list all alerts |
|ops alert list 5h             | list all alerts in the past 5 hours |
|ops alert ack 1               | acknowledge alert number 1 |
|ops alert ack all             | acknowledge all outstanding alert |
|ops policy list               | list all policies |
|ops policy test 1             | see what would match policy 1 (use 'ops policy list' to find the number) |
|ops policy test 1 5m          | see what would match policy 1 in the last 5 minutes |
|ops policy enable 1           | enable policy 1 |
|ops policy enable 1 1h        | enable policy 1 for 1 hour |
|ops policy disable 1          | enable policy 1 for 1 hour |
|ops filter your filter 1h30m  | create a policy and enable it for 1 hour and 30 minutes |
|ops help                      | your looking at it |

to configure add a opscli.config or .opscli config file with the api + team keys
