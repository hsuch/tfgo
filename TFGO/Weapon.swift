//
//  Weapon.swift
//  TFGO
//
//  Created by Sam Schlang on 2/28/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation

class Weapon {
    var name: String
    var clipSize: Int
    var clipFill: Int
    var shotReload: Float
    var clipReload: Float
    
    init () {
        name = ""
        clipSize = 0
        clipFill = 0
        shotReload = 0
        clipReload = 0
    }
}

class Sword: Weapon {
    
    override init() {
        super.init()
        name = "Sword"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Shotgun: Weapon {
    
    override init() {
        super.init()
        name = "Shotgun"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Pistol: Weapon {
    
    override init() {
        super.init()
        name = "Pistol"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Blaster: Weapon {
    
    override init() {
        super.init()
        name = "Blaster"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Crossbow: Weapon {
    
    override init() {
        super.init()
        name = "Crossbow"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Rifle: Weapon {
    
    override init() {
        super.init()
        name = "Rifle"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Boomerang: Weapon {
    
    override init() {
        super.init()
        name = "Boomerang"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Lightsaber: Weapon {
    
    override init() {
        super.init()
        name = "Lightsaber"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Spear: Weapon {
    
    override init() {
        super.init()
        name = "Spear"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Banhammer: Weapon {
    
    override init() {
        super.init()
        name = "BanHammer"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class BeeSwarm: Weapon {
    
    override init() {
        super.init()
        name = "BeeSwarm"
        clipSize = 1337
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}
