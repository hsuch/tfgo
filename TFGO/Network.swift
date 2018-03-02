//
//  Network.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import SwiftSocket
import SwiftyJSON
import MapKit

class Connection {
    private var servadd: String = "10.150.236.157" // to be replaced with real server ip
    private var servport: Int32 = 9265
    private var client: TCPClient

    func sendData(data: Data) -> Result{
        return client.send(data: data)
    }

    func recvData() -> Data {
        var response = Data()
        while true {
            guard let data = client.read(1024*10, timeout: 1)
                else { break }
            response += data
        }
        return response
    }

    init() {
        client = TCPClient(address: servadd, port: servport)
        DispatchQueue.global(qos: .background).async {
            switch self.client.connect(timeout: 10) {
            case .failure:
                print("Connection failed")
            default:
                print("Successful connection")
            }
        }
    }
}

/* handleMsgFromServer(): takes the incoming message and splits it if there are multiple messages in buffer,
 * then parses all of them
 */
func handleMsgFromServer() -> Bool {
    print("just got in handle")
    let conn = gameState.getConnection()
    var received: Data? = nil
    while(received == nil)
    {
        print("whiling")
        received = conn.recvData()
    }
    print("got message from server")
    let recvStr: String = String(data: received!, encoding: .utf8)!
    var strArray: [String] = []
    recvStr.enumerateLines { line, _ in
        strArray.append(line)
    }
    for line in strArray {
        var data = try! JSONSerialization.jsonObject(with: line.data(using: .utf8)!, options: []) as! [String: Any]
        let type = data.removeValue(forKey: "Type") as! String
        print("starting parsing message")
        if (!parse(data: data, type: type)) {
            return false
        }
    }
    return true
}

/* parse(): convert data array into appropriate data struct depending on message type */
/* URL for the structure of the messages : https://github.com/hsuch/tfgo/wiki/Network-Messages */
func parse(data: [String: Any], type: String) -> Bool {
    switch type {
    case "PlayerListUpdate":
        return parsePlayerListUpdate(data: data)
    case "AvailableGames":
        return parseAvailableGames(data: data)
    case "GameInfo":
        return parseGameInfo(data: data)
    case "JoinGameError":
        return parseJoinGameError(data: data)
    case "LeaveGame":
        return parseLeaveGame()
    case "GameStartInfo":
        return parseGameStartInfo(data: data)
    case "GameUpdate":
        return parseGameUpdate(data: data)
    case "StatusUpdate":
        return parseStatusUpdate(data: data)
    case "VitalsUpdate":
        return parseVitalsUpdate(data: data)
    case "GameOver":
        return parseGameOver(data: data)
    case "AcquireWeapon":
        return parseAcquireWeapon(data: data)
    case "PickupUpdate":
        return parsePickupUpdate(data: data)
    default:
        return false
    }
}

/* The parser for the PlayerListUpdate message */
func parsePlayerListUpdate(data: [String: Any]) -> Bool {
    print("starting parsePlayerListUpdate")
    if let info = data["Data"] as? [[String: Any]] {
        
        // we want to build a list of the players from the given JSON message
        var players = [Player]()
        
        for player in info {
            if let name = player["Name"] as? String, let icon = player["Icon"] as? String {
                // build a list of the passed in players
                players.append(Player(name: name, icon: icon))
            }
        }
        
        // set the gameState's current game's player list to the built list
        gameState.getCurrentGame().setPlayers(toGame: players)
        print("finished updating players")
        return true
    }
    return false
}

/* The parser for the AvailableGames message */
func parseAvailableGames(data: [String: Any]) -> Bool {

    if let info = data["Data"] as? [[String: Any]] {
        
        // Similar to parsePlayerListUpdate, we want to build a list of the parsed games
        var games = [Game]()
        
        for game in info {
            if let id = game["ID"] as? String {
                if let name = game["Name"] as? String, let mode = game["Mode"] as? String, let loc = game["Location"] as? [String: Any], let players = game["PlayerList"] as? [[String: Any]], let hasPassword = game["HasPassword"] as? Bool {
                        
                    // build the new game given the passed in information
                    let newGame = Game()
                    newGame.setID(to: id)
                    if newGame.setName(to: name) {}  // these games will always give a valid name
                    newGame.setMode(to: Gamemode(rawValue: mode)!)
                        
                    if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                        newGame.setLocation(to: MKMapPointMake(x, y))
                    }
                        
                    newGame.setHasPassword(to: hasPassword)
                        
                    for player in players {
                        if let name = player["Name"] as? String, let icon = player["Icon"] as? String {
                            // add the player to the game's list of players
                            newGame.addPlayer(toGame: Player(name: name, icon: icon))
                        }
                    }
                    // append the created game to our new array of found games
                    games.append(newGame)
                }
            }
        }
        
        // set the currentGame's array of found games to the built array
        gameState.setFoundGames(to: games)
        return true
    }
    return false
}

/* The parser for the GameInfo message */
func parseGameInfo(data: [String: Any]) -> Bool {
    
    if let info = data["Data"] as? [String: Any] {
        if let desc = info["Description"] as? String, let playerNum = info["PlayerLimit"] as? Int, let pointLim = info["PointLimit"] as? Int, let timeLim = info["TimeLimit"] as? Int, let id = info["ID"] as? String {
            
            // build a new game with the passed in information
            let newGame = Game()
            newGame.setMaxPlayers(to: playerNum)
            if newGame.setDescription(to: desc) {}
            newGame.setMaxPoints(to: pointLim)
            newGame.setTimeLimit(to: timeLim)
            newGame.setID(to: id)
            
            // since this game was selected from the abailable game list, there must be an entry
            // in the available game list that corresponds to this game, and that has information
            // on the game that we need
            let foundGame = gameState.findGame(to: id)
            if newGame.setName(to: (foundGame?.getName())!) {}
            newGame.setMode(to: (foundGame?.getMode())!)
            newGame.setLocation(to: (foundGame?.getLocation())!)
            newGame.setHasPassword(to: (foundGame?.getHasPassword())!)

            // set the currentGame to the built game
            if gameState.setCurrentGame(to: newGame) {}
        }
        
        // we want to set the passed in boundaries as well
        if let bounds = info["Boundaries"] as? [[String: Any]] {
            var gameBounds: [MKMapPoint] = []
            for bound in bounds {
                if let x = bound["X"] as? Double, let y = bound["Y"] as? Double {
                    let newBound = MKMapPoint(x: x, y: y)
                    gameBounds.append(newBound)
                }
            }
            gameState.getCurrentGame().setBoundaries(gameBounds)
        }
        
        // finally, we also want to set the list of players in the game
        if let players = info["PlayerList"] as? [[String: Any]] {
            for player in players {
                if let name = player["Name"] as? String, let icon = player["Icon"] as? String {
                    gameState.getCurrentGame().addPlayer(toGame: Player(name: name, icon: icon))
                }
            }
            return true
        }
    }
    return false
}

/* The parser for the JoinGameError message */
func parseJoinGameError(data: [String: Any]) -> Bool {
    return gameState.setCurrentGame(to: Game())
}

/* The parser for the LeaveGame message */
func parseLeaveGame() -> Bool {
    return true   //TODO have the player leave the game
}

/* The parser for the GameStartInfo message */
func parseGameStartInfo(data: [String: Any]) -> Bool {
    
    if let info = data["Data"] as? [String: Any] {
        if let players = info["PlayerList"] as? [[String: Any]] {
            for player in players {
                // we want to update the playerList of the game, i.e. we want to make sure that
                // every member of the game has the completed playerList
                if let name = player["Name"] as? String, let team = player["Team"] as? String {
                    let index = gameState.getCurrentGame().findPlayerIndex(name: name)
                    if index > -1 {
                        gameState.getCurrentGame().getPlayers()[index].setTeam(to: team)
                        if gameState.getCurrentGame().getPlayers()[index].getName() == gameState.getUser().getName() {
                            gameState.getUser().setTeam(to: team)
                        }
                    }
                }
            }
        }
        
        // we also have to set information on the Game's objectives
        if let objectives = info["Objectives"] as? [[String: Any]] {
            for objective in objectives {
                if let loc = objective["Location"] as? [String: Any], let radius = objective["Radius"] as? Double, let id = objective["ID"] as? String {
                    if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                        gameState.getCurrentGame().addObjective(toObjective: Objective(x: x, y: y, radius: radius, id: id))
                    }
                }
            }
        }
        
        // additionally, we store information on every Pickup that exists within the game
        if let pickups = info["Pickups"] as? [[String: Any]] {
            var gamePickups: [Pickup] = []
            for pickup in pickups {
                if let loc = pickup["Location"] as? [String: Any], let type = pickup["Type"] as? String, let amount = pickup["Amount"] as? Int {
                    if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                        let point = MKMapPoint(x: x, y: y)
                        gamePickups.append(Pickup(loc: point, type: type, amount: amount))
                    }
                }
            }
            gameState.getCurrentGame().setPickups(toPickup: gamePickups)
        }
        
        // furthermore, we have to set the input information on the Game's Red Base and Blue Base
        if let bBase = info["BlueBase"] as? [String: Any], let rBase = info["RedBase"] as? [String: Any] {
            if let bLoc = bBase["Location"] as? [String: Any], let bRad = bBase["Radius"] as? Double, let rLoc = rBase["Location"] as? [String: Any], let rRad = rBase["Radius"] as? Double {
                if let bX = bLoc["X"] as? Double, let bY = bLoc["Y"] as? Double, let rX = rLoc["X"] as? Double, let rY = rLoc["Y"] as? Double {
                    gameState.getCurrentGame().setBlueBaseLoc(to: MKMapPoint(x: bX, y: bY))
                    gameState.getCurrentGame().setBlueBaseRad(to: bRad)
                    gameState.getCurrentGame().setRedBaseLoc(to: MKMapPoint(x: rX, y: rY))
                    gameState.getCurrentGame().setRedBaseRad(to: rRad)
                }
            }
        }
        
        // we also want to se the bounds again so that the host can get the array of
        // boundaries in the proper order, since the host never receives a gameInfo message
        if let bounds = info["Boundaries"] as? [[String: Any]] {
            var gameBounds: [MKMapPoint] = []
            for bound in bounds {
                if let x = bound["X"] as? Double, let y = bound["Y"] as? Double {
                    let newBound = MKMapPoint(x: x, y: y)
                    gameBounds.append(newBound)
                }
            }
            gameState.getCurrentGame().setBoundaries(gameBounds)
        }
        
        // finally, we want to set the StartTime of the game
        if let startTime = info["StartTime"] as? String {
            gameState.getCurrentGame().setStartTime(to: startTime)
        }
        return true
    }
    return false
}

/* The parser for the GameUpdate message */
func parseGameUpdate(data: [String: Any]) -> Bool {
    
    if let info = data["Data"] as? [String: Any] {
        if let players = info["PlayerList"] as? [[String: Any]], let points = info["Points"] as? [String: Any], let objectives = info["Objectives"] as? [[String: Any]] {
            
            // Update the liest of players in case anyone left the game
            for player in players {
                if let name = player["Name"] as? String, let orientation = player["Orientation"] as? Double, let loc = player["Location"] as? [String: Any] {
                    
                    let index = gameState.getCurrentGame().findPlayerIndex(name: name)
                    if index > -1 {
                        gameState.getCurrentGame().getPlayers()[index].setOrientation(to: Float(orientation))
                        if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                            gameState.getCurrentGame().getPlayers()[index].setLocation(to: x, to: y)
                        }
                    }
                }
            }
            
            // Update the scores of the red team and the blue team
            if let redPoints = points["Red"] as? Int, let bluePoints = points["Blue"] as? Int {
                gameState.getCurrentGame().setRedPoints(to: redPoints)
                gameState.getCurrentGame().setBluePoints(to: bluePoints)
            }
            
            // Update the status of the objectives in the game, such as which team they belong to
            for objective in objectives {
                if let id = objective["ID"] as? String, let occupants = objective["Occupying"] as? [String], let owner = objective["BelongsTo"] as? String, let progress = objective["Progress"] as? Int {
                    
                    let objIndex = gameState.getCurrentGame().findObjectiveIndex(id: id)
                    if objIndex > -1 {
                        gameState.getCurrentGame().getObjectives()[objIndex].setProgress(to: progress)
                        gameState.getCurrentGame().getObjectives()[objIndex].setOccupants(to: occupants)
                        if gameState.getCurrentGame().getObjectives()[objIndex].getOwner() != owner {
                            gameState.getCurrentGame().getObjectives()[objIndex].setOwner(to: owner)
                            gameState.getCurrentGame().getObjectives()[objIndex].setRedraw(to: true)
                        }
                    }
                }
            }
        }
        return true
    }
    return false
}

/* The parser for the GameOver message */
func parseGameOver(data: [String: Any]) -> Bool {
    // no data to parse, we just need for the player to receive a notification that the game is over
    gameState.getCurrentGame().setGameOver(to: true)
    return true
}

/* The parser for the StatusUpdate message */
func parseStatusUpdate(data: [String: Any]) -> Bool {
    // updates the player's status from the given message
    if let status = data["Data"] as? String {
        gameState.setUserStatus(to: status)
        return true
    }
    return false
}

/* The parser for the VitalsUpdate message */
func parseVitalsUpdate(data: [String: Any]) -> Bool {
    if let info = data["Data"] as? [String: Any] {
        if let health = info["Health"] as? Int, let armor = info["Armor"] as? Int {
            //update the players health and armor accordingly
            gameState.setUserHealth(to: health)
            gameState.setUserArmor(to: armor)
            return true
        }
    }
    return false
}

/* The parser for the AcquireWeapon message */
func parseAcquireWeapon(data: [String: Any]) -> Bool {
    if let weapon = data["Data"] as? String {
        // add the given weapon to the user's inventory
        gameState.getUser().addWeapon(to: weapon)
        return true
    }
    
    return false
}

/* The parser for the PickupUpdate message */
func parsePickupUpdate(data: [String: Any]) -> Bool {
    if let info = data["Data"] as? [String: Any] {
        if let loc = info["Location"] as? [String: Any], let available = info["Available"] as? Bool {
            if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                // if the pickup exists in our list of pickups, we want to update its status
                let index = gameState.getCurrentGame().findPickupIndex(x: x, y: y)
                if index > -1 {
                    gameState.getCurrentGame().getPickups()[index].setAvailability(to: available)
                    return true
                }
            }
        }
    }
    
    return false
}

/* MsgToServer
 * A class that stores data for messages sent to the server, and implements a method to convert them to json
 */
class MsgToServer {
    private var action: String
    /* possible message actions:
        case CreateGame, ShowGames, ShowGameInfo, JoinGame, StartGame, LocationUpdate, Fire
    */

    private var data: [String: Any]

    /* toJson(): convert action type and data array into server-readable json */
    func toJson() -> Data {
        let message = ["Action": action, "Data": data] as [String : Any]
        let retval = try! JSONSerialization.data(withJSONObject: message, options: JSONSerialization.WritingOptions.prettyPrinted)
        return retval
    }

    init(action: String, data: [String: Any]) {
        self.action = action
        self.data = data
    }
}

private func boundariesToArray(boundaries: [MKMapPoint]) -> [[String: Any]] {
    let bound1 = ["X": boundaries[0].x, "Y": boundaries[0].y]
    let bound2 = ["X": boundaries[1].x, "Y": boundaries[1].y]
    let bound3 = ["X": boundaries[2].x, "Y": boundaries[2].y]
    let bound4 = ["X": boundaries[3].x, "Y": boundaries[3].y]
    let retval = [bound1, bound2, bound3, bound4]
    return retval
}

/* Message generators: the following functions generate messages that can be directly sent to the server via Connection.sendData()*/

func RegisterPlayerMsg() -> Data {
    let payload = ["Name": gameState.getUserName(), "Icon": gameState.getUserIcon()] as [String : Any]
    return MsgToServer(action: "RegisterPlayer", data: payload).toJson()
}
func CreateGameMsg(game: Game) -> Data {
    let timelimit = game.getTimeLimit()
    let payload = ["Name": game.getName()!, "Password": game.getPassword() ?? "", "Description": game.getDescription(), "PlayerLimit": game.getMaxPlayers(), "PointLimit": game.getMaxPoints(), "TimeLimit": timelimit, "Mode": game.getMode().rawValue, "Boundaries": boundariesToArray(boundaries: game.getBoundaries()), "NumCP": game.getMaxObjectives()] as [String : Any]
    return MsgToServer(action: "CreateGame", data: payload).toJson()
}
func ShowGamesMsg() -> Data {
    return MsgToServer(action: "ShowGames", data: [:]).toJson()
}
func ShowGameInfoMsg(IDtoShow: String) -> Data {
    return MsgToServer(action: "ShowGameInfo", data: ["GameID": IDtoShow]).toJson()
}
func JoinGameMsg(IDtoJoin: String, password: String) -> Data {
    return MsgToServer(action: "JoinGame", data: ["GameID": IDtoJoin, "Password": password]).toJson()
}
func LeaveGameMsg() -> Data {
    gameState.getUser().setStatus(to: "")
    return MsgToServer(action: "LeaveGame", data: [:]).toJson()
}
func StartGameMsg() -> Data {
    return MsgToServer(action: "StartGame", data: [:]).toJson()
}
func LocUpMsg() -> Data {
    let location = ["X": gameState.getUserLocation().coordinate.latitude, "Y": gameState.getUserLocation().coordinate.longitude]
    let payload = ["Location": location, "Orientation": gameState.getUserOrientation()] as [String : Any]
    return MsgToServer(action: "LocationUpdate", data: payload).toJson()
}
func FireMsg() -> Data {
    let weapon = gameState.getUserWeapon().name
    let direction = gameState.getUserOrientation()
    let payload = ["Weapon": weapon, "Direction": direction] as [String: Any]
    return MsgToServer(action: "Fire", data: payload).toJson()

}
