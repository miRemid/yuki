import React, { useState, useEffect } from 'react'
import { makeStyles } from '@material-ui/core'

import { RuleList } from './components/RuleList'

import { get, post, del, update } from '../../request'
import cogoToast from 'cogo-toast'
import { confirmAlert } from 'react-confirm-alert'

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

  const [rules, setRules] = useState([])
  const [nodes, setNodes] = useState([])

  useEffect(async () => {
    const res = await get('/api/node')
    setNodes(res.data.nodes)
    const rules = await get('/api/rule')
    setRules(rules.data.rules)
  }, [])

  const addRule = async (cmd, remote) => {
    const res = await post('/api/rule', {
      cmd: cmd,
      remote_addr: remote
    })
    if (res.code === 0) {
      setRules((rules) => [...rules, { cmd: cmd, remote_addr: remote }])
      cogoToast.success('Add rule success')
    }
  }

  const removeRule = async (cmd) => {
    confirmAlert({
      title: 'Confirm to delete',
      message: 'Are you sure to do this?',
      buttons: [
        {
          label: 'Yes',
          onClick: async () => {
            const res = await del('/api/rule', {
              cmd: cmd
            })
            if (res.code === 0) {
              setRules((rules) => rules.filter(item => item.cmd !== cmd))
              cogoToast.success('Remove rule success')
            }
          }
        },
        {
          label: 'No'
        }
      ]
    })
  }

  return (
    <div className={classes.root}>
      <br />
      <h4 className="title">命令规则配置</h4>
      <br />
      {/* Rules List */}
      <RuleList
        nodes={nodes}
        rules={rules}
        addRule={addRule}
        removeRule={removeRule}
      />
    </div>
  )
}
