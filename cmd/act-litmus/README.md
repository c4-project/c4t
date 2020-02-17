# act-litmus

`act-litmus` is a wrapper around [`litmus`](https://github.com/herd/herdtools7)
that handles various corner cases that arise when using `litmus` in ACT tests.  Specifically, it:

- automatically enables `-ascall true` if the input Litmus test has return values;
- patches the test harness

It depends on both `litmus7` and `act-c` being in `PATH`.

## Usage

`act-litmus -carch ARCH -o DIR FILE`

Both `-carch` and `-o` are mandatory.  `act-litmus` always runs in C11 harness outputting mode,
and doesn't support any other operating mode.