/*
(c) 2020 Stan Drozdov https://github/rudestan

Broadlink Api Engine cli application that can control the Broadlink devices using pre-configured JSON file. Can be run
in two modes:

	1) web server mode that listens the incoming requests and executes commands via corresponding requests or
       parses the Alexa API request JSON and executes matched scenarios or commands if any.
	2) command line mode for direct execution of the command or scenario
*/

package main
