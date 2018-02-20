//
//  InventoryViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/20/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class PickupViewCell: UICollectionViewCell {
    
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
        return cell
    }
}
