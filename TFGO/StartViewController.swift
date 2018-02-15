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

class StartViewController: UIViewController, UICollectionViewDelegate, UICollectionViewDataSource, UITextFieldDelegate {
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return icons.count
    }
    
    private var icons = ["🦄","🦊","🐙","🐼","🐵","🐥","🐳","🐸","🐧","🦀","🐍","🐡","🐚","🦁","🐨","🦇","🐶","🐭","🦋","🐝","🤔","😎","🤑","😇","😜","😴","😱","😍","🙃","🤠","👻","💀","🤖","💩","😈","🍁","🌵","🍄","🌺","🥀","🌪","💦","🌈","❄️","☃️","🍎","🍑","🥑","🍌","🍉","🌽","🍑","🥝","🌶","🍒","⚽️","🏀","🏈","⚾️","🎾"]
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "icon", for: indexPath) as! IconViewCell
        
        cell.label.backgroundColor = randomColor()
        cell.label.text = icons[indexPath.row]
        return cell
    }
    
    @IBOutlet weak var iconCollection: UICollectionView!
    
    func collectionView(_ collectionView: UICollectionView, didSelectItemAt indexPath: IndexPath) {
        let cell = iconCollection.cellForItem(at: indexPath) as! IconViewCell
        gameState.setUserIcon(to: cell.label.text!)
    }
    
    private func randomColor() -> UIColor {
        let color = UIColor(red: CGFloat(6), green: CGFloat(96), blue: CGFloat(200), alpha: 1)
        return color
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        return gameState.getUser().isValid()
    }
    
    @IBOutlet weak var nameField: UITextField!
    
    func textField(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool {
        var name = gameState.getUserName()
        name.append(string)
        gameState.setUserName(to: name)
        return true
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        return textField.resignFirstResponder()
    }
    
    @IBAction func playerAlert(_ sender: UIButton) {
        if gameState.getUserName() == "" {
            let alertController = UIAlertController(title: "Invalid Username", message:
                "Please input a username", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "I didn't mean it", style: UIAlertActionStyle.default,handler: nil))
            self.present(alertController, animated: true, completion: nil)
        } else if gameState.getUserIcon() == "" {
            let alertController = UIAlertController(title: "Invalid User Icon", message:
                "Please select an icon", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "I will do better next time", style: UIAlertActionStyle.default,handler: nil))
            self.present(alertController, animated: true, completion: nil)
        }
    }
    
}
