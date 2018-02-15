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
    
    override func viewDidLoad() {
        super.viewDidLoad()

        let game = gameState.getCurrentGame()
        table.numberOfRows(inSection: game.getPlayers().count)

        // Do any additional setup after loading the view.
    }
    
    func tableView(_ tableView: UITableView, willDisplay cell: UITableViewCell, forRowAt indexPath: IndexPath) {
        let cell = tableView.dequeueReusableCell(withIdentifier: "Player", for: indexPath) as! WaitingViewCell
        let player = gameState.getCurrentGame().getPlayers()[indexPath.row]
        cell.icon.text = player.getIcon()
        cell.icon.backgroundColor = #colorLiteral(red: 1, green: 0.5212053061, blue: 1, alpha: 1)
        cell.name.text = player.getName()
    }

    @IBOutlet weak var table: UITableView!
    
    private var playerTimer = Timer.scheduledTimer(timeInterval: 1, target: self, selector: #selector(WaitingViewController.updatePlayers), userInfo: nil, repeats: true)
    
    @objc private func updatePlayers() {
        table.reloadData()
    }

}
