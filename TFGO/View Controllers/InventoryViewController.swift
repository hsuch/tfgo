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
    
    var name = ""
}

class InventoryViewController: UIViewController, UICollectionViewDelegate, UICollectionViewDataSource{

    override func viewDidLoad() {
        super.viewDidLoad()
        
    }
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return 0
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "Item", for: indexPath) as! PickupViewCell
        cell.backgroundColor = randomColor()
        cell.layer.cornerRadius = 8.0
        cell.clipsToBounds = true
        return cell
    }
    
    @IBAction func selectItem(_ sender: UIButton) {
        let pickup = sender.superview as! PickupViewCell
        let actionController = UIAlertController(title: nil, message:
            pickup.name, preferredStyle: UIAlertControllerStyle.actionSheet)
        actionController.addAction(UIAlertAction(title: "Equip", style: UIAlertActionStyle.default,handler: nil))
        actionController.addAction(UIAlertAction(title: "Discard", style: UIAlertActionStyle.default,handler: nil))
        actionController.addAction(UIAlertAction(title: "Consume?!?", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            print("You Chose Poorly")
            gameState.setUserHealth(to: 0)
        }))
        
        actionController.addAction(UIAlertAction(title: "Cancel", style: UIAlertActionStyle.cancel,handler: nil))
        self.present(actionController, animated: true, completion: nil)
    }
}
