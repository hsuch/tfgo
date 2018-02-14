//
//  WaitingViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/13/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class WaitingViewCell: UITableViewCell {
    @IBOutlet weak var icon: UILabel!
    @IBOutlet weak var name: UILabel!
}

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
        let cell = tableView.dequeueReusableCell(withIdentifier: "Player", for: indexPath) as! WaitingViewCell
        let player = state?.getCurrentGame().getPlayers()[indexPath.row]
        cell.icon.text = player?.getIcon()
        cell.icon.backgroundColor = #colorLiteral(red: 1, green: 0.5212053061, blue: 1, alpha: 1)
        cell.name.text = player?.getName()
    }

    @IBOutlet weak var table: UITableView!
    

}
