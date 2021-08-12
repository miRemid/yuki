import React, {
    useEffect
} from 'react'
import {
    NavLink
} from 'react-router-dom'

import './index.less'

export default function Header() {

    const nav_items = [
        {
            name: 'Home',
            path: '/'
        },
        {
            name: 'Node',
            path: '/node'
        },
        {
            name: 'Rule',
            path: '/rule'
        },
        {
            name: 'System',
            path: '/system'
        }
    ]

    // mount
    useEffect(() => {
        // add windows listener
    }, [])

    return (
        <header className="header">
            {/* This is the Login */}
            <div className="logo">
                YUKI
            </div>
            <nav >
                <ul className="navs">
                    {
                        nav_items.map(item => {
                            return <li className="nav-item" key={item.name}>
                                <NavLink exact to={item.path} activeClassName="active">{item.name}</NavLink>
                            </li>
                        })
                    }
                </ul>
            </nav>
        </header>
    )
}