# cbot
Cabify has become self-aware

#### Commands

The `commands` directory contains executables which are run in response to chat messages (these must have the executable bit set, ie `chmod +x`).

The executable is passed the chat message as its argument(s). Any output on STDOUT and STDERR are sent as Flowdock replies to the message. 

Blank newlines divide multiple replies. Eg

```
1  hello
2  same message
3  
4  a new reply
```

Lines 1-2 will be sent as a reply, then line 4 as a separate reply.

##### There are two types of commands:

 * Executables which begin with `_`. These are run on **each** message received and are passed the entire message string as arg1.

 * Executables which don't begin with `_` are only run when cbot receives a command with the executable filename. Eg there's an executable called `deploy`, so if the chat message `cbot deploy first second third` is seen, the `deploy` executable will be run with arg1 = first, arg2 = second, etc

 #### Environment
  Every command is called with some environment
  You can get some variables with the message
  ```
  CURRENT_FLOW
  CURRENT_USER_AVATAR
  CURRENT_USER_EMAIL
  CURRENT_USER_NICK
  CURRENT_USER_NAME
  CURRENT_USER_ID
  ```