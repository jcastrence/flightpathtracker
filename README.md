# Flight Path Tracker
---
## Running Application
From root directory:
```make```
Application will run on localhost:8080 by default

---
## Endpoints

### /calculate
Example Request
```
curl --location --request GET 'localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data '[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]'
```
Response
```
[
    "SFO",
    "EWR"
]
```
### Usage

The calculate endpoint is used to determine a person's originating source and final destination airports. For example, a person with the following flights
```
[["ATL", "EWR"], ["SFO", "ATL"]]
```
would have a resulting source and destination of
```
["SFO", "EWR"]
```
since it can be inferred that the originating source was SFO and the final destination was EWR.


### Input

The endpoint accepts a GET request where the body of the request is raw JSON in the form of a list of string lists. Make sure to set the Content-Type header to application/json.

Each string represents an airport code. Each of the string lists must be of length 2, as any other length would make it ambiguous as to which airports are the source and the destination. The following input would be considered invalid
```
[["ATL", "EWR"], ["SFO", "ATL", "DEN"]]
```
as would
```
[["ATL"], ["SFO", "ATL"]]
```
Lists may also not contain duplicate airports as it is assumed to be a user error (one typically does not originate and land from the same airport)
```
[["ATL", "ATL"], ["SFO", "ATL"]]
```
Airport codes must also be of a fixed string length of 3 containing only capital letters from A to Z. The following inputs would be invalid
```
[["ATL", "SF0"]]
[["atl", "SFO"]]
[["ATL", "SF"]]
```
### Dependencies
At it's core the microservice is a REST API. The goal was to keep it simple and performant. [Echo](https://echo.labstack.com/) was chosen as the web framework due to its lightweight, flexibility, extensibility, ease of use, scalability, and performance. [Testify](https://github.com/stretchr/testify) was chosen as the testing framework due to its proven robustness and effectiveness. Its philosophy of simplicity make it easy to write a lot of tests, quicker, providing better coverage. Both frameworks are actively and well maintained. This is perhaps an even more important aspect to consider when choosing to utilitize a framework.

### Implementation
The crux of the application lies in its implementation of the flight path tracking algorithm. The service is essentially a topological sort. The list of flights can be considered a graph (even more specifically as a tree) represented by a directed edge list, where each flight is an edge and each airport is a node.

Though only logical flight paths were provided in the examples, an environment where users exclusively provide logical flight paths is naive, and it is assumed that the user may provide flight paths that would make finding a definitive source and destintation ambiguous. While checks are in place to prevent nodes from having multiple children, essentially constraining the tree to a linked list, it is still important to consider the core algorithm of a topological sort as it prevents cycles. Imagine a flight path like
```
[["JFK", "ATL"], ["ATL", "LAX"], ["LAX", "JFK"]]
```
Without further information like time of flights, it is impossible to determine which of these airports was the original source and final destination. These situations should return an error.

Another relevant consideration is that unlike a standard topological sort where a node may have multiple children and parents, either of these circumstances would also make it ambiguous.
```
[["JFK", "LAX"], ["ATL", "LAX"]]
```
In the above example it is impossible to determine a source as there are two possible sources.
```
[["JFK", "LAX"], ["JFK", "ATL"]]
```
Conversely, the above example would make it impossible to determine a destination.

Two final considerations are disconnected flight paths
```
[[JFK, ATL], [ATL, LAX], [CLT, DCA], [DCA, DEN]]
```
and repeated flights
```
[[JFK, ATL], [ATL, LAX], [JFK, ATL]]
```
While these are less ambiguous (could decide to choose the first flight path seen in a set of disconnected flights and could choose to ignore repeated flights) in the context of a flight path tracking system, they still appear illogical and are considered errors to protect users from unintended inputs.

Ultimately, a hashmap approach was chosen to maximize time and space complexity (linear complexity for both). Differing from a typical topological sort however, an optimization was made by assuming each node will have at most one child (more explanation will be provided further). This allows forgoing a hashset, commonly used to track which nodes have been visited to prevent cycles.

The algorithm works in three linear time phases:

1. Converting the flight list to a hashmap where sources point to their destination

2. Iterating through the sources to remap and reduce each source to its final destination, potentially creating multiple reduced sub-mappings

3. Iterating through the reduced sub-mappings to determine if a valid result is present

Of the three phases the most critical is phase two as determining a way to mark nodes as either visited or sub-mappings is not necessarily trivial. This is also where the most optimization lies.

A DFS like approach was implemented to traverse a range of nodes, by continually scanning the value of each key until either an existing sub-mapping is found (the only sub-mapping in the case of a happy path) or the final destination, which would not have an entry in the hashmap.

To forgo using an additional data structure like a hashset to keep track of visited nodes, the value (destination) of each key-value pair (source -> destination) is altered. Prefixing an airport code string with a "-" is used to denote a node that was already visited. This allows the algorithm to detect cycles during iteration and skip redundant iterations.

When traversing the nodes, if the source ever finds a source that reduces to the final destination, there is no need to continue iterating. A "*" is used to denote a new 'head' source, simply a source that has found the final destination. This strategy allows the algorithm to easily detect an erroneous input in the form of multiple heads.

Though this document describes the exact implementation of the visited and head marking strategy, an effort was made to decouple this implementation from the core algorithm. This is because of a final language specific consideration; strings in Go are immutable. This implies that each marking of a visited or head node would require a completely new string. Since reasonable input data assumptions were made that every airport code would consist of a very small fixed size, and since the overall complexity of marking a visited or head node would still follow a linear pattern, it was decided that string rewrites would not be a significant cost. However, it would be wise to consider that, should the assumptions about the input data change (like longer / more flexible airport codes), a different strategy might be more appropriate, i.e. converting the strings to a mutable or int-based data type.