# TFGO: The Fun Game Online
 
## How to Compile
To compile the server, make sure that Go is installed on your computer. Create a folder go/src, and clone the Git repository there. Navigate to the TFGOServer directory, then run “go build” from the command line to compile.
 
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
  -v	print sent and received JSON
  ```
 
To properly run the server and make it accessible from devices on the same network, determine the IP address of your host machine (this can be found on MacOS by running `ifconfig` and using the `inet` field of `en0`), and supply this as the value of `-host` when entering the command. The port should not need to be changed from 9265. 
 
`-v` enables verbose server-side logging, which displays all JSON messages received and sent by the server
 
Example usage: `./TFGOServer -host 192.168.0.101 -port 9265 -v`
 
 
### Client
 
Open Xcode and select TFGO.xcworkspace. Select the TFGO file (the outermost TFGO file with the blue icon), and under the ‘General’ tab, go to the ‘Signature’ section. Go to the ‘Team’ drop-down bar and sign into your Apple ID, and then select your Apple ID from the drop-down bar. If you want to run it on your simulator, just pick a simulator to run on and click the run button. You can also run the game on your phone by pressing the Play button with your phone connected to the computer. 
Once the simulator opens (or you open the app on your phone), enter a valid name and then select an emoji (it won't be highlighted, but will still be selected). Go to host game, give the game a name and change any other settings as you wish (boundary is currently hardcoded, which will be implemented in next iteration), then scroll down to create the game. This will take you to the game lobby, and you can press Start Game to enter the game (note that if there are two players in the game, both will have to press the Start Game button to enter the game).
 
As for joining a game, there is a...not so conspicuous pink button on the welcome page following a login. For iteration 1, since we will only be testing 1 game, we will not be using lobbies, so we decided to bypass the lobby step using the pink button. So if you’re trying to join the single game, press the pink button to automatically join the game.
 
If your phone is not up to the current build level, you can lower the build level ‘Deployment Info’ section, which is just below the ‘Signature’ section.
 
 
## How to Run Tests
### Server
After compiling, navigate to the TFGOServer directory and run `go test`--this will automatically run all tests in the \*_test.go files. As is the convention in Go, the server test files are all located in the TFGOServer folder along with the code which they test.
 
### Client
In XCode, go to the “TFGOTests” Folder and find “TFGOTests.swift”, then single click the diamond at the very beginning of line 12 (between line 11 and 13 and the number is replaced by a diamond), and this will compile all the code and run the test cases.
 
To test individual functions (such as `testGame()` and `testParse()`), go to the diamond on the line that the function starts at and press it. This will only run the tests in that specific function.
 
## Suggested Acceptance Tests
Since the application is a multiplayer game, it would be a bit difficult to perform any Acceptance Tests with only one phone. If you have two location-enabled iPhones, here’s what you can do:
1. Start the TFGOServer on the IP and port of your choice according to the instructions above.
2. Open the TFGO.xcworkspace file in Xcode. Navigate to the file named “network.swift” and locate the two private variables within the `Connection` class named `address` and `port`. Replace these with the IP and port of your choice.
3. Sign, build, and install the code onto both iPhones according to the instructions above.
4. Start the game on both iPhones, and pick a different name on both.
5. On one iPhone, tap “Host Game” and fill out the “Name” field; others have default values. After hitting “Create Game”, you should see the “Waiting for Players” screen.
6. On the other iPhone, tap the big ugly pink button on the bottom. It should also bring you to the “Waiting for Players” screen. If you’re running the server in verbose mode, you should be able to see the server send PlayerListUpdate messages to both phones.
7. Have the host phone tap “Start Game”. Both phones should begin transmitting LocationUpdate messages with their location and orientation. Pressing the red button will cause the app to transmit a Fire message. If you have been hit, you will receive a notification.
 
 
## Text Description of What is Implemented
### Overview
As planned in our Milestone 2 document, we have focused on laying the groundwork upon which everything is built, and implementing two use cases: shooting other players and capturing control points. To make these actions possible, we have also had to implement creating and joining a game and taking damage/being killed.
 
### Server - TFGOServer folder
On the server side, we have written a series of functions to connect to clients and exchange JSON messages with clients, contained in main.go, clienthandler.go, and clientmessage.go. Struct.go contains the definitions of the structs we use, the main ones being Game, Player, Team, ControlPoint, and Weapon. Go does not support global non-primitive constants, so we have also defined several functions that return important numbers like BASERADIUS and TICK, to avoid the use of “magic numbers.” The client uses latitude-longitude degrees as the unit of measurement for locations, but for the math related to game logic, it makes the code much more readable to use meters as the base unit, so there are also degreeToMeter and meterToDegree conversion functions. Setup.go contains the code for registering a new player, creating a new game, joining an existing game, and starting and stopping the game. Movement.go contains handleLoc(), which is called whenever a player’s location changes. Depending on the result of the movement, handleLoc() calls helper functions to deal with a player being out-of-bounds and a change in control point capture status. Fire.go contains functions related to shooting another player. Fire() determines whether the shot hits another player, and if so calls takeHit(). TakeHit() reduces the affected player’s health and, if the player is now dead, calls awaitRespawn(), forcing the player to return to their base and wait to respawn.
 
### Client
 
## Who Did What
### This Document
Jenny wrote most of the server-related documentation in this README file; Oliver wrote the compilation and running instructions (Anders documented the command line arguments); ____ wrote most of the client-related documentation; everyone read through the document to confirm that the information was correct and contributed small edits.
 
### Server Team
Oliver wrote the code which deals with setting up the server and communicating with clients, and worked with the App Team to agree on the format and contents of the JSON messages exchanged by the client and server. He also wrote the code associated with parsing JSON messages from clients and taking the appropriate action. Jenny wrote `fire()` and its helper function, `canHit()`. She also designed the game’s coordinate space and handled all of the linear algebra-related tasks, such as setting game boundaries, choosing team base locations, choosing control point locations, and determining whether a location is within the game boundaries. Brad primarily dealt with writing tests and debugging. In addition to fixing bugs, he also made sure that the game logic was sound (for instance, in dealing with weapon spread). He wrote the update functions `updateLocation()` and `updateStatus()` as well. Anders looked into solutions for hosting the server in a test environment, eventually arriving at the method currently used for testing, and looked into production hosting, ruling out shared web hosting as a viable platform due to poor performance. He also wrote `takeHit()`, and contributed to `respawn()`. Whenever anyone has run into problems, we communicate via Messenger, and we have all looked over each other’s code and helped talk through issues that have arisen. Generally, we notified / consulted each other when making decisions that impact integration, such as function prototypes, return values, and gameplay. 
 
### App Team
Sam did pretty much all of the UI-related stuff, whether it came to building pages/managing buttons and the like. He also handled all of the ViewControllers aside from GPS tracking, and created the bases for a lot of the classes. Additionally, he wrote the tests for the game. Connie made a mock-up design early on that was very helpful for designing the UI, and she additionally worked with Oliver and the rest of the Server Team to decide the layout of network messages sent between the server and client, and did pretty much all network related coding on the client side. Zach worked on map-related content, from tracking the user’s location on the Game Map, in addition to constantly updating the user’s location and orientation so that the values could be sent to the server. Additionally, he wrote all the functions used for parsing JSON files sent from the server to the client, and updated game and user information using the data received from server to client JSON files, and added to structs so that this information could be stored. Jason worked on building and sending JSON messages from the client to the server, and also implemented some part of the definitions and also worked with Zach to write the unit tests for the parsing functions. In short, Sam and Connie worked on setting the groundwork for the app and designing the UI, Connie worked with Oliver and other members of the server side to setup the client network connection, Zach worked on maps in addition to tracking the user’s location and updating their position and orientation, and Zach and Jason worked on formatting messages between the server and the client.
 
## Changes from Earlier Milestones
**Location struct**: Since Go is not strictly object-oriented and does not have public and private fields or functions, there was no reason to write a getter or setter function.

**Direction struct**: To aid in some of the geometry/linear algebra-related calculations, we created the Direction struct, which contains two floats, X and Y, and represents a 2D direction vector.

**Border struct**: To enable us to easily test whether a point is within the game boundaries even with an arbitrary number of boundaries that are not necessarily axis-aligned and don’t have to form any specific shape, we created the Border struct to represent game boundaries. It is the parametric representation of a line segment, and contains a Location (P) and a Direction (D). Because of the way in which D is calculated during setup, the line segment is always defined by t-values between 0 and 1.

**Game struct**: We have added a Description, a PlayerLimit, a PointLimit, and a TimeLimit to the struct. We have moved StartPoints to the Team struct (they are the base locations). As described above, Boundaries is now a slice of Borders rather than Locations, and it need not be of length 4.

**Player struct**: We have added the Conn, Chan, and Encoder fields to link a player with a specific client and enable message-sending between the server and that client. We have also added an Orientation field, which indicates the direction in which the player is currently facing, and an OccupyingPoint field, which is a pointer to the ControlPoint that the player is currently occupying, if any. The Player struct also no longer contains an ActiveWeapon field; that is left up to the client, so whenever a shot is fired, the client informs the server which weapon was used.

**Team struct**: Teams no longer have a name, icon, or color; there are always exactly two and they are simply referred to as RedTeam and BlueTeam. They also no longer have a list of Players, because we found that Teams do not need to reference players, and this avoids mutual references. We also found that we did not need the getTeamLocs() function, which was intended as a helper for canHit(). And Team now includes a Base and BaseRadius, indicating the position and size of that team’s home base.

**ControlPoint struct**: We have added a PayloadLoc field to indicate the payload’s current location.

**Pickup interface**: We have decided to make Pickup an interface rather than a struct, and because of this it is only required to have a use() function, and we will not be implementing this until Iteration II.

**Weapon struct**: We have added a Name field for ease of identification. Since `nearestHit()` would only have been a few lines of code and would only have been called once, we decided to remove it. `canHit()` now returns a float64 (the distance to the target or MaxFloat64 if the shot does not hit the target) rather than a bool, because this was more useful to the `fire()` function.
 
## Testing Changes since Milestone 3a
### Server
The server team has added unit tests for all of the major functions that do not involve server-client network interaction (such as those in `clienthandler.go` or `clientmessage.go`) or a random element (in `setup.go`). Our tests are written in the various `*_test.go` files. Network related and random element containing functions were purposefully excluded not only because it is difficult to test them, but also because they should be covered by the acceptance tests. Certain other functions, such as those in `struct.go`, were also not tested due to their straightforward nature (generally acting as constants or performing basic and arbitrary conversions that can’t really be tested for "correctness").
 
### Client
There are two parts of the testing we did in the App team. The first part is to test the functionalities in the code and the second part is for UI Testing. As for the unit tests part for the functionalities, we have created tests to test setters such as “setID”, “setName”, “setMode”, “setTimeLimit”, “setMaxPoints”, “setMaxPlayers”, “setDescriptions”, as well as testing if the game is valid. We also test the functions that parser Server to Client JSONs by checking if they properly update the gameInfo or not. As for the UI Testing, since the whole UI suite is not completed yet, so we decided to postpone UI Testing to later iterations.
 
## Other
Although they are shown in the class diagram in Milestone 2, as per the description of Iteration I given in the Milestone 2 document, the functions `setWinner()`, `heal()`, `checkInventory()`, `addItem()`, and `removeItems()`, and everything related to Pickups, have not been implemented as part of Iteration I. They will be handled in Iteration II.
 
