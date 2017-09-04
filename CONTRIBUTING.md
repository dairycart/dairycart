# Guidance on how to contribute

> All contributions to this project will be released under the MIT License. By submitting
> a pull request or filing a bug, issue, or feature request, you are agreeing to comply
> with this waiver of copyright interest.

There are two primary ways to help:

- Using the issue tracker, and
- Changing the code-base.

Contributions to this project that do not adhere to the [Code of Conduct](CODE_OF_CONDUCT.md) will be rejected.

## Using the issue tracker

Use the issue tracker to suggest feature requests, report bugs, and ask questions.
This is also a great way to connect with the developers of the project as well
as others who are interested in this solution.

Use the issue tracker to find ways to contribute. Find a bug or a feature, mention in
the issue that you will take on that effort, then follow the _Changing the code-base_
guidance below.

## Changing the code-base

Generally speaking, you should fork this repository, make changes in your
own fork, and then submit a pull-request. All new code should have associated unit
tests that validate implemented features and the presence or lack of defects.
Additionally, the code should follow any stylistic and architectural guidelines
prescribed by the project. In the absence of such guidelines, mimic the styles
and patterns in the existing code-base.

## Best practices

- Route handlers should be returned by functions that accept things like cookie stores and database connections.
- Inner functions called by route handlers should wrap errors using [errors.Wrap](https://godoc.org/github.com/pkg/errors#Wrap)
- Functions that accept IDs as arguments should always accept them as `uint64`s. They may more conveniently arrive as strings, but [strings can be converted to `uint64` fairly easily](https://golang.org/pkg/strconv/#ParseUint) (so long as the string adequately represents a `uint64`).

note: The following bullet points are things I always intend to follow, but frequently forget. This list is not meant as a shame-on-you-if-you-don't-do-these-things list in any way, simply a list of things I might ask you to fix or implement when you issue a PR. It's also a convenient list of things you can look for in the existing codebase and issue a PR to fix.

## Ways you can always help

### Documentation

Adding or updating documentation cannot have its importance understated. There are a ton of functions in Dairycart, and not all of them have documentation at all, much less great documentation (because I'm bad at it). Having great documentation is probably the greatest feature any software can ever have, and writing great documentation takes as much skill and effort as any software takes to write.

### Direct unit tests

Ensuring that every function has a direct unit test. I'll be working on a tool for this separately, but say you have the following code:

```go
func A() string {
    return "hello"
}

func B() string {
    return "hi"
}

func C() string {
    return "greetings"
}

func sayHello() {
    fmt.Println(A())
    fmt.Println(B())
    fmt.Println(C())
}
```

and you build tests for `A`, `C`, and `sayHello`. `B` will show up in the coverage report as tested, but the moment `sayHello` stops calling `B` you have a new coverage gap that wouldn't have been there if you had just built the appropriate unit tests in the first place. Again, the tool I plan on building will detect all this, and eventually even fail the Travis builds.

A big thanks to the [CFPB](https://github.com/cfpb/open-source-project-template) for the template this file is based on.