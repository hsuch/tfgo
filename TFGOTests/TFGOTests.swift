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
    
    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }
    
}
