//
//  InventoryViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/20/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class PickupViewCell: UICollectionViewCell {
    @IBOutlet weak var image: UIImageView!
    @IBOutlet weak var button: UIButton!
}

class InventoryViewController: UIViewController, UICollectionViewDelegate, UICollectionViewDataSource{

    private var player = gameState.getUser()
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return player.getWeaponsList().count
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "Item", for: indexPath) as! PickupViewCell
        let name = player.getWeaponsList()[indexPath.row]
        let color = randomColor()
        cell.image.image = UIImage(named: name)
        cell.button.setTitle(name, for: .normal)
        cell.button.setTitleColor(color, for: .normal)
        cell.backgroundColor = color
        cell.layer.cornerRadius = 8.0
        cell.clipsToBounds = true
        return cell
    }
    
    @IBAction func selectItem(_ sender: UIButton) {
        let actionController = UIAlertController(title: nil, message:
            sender.currentTitle ?? "halp", preferredStyle: UIAlertControllerStyle.actionSheet)
        actionController.addAction(UIAlertAction(title: "Equip", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            self.player.setWeapon(to: sender.currentTitle ?? "BeeSwarm")
            self.performSegue(withIdentifier: "back", sender: nil)
        }))
        actionController.addAction(UIAlertAction(title: "Discard", style: UIAlertActionStyle.default,handler: nil))
        actionController.addAction(UIAlertAction(title: "Consume?!?", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            print("You Chose Poorly")
            gameState.setUserHealth(to: 0)
            self.performSegue(withIdentifier: "back", sender: nil)
        }))
        actionController.addAction(UIAlertAction(title: "Cancel", style: UIAlertActionStyle.cancel,handler: nil))
        self.present(actionController, animated: true, completion: nil)
    }
}
