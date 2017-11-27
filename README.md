# Typy Bird

## Overview

Typy Bird is a multiplayer web-based typing game. Players up to 4 can share a
single game room and play against each other. Each player will control a
bird-like game object called **Typie** and make it move forward by typing the
whole word prompted to the player. If one of the Typies reaches the finishing
line, the player who controls that Typie is then considered as the only winner
in this game.

## Core Technologies

* HTML5 Canvas: display game UI.
* JavaScript: build a simple game rendering engine.
* Aseprite: create texture atlas and spritesheet.
* Websocket: notify players about game state.
* Golang: build a game API server.

## Work Breakdown Structure

1. Build a simple game rendering engine (Zico Deng).
2. Design a database schema that stores players information and their best
   records. Each player is uniquely identified by the Typie's name.
3. Build a game waiting room page that allows players to create a new Typie and
   join the game. New players can only join the game when the game room is
   marked as **get ready** and **capacity (4 players max) is not reached yet**.
   This page should also contain a table that shows all previous players
   (Typies) information and their best records. The table should also be able to
   switch between chronological order or order by historical records.
4. Build a game page that consumes the data sent by the game API server, renders
   the updated game state, and allows players to send their typed words to the
   game API server (Zico Deng).
5. Build a game API server that stores players information and their best
   records to MongoDB, validates typed word sent by players, and broadcast
   updated game state to client through websocket.
6. Design game environment and entities (Zico Deng).
7. Create texture atlas and spritesheet (Zico Deng).
8. Design game dictionary. Random words will be picked from this dictionary
   during the game. Players need to type them correctly in order to make their
   controlled Typies move forward.

Note: this work breakdown structure only shows a very high-level idea of work
breakdown pieces and who will be building which piece. The details of each task
is hard to predict and measure, and the entire WBS is subjective to change.

## Contributors

* Zico Deng
* Eric Jacobson
* Matthew Bond
* Emily Zhai

# Development

## Collaboration

* Always create a new work branch for different features or tasks.
* Name your work branch as your name + feature you are implementing. e.g.
  `zico-game-ui`.
* Before you create a new work branch, always make sure to pull the latest code.
* Never push your code to master branch directly.
* Each commit should only concern its primary task. For example, if you
  encounter a bug and you know a quick fix for it, do not fix it in your current
  commit. Create a new commit for fixing this bug instead.
* Capitalize the first letter in your commit message.

## Client-Side

## Server-Side
