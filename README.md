# Detection-Validation

## Purpose

The tool automate the process of simulating malicious process events without need to go through setup of real processes. 

Suppose you want to test w3wp.exe spawning Powershell, you will need to go through exchange setup to simulate w3wp.exe spawning Powershell, since detection engines work based simple string matching from telemetry collection tools such as Sysmon. 
Any binary with the same name and path can be used to test the logic, hence no need to setup real exchange to simulate the behavior.

![w3wp_powershell.png](img/w3wp_powershell.png)

The tool allow you to create a child process with a custom parent and path. In addition to couple of other events such as file create from specific process and path, DNS query, Process connections. 

```
NAME:
   Malware Cli - A new cli application

USAGE:
   main.exe [global options] command [command options] [arguments...]

DESCRIPTION:
   Detection validation tool.
   The objective is to generate event with specific conditions to validate detection rule.
   You can execute commands such as w3wp.exe spawning shell or winword creating file or making DNS queries.

COMMANDS:
   argsfree    Accept any commandline
   connect     Connect to host
   dnsquery    Resolve DNS
   execute     Execute command with custom commandline and parent process
   createfile  Create file at a spcific path
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## Examples

**winword.exe spawning cscript.exe**  

```
 mcli.exe execute --parent winword.exe --command cscript.exe
```

**explorer.exe making DNS request** 

```
mcli.exe dnsquery --binpath c:\temp\explorer.exe --host [malicious.com](http://malicious.com/)
```

**w.exe creating file from path C:\temp**  

```
mcli.exe createfile --path f.dat --binpath c:\temp\w.exe
```
