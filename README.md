# edgecov
[WIP] Measuring coverage of integration tests.

##### Table of Contents  
[The Idea](#idea)  
[Go + Cucumber](#gocumber)  

<a name="idea"/>

## Idea

Split your functionality into components (not easy, but some approximation will exist - true?). Say we have a machine and we can manage it via the following four commands: `setup`, `start`, `stop`, and `delete`. These are not independent in the sense that executing `stop` after `delete` is useless because there is no machine to stop if it had already been deleted. So we can't expect `stop` to succeed if executed immediately after `delete`. We also assume that we only test for exit codes, i.e. {0, 1}. So either a command succeeds or it does not. We do not check why, how, or any other output.

Given this situation, one can design a digraph `G` with the following nodes: `setup`, `start`, `stop`, `delete` and (directed) edges: `(setup, start)`, `(start, stop)`, `(stop, delete)`, `(delete, setup)`, `(stop, start)`. We intend to have an edge for a step that returns exit code 0 and non-edge for exit code 1 steps. We also need to designate a **source** and a **sink**. This is because: `setup` could be required before `start` is run, so testing the edge `(start, stop)` before `(setup, start)` will return 1 even though it is an edge. If, however, we require that every walk starts at the **source**, we prevent this problem.

Every test run (call it walk) starts at the source and ends when there are no unexplored edges or it hits a dead end. In our example, there are only 2 options:
1. `setup`, `start`, `stop`, `start`, `stop`, `delete`, `start`
2. `setup`, `start`, `stop`, `delete`, `start`, `stop`, `start`
Clearly, both of them visit every edge and are of the same length. Ideally, we would not care which one our test chooses. After execution of such a walk, we would consider the coverage to be 100% (5/5). 
If, instead, the test executed only the following path:
3. `setup`, `start`, `stop`, `delete`, `start`
the edge coverage would be 80% (4/5). 

Edges represent relationships between functionalities, whereas the nodes represent functionalities themselves. Ideally, unit tests would have tested all the nodes by themselves. What integration tests aim for is to execute combinations of them: We know what *A* does, we know what *B* does, but do we know what *B* after *A* does? The example above tries to mimic this and measure how well we tested **dependencies** between functionalities. 

<a name="gocumber"/>

## Go and Cucumber

The tests I work with are written in Go and use some Cucumber framework or whatnot. Anyway, there are Feature files which contain scenarios which in turn contain steps. Each Feature file is supposed to test one "feature" -- whatever that means. The way I imagine this is that there is a library of Scenarios such as:

```
Scenario: (start, stop)
  Given executing "machine status" succeeds
  And stdout contains "Stopped"
  When executing "machine start" succeeds
  Then stdout should contain "Running"
```
Of course, we could simplify it to check only *success* vs *failure* (exit codes). But here it's possible to generically test more, so the scenario can be more detailed without impacting reusability.

Given the above scenario for the edge `(start, stop)`, we can have analogous scenarios for the rest of the edges. Then, the test engine would comprise of the following tasks:

1. Produce a list of paths (Features)
2. Construct feature files from the library of Scenarios
3. Select a subset of the feature files (or all of them) and run integration tests
4. Record info: whenever a scenario is executed - increment that edge's counter in the heatgraph and in the list of edges
5. Produce the statistic: edgecov in % (nonzero entries in list of edges normalized by num of edges)

Over multiple runs, the heatgraph will show that the edges starting at the source are visited most often while some other edges are not used that frequently. This can inform writing of additional manual tests.

  
