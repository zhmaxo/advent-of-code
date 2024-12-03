# Advent of Code

My Advent of Code solutions.
I want to play around some quality of life features.

## Motivation

In addition to solving advent puzzles, I'm interested in maintaining a fairly
nice architecture of the project. I don't like scattered single solutions that
run separately, so I want to try the option with a single entry point and
parameterization of the day and input.

## Tools

- Neovim v0.10.2 + LazyVim distribution with extras, especially:
  - lang/go
  - util/rest (thx kulala.nvim)
- Go 1.23.1
- Wezterm
- Nushell with some cli tools installed, including:
  - ripgrep
  - curl

## Roadmap

- [x] Get input via API (see ./advent.http)
- [x] Single entry-point
- [x] Dynamic solution usage
- [x] Dynamic testing without new test func every single day
- [ ] Load input from adventofcode.com using session token
- [ ] Simple HTTP server
- [ ] Multi-year support
