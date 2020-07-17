# stdtimer

Run a command, prefixing each line with the current duration.

Example:

```
$ stdtimer bash -c 'echo Beginning; sleep 2; echo Middle; sleep 2; echo End'
[0:00] Beginning
[0:02] Middle
[0:04] End
```

## License

Apache 2.0. See [LICENSE](LICENSE).
