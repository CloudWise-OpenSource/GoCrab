// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

// package GoCrab provides a minimal framework for creating and organizing command line
// Go applications. cli is designed to be easy to understand and write, the most simple
// cli application can be written as follows:
//   func main() {
//     cli.NewApp().Run(os.Args)
//   }
//
// Of course this application does not do much, so let's make this an actual application:
//   func main() {
//     app := cli.NewApp()
//     app.Name = "greet"
//     app.Usage = "say a greeting"
//     app.Action = func(c *cli.Context) {
//       println("Greetings")
//     }
//
//     app.Run(os.Args)
//   }
package GoCrab

import (
	"os"
)

func Run() {
	CrabApp.Run(os.Args)
}
