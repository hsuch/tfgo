# TFGO: The Fun Game Online

## Team Members
App Team: Sam Schlang, Connie Hsu, Zach Saunders, Jason Guo
Server Team: Oliver Zou, Jenny Haar, Brad Lee, Anders Segerberg
 
## How to Compile
To compile the server, make sure that Go (ver 1.8 or higher) is installed on your computer.  Navigate to go/src, and clone the Git repository there.  Navigate to the TFGOServer directory, then run “go build” from the command line to compile.
 
To compile the xcode project, you must have cocoapods installed (instructions can be found online).  After installing cocoapods and cloning the Git repository, navigate to the “tfgo” folder and run “pod install”.  Then open the xcode project via TFGO.xcworkspace; opening it via TFGO.xcodeproj won’t work.
 
## How to Run the Code
### Server
 
To start the server, run `./TFGOServer` with optional command line flags as specified below:
 
Usage of `./TFGOServer`:
```
  -host string
    	ip address on which to run server (default "127.0.0.1")
  -port string
    	port on which to run server (default "9265")
  -v	int
     verbosity level
     0: print only high-level information
     1: also print all non-periodic JSON messages between client and server
     2: also print periodic JSON messages between client and server
  ```
 
To properly run the server and make it accessible from devices on the same network, determine the IP address of your host machine (on MacOS, running `ifconfig` and use the `inet` field of `en0`; on Windows, run `ipconfig` and look for your IPv4 Address), and supply this as the value of `-host` when entering the command.  The port should not need to be changed from 9265.  
 
Example usage: `./TFGOServer -host 192.168.0.101 -port 9265 -v 1`
 
 
### Client
 
If you want to run the game on your phone, you must sign it.  To sign it, open Xcode and select TFGO.xcworkspace.  Select the TFGO file (the outermost TFGO file with the blue icon), and under the ‘General’ tab, go to the ‘Signature’ section.  Go to the ‘Team’ drop-down bar and sign into your Apple ID, and then select your Apple ID from the drop-down bar.  Create a new unique bundle identifier if it doesn’t let you sign it, then the game should sign itself.  Please note that if you cannot install the app due to a antique version of iOS on your phone, you should also lower the compiling version accordingly.

Prior to running the game, it is necessary to update the server IP (servadd) and port (servport) variables with the IP and port that you use for running the server.  These variables are located within the "Network.swift" file in the list of variables for the Connection class at the very top of the page.  If you want to run the application on a simulator, just pick a model of iPhone to run on from the drop down menu on the top left of the Xcode page and click the run button.  You can also run the game on your phone by pressing the Play button with your phone selected from the device menu, as long as you signed it earlier.  If your phone is not up to the current build level, you can lower the build level in the ‘Deployment Info’ section of the TFGO file, which is just below the ‘Signature’ section from earlier.
  
## Text Description of What is Implemented
### Overview
As planned in our Milestone 4a document, we have built upon the groundwork established in Milestone 3, and implemented most of the remaining features of the game, including lobbies, game boundaries set by the user, creating and capturing control points, creating and getting pickups, and choosing from a variety of different weapons.  The server-side code for payload mode games with a single moving control point exists and has unit tests, but we did not have time to implement the app-side of payload games.

### Server - TFGOServer folder
On the server side, we have written a series of functions to connect to clients and exchange JSON messages with clients, contained in "main.go", "clienthandler.go", and “clientmessage.go.” "struct.go" contains the definitions of the structs we use, the main ones being Game, Player, Team, ControlPoint, Weapon, and PickupSpot.  Go does not support global non-primitive constants, so we have also defined several functions that return important numbers like `BASERADIUS` and `TICK`, to avoid the use of “magic numbers.” The client uses latitude-longitude degrees as the unit of measurement for locations, but for the math related to game logic, it makes the code much more readable to use meters as the base unit, so there are also `degreeToMeter` and `meterToDegree` conversion functions. "setup.go" contains the code for registering a new player, creating a new game, joining an existing game, and starting and stopping the game.  "update.go" contains code which handles all game status updates based on information received from the client, i.e., updating player locations, updating control point status, and applying pickups.  “fire.go” contains the functions related to firing weapons, taking a hit, dying, and respawning.  “pickup.go” contains `consumePickup ()`, which is called when a player is in range of a pickup that is available for consumption.  `consumePickup ()` applies the pickup to the player, sends an update message, marks the pickup as unavailable, and starts a respawn timer, after which the pickup will again be available.  “pickup.go” also contains helper functions implementing this functionality, and a function for generating pickups at game creation.
 
### Client
The app team has written functions to connect to the server as well as to generate and parse JSON messages, located within "Network.swift".  This information is used to populate the GameState, Game, Player, and Objective classes included within the "Game.swift" file.
We have implemented an extensive UI containing several distinct views.  Each view that was fully implemented has a unique ViewController class, which includes all relevant information to each view including outlets and actions for UI elements, delegate functions to control and monitor the map displays and text fields, timers, and send and receive message calls with their corresponding threads.  These features, in addition to what was implemented for the server, enable player identification through the initial screen, game setup through the "Host Game" Button, a join game feature through the “Join Game” button, and includes gameplay such as obtaining and displaying user location along with orientation, displaying objective location, and a simplified tagging system that notifies a user when they are within the spread of an opponent's weapon.

As of iteration 2, the app team has fully implemented both public and private lobbies, setting game boundaries, game passwords, as well as gameplay mechanics such as displaying other players’ locations and the locations of relevant points of interest, such as objectives and pickups.  Several users can now play the complete game as outlined in the “How To Play” section on the game menu!  This will include the ability to capture control points, pick up pickups, use/discard/consume weapons, attack enemies, travel out of bounds, die, respawn, preemptively leave the game, and throw your phone with real world consequences.  For additional convenience, a link to our github page has been included within the apps main menu.
 
##How to Play
Once the simulator opens (or you open the app on your phone), enter a name and select an emoji.  Go to host game, give the game a name and change any other settings as you wish, then scroll down to choose the game boundaries by tapping the map.  After you create the game, you’ll be at a screen waiting for players, and you can press Start Game to enter the game.  Any other players in the game will have their games start simultaneously upon the host starting the game.

If you want to test the game with friends, have them wait until you’ve created a game and are waiting on the “Waiting for Players” screen, then have them go to “Join Game” and tap your game.  Once they tap “Join” in the upper right, they’ll show up in your player lobby, and you can start the game at any time.

For more information on the rules of the game, please select the “How to Play” button from the main menu of the app.

## Invalid Inputs and Other Known Issues
TFGO should be able to be run on any Apple mobile device with gps functionality including iPads, however it is optimized for iPhone 8 and iPhone X.  Older models of iPhones may have less accurate location data including orientation information, which may complicate gameplay.

Since the server is being run locally, it is necessary to change the IP address inside the network.swift file found inside the TFGO folder to the IP address of the machine that the server is being run on.  If the server is being run on a private network such as the UChicago network, it is necessary to use the private IP address.  Due to this fact, it is important to ensure that any mobile device running TFGO is connected to the same network as the server.  Disconnection from the network or the use of cellular data will prevent gameplay.  Reconnection is only possible by restarting the app, either from the computer or by double tapping the home button to open up the active apps view and manually closing the TFGO app.  On a successful connection, the server will print that a new connection has been received.  If a button is pressed and it does not perform an action within five seconds, the app has lost connection to the server and needs to be restarted.

Make sure that you allow TFGO to use location services.  Otherwise the game will not be able to function on a fundamental level.

Lobbies will occasionally not appear in the menu, especially if a user joins the lobby prior to having a game be created.  Returning to the main menu resolves this issue.

Player icons appear incorrectly in the game lobby, but appear correctly in both the game info screen and the waiting for players lobby.

Upon starting a game, the app will prompt the user to head to their team’s base.  In an effort to make testing easier, we have disabled the auto-kill feature if a player fails to reach the base prior to the actual game start.  Likewise, we have set the start countdown to 30 seconds, but in a real game this value would be proportional to the size of the field so that players could reasonably reach their base in time.

Selecting a game with an extremely high time limit or point limit may cause the display to appear incorrectly.  The values are still tracked correctly and this does not affect gameplay.

Melee weapons may appear to not be functioning.  This is not an error, but rather a large clip size that when displayed does not generate much of a visual difference.  Repeated firings can be used to see that the clip is being decremented.

In an attempt to make the game needlessly difficult we elected to not implement any way to view information about the weapons.  They are also in no way balanced.

Using weapons with small clip sizes may appear to not reload properly.  This is a graphical error due to the time it takes for the animation to show that ammo has been used up.

It is possible to make a game whose boundaries are so small that there is not enough room for the objectives such that a portion or even all of the bases are within the objectives boundaries.  This is still a valid game and the objectives and bases will work as normal.  The same goes for if the boundaries are so small that a portion of the objective is outside of the boundaries.

Stepping on an objective does not immediately capture it.  You must be present for a short while and then you capture it and begin to accumulate points.  Once you have captured a point you no longer need to stand on it and it will passively accumulate points for your team.
