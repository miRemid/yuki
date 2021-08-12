import React, { useEffect, useState } from 'react'
import TextField from '@material-ui/core/TextField'
import { makeStyles } from '@material-ui/core/styles'
import { Button } from '@material-ui/core'

import './index.less'
import { get, post } from '../../request'
const useStyles = makeStyles((theme) => ({
    root: {
        '& .MuiTextField-root': {
            'margin-right': 8,
            'margin-bottom': 8,
            width: 200,
        },
    },
}));
export default function System() {

    const [cqhttp, setCqhttp] = useState('')
    const [admin, setAdmin] = useState('')
    const [secret, setSecret] = useState('')
    const [prefix, setPrefix] = useState('')

    const updateConfig = async () => {
        try {
            const arrs = prefix.split(';')
            arrs.forEach((item, i) => {
                if (item === '') {
                    arrs.splice(i, 1)
                }
            })
            const data = {
                'cqhttp_address': cqhttp,
                'admin_qq': admin,
                'secret': secret,
                'prefix': arrs,
            }
            console.log(data)
        } catch (error) {
            console.log(error)
        }
    }

    const handleInput = (name) => {
        return e => {
            const value = e.target.value
            console.log('name: ', name)
            console.log('value: ', value)
            switch (name) {
                case 'cqhttp':
                    setCqhttp(value)
                    break
                case 'admin':
                    setAdmin(value)
                    break
                case 'secret':
                    setSecret(value)
                    break
                case 'prefix':

                    setPrefix(value)
                default:
                    break;
            }
        }
    }

    const classess = useStyles()

    useEffect(() => {
        const fetchData = async () => {
            const response = await get('/api/config/get')
            setAdmin(response.data.data.admin_qq)
            setCqhttp(response.data.data.cqhttp_address)
            setPrefix(response.data.data.prefix.join(';'))
            setSecret(response.data.data.secret)
        }
        fetchData()
    }, [])

    return (
        <div className="container system">
            <br />
            <h4 className="title">系统配置</h4>
            {/* This is the system's configuration form */}
            <br />
            <div className={classess.root + " form-box"}>
                <TextField
                    id="cqhttp_address"
                    label="CQHTTP_ADDRESS"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('cqhttp')}
                    value={cqhttp}
                />
                <TextField
                    id="admin_qq"
                    label="ADMIN_QQ"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('admin')}
                    value={admin}
                />
                <TextField
                    id="secret"
                    label="SECRET"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('secret')}
                    value={secret}
                />
                <TextField
                    id="prefix"
                    label="PREFIX"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('prefix')}
                    value={prefix}
                />
            </div>
            <Button variant="outlined" color="primary" onClick={updateConfig}>
                Update
            </Button>
        </div>
    )
}
