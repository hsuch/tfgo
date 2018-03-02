# TFGO: The Fun Game Online
 
## How to Compile
To compile the server, make sure that Go (ver 1.8 or higher) is installed on your computer. Navigate to go/src, and clone the Git repository there. Navigate to the TFGOServer directory, then run “go build” from the command line to compile.
 
To compile the xcode project, you must have cocoapods installed (instructions can be found online). After installing cocoapods and cloning the Git repository, navigate to the “tfgo” folder and run “pod install”. Then open the xcode project via TFGO.xcworkspace; opening it via TFGO.xcodeproj won’t work.
 
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
 
To properly run the server and make it accessible from devices on the same network, determine the IP address of your host machine (on MacOS, running `ifconfig` and use the `inet` field of `en0`; on Windows, run `ipconfig` and look for your IPv4 Address), and supply this as the value of `-host` when entering the command. The port should not need to be changed from 9265. 
 
Example usage: `./TFGOServer -host 192.168.0.101 -port 9265 -v 1`
 
 
### Client
 
If you want to run the game on your phone, you must sign it. To sign it, open Xcode and select TFGO.xcworkspace. Select the TFGO file (the outermost TFGO file with the blue icon), and under the ‘General’ tab, go to the ‘Signature’ section. Go to the ‘Team’ drop-down bar and sign into your Apple ID, and then select your Apple ID from the drop-down bar. Create a new unique bundle identifier if it doesn’t let you sign it, then the game should then sign itself. 

Prior to running the game, it is necessary to update the server IP (servadd) and port (servport) variables with the IP and port that you use for running the server.  These variables are located within the "Network.swift" file in the list of variables for the Connection class at the very top of the page.  If you want to run the application on a simulator, just pick a model of iPhone to run on from the drop down menu on the top left of the Xcode page and click the run button. You can also run the game on your phone by pressing the Play button with your phone selected from the device menu, as long as you signed it earlier.

Once the simulator opens (or you open the app on your phone), enter a name and select an emoji. Go to host game, give the game a name and change any other settings as you wish, then scroll down to choose the game boundaries by tapping the map. After you create the game, you’ll be at a screen waiting for players, and you can press Start Game to enter the game. Any other players in the game will have their games start simultaneously upon the host starting the game.
If you want to test the game with friends, have them wait until you’ve created a game and are waiting on the “Waiting for Players” screen, then have them go to “Join Game” and tap your game. Once they tap “Join” in the upper right, they’ll show up in your player lobby, and you can start the game at any time.
 
If your phone is not up to the current build level, you can lower the build level in the ‘Deployment Info’ section of the TFGO file, which is just below the ‘Signature’ section from earlier.
 
 
## How to Run Tests
### Server
After compiling, navigate to the TFGOServer directory and run `go test`--this will automatically run all tests in the \*_test.go files. As is the convention in Go, the server test files are all located in the TFGOServer folder along with the code which they test.
 
### Client
In XCode, go to the “TFGOTests” Folder and find “TFGOTests.swift”, then single click the diamond at the very beginning of line 12 (between line 11 and 13 and the number is replaced by a diamond), and this will compile all the code and run the test cases.
 
To test individual functions (such as `testGame()` and `testParse()`), go to the diamond on the line that the function starts at and press it. This will only run the tests in that specific function.

Note that since the structure of some messages changed since Milestone 4a, the unit tests that were made at that time have changed.
 
## Suggested Acceptance Tests
Since the application is a multiplayer game, acceptance testing is best performed with multiple devices. If you have two location-enabled iPhones, or a friend to bother, here’s what you can do:
1. Start the TFGOServer on the IP and port of your choice according to the instructions above.
2. Open the TFGO.xcworkspace file in Xcode. Navigate to the file named “network.swift” and locate the two private variables within the `Connection` class named `address` and `port`. Replace these with the IP and port of your choice.
3. Sign, build, and install the code onto both iPhones according to the instructions above.
4. Start the game on both iPhones. Pick your names and icons.
5. On one iPhone, tap “Host Game” and fill out the “Name” field; others have default values. Remember to set your boundaries appropriately. We are not responsible if any game objectives happen to be in bodies of water or private properties or whatever.
6. Create the game and have the other user join your game.
7. Press “Start Game” to begin. The game will give you a short period of time to move to your base. Once the game has begun, you can attempt to capture the control point, fire your weapon, or gather weapons and other helpful items from pickups scattered around the map.
8. The game ends once time is up or one team has reached the maximum number of points to win.

It is also possible to do some simple acceptance testing with only one phone. In particular, it is possible to perform everything described above except for joining games and hitting other players (although if you're running the server with at least -v=1, you can see the fire messages being sent from client to server every time you fire your weapon). We suggest hosting two games - one where the boundaries are relatively small and centered on you so you can test capturing points and acquiring pickups, and another where the boundaries are entirely outside your location so you can test being out of bounds and needing to respawn.
 
## Text Description of What is Implemented
### Overview
As planned in our Milestone 4a document, we have built upon the groundwork established in Milestone 3, and implemented most of the remaining features of the game, including lobbies, game boundaries set by the user, creating and capturing control points, creating and getting pickups, and choosing from a variety of different weapons. The server-side code for payload mode games with a single moving control point exists and has unit tests, but we did not have time to implement the app-side of payload games.

### Server - TFGOServer folder
On the server side, we have written a series of functions to connect to clients and exchange JSON messages with clients, contained in "main.go", "clienthandler.go", and “clientmessage.go.” "struct.go" contains the definitions of the structs we use, the main ones being Game, Player, Team, ControlPoint, Weapon, and PickupSpot. Go does not support global non-primitive constants, so we have also defined several functions that return important numbers like `BASERADIUS` and `TICK`, to avoid the use of “magic numbers.” The client uses latitude-longitude degrees as the unit of measurement for locations, but for the math related to game logic, it makes the code much more readable to use meters as the base unit, so there are also `degreeToMeter` and `meterToDegree` conversion functions. "setup.go" contains the code for registering a new player, creating a new game, joining an existing game, and starting and stopping the game. "update.go" contains code which handles all game status updates based on information received from the client, i.e. updating player locations, updating control point status, and applying pickups. “fire.go” contains the functions related to firing weapons, taking a hit, dying, and respawning. “pickup.go” contains `consumePickup ()`, which is called when a player is in range of a pickup that is available for consumption.  `consumePickup ()` applies the pickup to the player, sends an update message, marks the pickup as unavailable, and starts a respawn timer, after which the pickup will again be available. “pickup.go” also contains helper functions implementing this functionality, and a function for generating pickups at game creation.
 
### Client
The app team has written functions to connect to the server as well as to generate and parse JSON messages, located within "Network.swift".  This information is used to populate the GameState, Game, Player, and Objective classes included within the "Game.swift" file.
We have implemented an extensive UI containing several distinct views.  Each view that was fully implemented has a unique ViewController class, which includes all relevant information to each view including outlets and actions for UI elements, delegate functions to control and monitor the map displays and text fields, timers, and send and receive message calls with their corresponding threads.  These features, in addition to what was implemented for the server, enable player identification through the initial screen, game setup through the "Host Game" Button, a join game feature through the “Join Game” button, and includes gameplay such as obtaining and displaying user location along with orientation, displaying objective location, and a simplified tagging system that notifies a user when they are within the spread of an opponent's weapon.

As of iteration 2, the app team has fully implemented both public and private lobbies, setting game boundaries, game passwords, as well as gameplay mechanics such as displaying other players’ locations and the locations of relevant points of interest, such as objectives and pickups. Several users can now play the complete game as outlined in the “How To Play” section on the game menu!  This will include the ability to capture control points, pick up pickups, use/discard/consume weapons, attack enemies, travel out of bounds, die, respawn, preemptively leave the game, and throw your phone with real world consequences.  For additional convenience, a link to our github page has been included within the apps main menu
 
## Who Did What
### This Document
Jenny wrote most of the server-related documentation in this README file; Oliver wrote the compilation and running instructions (Anders documented the command line arguments); Jason, Sam, and Zach wrote most of the client-related documentation (Connie wrote the setup instructions and acceptance tests); everyone read through the document to confirm that the information was correct and contributed small edits.
 
### Server Team
Oliver wrote the code which deals with setting up the server and communicating with clients, and worked with the App Team to agree on the format and contents of the JSON messages exchanged by the client and server. He also wrote the code associated with parsing JSON messages from clients and taking the appropriate action. As in the previous iteration, Jenny was in charge of all linear algebra-related functions, which for this iteration involved generating the correct number of pickups distributed evenly but semi-randomly throughout the arena, and writing an algorithm which sorts the boundary vertices into an order which, when connected, produces a planar shape, regardless of the order in which the vertices are received from the client. She also wrote the descriptions for most of the new weapons that have been implemented. Brad primarily dealt with writing tests and debugging. In addition to fixing bugs, he also made sure that the game logic was sound (for instance, in dealing with weapon spread). He wrote the update functions `updateLocation()` and `updateStatus()` as well. For iteration two, Brad and Anders contributed to the logic, implementation, and testing of pickups, working in “pickup.go”, “pickup_test.go”, and "struct.go". Whenever anyone has run into problems, we communicate via Messenger, and we have all looked over each other’s code and helped talk through issues that have arisen. Generally, we notified and consulted each other when making decisions that impact integration, such as function prototypes, return values, and gameplay logic. Everyone contributed to collectively debugging test cases and reviewing each other's code once testing began. 
 
### App Team
Like last time, Sam and Connie did pretty much all of the UI-related stuff, Sam implemented the user interface and Connie drew all the image assets needed for the app. Sam led the app side development and planned all the meet-ups for the team. In terms of functions, Sam programmed the public and private game lobbies and wrote the how to play section.  He implemented the health bar, armor bar, and reload indicator on the game screen and created all of the app side weapon logic in addition to the inventory system. Sam also handled concurrency and the timers used for repeated messages to the server.  Connie, in addition to all the images, also further improved the client-server communication logic as well as planning and implementing further communication messages for the iteration 2 functions. She also implemented the timer for the game’s UI. Zach further improved the map system by displaying other players and objects on the map using Connie’s images, as well as implementing the parsing functions for the app. He drew overlays such as circles to represent objective areas and polygons to represent boundaries in order to further enhance the map. Since we added more functionalities such as capturing points, weapon pick-ups as well as the team system, Zach also implemented more parsing functions and edited previous parsing functions to assist server-client communication as well as working with the rest of the team on improving the game structs for the app. Jason worked on setting boundaries for the host map as well as part of the parsing functions for the communication. He designed the pin dropping mechanism and improved it to a 4-pin system so that the boundary can be set property and then used the parsing functions as well as setBoundaries functions to send the location points to the server to determine the game boundary. In addition, he also worked on part of the unit tests with Zach as well as helped Sam with the “How to Play” Page.
 
## Changes from Earlier Milestones
### Server
**Game struct**: We have added an array of PickupSpots to keep track of the pickups in the game. We have also added PayloadPath and PayloadSpeed. PayloadPath is a Direction representing the direction vector of the payload in a payload mode game, from Red base to Blue Base. PayloadSpeed is the speed at which the payload can move, in meters/second. These are included in the Game struct because in Payload mode games, they do not change throughout the game and they are used every time the server and client communicate, so it doesn’t make sense to recalculate them every time.

**ControlPoint struct**: We have removed the PayloadLoc field, because we realized that we can just use the Location field for that.

**PickupSpot struct**: Since we ended up needing pickups to have some attributes in common as well as a use() function, we wrapped the Pickup interface in a PickupSpot struct. This struct includes a Location, a Pickup interface with the appropriate use() function, an Available bool indicating whether or not the pickup is currently available, and a SpawnTimer, which is started whenever someone picks up the pickup and tells the server how long to wait before the pickup respawns.

### Client
**Message protocol**: Some changes have been made to the message protocol since the last iteration. These changes are fully documented on the Network Messages wiki page: https://github.com/hsuch/tfgo/wiki/Network-Messages

**Main menu design**: As users are able to set their name and icon every time they enter the app, we deemed a “Player Info” page to be unnecessary, and removed it. Instead, we’ve replaced that page with a “How to Play” page, which will be tweaked up until final deployment. We’ve also replaced the “Settings” page, because we couldn’t come up with any meaningful settings to add. It just links to the GitHub now.

**Network functions**: We came across a bug when testing, and realized that our `conn.recvData()` function was not receiving all the data in the buffer if the data was especially long. This was not caught in earlier iterations because messages back then were short and didn’t contain much data; after adding pickup information to messages, this error began to pop up. Oliver offered a lot of help to fix that bug, and ended up implementing a version of recvData() which would read until it saw a newline, which has largely fixed that bug.

**Weapon class**: Exists now on the app side.  All weapons classes inherit from the weapon class and contains information relevant only to the app such as cooldown times.

**Item actions**: You can eat your weapon now.
 
## Testing Changes since Milestone 4a
### Server
To facilitate our work with the app team, we have made a number of small changes to the setup functions, such as setBoundaries() and generateObjectives(), and we have altered our tests to account for these changes. In particular, for setBoundaries(), since ensuring that the boundary vertices are correctly ordered is now the server’s responsibility, we have added a test which involves receiving out-of-order vertices from the client. And we have changed the pickup generation aspect of the generateObjectives() test to account for the maximum pickups per game limit that we have imposed. We have also added tests for several new helper functions, testBorders() and connectTheDots().
 
### Client
For the unit tests, we have created tests to test setters such as “setID”, “setName”, “setDescriptions”, as well as testing if the game is valid. “setMode”, “setTimeLimit”, “setMaxPoints”, and “setMaxPlayers” were updated to no longer return a bool since the UI automatically limited their inputs to only valid values.  We also test the functions that parse Server to Client JSONs by checking if they properly update the gameInfo or not.

We also updated the unit tests based on the tests we had before as well as the updated game structs and parsing functions. [TBD, maybe some specifics] Like last time, since we have assumed that all of inputs we have received from the server side would have been checked thoroughly, so we checked mostly valid data to test if the functions truly works. (We also added some invalid inputs as well just in case to make sure that our functions can handle some unexpected inputs.)

We also performed our own acceptance tests on top of the unit tests, just to ensure the game would run smoothly. Other than a base being unfortunately placed behind a fence, it worked out pretty well and was unexpectedly(?) fun.
 
## Other

 


