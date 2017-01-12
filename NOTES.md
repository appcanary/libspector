# Notes

Abstractions:

* lsof, pmap, pldd
* ldd, otool (osx)
* dpkg/dpkg-query


## Ubuntu

Find process:

```
$ pgrep -f "nginx: master"
10546
```

List open files of pid:

```
$ lsof -p 10546
COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF     NODE NAME
nginx   10546 root  cwd    DIR              253,1     4096        2 /
nginx   10546 root  rtd    DIR              253,1     4096        2 /
nginx   10546 root  txt    REG              253,1   892392   925590 /usr/sbin/nginx
nginx   10546 root  mem    REG              253,1    52160  1049455 /lib/x86_64-linux-gnu/libnss_files-2.17.so
nginx   10546 root  mem    REG              253,1    47760  1049453 /lib/x86_64-linux-gnu/libnss_nis-2.17.so
nginx   10546 root  mem    REG              253,1    97296  1049441 /lib/x86_64-linux-gnu/libnsl-2.17.so
nginx   10546 root  mem    REG              253,1    35728  1049444 /lib/x86_64-linux-gnu/libnss_compat-2.17.so
...
```

Or by command substring:

```
$ lsof -c nginx
...
```

See: [lsof quickstart](http://www.akadia.com/services/lsof_quickstart.txt)

Don't see a way to filter (ideally FD=mem TYPE=REG)


List of dynamically linked libraries in a binary:

```
$ ldd /usr/sbin/nginx
    linux-vdso.so.1 =>  (0x00007fffb9da7000)
    libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f2d51cb9000)
    libcrypt.so.1 => /lib/x86_64-linux-gnu/libcrypt.so.1 (0x00007f2d51a80000)
    ...
```

List of dynamically linked libraries of a *process*:

```
$ pldd 10546
10546:  /usr/sbin/nginx
linux-vdso.so.1
/lib/x86_64-linux-gnu/libpthread.so.0
/lib/x86_64-linux-gnu/libcrypt.so.1
/lib/x86_64-linux-gnu/libpam.so.0
/lib/x86_64-linux-gnu/libexpat.so.1
/lib/x86_64-linux-gnu/libpcre.so.3
/lib/x86_64-linux-gnu/libssl.so.1.0.0
/lib/x86_64-linux-gnu/libcrypto.so.1.0.0
/lib/x86_64-linux-gnu/libz.so.1
/usr/lib/x86_64-linux-gnu/libxml2.so.2
/usr/lib/x86_64-linux-gnu/libxslt.so.1
...
$ pmap 10546
10546:   nginx: master process /usr/sbin/nginx
0000000000400000    788K r-x--  /usr/sbin/nginx
00000000006c4000      4K r----  /usr/sbin/nginx
00000000006c5000     84K rw---  /usr/sbin/nginx
00000000006da000     60K rw---    [ anon ]
0000000000cd3000    504K rw---    [ anon ]
0000000000d51000    844K rw---    [ anon ]
00007fee20e23000     48K r-x--  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee20e2f000   2044K -----  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee2102e000      4K r----  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee2102f000      4K rw---  /lib/x86_64-linux-gnu/libnss_files-2.17.so
...
```

pmap seems to give resolved paths.

Package owns file:

```
$ dpkg -S /lib/x86_64-linux-gnu/libcrypt.so.1
libc6:amd64: /lib/x86_64-linux-gnu/libcrypt.so.1
```

Also works with substrings.

Version of package installed:

```
$ dpkg -s libc6:amd64 | grep "Version:"
Version: 2.17-93ubuntu4
$ dpkg-query --showformat='${Version}\n' --show libc6
2.17-93ubuntu4
```

Alternatively:

```
$ apt-cache policy libc6:amd64
libc6:
  Installed: 2.17-93ubuntu4
  Candidate: 2.17-93ubuntu4
  Version table:
 *** 2.17-93ubuntu4 0
        500 http://archive.ubuntu.com/ubuntu/ saucy/main amd64 Packages
        100 /var/lib/dpkg/status
```

Reverse depends:

```
$ apt-cache rdepends --installed libc6
...
```

All files owned by a package (does not include binaries?):

```
$ dpkg-query -L libc6
...
```

When was a process started:

```
$ ps -p 10546 -o lstart=
Mon Sep 14 20:02:43 2015
```

When was a library last modified:

```
$ stat --format="%z" /lib/x86_64-linux-gnu/libcrypt.so.1
2014-04-09 06:57:09.555173000 +0000
```


## OSX

* `otool -L`
* `vmmap $PID`
