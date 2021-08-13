import React from 'react'
import {makeStyles} from '@material-ui/core'

import './index.less'

const useStyle = makeStyles((theme) => ({
    root: {
        width: '100%',
        height: '90vh',
        padding: '0 10%',
        margin: '0 auto',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignContent: 'center',
        textAlign: 'center'
    }
}))

export default function index() {
    const classes = useStyle()
    return (
        <div className={classes.root}>
            <h1>YUKI</h1>
            <br />
            <p>A proxy reverse gateway for cqhttp</p>
            <br />
            <span style={{
                fontSize: '12px',
                color: '#b2b2b2'
            }}>developing...</span>
        </div>
    )
}
