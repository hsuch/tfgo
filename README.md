# TFGO: The Fun Game Online

## Team Members
App Team: Sam Schlang, Connie Hsu, Zach Saunders, Jason Guo

Server Team: Oliver Zou, Jenny Haar, Brad Lee
 
## How to Compile
Begin by cloning this github repo to your computer.

To compile the server, begin by installing Go (https://golang.org/doc/install).  If you already have Go installed, verify that it is version 1.8 or higher.  Move either the entire github repo or the TFGOServer directory inside your go/src directory.  From within the TFGOServer directory, run “go build” from the command line.

To compile the Xcode project, begin by installing Xcode (https://developer.apple.com/xcode/) and CocoaPods (https://guides.cocoapods.org/using/getting-started.html). Navigate to the “tfgo” folder (the main github repo folder) and run “pod install”.  Then open the Xcode project via TFGO.xcworkspace; opening it via TFGO.xcodeproj won’t work.
 
## How to Run the Code
### Server
 
To start the server, run `./TFGOServer` with command line flags as specified below:
 
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
 
To properly run the server and make it accessible from devices on the same network, determine the IP address of your host machine (on MacOS, run `ifconfig` and use the `inet` field of `en0`), and supply this as the value of `-host` when entering the command.  The port should not need to be changed from 9265.  
 
Example usage: `./TFGOServer -host 192.168.0.101 -port 9265 -v 1`
 
 
### Client
 
To run the game on your phone, you must sign it.  To do so, open Xcode and select TFGO.xcworkspace.  Select the TFGO file (the outermost TFGO file with the blue icon), and under the ‘General’ tab, go to the ‘Signature’ section.  Go to the ‘Team’ drop-down bar and sign into your Apple ID, and then select your Apple ID from the drop-down bar.  Create a new unique bundle identifier if it doesn’t let you sign it, then the game should sign itself.  Please note that if you cannot install the app due to a antique version of iOS on your phone, you should also lower the compiling version accordingly.

Prior to running the game, it is necessary to update the server IP (servadd) and port (servport) variables with the IP and port that you use for running the server.  These variables are located within the "Network.swift" file in the list of variables for the Connection class at the very top of the page.  If you want to run the application on a simulator, just pick a model of iPhone to run on from the drop down menu on the top left of the Xcode page and click the run button.  You can also run the game on your phone by pressing the Play button with your phone selected from the device menu, as long as you signed it earlier.  If your phone is not up to the current build level, you can lower the build level in the ‘Deployment Info’ section of the TFGO file, which is just below the ‘Signature’ section from earlier.
  
## How to Play
The “How to Play” button on the app's main menu contains detailed instructions on setting up, joining, and starting games, as well as the game rules and functionalities of various game components.

In brief, pick a name and icon from the starting screen, which then brings you to the main menu. From there, you can either host or join a game. Once the host starts a game, you can contest control points to earn points for your team, fire at enemy players, acquire pickups scattered throughout the map, swap weapons, or leave the game.

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

It is possible to make a game whose boundaries are so small that there is not enough room for the number of requested objectives, in which case a reduced number of objectives will be generated. Particularly small games might also cause a portion or even all of the bases to be within the objectives boundaries.  This is still a valid game and the objectives and bases will work as normal.  The same goes for if the boundaries are so small that a portion of the objective is outside of the boundaries.

Stepping on an objective does not immediately capture it.  You must be present for a short while and then you capture it and begin to accumulate points.  Once you have captured a point you no longer need to stand on it and it will passively accumulate points for your team.

Occasionally, the boundaries on the Game Info screen's map will not be displayed. This seems to happer more often in private games as opposed to public games, but it happens fairly frequently with both.

Additionally, while it is very rare, a player's icon on the Game Map may change color (from red to blue or blue to red). In the case that this happens to the user, at the beginning of each game a message will be displayed, telling you to go to a certain base; that base is the user's team color, so remembering it will help prevent confusion. In the case that it happens to another player, you can always press on any player icon aside from your own to see a player's name and team color; their team color will always be what team they are on.

In the Lobby screen where it lists games you can join, each game shows the user's distance from the game. Unfortunately, we were not able to implement this calculation, so it will always read "0.0 units".

When counting down to the start of a game, the user's player annotation will appear momentarily and then disappear. It will reappear the moment the game starts (you can verify your specific annotation by pressing an annotation; if it says "My Location", that's your icon).
