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

class WaitingViewController: UIViewController, UITableViewDelegate, UITableViewDataSource {
    
    private var startGame = false
    
    private let game = gameState.getCurrentGame()
    
    @IBOutlet weak var startButton: UIButton!
    
    override func viewWillAppear(_ animated: Bool) {
        startButton.isHidden = !gameState.getUser().isHost()
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        self.table.dataSource = self;
        self.table.delegate = self;
        print("tick")
        runTimer()
        print("tock")
        // Do any additional setup after loading the view.
        print("begin")
        DispatchQueue.global(qos: .userInitiated).async {
            print("end")
            while !self.startGame {
                print("trying to handle")
                if handleMsgFromServer() { }
            }
        }
    }
    
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return game.getPlayers().count
    }
    
    
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "Player", for: indexPath) as! WaitingViewCell
        let player = gameState.getCurrentGame().getPlayers()[indexPath.row]
        cell.icon.text = player.getIcon()
        cell.name.text = player.getName()
        cell.icon.backgroundColor = randomColor()
        cell.icon.layer.cornerRadius = 8.0
        cell.icon.clipsToBounds = true
        cell.layer.cornerRadius = 8.0
        cell.backgroundColor = randomColor()
        return cell
    }

    @IBOutlet weak var waitingLabel: UILabel!
    @IBOutlet weak var table: UITableView!
    
    private var playersTimer = Timer()
    
    func runTimer() {
        playersTimer = Timer.scheduledTimer(timeInterval: 1, target: self,   selector: (#selector(checkPlayers)), userInfo: nil, repeats: true)
    }
    
    private var dots = "..."
    
    @objc private func checkPlayers() {
        self.table.reloadData()
        waitingLabel.textColor = randomColor()
        waitingLabel.text = "Waiting For Players\(dots)"
        if dots == "..." {
            dots = "."
        } else if dots == "." {
            dots = ".."
        } else {
            dots = "..."
        }
        if game.getStartTime() != [] {
            performSegue(withIdentifier: "startGame", sender: nil)
        }
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier != "startGame" {
            if gameState.getConnection().sendData(data: LeaveGameMsg()).isSuccess {
                return true
            }
        }
        if gameState.getUser().isHost(), gameState.getConnection().sendData(data: StartGameMsg()).isSuccess {
           // DispatchQueue.global(qos: .background).async {
                //while gameState.getCurrentGame().getStartTime() == "" {
            if handleMsgFromServer() {
                if gameState.getCurrentGame().getStartTime() != [] {
                    startGame = true
                    return true
                }
                
            }
                //}
            //}
            
        } else {
            startGame = true
            return true
        }
        return false
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        playersTimer.invalidate()
    }

}
