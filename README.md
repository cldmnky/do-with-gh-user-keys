# do-with-gh-user-keys

A small tool writen in GO to loop though ssh keys for users in an GitHub organization and pipe:ing them to a command's stdin.

I use it for adding SSH keys to a docker image running https://github.com/progrium/gitreceive.

````
./do-with-gh-user-keys-linux-amd64 -o <your GH organization> \
  -t <GH OATH2 Key (PAT)>  \
  -c "/usr/local/bin/gitreceive" \
  -a "upload-key" -u
````

Invoking:
```
NAME:
   do-with-gh-user-keys - Runs a program for each users ssh key in an organization

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value         A GitHub personal access token to be used for authentication (default: "<unset>") [$DOW_GH_TOKEN]
   --command value, -c value       Command that will reveive the piped key(s)
   --userarg, -u                   Add github username as last argument to command
   --args value, -a value          Args to command that will reveive the piped key(s)
   --organization value, -o value  List member keys in this organization
   --help, -h                      show help
   --version, -v                   print the version
   ```


