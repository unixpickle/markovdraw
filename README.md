# markovdraw

This is an attempt to use a Markov chain to draw images. You give it a bunch of things you've drawn (e.g. letters and numbers) and it uses line segments from those drawings to form a Markov chain. It produces a panel of random drawings that looks something like this:

![Example Output](samples/output.png)

# Usage

First and foremost, this project is written in Go, so you must [install Go](https://golang.org/doc/install).

In order to train the Markov chain, you need a JSON file full of line drawings. Each line drawing is an array of points (with x/y coordinates between 0 and 120). For example, the file might look like:

```json
[
  [{"x": 15, "y": 15}, {"x": 15, "y": 20}],
  [{"x": 10, "y": 15}, {"x": 15, "y": 20}]
]
```

You can automatically generate such a JSON file using [linetrace/datadraw](https://github.com/unixpickle/linetrace/tree/master/datadraw). You can also find a sample file in [samples/](samples).

Now, run this command as follows:

```
$ go run *.go paths.json output.png
```

Replace `paths.json` and `output.png` with the JSON file of sketches and the output image path, respectively.
