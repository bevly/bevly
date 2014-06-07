# Bevly Server

Serves a beverage menu for any configured restaurant.

The server will eventually have two parts:
- A background sync job that updates the menu
- A JSON frontend api

## Running the server

To run the server, install the dependencies, then

     $ go install github.com/bevly/bevly/bevly-server
     $ bevly-server

The server binds to localhost:3000 by default, you may alter this by
setting the PORT environment variable to change the port. The server
is intended to be reverse-proxied behind Apache or nginx, and not
directly exposed to the internet.