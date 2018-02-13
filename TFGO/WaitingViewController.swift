//
//  WaitingViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/13/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class WaitingViewController: UIViewController, UITableViewDelegate {

    var state: GameState?
    
    override func viewDidLoad() {
        super.viewDidLoad()

        if let game = state?.getCurrentGame() {
            table.numberOfRows(inSection: game.getPlayers().count)
        }
        // Do any additional setup after loading the view.
    }
    
    func tableView(_ tableView: UITableView, willDisplay cell: UITableViewCell, forRowAt indexPath: IndexPath) {
        tableView.dequeueReusableCell(withIdentifier: "Player", for: indexPath)
        
    }

    @IBOutlet weak var table: UITableView!
    

}
