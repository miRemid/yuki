import React, { useEffect, useState } from 'react'
import TextField from '@material-ui/core/TextField'
import { makeStyles } from '@material-ui/core/styles'
import { Button } from '@material-ui/core'
import cogoToast from 'cogo-toast'
import { get, update } from '../../request'


const useStyles = makeStyles((theme) => ({
    root: {
        padding: '0 10%',
        fontSize: 25,
        '& .MuiTextField-root': {
            'margin-right': 8,
            'margin-bottom': 8,
            width: 200,
        },
    },
    formBox: {
        display: 'flex'
    }
}));
export default function System() {

    const [data, setData] = useState({
        'cqhttp': '',
        'admin': '',
        'secret': '',
        'prefix': '',
        'format': ''
    })

    const updateConfig = async () => {
        try {
            const arrs = data.prefix.split(';')
            arrs.forEach((item, i) => {
                if (item === '') {
                    arrs.splice(i, 1)
                }
            })
            const res = await update('/api/config', {
                'cqhttp_address': data.cqhttp,
                'admin_qq': data.admin,
                'secret': data.secret,
                'prefix': arrs,
                'format': data.format
            })
            if (res.code === 0) {
                cogoToast.success('Update config success')
            } else {

            }
        } catch (error) {
            console.log(error)
        }
    }

    const handleInput = (name) => {
        return e => {
            const value = e.target.value
            console.log(name, value);
            switch (name) {
                case 'cqhttp':
                    setData({ ...data, 'cqhttp': value })
                    break
                case 'admin':
                    setData({ ...data, 'admin': value })
                    break
                case 'secret':
                    setData({ ...data, 'secret': value })
                    break
                case 'prefix':
                    setData({ ...data, 'prefix': value })
                    break
                case 'format':
                    setData({ ...data, 'format': value })
                    break
                default:
                    break;
            }
            console.log(data);
        }
    }

    const classess = useStyles()

    useEffect(() => {
        const fetchData = async () => {
            const res = await get('/api/config')
            if (res.code === 0) {
                setData({
                    ...data,
                    'cqhttp': res.data.cqhttp_address,
                    'admin': res.data.admin_qq,
                    'secret': res.data.secret,
                    'prefix': res.data.prefix.join(';'),
                    'format': res.data.format
                })
            }
        }
        fetchData()
    }, [])

    return (
        <div className={classess.root}>
            <br />
            <h4 className="title">系统配置</h4>
            {/* This is the system's configuration form */}
            <br />
            <div className={classess.formBox}>
                {
                    Object.keys(data).map((v, i) => {
                        return <TextField
                            key={v}
                            id={v}
                            label={
                                v === 'cqhttp' ? 'CQHTTP_ADDRESS' :
                                    v === 'admin' ? 'ADMIN_QQ' :
                                        v === 'prefix' ? 'PREFIX' :
                                            v === 'secret' ? 'SECRET' :
                                                'Command 404 format'
                            }
                            variant="outlined"
                            color="primary"
                            onChange={handleInput(v)}
                            value={data[v]}
                            style={{
                                width: v === 'format' ? '300px' : '200px'
                            }}
                        />
                    })
                }
                {/* <TextField
                    id="cqhttp_address"
                    label="CQHTTP_ADDRESS"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('cqhttp')}
                    value={data.cqhttp}
                />
                <TextField
                    id="admin_qq"
                    label="ADMIN_QQ"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('admin')}
                    value={data.admin}
                />
                <TextField
                    id="secret"
                    label="SECRET"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('secret')}
                    value={data.secret}
                />
                <TextField
                    id="prefix"
                    label="PREFIX"
                    variant="outlined"
                    color="secondary"
                    onChange={handleInput('prefix')}
                    value={data.prefix}
                /> */}
            </div>
            <Button variant="outlined" color="primary" onClick={updateConfig}>
                Update
            </Button>
        </div>
    )
}
