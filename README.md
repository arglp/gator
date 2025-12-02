# gator

gator is a cli tool that lets users register feeds from the internet, aggregates them and stores them in a postgres database.

## prerequisites
- go installed
- postgres installed and running

## install
- use go install github.com/arglp/gator to install gator

## config
- you need a config file in the home directory called ".gaterconfig/json"

## usage
- to run gator simple type gator into your console followed by the command

### commands

#### login
accepts a username as an argument, sets the provided user as the active user.
#### register
accepts a username as an argument, registers a new user with the provided username and sets it as the active user.
#### reset
accepts no argument, delets all users
#### agg
accepts a timestring (f.e. 3s) scrapes through all the feeds after a certain time set by the provided timestring and stores the items in the database
#### addfeed
accepts a feed name and a url as arguments and registers said feed in the database
#### follow
accepts the url of a feed as an argument and registers it as followed for the active user in the database
#### feeds
accepts no argument, shows all registered feeds
#### following
accepts no argument, shows all feeds followed by active user
#### users
accepts no argument, shows all registered users
#### unfollow
accepts the url of a feed as an argument and deregisters it as followed for the active user in the database
#### browse
accepts an optional numbered argument, it shows the most recent posts of the followed feeds of the active user limited by the number given as an argument
