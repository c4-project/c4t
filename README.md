# `act-tester`

`act-tester` is a work-in progress rewrite of the top-level testing framework
in the [ACT](https://github.com/MattWindsor91/act) project.  The idea is that the new tester:

- is in Go, not Python (and so is a little faster, more maintainable, and easier to deploy);
- has machine-level parallelism;
- takes control of more parts of the test process, such as fuzzing and optimiser levels;
- handles test failures more gracefully.

## Licence

As with the rest of ACT, `act-tester` uses the MIT licence.  (See `LICENSE` for details.)

This part of ACT _doesn't_ use any code from the
[Herdtools7](https://github.com/herd/herdtools7) project, but other parts do,
with all the licensing consequences that entails.

## Acknowledgements

`act-tester` is part of work carried out under the umbrella of the
[IRIS](https://interfacereasoning.com) project.
