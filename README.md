# stdtimer

Run a command, prefixing each line with the current duration.


## Example command

```
stdtimer bash -c '
    echo Beginning; sleep 2
    echo Middle;    sleep 2
    echo End
'
```

Output:

```
[0:00] Beginning
[0:02] Middle
[0:04] End
```


## Installation

Download from the [releases](https://github.com/nathforge/stdtimer/releases/latest) page.
Extract with `tar xzf FILENAME.tar.gz`, and move `stdtimer` to a directory
included in your system's `PATH`.


## License

Apache 2.0. See [LICENSE](LICENSE).
