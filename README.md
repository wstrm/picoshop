picoshop
========
Minimalistic web shop.

## Project status
| Build status | Test coverage |
|:------------:|:-------------:|
| [![Build Status](https://travis-ci.org/willeponken/picoshop.svg?branch=master)](https://travis-ci.org/willeponken/picoshop) | [![Coverage Status](https://coveralls.io/repos/github/willeponken/picoshop/badge.svg?branch=master)](https://coveralls.io/github/willeponken/picoshop?branch=master) |

## Project structure

### Code
![Code project structure tree](https://github.com/willeponken/picoshop/blob/master/doc/patterns/picoshop-project-structure_rev1.png)

<details>
<summary>Description</summary>

 * /cmd - main entry points for each binary
 * /controller - routes according to MVC pattern
 * /doc - documentation
 * /middleware - interceptors for routes
 * /model - interact with database according to MVC pattern
 * /static - content that is served by the web server
 * /tool - developer utilities
 * /view - HTML views that are rendered for each web page

</details>

### Database
![Database ER scheme](https://github.com/willeponken/picoshop/blob/master/doc/database/picoshop_sql-er-diagram_rev2.png)
