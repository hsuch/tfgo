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
        
        background.image = UIImage(named: "bg2")
    }
    
    override func viewWillAppear(_ animated: Bool) {
        self.navigationController?.setNavigationBarHidden(true, animated: true)
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        self.navigationController?.setNavigationBarHidden(false, animated: true)
    }
    
}

