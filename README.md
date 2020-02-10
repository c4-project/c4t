# `act-tester`

`act-tester` is a work-in progress rewrite of the top-level testing framework
in the [ACT](https://github.com/MattWindsor91/act) project.  The idea is that the new tester:

- is in Go, not Python (and so is a little faster, more maintainable, and easier to deploy);
- has machine-level parallelism;
- takes control of more parts of the test process, such as fuzzing and optimiser levels;
- handles test failures more gracefully.

## Components

We envision that `act-tester` will have the following components:

- [ ] `act-tester-plan`, which creates an initial test plan over some Litmus tests and compilers;
- [ ] `act-tester-fuzz`, which runs `act-fuzz` over a test plan to create a more useful plan;
- [ ] `act-tester-litmus`, which runs `litmus` over a test plan to generate compilable harnesses;
- [ ] `act-tester-cp`, which copies a test plan to one of its remote machines;
- [ ] `act-tester-run`, which runs compilers over Litmus harnesses to produce results;
- [ ] `act-tester-direct`, which combines the above into a looping test campaign.

## Licence

As with the rest of ACT, `act-tester` uses the MIT licence.  (See `LICENSE` for details.)

This part of ACT _doesn't_ use any code from the
[Herdtools7](https://github.com/herd/herdtools7) project, but other parts do,
with all the licensing consequences that entails.

## Acknowledgements

`act-tester` is part of work carried out under the umbrella of the
[IRIS](https://interfacereasoning.com) project.