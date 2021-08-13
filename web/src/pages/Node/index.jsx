import { makeStyles } from '@material-ui/core'
import React, { useState, useEffect } from 'react'
import {
  Button,
  Select, MenuItem, FormControl, InputLabel,
} from '@material-ui/core'
import { confirmAlert } from 'react-confirm-alert'
import cogoToast from 'cogo-toast'

import {
  NodeList,
} from './components/NodeList'

import {
  get, post
} from '../../request'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: '0 10%',
    fontSize: 25,
  },
  formControl: {
    minWidth: 200,
    '& .MuiButton-root': {
      marginTop: 10,
    },
  }
}))

export default function index() {
  const classes = useStyles()

  const [nodes, setNodes] = useState([])

  const [method, setMethod] = useState('random')

  const handleMethod = (e) => {
    setMethod(e.target.value)
  }

  const updateMethod = async () => {
    const res = await post('/api/node/modifySelector', {
      'func_name': method
    })
    if (res.code === 0) {
      cogoToast.success('Update success')
    }
  }

  const addProxyNode = async (remote) => {
    const res = await post('/api/node/add', {
      "remote_addr": remote
    })
    if (res.code === 0) {
      setNodes((nodes) => [...nodes, res.data])
      cogoToast.success('Add proxy node success')
    }
  }

  const removeProxyNode = (remote) => {
    confirmAlert({
      title: 'Confirm to delete',
      message: 'Are you sure to do this?',
      buttons: [
        {
          label: 'Yes',
          onClick: async () => {
            const res = await post('/api/node/remove', {
              "remote_addr": remote
            })
            if (res.code === 0) {
              setNodes((nodes) => nodes.filter(item => item.remote_addr !== remote))
              cogoToast.success('Remove proxy node success')
            }
          }
        },
        {
          label: 'No',
        }
      ]
    });

  }

  useEffect(async () => {
    // fetch from backend
    const res = await get('/api/node/getAll')
    setNodes(res.data.nodes)
    setMethod(res.data.method)
  }, [])


  return (
    <div className={classes.root}>
      <br />
      <h4 className="title">代理节点配置</h4>
      <br />
      {/* Load balance method */}
      <FormControl variant="outlined" className={classes.formControl}>
        <InputLabel id="demo-simple-select-outlined-label">Load balance method</InputLabel>
        <Select
          labelId="demo-simple-select-outlined-label"
          id="demo-simple-select-outlined"
          value={method}
          onChange={handleMethod}
          label="Load balance method"
        >
          <MenuItem value={'random'}>Random</MenuItem>
          <MenuItem value={'round_robin'}>Round_Robin</MenuItem>
          <MenuItem value={'weight'}>Weight</MenuItem>
          {/* <MenuItem value={'hash'}>Hash</MenuItem> */}
        </Select>
        <Button variant="outlined" color="primary" onClick={updateMethod}>
          Update
        </Button>
      </FormControl>
      <br />
      <br />
      {/* List */}
      <NodeList
        nodes={nodes}
        addProxyNode={addProxyNode}
        removeProxyNode={removeProxyNode}
      />
    </div>
  )
}
