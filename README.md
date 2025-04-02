# Concurrent-Programming

table of contents
- [Concurrent-Programming](#concurrent-programming)
  - [Repository Description](#description)
  - [Outline](#outline)
  - [Usage](#usage)

[//]: # (  - [Course Projects]&#40;#course-projects&#41;)


## Description

This course is designed to provide students with a comprehensive understanding of concurrent programming. 
The course will cover certain concepts of concurrent programming with a more nuanced minute approach, 
including threads, goroutines and go runtime, various synchronization techniques, locks, and condition variables, deadlock, and starvation. 
It will include hands-on programming assignments that will make the content of this repository.

## Outline
| capiton                            | synopsis                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | sources                                                            |
|------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| market chain-of-events simulation  | We simulate a simplified scheme of a random chain of market <br> decision process pertaining <br> manufacturing and retail, where: <br> &nbsp; • the manufacturer simulated by boss and workers decides with some probability to produce a certain number of products <br> &nbsp; • the market simulated by clients decides with some probability to buy a certain number of products                                                                                                                                                                                                                                                            | [GOLANG](./introduction)                                           |
| grid traveler - visualisation tool | Web-app tool visualisation for random trace process data, written in React. Provides communication via sockets and web page UI automatically listening for batch updates. Provides an exemplary golang backend                                                                                                                                                                                                                                                                                                                                                                                                                                   | [REACT & GO](./grid-travelers)                                     |
| grid traveler                      | **_logic 1:_**<br>&nbsp;&nbsp; Travelers enter the grid as ghosts, unaware of other travelers<br>**_logic 2_:**<br>&nbsp;&nbsp; each grid tile has a capacity of the amount of travelers that reside at it simultaneously at any point of time. If a traveler embarks upon a tile that has reached its capacity for the whole duration of his "maxDelay" parameter, they mark that with an lowercase letter and stops traveling<br>**_logic 3:_**<br>&nbsp;&nbsp; modifies logic 2 such that each traveler starts at the diagonal line and travels only either up or down or left or right depending on parity of its index | [GOLANG](./grid-travelers/golang) <br> [ADA](./grid-travelers/ada) |
| ...                                | ...                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | ...                                                                |


## Usage

To use any of the projects,
clone the repository and run the following command to install the dependencies

This repository is a collection of projects that are written in different programming languages.
Its main purpose is educational. It provides a comprehensive understanding of concurrent programming.

for example, to run the first project, you can run the following commands:

```bash
  $ git clone <url>
  $ cd <project>
  $ go mod download
  $ go run main.go
```