//
//  ViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 1/30/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

@IBDesignable
class ViewController: UIViewController {
    
    @IBOutlet weak var background: UIImageView!

    override func viewDidLoad() {
        super.viewDidLoad()
        
//        background.image = UIImage(named: "dark-triangles")
        // Do any additional setup after loading the view, typically from a nib.
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "temp" {
            if gameState.getConnection().sendData(data: ShowGamesMsg()).isSuccess {
                while true {
                    if MsgFromServer().parse() {
                        if gameState.findPublicGames().count > 0 {
                            return gameState.setCurrentGame(to: gameState.findPublicGames()[0])
                        }
                    }
                }
            }
            return false
        } else {
            return true
        }
    }
    
}

