# GoGoPing

This is a program that utilizes Go's concurrency features i.e goroutines and channels to ping a URL N times.

The aim was to learn worker pools and fan-in & fan-out pattern. The idea is: A worker pushes N tasks to a `tasks` channel, while M workers pick from it; process it and store it in the results channel, while the `results` get printed to stdout.

Here's an example with `const MAX_GOROUTINES = 5`:
```bash
➜ go run . -url https://youtube.com -n 5
Status code: 200 | Error: <nil> | TimeTaken: 430.81525ms
Status code: 200 | Error: <nil> | TimeTaken: 451.779375ms
Status code: 200 | Error: <nil> | TimeTaken: 482.215375ms
Status code: 200 | Error: <nil> | TimeTaken: 575.35525ms
Status code: 200 | Error: <nil> | TimeTaken: 592.054916ms
Total time: 592.264958ms
```