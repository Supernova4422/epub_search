# Epub Search
Epub search runs a simple web server that provides search results from files.
This project is used as a webserver to 
[boby](https://github.com/BKrajancic/boby)

# Running Environment
To use to this project, the only required software is a working go environment.
For installation instructions, see [this page.](https://golang.org/doc/install)

This project includes third party dependencies, so be sure to run
`go get -d -v ./..` to install those dependencies.

Alternatively, docker can be used, and there's a python script that helps run
docker.

# How to use
1. Extract .xhtml files from an epub file that you would like to query from.
This can be done by renaming an epub file to .zip, then extracting files.
2. Ensure that the files contain tables with only two columns. 
3. Place them into a folder. 
4. Use run_in_docker.py to run a server that can be queried, or do usual
compilation things to run this project.

The .xhtml files must contain tables with only two columns. A user searches the
dictionary by making a query such as
http://localhost:8080/?q=Hello&f=foldername the server
will search through all files under foldername until it finds a cell with
"Hello" then it will give the result that's next to it.

The unit tests will help with understanding.

