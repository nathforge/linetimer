# stdtimer

Run a command, prefixing each line with the current duration.

Example:

```
$ stdtimer bash -c 'echo Beginning; sleep 2; echo Middle; sleep 2; echo End'
[0:00] Beginning
[0:02] Middle
[0:04] End
```


## Installation

Download a binary for your OS from the
[releases](https://github.com/nathforge/stdtimer/releases/latest) section.
Extract with `tar xvf FILENAME.tar.gz`, and move `stdtimer` to a directory
included in your system's `PATH`.


## License

Apache 2.0. See [LICENSE](LICENSE).
