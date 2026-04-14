# oapi-codegen internal tests

Please note, we have a number of directories here, testing various topics.
Try to find an existing topic to extend when making new tests, or if the
topic is unrelated to what we have, make a new one.

Wherever tests depend on the runtime package, we import our pre-generated
runtime package from the repo root, except parameter tests, where
we generate it inline, since we want to test that code path and dead code elimination,
both of which are most heavily exercised in parameters code.
