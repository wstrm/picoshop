picoshop
========
Minimalistic web shop.

## Project status
| Build status | Test coverage |
|:------------:|:-------------:|
| [![Build Status](https://travis-ci.org/willeponken/picoshop.svg?branch=master)](https://travis-ci.org/willeponken/picoshop) | [![Coverage Status](https://coveralls.io/repos/github/willeponken/picoshop/badge.svg?branch=master)](https://coveralls.io/github/willeponken/picoshop?branch=master) |

## Code structure
 * cmd - main entry points for each binary
 * controller - routes according to MVC pattern
 * middleware - interceptors for routes
 * model - interact with database according to MVC pattern
 * view - HTML views that are rendered for each web page
