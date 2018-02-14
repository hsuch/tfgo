//
//  StartViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/14/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class IconViewCell: UICollectionViewCell {
    @IBOutlet weak var label: UILabel!
}

class StartViewController: UIViewController, UICollectionViewDelegate, UICollectionViewDataSource {
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return icons.count
    }
    
    private var icons = ["🦄","🦊","🐙","🐼","🐵","🐥","🐳","🐸","🐧","🦀","🐍","🐡","🐚","🦁","🐨","🦇","🐶","🐭","🦋","🐝","🤔","😎","🤑","😇","😜","😴","😱","😍","🙃","🤠","👻","💀","🤖","💩","😈","🍁","🌵","🍄","🌺","🥀","🌪","💦","🌈","❄️","☃️","🍎","🍑","🥑","🍌","🍉","🌽","🍑","🥝","🌶","🍒","⚽️","🏀","🏈","⚾️","🎾"]
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "icon", for: indexPath) as! IconViewCell
        
        cell.backgroundColor = randomColor()
        cell.label.text = icons[indexPath.row]
        return cell
    }
    

    override func viewDidLoad() {
        super.viewDidLoad()
        

        // Do any additional setup after loading the view.
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    private func randomColor() -> UIColor {
        return UIColor(red: CGFloat(arc4random_uniform(100)), green: CGFloat(arc4random_uniform(100)), blue: CGFloat(arc4random_uniform(100)), alpha: 1)
    }

}
