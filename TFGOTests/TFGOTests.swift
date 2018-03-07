//
//  TFGOTests.swift
//  TFGOTests
//
//  Created by Sam Schlang on 1/30/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import XCTest
@testable import TFGO

class TFGOTests: XCTestCase {
    
    override func setUp() {
        super.setUp()
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }
    
    func testGame() {
        // Tests Game object setup for invalid inputs by user
        let game = Game()
        XCTAssertFalse(game.isValid())
        game.setID(to: "abc")
        XCTAssertFalse(game.setName(to: ""))
        XCTAssertFalse(game.setName(to: "abcdefghijklmnopqrstuvwxyzabcd"))
        XCTAssertTrue(game.setName(to: "def"))
        XCTAssertFalse(game.setDescription(to: "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"))
        XCTAssertTrue(game.setDescription(to: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"))
        XCTAssertTrue(game.isValid())
        
    }
    
    func testParse() {
        // Tests Parse Functions
        //let game = Game()
        //let testAvailGame = "\"Data\": [{\"ID\": \"Abc\",\"Name\": \"Abc\",\"Mode\": \"SingleCapture\",\"Location\": {\"X\": 1.23, \"Y\": 12.3},\"PlayerList\": [{\"Name\": \"Abc\", \"Icon\": \"Abc\"}, {\"Name\": \"Abc\", \"Icon\": \"Abc\"}]}"
        let location = ["X":1.23, "Y":12.3]
        let player1 = ["ID": "123", "Name": "ABC", "Icon": "Abc", "Orientation": Float(1.23), "Location" : location] as [String : Any]
        let player2 = ["ID": "456", "Name": "BCD", "Icon": "Bcd", "Orientation": Float(1.23), "Location" : location] as [String : Any]
        let playerList = [player1, player2]
        let testAvailGame = ["Data": [["ID": "ABC", "Name": "DEF", "Mode": "SingleCapture", "Location": location, "PlayerList": playerList, "HasPassword": true]]] as [String : Any]
        let testAvailGameInvalid = ["Data": ["ID": "ABC", "Name": "DEF", "Mode": "SingleCapture", "Location": location, "PlayerList": playerList]] as [String : Any]
        //        let testGameInfo = "\"Data\": {\"Description\": \"Abc\",\"PlayerLimit\": 123,\"PointLimit\": 123,\"TimeLimit\": \"1h23m3s\",\"Boundaries\": [{\"X\": 12.3, \"Y\": 1.23}, {\"X\": 23.4, \"Y\": 2.34}],\"PlayerList\": [{\"Name\": \"Abc\", \"Icon\": \"Abc\"}, {\"Name\": \"Abc\", \"Icon\": \"Abc\"}],}"
        let boundary = [location, location, location, location]
        let testGameInfo = ["Data": ["ID": "ABC", "Description": "Abc", "PlayerLimit": 50, "PointLimit": 100, "TimeLimit": 80, "Boundaries": boundary, "PlayerList": playerList]]
        
        let testGameInfoInvalid = ["Data": ["PlayerLimit": 50, "PointLimit": 100, "TimeLimit": "1h23m3s", "Boundaries": boundary]]
        
        //        let testJoinGameError = "\"Data\": \"GameFull\""
        let testJoinGameErrorFull = ["Data": "GameFull"]
        let testJoinGameErrorStarted = ["Data": "GameStarted"]
        let player1Start = ["Name": "ABC", "Team": "Red", "ID": "123"]
        let player2Start = ["Name": "BCD", "Team": "Blue", "ID": "456"]
        let playerListStart = [player1Start, player2Start]
        let pickUps = ["Location": location, "Type": "Health", "Amount": 40] as [String: Any]
        let pickUpsArray = [pickUps]
        let pickUps2 = ["Location": location, "Available": false] as [String: Any]
        let Objective1 = ["Location": location, "Radius": 1.23, "ID": "HHH"] as [String : Any]
        let testGameStartInfo = ["Data": ["PlayerList": playerListStart, "Pickups": pickUpsArray, "Objectives": [Objective1], "StartTime": "2017-02-14 17:24:56"]] as [String : Any]
        let testGameStartInfoInvalid = ["Data1": ["PlayerList": playerListStart, "Objectives": [Objective1], "StartTime": "2017-02-14 17:24:56"]] as [String : Any]
        let redBase = ["Location": location, "Radius": 1.23] as [String : Any]
        
        
        let Objective0 = ["Location": location, "ID": "HHH", "Occupying": ["ABC"], "BelongsTo": "Red", "Progress" : 75] as [String : Any]
        let Objective = [Objective0]
        let points = ["Red": 100, "Blue": 105] as [String: Any]
        let testGameUpdate = ["Data": ["PlayerList": playerList, "Points":points, "RedBase": redBase, "BlueBase": redBase, "Objectives" : Objective, "StartTime": "2017-02-14 17:24:56"]]
        let weapon = ["Data": "Blaster", "Type": "AcquireWeapon"]
        
        
        //        let testGameStartInfo = "\"Data\": {\"PlayerList\": [{\"Name\": \"Abc\", \"Team\": \"Red\"}, {\"Name\": \"bcd\", \"Team\": \"Blue\"}],\"RedBase\": {\"Location\": {\"X\": 1.23, \"Y\": 12.3}, \"Radius\": 1.23},\"BlueBase\": {\"Location\": {\"X\": 1.23, \"Y\": 12.3}, \"Radius\": 1.23},\"Objectives\":[{\"Location\": {\"X\": 1.23, \"Y\": 12.3}, \"Radius\": 1.23}],\"StartTime\": \"2017-02-14 17:24:26\""
        //        let testGameUpdate = "\"Data\": {\"PlayerList\": [{\"Name\": \"Abc\", \"Orientation\": 1.23, \"Location\": {\"X\": 1.23, \"Y\": 12.3}, {\"Name\": \"Bcd\", \"Orientation\": 1.23, \"Location\": {\"X\": 1.23, \"Y\": 12.3}],\"Points\": {\"Red\": 100, \"Blue\": 250},\"Objectives\": [{\"Location\": {\"X\": 1.23, \"Y\": 12.3},\"Occupying\": [\"Abc\", \"Def\"],\"BelongsTo\": \"Red\",\"Progress\": 75}"
        //        let testOutofBound = "\"Data\": \"OutOfBounds\""
        //        let testBackInBounds = "\"Data\": \"BackInBounds\""
        //        let testRespawn = "\"Data\": \"Respawn\""
        //        let testRespawned = "\"Data\": \"Respawned\""
        //        let testTakeHit = "\"Data\": {\"Health\": 123, \"Armor\": 123}"
        if parseAvailableGames(data: testAvailGame) {}
        if parseGameInfo(data: testGameInfo) {}
        //XCTAssertTrue(parseJoinGameError(data: testGameInfo))
        XCTAssertTrue(parseGameStartInfo(data: testGameStartInfo))
        
        
        /***
         
         For the testing we only tested with valid data, given that all the data we
         receive from the server will be filtered already; otherwise it won't get
         passed the first step
         
         ***/
        
        XCTAssertTrue(gameState.getFoundGames()[0].getID() == "ABC")
        XCTAssertFalse(gameState.getFoundGames()[0].getID() == "BCD")
        XCTAssertTrue(gameState.getFoundGames()[0].getName() == "DEF")
        XCTAssertTrue(gameState.getFoundGames()[0].getMode().rawValue == "SingleCapture")
        XCTAssertTrue(gameState.getFoundGames()[0].getLocation().x == 1.23)
        XCTAssertFalse(gameState.getFoundGames()[0].getLocation().x == 12.3)
        XCTAssertTrue(gameState.getCurrentGame().getDescription() == "Abc")
        XCTAssertTrue(gameState.getCurrentGame().getMaxPoints() == 100)
        XCTAssertTrue(gameState.getCurrentGame().getMaxPlayers() == 50)
        //XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getXLoc() == 1.23)
        //XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getYLoc() == 12.3)
        //XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getRadius() == 1.23)
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[0].getName() == "ABC")
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[1].getName() == "BCD")
        XCTAssertTrue(gameState.getCurrentGame().getStartTime()[0] == "2017")
        XCTAssertTrue(gameState.getCurrentGame().getStartTime()[1] == "02")
        XCTAssertTrue(gameState.getCurrentGame().getStartTime()[2] == "14")
        XCTAssertTrue(gameState.getCurrentGame().getStartTime()[3] == "17")
        XCTAssertFalse(gameState.getCurrentGame().getStartTime()[0] == "2018")
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[0].getName() == "ABC")
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[1].getName() == "BCD")
        
        if parseGameUpdate(data: testGameUpdate) {}
        XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getProgress() == 75)
        XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getOwner() == "Red")
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[0].getOrientation() == 1.23)
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[1].getLocation().coordinate.latitude == 1.23)
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[1].getLocation().coordinate.longitude == 12.3)
        XCTAssertTrue(gameState.getCurrentGame().getRedPoints() == 100)
        XCTAssertTrue(gameState.getCurrentGame().getBluePoints() == 105)
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[0].getID() == "123")
        XCTAssertTrue(gameState.getCurrentGame().getPlayers()[1].getID() == "456")
        XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getOccupants()[0] == "ABC")
        XCTAssertTrue(gameState.getCurrentGame().getObjectives()[0].getOccupants().count == 1)
        
        if parseStatusUpdate(data: ["Data": "Dead"]) {}
        XCTAssertTrue(gameState.getUser().getStatus() == "Dead")
        if parseStatusUpdate(data: ["Data": "Alive"]) {}
        XCTAssertTrue(gameState.getUser().getStatus() == "Alive")
        
        if parseVitalsUpdate(data: ["Data" : ["Health": 58, "Armor": 44]]) {}
        XCTAssertTrue(gameState.getUser().getHealth() == 58)
        XCTAssertTrue(gameState.getUser().getArmor() == 44)
        
        XCTAssertTrue(gameState.getCurrentGame().getPickups()[0].getAvailability() == true)
        if parsePickupUpdate(data: ["Data" : pickUps2]) {}
        XCTAssertTrue(gameState.getCurrentGame().getPickups()[0].getAvailability() == false)
        
        if parsePlayerID(data: ["Data" : "qtpi55"]) {}
        XCTAssertTrue(gameState.getUser().getID() == "qtpi55")
        if parsePlayerID(data: ["Data" : "mrtw33tums"]) {}
        XCTAssertTrue(gameState.getUser().getID() == "mrtw33tums")
        
        XCTAssertTrue(weaponByName(name: "Blaster").name == "Blaster")
        
        /***
         
         Here we picked some casese where the parse function will "break" in the
         very beginning so here they will all return false
         
         ***/
        
        XCTAssertFalse(parseAvailableGames(data: testAvailGameInvalid)) // missing the right struct
        XCTAssertFalse(parseGameInfo(data: testGameInfoInvalid)) // missing playerlist
        
        XCTAssertFalse(parseGameStartInfo(data: testGameStartInfoInvalid)) // missing the data bracket
        //XCTAssertTrue(gameState.getFoundGames()[0].get == "ABC")
        
    }
    
    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }
    
}
