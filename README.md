# linetimer

Run a command, prefixing each line with the current duration.


## Example

```
$ linetimer bash -c "echo Beginning; sleep 1
                     echo Middle;    sleep 2
                     echo End"
[0:00] Beginning
[0:01] Middle
[0:03] End
```


## Installation

If you have [Go](https://golang.org/) you can install from source:

```
go install github.com/nathforge/linetimer/cmd/linetimer
```

Otherwise, download a file from the
[releases](https://github.com/nathforge/linetimer/releases/latest) page.
Extract with `tar xzf FILENAME.tar.gz` and move `linetimer` to a directory
included in your system's `PATH`.


## License

Apache 2.0. See [LICENSE](LICENSE).
